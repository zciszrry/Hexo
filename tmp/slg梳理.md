#### 发起行军？

​	首先是major到slg的innerRpc。slg中，需要构造路径，获得hero信息等等，然后去构造march对象，如果时玩家被行军，会用Me推送一下。

​	构造对象时，先向db插入值为0的数据行（在获得obj里调Find时init应该时一样的效果），再获得obj，获得obj时初始化了committer（此时committer的initExit值为false），然后用obj填充march的值（此时设置了CreatTime，故march的isExit值为true）。

​	在doAfterUpdate中，触发commit。

​	commit的add会将march转为March对象，将行军计时器刷新后，存入到SlgMap（地图实例）的Service里。

​	计时器结束后,触发行军到达事件，处理行军到达逻辑