[TOC]



# **后端工程师 Linux 生存与效率提升完全手册**

**核心哲学：** 你不是在学命令，而是在学习与机器对话、高效解决问题的思维。

## 一、 基础生存技能（必须掌握）

这些命令是你登录服务器后第一时间就会用到的。

### 1. 导航与文件管理
*   `pwd` (**P**rint **W**orking **D**irectory): 显示当前所在目录的绝对路径。
*   `ls` (**L**i**s**t): 列出目录内容。
    *   `ls -l`: 以长格式显示（权限、所有者、大小、时间）。
    *   `ls -la`: 显示**所有**文件（包括隐藏文件，如 `.bashrc`）。
    *   `ls -lh`: 以人类可读的形式显示文件大小（K, M, G）。
*   `cd` (**C**hange **D**irectory): 切换目录。
    *   `cd ~`: 回到家目录。
    *   `cd ..`: 返回上一级目录。
    *   `cd -`: 回到上一个所在的目录。
*   `mkdir` (**M**a**k**e **Dir**ectory): 创建新目录。`mkdir -p project/{src,log,bin}` (一次性创建嵌套目录)。`mkdir -p project{src,log,bin}`也可以运行看看哦。
*   `rm` (**R**e**m**ove): 删除。**危险命令！**
    *   `rm file.txt`: 删除文件。
    *   `rm -r directory/`: **递归删除**目录及其所有内容。
    *   `rm -rf directory/`: **强制递归删除**，不提示。**使用此命令前请再三确认！**
*   `cp` (**C**o**p**y): 复制。
    *   `cp source.txt dest.txt`: 复制文件。
    *   `cp -r source_dir/ dest_dir/`: 递归复制目录。
*   `mv` (**M**o**v**e): 移动或重命名。
    *   `mv old.txt new.txt`: 重命名。
    *   `mv file.txt /tmp/`: 移动文件。
*   `find`: 强大的文件查找工具。
    *   `find /home/app -name "*.log"`: 在 `/home/app` 目录下按名称查找 `.log` 文件。
    *   `find . -mtime -1`: 查找当前目录下最近1天内修改过的文件。
    *   `find . -size +100M`: 查找大于100MB的文件（清理磁盘时非常有用）。

### 2. 文件查看与编辑
*   `cat`: 连接并打印文件内容（适合小文件）。
*   `less` / `more`: 分页查看文件内容（适合大文件），**优先用 `less`**。
*   `head` / `tail`: 查看文件头/尾部。
    *   `tail -f app.log`: **实时追踪**日志输出，**后端必备神技**。
*   `vim` / `nano`: 文本编辑器。**强烈建议至少学会 `vim` 的基础操作**：
    *   `vim file.txt`: 打开文件。
    *   `i`: 进入插入模式（开始编辑）。
    *   `ESC`: 退出插入模式。
    *   `:wq`: 保存并退出。
    *   `:q!`: 不保存强制退出。

### 3. 系统状态与进程管理
*   `ps` (**P**rocess **S**tatus): 查看进程快照。
    *   `ps aux | grep java`: 查看所有 Java 进程。`aux` 参数组合最常用。
*   `top` / `htop`: 动态、交互式地查看系统进程和资源（CPU、内存）使用情况。**`htop` 更强大，需安装**。
*   `kill`: 终止进程。
    *   `kill <PID>`: 优雅地终止进程（发送 TERM 信号）。
    *   `kill -9 <PID>`: **强制杀死**进程（发送 KILL 信号）。终极手段，慎用。
*   `df` (**D**isk **F**ree): 查看磁盘空间使用情况。
    *   `df -h`: 以人类可读格式显示。
*   `du` (**D**isk **U**sage): 查看文件或目录的磁盘使用量。
    *   `du -sh /path/to/dir`: 汇总 (`-s`) 显示目录的总大小，人类可读格式 (`-h`)。
*   `free`: 查看内存使用情况。`free -h`。

### 4. 网络工具
*   `ping`: 测试与目标主机的网络连通性。`ping example.com`。
*   `netstat` / `ss`: 查看网络连接、路由表、接口统计等。
    *   `netstat -tulnp | grep 8080` 或 `ss -tulnp | grep 8080`: 查看谁在监听 8080 端口。**排查端口冲突必备**。
