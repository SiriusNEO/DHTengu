package chord

import (
	"github.com/sirupsen/logrus"
	"math/big"
	"net"
	"net/rpc"
)

//ReceiverType is used in RPC Call
//All methods should conform the format

type ReceiverType struct {
	Node      *NodeType
	Server    *rpc.Server
	Listener  net.Listener
}

func NewReceiver(ip string) *ReceiverType{
	ret := &ReceiverType{
		Node: NewNode(ip),
		Server: rpc.NewServer(),
		Listener: nil,
	}
	var err error
	ret.Listener, err = net.Listen("tcp", ip)

	if err != nil {
		Log.WithFields(logrus.Fields{
		}).Error("Listen Error")
	}

	return ret
}

func (this *ReceiverType) Test(args int, reply *string) error {
	*reply = this.Node.Addr.Ip
	return nil
}

func (this ReceiverType) GetPredecessor(_ int, reply *AddrType) error {
	this.Node.mux.Lock()
	defer this.Node.mux.Unlock()

	*reply = this.Node.predecessor
	return nil
}

func (this ReceiverType) GetSuccList(_ int, reply *SuccListType) error {
	this.Node.mux.Lock()
	defer this.Node.mux.Unlock()

	*reply = this.Node.succList
	return nil
}

func (this ReceiverType) FindSuccessor(args *big.Int, reply *AddrType) error {
	*reply = this.Node.findSuccessor(args)
	return nil
}

func (this *ReceiverType) Notify(args *AddrType, _ *int) error {
	this.Node.notify(args)
	return nil
}

func (this *ReceiverType) Put(args StrPair, reply *bool) error {
	*reply = this.Node.Put(args.First, args.Second)
	return nil
}

func (this *ReceiverType) Get(args string, reply *BoolStrPair) error {
	founded, value := this.Node.Get(args)
	*reply = BoolStrPair{founded, value}
	return nil
}

func (this *ReceiverType) Delete(args string, reply *bool) error {
	*reply = this.Node.Delete(args)
	return nil
}

func (this *ReceiverType) DirectlyPut(args StrPair, reply *bool) error {
	founded, _ := this.Node.data.Load(args.First)

	if founded {
		*reply = false
	}

	this.Node.data.Store(args.First, args.Second)
	*reply = true
	return nil
}

func (this *ReceiverType) BackupDirectlyPut(args StrPair, reply *bool) error {
	founded, _ := this.Node.backup.Load(args.First)

	if founded {
		*reply = false
	}

	this.Node.backup.Store(args.First, args.Second)
	*reply = true
	return nil
}

func (this *ReceiverType) DirectlyDelete(args string, reply *bool) error {
	*reply = this.Node.data.Delete(args)
	return nil
}

func (this *ReceiverType) BackupDirectlyDelete(args string, reply *bool) error {
	*reply = this.Node.backup.Delete(args)
	return nil
}

func (this *ReceiverType) GetData(_ int, reply *map[string]string) error {
	this.Node.mux.Lock()
	defer this.Node.mux.Unlock()
	*reply = this.Node.data.hashMap
	return nil
}

func (this *ReceiverType) UpdateBackup(args map[string]string, reply *int) error {
	this.Node.mux.Lock()
	defer this.Node.mux.Unlock()
	this.Node.backup.hashMap = args
	return nil
}

func (this *ReceiverType) PredecessorUpdate(args *AddrType, _ *int) error {
	this.Node.mux.Lock()
	defer this.Node.mux.Unlock()
	this.Node.predecessor = *args
	return nil
}

func (this *ReceiverType) SuccListUpdate(args *AddrType, _ *int) error {
	this.Node.succListUpdate(args)
	return nil
}