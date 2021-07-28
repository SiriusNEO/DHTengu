package kadmelia

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

func (this *ReceiverType) Store() error {

}

func (this *ReceiverType) FindNode() error {

}

func (this *ReceiverType) FindValue() error {

}
