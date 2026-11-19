## major.actalliancechampionship 模块详解分析

> 模块路径：`majorserver/module/activity/actalliancechampionship`
> 模块定位：跨服「公会联盟争霸赛」（Web Alliance Championship）活动在 **majorserver（主服务器）** 端的实现，负责：玩家/联盟报名、线路（部署）管理、赛程与对阵数据持久化、段位星数计算、结算发奖、邮件通知、与 `generalserver`（中央调度服）的 RPC 协作以及若干 GM 工具接口。

---

### 一、整体目录结构与定位

```
actalliancechampionship/
├── actalliancechampionship_def.go         // 公共常量定义（邮件延迟等）
├── iactalliancechampionship.go            // 对外接口契约（数据/模块/RPC接口）
└── actalliancechampionshipimpl/           // 实现包（业务全部在此）
    ├── actalliancechampionship_activity.go     // 活动生命周期 & 阶段定时器
    ├── actalliancechampionship_conf.go         // 配置读取 + 战力/段位/积分计算
    ├── actalliancechampionship_logic.go        // 核心业务逻辑（17 个对外函数）
    ├── actalliancechampionship_rpc.go          // 客户端 RPC 入口 (17 个接口)
    ├── actalliancechampionship_in.go           // 模块内对外暴露 + GM 命令
    ├── actalliancechampionship_togeneral.go    // M -> G 的所有 RPC 出口
    ├── actalliancechampionship_me.go           // me 消息总线监听 G -> M
    ├── actalliancechampionship_settlement.go   // 结算流程
    ├── actalliancechampionship_mail.go         // 邮件构造与发送
    ├── actalliancechampionship_rank.go         // 段位排行榜（编/解 score）
    ├── actalliancechampionship_system_obj.go   // 全局系统对象
    ├── actalliancechampionship_guild_obj.go    // 联盟对象（核心数据载体）
    ├── actalliancechampionship_team_obj.go     // 小组（六联盟）对象
    ├── actalliancechampionship_record_obj.go   // 联盟操作记录对象
    ├── actalliancechampionship_user_obj.go     // 玩家个人参与对象
    ├── actalliancechampionship_util.go         // 工具函数
    └── *_auto.go / inpool_auto.go              // 自动生成的对象池/注册代码
```

实现包内部按照 majorserver 一贯的框架使用 `TDataObjInner` + DB 表 (`tables.TActivity*`) + Refresh() 模式，保证数据「按活动期重新刷新一次」。

---

### 二、对外接口（iactalliancechampionship.go）

文件用 5 个 interface 把模块与外部解耦：

1. `IInActAllianceChampionShipData`：查找系统、联盟、小组、操作记录、用户五个数据对象。
2. `IInActAllianceChampionShipModule`：内部模块层调用 + 大量 GM 调试接口（解散、踢人、强制分组、强制开战、改段位、测试战力等）。
3. `IRpcActAllianceChampionShipModule`：客户端可调用的 17 个 pb 接口（GetInfo / GetNotSignupList / Invite / Dispatch / ReplaceRouter / GetDeployment / GetFormation / GetTeamRank / GetCourse / GetLastSimpleInfo / SwitchUserRouter / ExchangeRouter / GetOperateRecord / Rollback / GetBattleData / GetDanRank / GetSuccessGuildList）。
4. `IInRpcActAllianceChampionShipModule`：generalserver -> majorserver 的内部 RPC（分组结果 / 一轮战斗结果），代码中实际已被 me 消息系统替代，函数体保留为占位。
5. `IActAllianceChampionShip*Data`：四个 DataObj 接口（System、Guild、Team、User、GuildRecord）。

---

### 三、活动生命周期与阶段（actalliancechampionship_activity.go）

通过 `activityconfmgr.RegisterActivityExecutor("WebAllianceChampionship", …)` 注册执行器 `TActAllianceChampionShipExecutor`，实现 `CheckConf / Start / End / Close / GetRunningInfo`。

`Start()` 时按当前时间和配置一次性安装 4 类定时器：

