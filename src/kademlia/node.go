package kadmelia

import "github.com/sasha-s/go-deadlock"

type NodeType struct {
	Running     bool
	Addr        AddrType
	buckets		[M]KBucketType
	data        LockMap
	mux         deadlock.Mutex
}

func NewNode(ip string) *NodeType {
	ret := &NodeType{
		Running: false,
		Addr: AddrType{ip, Hash(ip)},
	}
	ret.data.Init()
	return ret
}

func (this *NodeType) Join(ip string) {
	id := Hash(ip)
	this.buckets[cpl(&id, &this.Addr.Id)].Update(AddrType{Ip: ip, Id: id})

}

