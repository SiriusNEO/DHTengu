# Megumu - HashTable

PPCA 2021, final assignment, Distributed Hash Table (chord protocol)

### 简介

> 饭纲丸龙，鸦天狗的首领，具有操纵星空程度的能力。
>
> 大天狗身上承担着发展、维系天狗社会的义务。天狗社会每时每刻都在进行着纷繁的通信与调度，作为大天狗的它巧妙地利用分布式的原理去中心化。
>
> 网络中的 Node，与天魔之山的星星有着不少相似之处呢。
>



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



- DHT



```go
FileSystem
	utils.go //utils
	torrent.go //torrent-file related
	p2p.go //upload & download
	command.go //terminal command
```



### Problem

- [x] 进程抢占？（好吧是我自己写错了）
- [x] Quit 问题 （刷新 SuccList + finger 跳错特判）
- [x] Round 2 failed？（data move deadlock）
- [x] too many files？（~~finger 表跳转过慢问题~~  client Close 问题）
- [x] 疑似死锁？（锁管理）
- [x] Map 拷贝问题（统一采用 Copy）

