package kademlia

import (
	"fmt"
	"github.com/sasha-s/go-deadlock"
	"github.com/sirupsen/logrus"
	"math/big"
	"time"
)

type NodeType struct {
	Running     	bool
	Addr        	AddrType
	routeTable		[M]KBucketType
	data        	DataType
	mux         	deadlock.RWMutex
}

func NewNode(ip string) *NodeType {
	ret := &NodeType{
		Running: false,
		Addr: AddrType{ip, Hash(ip)},
	}
	ret.data.Init()
	return ret
}

func (this *NodeType) Display() {
	fmt.Println(" * Display *")
	fmt.Println("ip: ", this.Addr.Ip, "id: ", this.Addr.Id, " Running: ", this.Running)
	fmt.Println("data: ", this.data.hashMap)
	fmt.Println("KBucket: ")
	for i := 0; i < M; i++ {
		fmt.Println(this.routeTable[i].bucket)
	}
	fmt.Println()
}

func (this *NodeType) FindNode(tarID *big.Int) (closestList ClosestList) {
	this.mux.RLock()
	defer this.mux.RUnlock()

	closestList.Standard = *tarID

	for i := 0; i < M; i++ {
		for j := 0; j < this.routeTable[i].size; j++ {
			if Ping(this.routeTable[i].bucket[j].Ip) == nil {
				closestList.Insert(this.routeTable[i].bucket[j])
			}
		}
	}

	return
}

//FindValue return value or K-closest nodes
func (this *NodeType) FindValue(key string, hash *big.Int) FindValueRet {
	this.mux.RLock()
	defer this.mux.RUnlock()

	founded, value := this.data.Load(key)

	if founded {
		//fmt.Println("hit!")
		return FindValueRet{ClosestList{}, value}
	}

	closestList := ClosestList{Standard: *hash}

	for i := 0; i < M; i++ {
		for j := 0; j < this.routeTable[i].size; j++ {
			if Ping(this.routeTable[i].bucket[j].Ip) == nil {
				closestList.Insert(this.routeTable[i].bucket[j])
			}
		}
	}

	return FindValueRet{closestList, ""}
}

func (this *NodeType) NodeLookup(tarID *big.Int) (closestList ClosestList) {
	closestList = this.FindNode(tarID)
	closestList.Insert(this.Addr) //robust: avoid 1-node can't find self
	updated := true
	diaged := make(map[string]bool)

	for updated {
		updated = false
		var closestListTmp ClosestList
		var removeList []AddrType
		for i := 0; i < closestList.Size; i++ {
			if diaged[closestList.List[i].Ip] == true {
				continue
			}
			this.kBucketUpdate(closestList.List[i])
			client, err := Diag(closestList.List[i].Ip)
			diaged[closestList.List[i].Ip] = true
			var ret ClosestList
			if err != nil {
				Log.WithFields(logrus.Fields{"from" : this.Addr.Ip, "to" : closestList.List[i].Ip}).Error("Diag Failed. " + err.Error())
				removeList = append(removeList, closestList.List[i])
			} else {
				err = client.Call("ReceiverType.FindNode", &FindNodeArg{TarID: *tarID, Sender: this.Addr}, &ret)
				for j := 0; j < ret.Size; j++ {
					closestListTmp.Insert(ret.List[j])
				}
				client.Close()
			}
		}

		for _, key := range removeList {
			closestList.Remove(key)
		}

		for i := 0; i < closestListTmp.Size; i++ {
			updated = updated || closestList.Insert(closestListTmp.List[i])
		}
	}

	return
}

