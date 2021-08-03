# DHT & Tengu

### 简介

PPCA 2021, final assignment, Distributed Hash Table (chord & kademlia protocol)

and its application, Tengu 

> 饭纲丸龙，鸦天狗的首领，具有操纵星空程度的能力。
>
> 大天狗身上承担着发展、维系天狗社会的义务。天狗社会每时每刻都在进行着纷繁的通信与调度，作为大天狗的它巧妙地利用分布式的原理去中心化。
>
> 网络中的 Node，与天魔之山的星星有着不少相似之处呢。



### Requirements

Chord & Kademlia 用到了如下第三方库

```
github.com/sirupsen/logrus
github.com/sasha-s/go-deadlock
```

Tengu App 在此基础上还用到了如下第三方库

```
github.com/jackpal/bencode-go
github.com/faiface/beep
github.com/nsf/termbox-go
```

此外，Tengu 的音乐播放功能需要声卡驱动支持

```
sudo apt install libasound2-dev
```



### DHT Interface

```go
NewNode(port int) dhtNode

Run()
Create()
Join(addr string) bool
Quit()
ForceQuit()
Ping(addr string) bool
Put(key string, value string) bool
Get(key string) (bool, string)
Delete(key string) bool  //Kademlia 不实现此方法
```



### Chord

**文件结构**

```go
node.go //
rpc.go //RPC Client & RPC Method 设计
utils.go //相关常数的规定以及工具函数
```

**类型设计**

```go
type PubNodeType struct {
    //实现 dhtNode interface
	receiver	*ReceiverType
}

type ReceiverType struct {
    //实现TCP通信，封装好的RPC客户端
    Node      *NodeType
	Server    *rpc.Server
	Listener  net.Listener
}

type NodeType struct { //真正节点
    ...
}
```

**关键部分算法**

- `FixFinger`：单开线程并行维护，每次修复一位，然后将位置加 1 准备对下个位置 Fix
- `FindSuccessor`：按照论文编写，跳 Finger 表，且要求 Ping 的通，若 Finger 表全部失效则访问后继
- `ForceQuit` 处理：每个节点备份自己**后继**的数据，同时维护后继列表，当维护后继列表过程中发现后继失效，将备份数据合并到后继列表第一个有效节点。



### Kademlia

**文件结构**

```go
dataType.go //实现带有有效时间的数据
kBucket.go //实现K桶以及一个有序队列（用于K-Closest）
node.go
rpc.go
utils.go
```

**类型设计**

同 Chord

**关键部分算法**

- `NodeLookUp` ：先自己进行一次 `FIND_NODE` 加入结果队列，之后每次选取结果队列中最近 K 个点发送 `FIND_NODE RPC`，直到结果队列不再能被更新。
- `Get`：同 `NodeLookUp`，只是将 `FIND_NODE` 换成 `FIND_VALUE`（即找到立即停止）
- `RePublish`：先遍历数据获取需要重新发布的键值对，对于每个键值对，以 `key` 为参数进行一次 `NodeLookUp`，然后对这 K 个点发送 `STORE RPC` 。优化：在一个点接收到 `STORE RPC` 后，说明另外 K-1 个点也会收到，因此记录下该键值，下一周期不再重布。



### Tengu

Tengu 是一个支持本地进行小文件共享的P2P文件系统。

此外，它还是一个简易的”共享歌单“，任何人都可以加入歌曲、播放歌单里的歌曲。

如果您拥有 Tengu 的可执行文件包，执行 `./tengu` 来运行它，目前它还是一个命令行软件。

目前支持的功能有：

- 上传文件，生成 torrent 种子与磁力链接
- 下载文件（使用种子或者磁力链接）
- 上传歌曲（比起上传文件，这是特别为歌曲定制的入口，它能将你的歌曲加入某一专辑。当然，它完全包含了上传文件的内容）
- 播放歌曲（播放歌曲先进行一个下载过程来获得歌曲临时文件，再调用 Tengu 的内置 Player 播放歌曲）



### Reference

[dht.pdf](./ref/dht.pdf) from [@xmhuangzhen](https://github.com/xmhuangzhen)

[Golang Net](https://pkg.go.dev/net)

[Chord: A Scalable Peer-to-peer Lookup Protocol for Internet Applications](./ref/paper-ton.pdf)

[Kademlia: A Peer-to-Peer Information System Based on the XOR Metric  ](./ref/2002_Book_Peer-to-PeerSystems.pdf)

[Building a BitTorrent client from the ground up in Go](https://blog.jse.li/posts/torrent/#putting-it-all-together)

[Kademlia 算法学习](https://shuwoom.com/?p=813)