*   `curl`: 强大的命令行 HTTP 客户端。**后端调试神器**。
    *   `curl http://api:8080/health`: 发起 GET 请求。
    *   `curl -X POST http://api:8080/order -H "Content-Type: application/json" -d '{"id": 123}'`: 发起带 JSON 体的 POST 请求。
    *   `curl -I http://api:8080/health`: 只获取响应头。
*   `wget`: 命令行下载工具。

---

## 二、 中级效率工具（强烈推荐掌握）

这些命令和技巧将极大提升你的工作效率。

### 1. 文本处理三剑客 (与日志分析强相关)
*   `grep`: **过滤**文本。`grep "ERROR" log.txt | grep -v "healthcheck"`。
*   `sed`: **流编辑器**，主要用于文本替换和过滤。
    *   `sed 's/foo/bar/g' file.txt`: 将文件中所有 `foo` 替换为 `bar`。
    *   `sed -n '10,20p' file.txt`: 只打印文件的第10到20行。
*   `awk`: **文本和报表生成器**，擅长处理结构化数据。
    *   `awk '{print $1}' access.log`: 打印空格分隔的第一列（如IP地址）。
    *   `awk -F ',' '{print $NF}' data.csv`: 使用逗号分隔，打印最后一列 (`$NF`)。

### 2. 管道与重定向 (Linux哲学的精华)
*   `|` (管道): 将一个命令的**输出**作为另一个命令的**输入**。
    *   `ps aux | grep nginx | head -n 5`
*   `>` 和 `>>` (重定向):
    *   `echo "hello" > file.txt`: 将输出**覆盖**写入文件（文件原有内容丢失！）。
    *   `echo "world" >> file.txt`: 将输出**追加**到文件末尾。

### 3. 权限管理
*   `chmod` (**Ch**ange **Mod**e): 改变文件权限。
    *   `chmod +x script.sh`: 给脚本添加可执行权限。
    *   `chmod 755 script.sh`: 用数字模式设置权限（用户：rwx，组：rx，其他：rx）。
*   `chown` (**Ch**ange **Own**er): 改变文件所有者。`chown user:group file.txt`。

### 4. 压缩与解压
*   `tar`: 打包和解包工具。
    *   `tar -czvf archive.tar.gz /path/to/dir`: 压缩 (`-c`) 为 gzip 格式，显示详细过程 (`-v`)，指定文件名 (`-f`)。
    *   `tar -xzvf archive.tar.gz`: 解压 (`-x`)。
*   `zip` / `unzip`: 处理 `.zip` 文件。

---

## 三、 高级生存与调试（问题排查专家）

### 1. 系统性能排查 (线上问题常用)
*   `top` / `htop`: 实时监控。
*   `iostat`, `vmstat`: 查看 CPU、IO 状态。
*   `lsof` (**L**i**st** **O**pen **F**iles): 列出被进程打开的文件。`lsof -p <PID>` 查看某个进程打开了哪些文件、端口。

### 2. 环境与变量
*   `env`: 查看所有环境变量。
*   `export`: 设置环境变量（当前会话有效）。`export JAVA_HOME=/usr/lib/jvm/java-11`。
*   `which` / `whereis`: 查找命令的安装路径。`which java`。

### 3. 后台运行与会话管理
*   `&`: 在命令结尾加上 `&`，使其在**后台运行**。
*   `nohup`: 忽略挂断信号，让命令在你退出登录后继续运行。`nohup java -jar app.jar &`。
*   `tmux` / `screen`: **终端复用器**。允许你在一个终端窗口中创建多个持久化会话，即使网络断开，任务仍在后台运行。**这是生产环境操作的必备保险**，防止长时间操作因断网而中断。

---

## 四、 思维与最佳实践

1.  **--help 和 man**: 遇到不熟悉的命令，第一反应是 `command --help` 或 `man command`（ manual ）。这是你最应该掌握的“命令”。
2.  **Tab 键补全**: 输入命令或路径时，多按 `Tab` 键，可以自动补全或给出提示，避免拼写错误。
3.  **命令历史**: 按 `↑` 键翻阅历史命令。使用 `history` 查看所有历史，`!<编号>` 可以快速执行历史中的某条命令。
4.  **知其然，知其所以然**: 不要死记硬背。理解每个命令的用途（它是用来查看什么的？修改什么的？），组合使用比单一命令更强大。
5.  **安全第一**:
    *   对 `rm -rf`, `chmod`, `chown` 等**危险命令**保持敬畏之心，执行前再三确认路径和参数。
    *   修改重要配置文件前，**先备份**！`cp config.yml config.yml.bak`。