| 定时器 | 触发时间 | 行为 |
|---|---|---|
| `beforeGroupingTimer` | 报名结束前 12 小时 | 给报名人数不足 10 人的联盟管理发提醒邮件 |
| `groupingTimer` | 匹配阶段开始 | RPC 调用 `alliancechampionship.GroupingStarByInner` 让 generalserver 进行分组 |
| `roundBattlesTimer` × 5 | 每个战斗轮开始 | RPC 调用 `RoundBattleStarByInner` 让 generalserver 计算第 N 轮三场对阵 |
| `settlementTimer` | 全部战斗结束 | 走 `settlement()` 发奖、写晋级邮件 |

阶段判定函数（同文件下）：
- `isInSignupTime` / `isInMatchTime` / `isInPrepareWarTime` / `isInPrepareBattleWarTime` / `isInAllWarTime` / `isSettlementTime`。
- 时间长度均来自配置项 `RegisterTime / MatchTime / ReadyTime / BattleTime`，轮数 = `GroupSize - 1`（默认 6 公会一组，5 轮循环赛）。

---

### 四、配置与战力公式（actalliancechampionship_conf.go）

- `getConfConfig` 读取 5 个时长 + 总轮数。
- 段位配置 `getConfDan` 返回「大段位下小段位个数 / 每小段位星数 / 大段位总数」，整体段位以 `dan*100 + dan2` 表示（例 302=白银2段，505=钻石5段）。
- `getConfRank` 把段位×名次拍平为 `map[dan]map[rankNo]conf`，结算时 O(1) 查奖励。
- 比赛结果 `getResultType`：以「夺旗数 → 三路剩余兵力差」依次比较，返回 0=赢 / 1=输 / 2=平。
- 积分映射 `getConfScoreByFlagNum`：0/3 → PointWinner，1 → PointLoser，2 → PointDraw。
- 战力公式 `getFighting1`：
  - `singleFight = sqrt(sqrt(∏atk × ∏def))`，按兵种类型 (1/2/3) 乘修正系数 1.13/1.05/1.00；
  - `fightingScore = Σ(num × singleFight × factor) / sqrt(totalNum)`；
  - 最终 `ret = ceil(0.1 × fightingScore^1 + 100)`。
  - 系数 A=0.1、B=1、C=100，注释明确「策划文档让写进代码，不走配置」。

---

### 五、核心业务（actalliancechampionship_logic.go）

订阅了两个事件：
- `EventGuildQuit` → `onUserExitGuild`：报名期允许退出（同步 generalserver 与本地数据），其它阶段已成功报名的会忽略。
- `EventGuildDisSolve` → `onGuildDisSolve`：报名/结算期允许解散，其它阶段调用 `logger.ThrowFake` 拒绝。

17 个对外函数与 RPC 一一对应，主要业务点：
- **Dispatch（派遣部队）**：限报名期 + 主城等级达到 `NeedPlayerCastleLevel`；保存到 `Squad` 模块（类型 `EST_AllianceChampionship`），写入联盟对象的左/中/右路并通知 generalserver。
- **ReplaceRouter / SwitchUserRouter / ExchangeRouter / Rollback**：分别是玩家自换、管理换他人、管理整路对调、管理回退。后三者要求备战阶段（`isInPrepareWarTime`），首次操作触发 `BackUp()` 备份当前部署，所有操作都写入 `TActivityAllianceChampionShipGuildRecord`。
- **GetTeamRank / GetCourse / GetLastSimpleInfo / GetBattleData**：从 `TActAllianceChampionShipTeamObj` 读取当组对阵、排名、上一轮简略信息和某一场战报回放。
- **GetDeployment / GetFormation**：返回某条线路当前阵容（按战力倒序）。

权限校验通过 `assertGuildMemberCanOperate(... 52)`（位掩码权限位 52 = R5/R4 管理权限）。

---

### 六、客户端 RPC 入口（actalliancechampionship_rpc.go）

固定模式：
1. 检查 `isActivityOpen`，否则 `ThrowFake`；
2. 取 uid → gid，校验非 0；
3. `LockGuildData / LockUserData / LockData(...)` 加锁；
4. 调 logic 层函数。

`Dispatch` 是唯一一个同时锁公会、世界（`EworldLockId_AllianceChampionShip_Id`）、用户三把锁的，因为既要写联盟部署、又要写全局 system obj 的 `gidArr`。

`GetSuccessGuildList` 在拿到 `successGuildArr` 后会把成员数、最大成员数、联盟旗帜、自定义 ID 都补全，作为客户端的「本服参赛公会名单」。

