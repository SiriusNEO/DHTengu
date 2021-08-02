# DHT Development Doc.

PPCA 2021, final assignment, Distributed Hash Table

其实是写给自己看的（



### Log

- Week1：爽摸
- Week2：chord + tengu 下载、上传部分
- Week3：kademlia
- Week4：调试，tengu music player



### TO DO

- [x] 基本的 Go 语法
- [x] OOP 练习：LinkList
- [x] 读懂测试代码（道阻且长...）
- [x] net 入门，实现简单多 server RPC
- [x] DHT框架设计 （粗略√）
- [x] debug.go & utils.go
- [x] node.go
- [x] rpc.go
- [x] ForceQuit
- [x] BT App
- [x] Magnet Support & Test
- [x] Kademlia Learning
- [x] Kademlia Implement
- [x] Music Player



### Draft

##### chord

使用 log 做一套完备的debug输出

```go
chord
 - utilities.go
 - node.go
 - rcvr.go
 - debug.go
```

靠谱的添加？并行更新

```go
Run: //单开线程执行
CheckPredecessor
Stabilization - Notify someone
Fix_Fingers


```

```go
type Node struct {
    running //true or false
    server
    listener
    data
    backup //pre or suc?
    finger
    predecessor
    succList
    many mutex...
    
}
```



Running 是 RPC 层面的，Node 中只读不可写（读：for Running ...）



开闭区间问题？

chord ring 上 (key1, key2] 这一段是由 key2 对应的 Node 管的

$$key \in (this, suc]$$

closest_preceding_node 中，finger is the predecessor of id

$$finger \in (n, id)$$

stabilize 里，更新后继，x 和 succ 重复没意义（不一样才要更新）

$$x \in (n, succ)$$



##### app

```go
FileSystem
	utils.go //utils
	torrent.go //torrent-file related
	torrentClient.go //upload & download
	fileOperation.go  //save & read to disk
	command.go //terminal command
```



##### kademlia

```go
[RPC] FindNode() //在当前节点路由表中找最近K个节点（异或值小：近），除非路由表不满K个节点否则一定要找完
[RPC] FindValue() //与FindNode一样，只不过如果这个点有存储该key值，返回value，否则找这个点的最近K个信息

NodeLookup() //在整个网络中找离目标ID最近的K个节点
//先找本节点最近K个，然后并发地（并发数α）向它们发送 FindNode
//接受FindNode信息，再找K个距离最近的，发送FindNode
//停止条件：本轮FindNode结果已经无法再更新任何节点，即节点已经最近

Store(Key, Value) 
//NodeLookup，找到K个
//把数据保存在这K个（发送Store RPC）
//RePublish?

Find(Key)
//类似NodeLookup，不过用FindValue代替FindNode

Join(bootstrap)
//bootstrap加入自己的桶
//自己执行NodeLookup
//刷新KBucket（?）
```



### Problem

- [x] 进程抢占？（好吧是我自己写错了）
- [x] Quit 问题 （刷新 SuccList + finger 跳错特判）
- [x] Round 2 failed？（data move deadlock）
- [x] too many files？（~~finger 表跳转过慢问题~~  client Close 问题）
- [x] 疑似死锁？（锁管理）
- [x] Map 拷贝问题（统一采用 Copy）

