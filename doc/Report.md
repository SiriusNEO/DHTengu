# Report

### Chord

**关键部分算法**

- `FixFinger`：单开线程并行维护，每次修复一位，然后将位置加 1 准备对下个位置 Fix
- `FindSuccessor`：按照论文编写，跳 Finger 表，且要求 Ping 的通，若 Finger 表全部失效则访问后继
- `ForceQuit` 处理：每个节点备份自己**后继**的数据，同时维护后继列表，当维护后继列表过程中发现后继失效，将备份数据合并到后继列表第一个有效节点。

### Kademlia

**关键部分算法**

- `NodeLookUp` ：先自己进行一次 `FIND_NODE` 加入结果队列，之后每次选取结果队列中最近 K 个点发送 `FIND_NODE RPC`，直到结果队列不再能被更新。
- `Get`：同 `NodeLookUp`，只是将 `FIND_NODE` 换成 `FIND_VALUE`（即找到立即停止）
- `RePublish`：先遍历数据获取需要重新发布的键值对，对于每个键值对，以 `key` 为参数进行一次 `NodeLookUp`，然后对这 K 个点发送 `STORE RPC` 。优化：在一个点接收到 `STORE RPC` 后，说明另外 K-1 个点也会收到，因此记录下该键值，下一周期不再重布。

### Tengu App

**上传下载**

制作种子参考了 [Building a BitTorrent client from the ground up in Go](https://blog.jse.li/posts/torrent/#putting-it-all-together)

思路是直接将上传文件切成小 Piece，并行 Put 到 DHT，并做种。种子中的 PieceHash 作为 key 值。

**音乐播放**

mp3 文件解析、播放使用了 [beep package](github.com/faiface/beep)，音乐的储存管理继承自上传下载的功能。