6.  **使用脚本自动化**: 如果某个复杂的操作需要重复执行，就把它写成一个 Shell 脚本 (`*.sh`)。

## 总结：每日命令清单

| 场景             | 常用命令                                  |
| :--------------- | :---------------------------------------- |
| **登录服务器**   | `ssh user@host`                           |
| **查看当前位置** | `pwd`                                     |
| **查看文件**     | `ls -la`, `less file.txt`                 |
| **查看日志**     | `tail -f app.log`, `grep "ERROR" app.log` |
| **查看进程**     | `ps aux \| grep java`, `top`              |
| **查看端口**     | `netstat -tulnp \| grep 8080`             |
| **调试API**      | `curl -X POST http://...`                 |
| **检查磁盘**     | `df -h`, `du -sh .`                       |
| **操作后台**     | `nohup ... &`, `tmux`                     |

**给你的挑战：**
不要一次性学完所有命令。每天尝试使用1-2个新命令来解决实际工作中的问题。很快，你就会发现与服务器交流的乐趣和力量。

Welcome to the Linux world!



# awk 命令详解

`awk` 是 Linux/Unix 系统中功能极其强大的文本处理工具，它不仅是一个命令，更是一门功能丰富的**模式扫描和处理语言**。其名称来源于三位创始人 Alfred **A**ho、Peter **W**einberger 和 Brian **K**ernighan 的姓氏首字母。

## 一、核心思想与工作流程

`awk` 的基本思想是：**逐行扫描文件（或输入流），根据给定的"模式（Pattern）"进行匹配，对匹配到的行执行对应的"动作（Action）"**。

### 工作流程：
1. **读取输入流的第一行**
2. **将行内容分割成"字段（Field）"**
   - 默认以空格（或制表符）作为分隔符
   - 第一个字段是 `$1`，第二个是 `$2`，以此类推
   - **整行内容**为 `$0`
3. **检查该行是否匹配"模式"部分**
   - 如果匹配，则执行对应的"动作"
   - 如果没有模式部分，则默认匹配所有行，执行动作
   - 如果没有动作部分，则默认动作为 `{ print }`，即打印整行
4. **读取下一行，重复上述过程，直到所有行处理完毕**

## 二、基本语法结构

```bash
awk '模式 { 动作 }' 输入文件
# 或者从管道读取
cat 输入文件 | awk '模式 { 动作 }'
```

- **单引号 `' '`**：将 awk 程序括起来，防止 shell 解释其中的特殊字符（如 `$`）
- **模式 (Pattern)**：可选部分，用于筛选行
- **动作 (Action)**：可选部分，包含在花括号 `{ }` 中

### 模式类型：
- **正则表达式**：例如 `/^root/` 匹配以 `root` 开头的行
- **关系表达式**：例如 `$1 == "foo"` 或 `$3 > 100`
- **特殊模式**：
  - `BEGIN`：在读取任何输入行**之前**执行一次（常用于初始化变量、打印表头）
  - `END`：在处理完所有输入行**之后**执行一次（常用于汇总、输出结果）
- **模式组合**：使用 `&&` (与), `||` (非), `!` (非) 组合多个模式

### 动作内容：
- 包含一系列语句（命令），用分号 `;` 分隔
- 常见动作为 `print`, `printf`, 计算，条件判断，循环等

## 三、常用内置变量

`awk` 提供了许多有用的内置变量，无需定义即可使用。

| 变量              | 描述                                                      | 默认值       |
| ----------------- | --------------------------------------------------------- | ------------ |
| `$0`              | 当前行的全部内容                                          |              |
| `$1, $2, ..., $n` | 当前行的第1，2，...，n个字段                              |              |
| `NF`              | **Number of Fields**，当前行的字段总数                    |              |
| `NR`              | **Number of Records**，当前处理的是第几行（所有文件累计） |              |
| `FNR`             | **File Number of Records**，当前处理的是当前文件的第几行  |              |
| `FS`              | **Field Separator**，输入字段分隔符                       | 空格或制表符 |
| `OFS`             | **Output Field Separator**，输出字段分隔符                | 空格         |
| `RS`              | **Record Separator**，输入行（记录）分隔符                | 换行符 `\n`  |
| `ORS`             | **Output Record Separator**，输出行（记录）分隔符         | 换行符 `\n`  |
| `FILENAME`        | 当前正在处理的文件名                                      |              |

## 四、常用命令与示例

