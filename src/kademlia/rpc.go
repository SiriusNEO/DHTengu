package kademlia

import (
	"github.com/sirupsen/logrus"
	"net"
	"net/rpc"
)

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
		Log.WithFields(logrus.Fields{}).Error("Listen Error")
	}

	return ret
}

func (this *ReceiverType) Store(args *StoreArg, _ *int) error {
	this.Node.data.Store(args.Key, args.Value)
	this.Node.kBucketUpdate(args.Sender)
	return nil
}

func (this *ReceiverType) FindNode(args *FindNodeArg, reply *ClosestList) error {
	*reply = this.Node.FindNode(&args.TarID)
	this.Node.kBucketUpdate(args.Sender)
	//fmt.Println(args.Sender.Ip, this.Node.Addr.Ip)
	return nil
}

func (this *ReceiverType) FindValue(args *FindValueArg, reply *FindValueRet) error {
	*reply = this.Node.FindValue(args.Key, &args.Hash)
	this.Node.kBucketUpdate(args.Sender)
	return nil
}