---

### 七、与 generalserver 的协作

#### M → G（actalliancechampionship_togeneral.go）
全部走 `majorserver.RpcLogicGeneral(aContext, "alliancechampionship.<Method>ByInner", req, rsp)`，包括：
- 玩家派遣 / 自换 / 管理换 / 整路对调 / 回退同步全量；
- 解散公会 / 玩家退出公会；
- 阶段触发：`GroupingStarByInner` / `RoundBattleStarByInner` / 重发结果 / 强制重生分组 / 强制重生轮战。

#### G → M（actalliancechampionship_me.go）
通过 me 消息总线（`me.Listen`）订阅四种事件：
| 事件 | 处理函数 |
|---|---|
| `GENNERAL_TO_MAJOR_ME_EVENT_GROUPING_RESULT` | `meListener_GroupingResult` → `doGroupingResult` |
| `GENNERAL_TO_MAJOR_ME_EVENT_BATTLE_RESULT` | `meListener_RoundResult` → `doRoundResult` |
| `GENNERAL_TO_MAJOR_ME_EVENT_SIGNUP_FAIL` | 发送报名失败邮件 4077 |
| `GENNERAL_TO_MAJOR_ME_EVENT_GROUPING_FAIL` | 写入 `failedGuildArr`，结算时按第一名计算（策划要求） |

`doGroupingResult` 把 6 个公会绑定到同一个 teamId（`SetTeamId`），统计 `signupNum / totalFightForce` 写入 `successGuildArr`，并触发「报名成功邮件 4076」。

`doRoundResult` 处理一轮中三场对阵：
- 先 `LockWorldData(teamId)` → `AddCourseResult` 持久化到 team 对象；
- 然后为每个本服公会构造战报邮件 (`sendMailBattle`)；
- 中间通过 `saveClientPB` / `saveClientPBUser2` 把战斗服回传的逐场结果转换成客户端展示用 pb，并统计每个 uid 的「胜利次数 / 轮空 / 仅报名 / 未参战」标签（通过 `EActAllianceChampionShipKillNum_*` 枚举区分）。

值得注意的特殊规则：
- 最后一场 uid1（或 uid2）若是首次出战且获胜，记 `BattleOnceAndWin`；
- 队列中尚未出场的玩家：在「线路兵力上限」内的标记 `BattleButNotWar`，超过的标 `SignButNotBattle`；轮空时则全部标 `RoundEmpty/SignButNotBattle`。

---

### 八、数据对象（DataObj）

#### 1. SystemObj（actalliancechampionship_system_obj.go）
全局唯一，按 `EworldLockId_AllianceChampionShip_Id` 加世界锁。维护：
- `gidArr`：本服已参与的公会 id 集合；
- `successGuildArr`：分组成功的公会列表（含报名人数、总战力）；
- `failedGuildArr`：分组失败的公会；
- `isSettlement / isRemindManger`：幂等开关。

`Refresh()` 会在每次活动新一期开启时清空字段，避免跨期脏数据。

#### 2. GuildObj（actalliancechampionship_guild_obj.go）
活动核心数据载体。字段：
- `LeftSide / Middle / RightSide`：3 条线路当前部署（玩家 + 阵容（英雄+士兵））；
- `*Backup`：备战期首次修改前的快照，用于 Rollback；
- `TeamId`：分组后所属小组；
- `Dan / Dan2 / StarNum` + `LastDan*`：当前段位与上期段位；
- `DanFirstStatusArr`：是否首次抵达过该大段位（用来判断首次晋级邮件）。

行为亮点：
- `UserDispatch` 把玩家从原线路移除并按战力倒序 `Insert` 到目标线路（`insertUser`）；
- `Rollback` / `BackUp` 通过 `copyTAllianceChampionShipRoute` 实现深拷贝；
- `UpdateStarNum` 调 `calcDan` 处理「升大段 / 升小段 / 触顶 / 降大段 / 降小段 / 触底」六种状态机分支，并返回 `danChange` 与 `firstDan` 标志；
- `GetUserFightForce` / `GetTotalFightForce` 通过 `calcArmyFightForce` 单点求战力。
- `GenUidToWinNumWhenRoundEmpty` 在轮空时按线路 squadSize 限制给玩家发「轮空胜利」或「仅报名」的标签。