假设我们有一个名为 `demo.txt` 的文件，内容如下：
```
Alice 90 95
Bob 85 80
Charlie 78 92
Dave 95 98
```

### 1. 打印整行或特定字段

```bash
# 打印每一行的全部内容（等同于 cat）
awk '{print $0}' demo.txt

# 打印每行的第一个字段（名字）
awk '{print $1}' demo.txt

# 打印每行的第一和第三个字段
awk '{print $1, $3}' demo.txt

# 自定义输出格式（使用逗号连接）
awk '{print $1 "的成绩是：" $2 "分"}' demo.txt
# 输出：Alice的成绩是：90分
```

### 2. 基于条件过滤行

```bash
# 打印第二个字段（数学成绩）大于90的行
awk '$2 > 90 {print $0}' demo.txt
# 输出：Dave 95 98

# 打印包含 "li" 的行（正则表达式匹配）
awk '/li/ {print $1}' demo.txt
# 输出：Alice
#       Charlie

# 打印第一个字段是 "Bob" 的行
awk '$1 == "Bob" {print}' demo.txt
```

### 3. 使用内置变量

```bash
# 打印每一行及其字段数量
awk '{print $0, "有", NF, "个字段"}' demo.txt

# 打印文件的行号及其第一列
awk '{print NR, $1}' demo.txt
# 输出：
# 1 Alice
# 2 Bob
# 3 Charlie
# 4 Dave

# 改变输入和输出分隔符（处理CSV文件）
echo "Alice,Bob,Charlie" | awk 'BEGIN{FS=","; OFS=" | "} {print $1, $2, $3}'
# 输出：Alice | Bob | Charlie
```

### 4. BEGIN 和 END 块

```bash
# 在开头打印表头，在结尾打印总结
awk 'BEGIN {print "开始处理成绩单..."} {print $1} END {print "处理结束"}' demo.txt

# 计算数学成绩（第二列）的平均分
awk '{sum += $2; count++} END {print "平均分是: ", sum/count}' demo.txt
# sum += $2：累加第二列的值
# count++：计数行数
# END：最后计算平均值
```

### 5. 结合管道与其他命令

```bash
# 查找当前目录下的文件，并打印大小和文件名
ls -l | awk '{print $5, $9}'

# 显示系统上所有用户的用户名和家目录（来自 /etc/passwd）
awk -F: '{print $1, $6}' /etc/passwd
# -F: 指定输入字段分隔符为冒号 `:`
```

## 五、高级用法

### 1. 自定义变量
```bash
awk '{total = $2 + $3; print $1, total}' demo.txt
```

### 2. 条件语句 (if-else)
```bash
awk '{if ($2 > 90) print $1, "优秀"; else print $1, "良好"}' demo.txt
```

### 3. 循环 (for, while)
```bash
# 打印每一行的每一个字段
awk '{for (i=1; i<=NF; i++) print "字段", i, ":", $i}' demo.txt
```

### 4. 数组 (关联数组)
```bash
# 统计不同分数出现的次数（假设第二列是分数）
awk '{count[$2]++} END {for (score in count) print score, "出现了", count[score], "次"}' demo.txt
```

## 六、常用选项

| 选项            | 描述                                                  |
| --------------- | ----------------------------------------------------- |
| `-F fs`         | 指定输入字段分隔符 `fs`（等同于在 BEGIN 中设置 `FS`） |
| `-v var=value`  | 定义一个 awk 变量并赋值，可以在 awk 程序中使用        |
| `-f scriptfile` | 从脚本文件（而不是命令行）中读取 awk 程序             |

### 示例：
```bash
# 使用 -F 处理冒号分隔的文件
awk -F: '{print $1}' /etc/passwd

# 使用 -v 传递外部变量
awk -v name="Bob" '$1 == name {print $2}' demo.txt
```

## 总结

`awk` 是一个极其强大和灵活的工具，尤其擅长处理**结构化的文本数据**（如日志、CSV、表格数据）。它的核心在于**"模式-动作"** 模型。

**主要应用场景：**
- **简单文本提取和打印**：使用 `{print $N}`
- **过滤**：使用关系表达式或正则表达式作为模式
- **计算和统计**：结合变量、数组和 `END` 块
- **格式化报告**：使用 `printf` 进行精细化的格式控制

掌握 `awk` 可以让你在命令行下的文本处理能力提升一个数量级，是系统管理员和开发者必备的技能。建议从简单的例子开始练习，逐步尝试更复杂的操作。