# Go 为什么是神

### 部署

##### go env

`GOPATH`  当前 Project 的 Go 路径，分为 GLOBAL / Project GOPATH

`GOROOT`  系统根路径，golang 的安装路径

```
go_project //GOPATH
-- bin //可执行文件
-- pkg //.a
-- src //.go
```

GO 的包搜索方式：

先去 GOROOT/src

再去 Project GOPATH/src

再去 Global GOPATH/src

##### package

package 在物理上对应一个文件夹（最好保持命名一致）

必须要有一个 main package

引入包

```go
import (
)
```

注意

```go
Demo //大写 public
demo //小写 protected

//命名规范：驼峰
```



### 类型

```go
var b bool = true
var b int = 1
var b uint = 1
var b float64 = 2
var b complex64 = 1+4i

:= //声明赋值，不用写 var, 不能在函数外用

/*
派生类型：
Pointer
Array
struct
Channel
func
Slice
interface
Map
*/

//Map、Slice、Channel这些都要make！（不然就是nil 什么都做不了） 
//自己的类型最好用指针，定义一个new方法来创建

_, b = 5, 7 //空白标识符：只写，表示抛弃变量
```



### 语句

```go
//for-each
for key, value := range oldMap {
    newMap[key] = value
}

//interface
type interface_name interface {
    
}

func (instance InstanceType) func_name(para) [return] {
    
}
```



###  并发

```go
go func //开启一个goroutine

//Channel
ch := make(chan int)
ch <- x
<-ch

ch := make(chan int, 100) //100 buffer

sync.Mutex //互斥锁，保证只有一个goroutine
sync.RWMutex //读写互斥锁，更高效：读的线程可以拿
```



### Net

客户端  Client / Workstation

服务端  Server

RPC（Remote Procedure Call） 远程过程调用，A 调用 B 的方法

```go
"net/rpc"

1. Server
func (server *Server) Register(rcvr interface{}) error //向Server注册RPC服务
func (server *Server) Accept(lis net.Listener) //监听器

//开一个Server, 带监听器

server := rpc.NewServer()
server.Register(new(Type)) //Register参数是interface{}, 要求实现符合5个条件的接口
                           //Type用于接口，重要的是Type的方法！

//重要的是：都要exported，方法参数两个： args reply，reply要指针，返回error

listener, err := net.Listen("tcp", addr) //注意都要时刻关注err 
server.Accept(listener) //阻塞，直到listener close

2. Client
client, err := rpc.Diag("tcp", ip) //Diag return *client
client.Call("Method_Name", args, &reply) //
```

P2P 没有服务器，用户之间传递信息

port 端口号，`ip : port`  唯一确定

```go
//一些常用的互联网接口含义

Ping //尝试与客户端连接

```



### 其它

Go 没有重载运算符，因此 package 给的一些类型的基本运算都是靠函数

> 比如 bigInt（怨念）
>
> bigInt.Cmp()   -1<  0=  1>

注意 big.Int 深浅拷贝问题！

big.Int unaddressable