#### 3. TeamObj（actalliancechampionship_team_obj.go）
绑定一个小组（最多 6 联盟）的赛程与排名：
- `Course`：5 轮 × 3 场对阵的全量持久化（含 guild1/guild2 信息、左中右每个 uid 的对战明细：HP百分比、出战次数、双方战力、`BattleResultId` 等）；
- `Rank`：每个公会本组的累计 score / flagNum / killNum；
- `AddCourseResult` 把 generalserver 回传的 pb 灌进 db；
- `GetBattleData` 客户端重播查询，参数为 `roundNo + gid1 + gid2`；
- `sortRank` + 同分判定（score / flag / kill 三键比较）保证名次稳定。

#### 4. GuildRecordObj（actalliancechampionship_record_obj.go）
- `RecordArr` 保存所有调整线路（switch / exchange / rollback）操作日志；
- `BackupStatusArr` 标记每轮是否已对部署做过快照；
- `InvitedUserArr` + 每天 `LastRefreshTimeToday` 0 点重置：限制 R5/R4 邀请同一玩家的次数；
- 保存「操作类型 + 管理 uid + 被操作 uid + from/to 路 + 时间戳」便于联盟内审计。

#### 5. UserObj（actalliancechampionship_user_obj.go）
存「每轮战斗邮件已发时间戳」与「结算邮件已发时间戳」，用于幂等：避免同一轮、同一次结算重复发奖发邮件。

---

### 九、段位排行榜（actalliancechampionship_rank.go）

利用 `aContext.GetInRankListModule()` 全局排行键 `activity_alliance_champion_ship_rank`：
- `enScore`：把 dan / dan2 / star / 时间 / serverId / gidPiece 编码到一个 int64 score：
  - `dan*1e17 + dan2*1e15 + star*1e14 + (2071289203 - now)/1e5*1e10 + serverIdSuffix*1e6 + gidSuffix`
  - 时间使用 `2071289203 - now`（≈2035-08-21）实现「先达到此段位者排名靠前」；
- `deScore` 反解出 dan/dan2/star；
- `getRank` 取 Top100；
- `delRank` 解散公会时移除（先 `GetItemByUid` 再 `DelById`，按错误约定 not-in-rank 不返回 err）。

---

### 十、结算流程（actalliancechampionship_settlement.go）

1. 拉取所有匹配过的 `teamId`（按 `lastRefreshTime BETWEEN aStartTime AND aEndTime` 分页查询，每批 1000）。
2. 对每个 team 排序后产生 `gid -> rankNo`（同分并列）。`util.ExtractServerIdFromGid` 过滤本服公会。
3. 把分组失败 (`failedGuildArr`) 并入，按第一名待遇。
4. 对每个 gid 执行 `UpdateStarNum(starChange)` 推进段位，并收集 settlementGuild 结构。
5. `sendRewardToUser`：
   - 按 `confDanRank.GetReward()` 给参赛玩家发邮件 + 奖励（`InEmailExModule`）；
   - 通过 `userObj.GetSettlementTime()` 幂等控制，每人本期只收一次；
   - 若 `firstDan == 1` 即首次晋级到新大段位，再向**整盟成员**广播一次「联盟段位晋升邮件 + PromotionReward」。
6. `sysDataObj.SetIsSettlement()` 标记结算完成。

---

### 十一、邮件体系（actalliancechampionship_mail.go）

通过 `pb.NewTActAllianceChampionShipMail*` 序列化为 ext 数据，调用 `InEmailExModule.SendEmail(<mailId>, uid, …)`：
- `4076` 报名成功；
- `4077` 报名失败；
- 战斗结果邮件（`sendMailBattle`，含 BattleInfo 重播信息）；
- 结算邮件（含旧段位/新段位/星数/名次/奖励）；
- 大段位首次晋级邮件；
- 报名邀请邮件（`sendMailInvitePerson`）；
- 报名提醒管理邮件（`sendMailInviteManager`）。
- `Reissue*` 常量定义在 `actalliancechampionship_def.go`，预留延迟补发邮件（5–10 秒）使用。

---

### 十二、并发与一致性