func (this *NodeType) Get(key string) (founded bool, value string) {
	keyID := Hash(key)
	diaged := make(map[string]bool)
	findValue := this.FindValue(key, &keyID)
	if findValue.Second != "" {
		return true, findValue.Second
	}
	closestList := findValue.First
	//fmt.Println("Get List", closestList.List)
	updated := true

	for updated {
		updated = false
		var closestListTmp ClosestList
		var removeList []AddrType
		for i := 0; i < closestList.Size; i++ {
			if diaged[closestList.List[i].Ip] == true {
				continue
			}
			client, err := Diag(closestList.List[i].Ip)
			diaged[closestList.List[i].Ip] = true
			var ret FindValueRet
			if err != nil {
				Log.WithFields(logrus.Fields{"from" : this.Addr.Ip, "to" : closestList.List[i].Ip}).Error("Diag Failed. " + err.Error())
				removeList = append(removeList, closestList.List[i])
			} else {
				err = client.Call("ReceiverType.FindValue", &FindValueArg{Key: key, Sender: this.Addr}, &ret)
				if ret.Second != "" {
					return true, ret.Second
				} else {
					for j := 0; j < ret.First.Size; j++ {
						closestListTmp.Insert(ret.First.List[j])
					}
				}
				client.Close()
			}
		}

		for _, key1 := range removeList {
			closestList.Remove(key1)
		}

		for i := 0; i < closestListTmp.Size; i++ {
			updated = updated || closestList.Insert(closestListTmp.List[i])
		}
	}

	secondList := this.NodeLookup(&keyID)

	for _, i := range secondList.List {
		client, _ := Diag(i.Ip)
		defer client.Close()
		var ret FindValueRet
		client.Call("ReceiverType.FindValue", &FindValueArg{Key: key, Sender: this.Addr}, &ret)
		if ret.Second != "" {
			return true, ret.Second
		}
	}

	return false, ""
}

func (this *NodeType) Put(key string, value string) bool {
	keyID := Hash(key)

	//sta := time.Now()
	closestList := this.NodeLookup(&keyID)
	//fmt.Println("Put Look-Up Time: ", time.Now().Sub(sta), key, value)

	closestList.Insert(this.Addr) //robust: avoid 1-node can't find self
	for i := 0; i < closestList.Size; i++ {
		client, err := Diag(closestList.List[i].Ip)
		if err != nil {
			Log.WithFields(logrus.Fields{"from" : this.Addr.Ip, "to" : closestList.List[i].Ip}).Error("Diag Failed. " + err.Error())
		} else {
			err = client.Call("ReceiverType.Store", &StoreArg{Key: key, Value: value, Sender: this.Addr}, nil)
			if err != nil {
				Log.WithFields(logrus.Fields{"from" : this.Addr.Ip, "to" : closestList.List[i].Ip}).Error("Call Failed. " + err.Error())
			}
			client.Close()
		}
	}

	return true
}

func (this *NodeType) Create() {
	LogInit()
}

func (this *NodeType) Join(ip string) bool {
	id := Hash(ip)
	this.kBucketUpdate(AddrType{Ip: ip, Id: id})

	client, err := Diag(ip)

	if err != nil {
		Log.WithFields(logrus.Fields{"from" : this.Addr.Ip, "to" : ip}).Error("Diag Failed. " + err.Error())
	} else {
		var ret ClosestList
		err = client.Call("ReceiverType.FindNode", &FindNodeArg{this.Addr.Id, this.Addr}, &ret)
		for j := 0; j < ret.Size; j++ {
			this.kBucketUpdate(ret.List[j])
		}
		client.Close()
	}

	closestList := this.NodeLookup(&this.Addr.Id)
	for i := 0; i < closestList.Size; i++ {
		this.kBucketUpdate(closestList.List[i])
		client, err = Diag(closestList.List[i].Ip)
		if err != nil {
			Log.WithFields(logrus.Fields{"from" : this.Addr.Ip, "to" : closestList.List[i].Ip}).Error("Diag Failed. " + err.Error())
		} else {
			var ret ClosestList
			err = client.Call("ReceiverType.FindNode", &FindNodeArg{this.Addr.Id, this.Addr}, &ret)
			for j := 0; j < ret.Size; j++ {
				this.kBucketUpdate(ret.List[j])
			}
			client.Close()
		}
	}
	return true
}

func (this *NodeType) Quit() {
	this.data.Init() //clear data
}

func (this *NodeType) RePublish() {
	for this.Running {
		for i := 0; i < M; i++ {
			this.routeTable[i].Reflesh()
		}

		this.mux.Lock()
		thisData := this.data.Copy()
		republishList := this.data.RePublishList()
		this.mux.Unlock()

		sta := time.Now()
		for _, key := range republishList {
			this.Put(key, thisData[key]) //RePublish To Closest Nodes.
		}
		if false {
			fmt.Println("RePublish Time: ", time.Now().Sub(sta))
		}

		this.data.Expire()

		time.Sleep(RePublishInterval)
	}
}

func (this *NodeType) kBucketUpdate(addr AddrType) {
	this.mux.Lock()
	defer this.mux.Unlock()

	if addr.Ip == "" || addr.Ip == this.Addr.Ip {
		return
	}

	this.routeTable[cpl(&this.Addr.Id, &addr.Id)].Update(addr)
}