- 所有写操作均显式加锁（`LockWorldData / LockGuildData / LockUserData / LockData`），尤其是 Dispatch 同时取三类锁；
- 跨服改动通过 `RpcLogicGeneral` 同步，本服保持源数据；
- 用 `AddSuccDeferCall` 在事务提交后再做异步邮件、再次开新事务发送，避免长事务；
- 大量 `Refresh()` 用一个 `LastRefreshTime` 与活动开始时间比较，实现跨期清零；
- `*_auto.go`（如 `actalliancechampionship_inpool_auto.go`、各 obj 的 `_auto.go`）由代码生成器产出，做对象池注册和接口导出，源文件本身不需要业务改动。

---

### 十三、GM/调试入口

在 `actalliancechampionship_in.go` 的 `TInActAllianceChampionShipModule` 中暴露的接口（仅供后台/GM 触发）：
- `GMSendToGeneralServer / GMGetInfo`：手动联调 RPC；
- `GMSendDisSolveGuild / GMSendQuitGuild`：模拟解散/退出；
- `GMExecGrouping / GMSendGroupingResult / GMGroupingStartAgain`：手动触发或重发分组流程；
- `GMExecRoundBattle / GMSendRoundBattle / GMRoundBattleStartAgain`：手动触发或重发某轮战斗；
- `GMTestFightForce`：跑一遍战力计算用于验证；
- `GMChangeDan`：直接修改公会大小段与星数（带配置合法性校验）。

---

### 十四、整体流程总结（按时间线）

1. **活动开启** → 注册四组定时器（提醒/分组/5 轮战斗/结算）。
2. **报名期** → 客户端调用 `Dispatch` / `ReplaceRouter` / `Invite` / `GetNotSignupList`，logic 写入 GuildObj 并同步 generalserver。事件：玩家退会/解散通过事件回调清理本服侧数据。
3. **匹配期** → 定时器触发 `GroupingStarByInner`，G 服分组完成后通过 me 总线下发 `GroupingResult` → `doGroupingResult` 持久化 Course、Rank，发 4076 邮件。
4. **备战 + 战斗 5 轮循环**：
   - 备战阶段允许 R5/R4 调整线路并自动 BackUp，玩家可看上一轮简单信息；
   - 战斗阶段开始时 G 服计算结果，me 下发 `RoundBattleResult` → `doRoundResult` 持久化战报、产出战斗邮件、统计胜利次数标签。
5. **结算期** → `settlementTimer` 调 `settlement`：拉所有 team → 生成 rank → `UpdateStarNum` → 发奖励邮件、晋升邮件、写 `IsSettlement = 1`。
6. **下一期开启** → DataObj `Refresh()` 全部清空。

---

### 十五、可关注的实现细节 / 风险点（仅做客观提示）

- 战力公式中 `fightingScore^B` 的 `B=1`，`^` 是 Go 的位异或而非乘方；按字面看表达式 `int64(fightingScore^B)` 会对最低位做异或，与注释「策划文档让写进代码」存在语义偏差，是否符合预期需要结合策划同步确认。
- `GetSignupRouterNo` 的注释误把变量名拼成 `rounterNo`（拼写问题，不影响逻辑）。
- `GroupingResultByInner` / `ResultByInner` 在 `actalliancechampionship_rpc.go` 内已被注释掉的实现，实际依赖 me 异步事件，老 RPC 接口仍被保留，需注意防止误用。
- `getNotSignupList` 在阶段判断错误时通过 `ThrowFake` 抛错而非返回错误码，调用方需要捕获。
- `delRank` 中「不在排行榜中也不返回 err」的注释依赖 `GetItemByUid` 的实现约定，迁移底层排行实现时需注意保持语义。
- `ReissueEmailDelayMin/MaxSecs = 5/10` 常量已声明但当前可见代码未使用，可能用于补发邮件的将来扩展。

---

总结：该模块是一个相当典型的「主服 + 中央服 + 跨服跨公会」的活动闭环实现。**主服只保留本服公会/玩家相关的数据与展示**，**真正的匹配与战斗仿真都在 generalserver**，两者通过 `RpcLogicGeneral`（M→G）+ `me.Listen`（G→M）的双向通道协同；持久化层用 5 个 `tables.TActivityAllianceChampionShip*` 表分别承载 system / guild / team / record / user 五类数据，并通过 `Refresh()` 实现按活动期幂等。
