package chord

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"math/big"
	"sync"
	"time"
)

//NodeType (true node in the network)
//Addr and Running are visible

type NodeType struct {
	Addr         AddrType
	Running      bool
	data         LockMap
	backup       LockMap
	finger       [M]AddrType
	predecessor  AddrType
	succList     SuccListType

	mux           sync.Mutex
}

func NewNode(ip string) *NodeType{
	ret := &NodeType{
		Running: false,
		Addr: AddrType{ip, Hash(ip)},
	}
	ret.data.Init()
	ret.backup.Init()
	return ret
}

//Display for debug
func (this *NodeType) Display() {
	fmt.Println(" * Display *")
	fmt.Println("ip: ", this.Addr.Ip, "id: ", this.Addr.Id, " Running: ", this.Running)
	fmt.Println("data: ", this.data)
	fmt.Println("finger table: ", this.finger)
	fmt.Println("predecessor: ", this.predecessor)
	fmt.Println("succList: ", this.succList)
	fmt.Println()
}

//Create the Network in NodeType: pre and suc <- self
func (this *NodeType) Create() {
	this.predecessor = this.Addr

	for i := 0; i < SuccListLen; i++ {
		this.succList[i] = this.Addr
	}

	for i := 0; i < M; i++ {
		this.finger[i] = this.Addr
	}
}

func (this *NodeType) findSuccessor(keyId *big.Int) AddrType {
	this.succListFlush()

	this.mux.Lock()
	succ := this.succList[0]
	this.mux.Unlock()

	if IsIn(&this.Addr.Id, &succ.Id, keyId, false, true) {
		return succ
	}

	cpn := this.closestPrecedingNode(keyId)
	if cpn.Ip == this.Addr.Ip || cpn.Ip == "" {
		cpn = this.succList[0]
	}

	//cpn := this.succList[0]

	client, err := Diag(cpn.Ip)
	if err != nil {
		/*Log.WithFields(logrus.Fields{
			"from" : this.Addr.Ip,
			"to" : cpn.Ip,
			"key" : keyId,
		}).Error("Diag Failed. " + err.Error())*/
		return AddrType{}
	} else {
		/*
		Log.WithFields(logrus.Fields{
			"from" : this.Addr.Ip,
			"to" : cpn.Ip,
			"key" : keyId,
		}).Info("Diag Success. ")
		*/
	}
	defer client.Close()
	var ret AddrType
	err = client.Call("ReceiverType.FindSuccessor", keyId, &ret)
	return ret
}

//Search in the finger table, find the index-largest pos which is the predecessor of key
func (this *NodeType) closestPrecedingNode(keyId *big.Int) AddrType {
	this.mux.Lock()
	defer this.mux.Unlock()

	pingFailed := make(map[string]bool)

	for i := M-1; i > 0; i-- {
		// if finger[i] \in (n, Id)
		if this.finger[i].Ip != "" && IsIn(&this.Addr.Id, keyId, &this.finger[i].Id, false, false) {
			if pingFailed[this.finger[i].Ip] {
				continue
			}
			err := Ping(this.finger[i].Ip)
			if err == nil {
				return this.finger[i]
			} else {
				Log.WithFields(logrus.Fields{
					"from" : this.Addr.Ip,
					"to" : this.finger[i].Ip,
				}).Error("Ping Failed. " + err.Error())
				pingFailed[this.finger[i].Ip] = true
			}
		}
	}
	return this.Addr
}

//Join , GuIde, Ip
func (this *NodeType) Join(ip string) bool {
	client, err := Diag(ip)
	if err != nil {
		Log.WithFields(logrus.Fields{
			"from" : this.Addr.Ip,
			"to" : ip,
		}).Error("Diag Failed. " + err.Error())
		return false
	}
	err = client.Call("ReceiverType.FindSuccessor", &this.Addr.Id, &this.succList[0])
	client.Close()

	Log.WithFields(logrus.Fields{
		"from" : this.Addr.Ip,
		"to" : ip,
	}).Info("Tracing1... Join")

	client, err = Diag(this.succList[0].Ip)

	if err != nil {
		Log.WithFields(logrus.Fields{
			"from" : this.Addr.Ip,
			"to" : this.succList[0].Ip,
		}).Error("Diag Failed. " + err.Error())
		return false
	}

	var succData map[string]string
	err = client.Call("ReceiverType.GetData", 0, &succData)

	Log.WithFields(logrus.Fields{
		"from" : this.Addr.Ip,
		"to" : this.succList[0].Ip,
	}).Info("Tracing2... Join")

	if err != nil {
		Log.WithFields(logrus.Fields{
			"from" : this.Addr.Ip,
			"to" : this.succList[0].Ip,
		}).Error("Call Failed. " + err.Error())
		return false
	}

	for key := range succData {
		keyID := Hash(key)
		fmt.Println("Data Trans", keyID)
		if !IsIn(&this.Addr.Id, &this.succList[0].Id, &keyID, false, true) {
			this.data.Store(key, succData[key])
			var ret bool
			err = client.Call("ReceiverType.DirectlyDelete", key, &ret)
			if err != nil {
				Log.WithFields(logrus.Fields{
					"from" : this.Addr.Ip,
					"to" : this.succList[0].Ip,
				}).Error("Call Failed. " + err.Error())
				return false
			}
		}
	}
	return true
}

//Quit tell pre and suc, move its data to succ
func (this *NodeType) Quit() {
	client, err := Diag(this.predecessor.Ip)
	if err != nil {
		/*Log.WithFields(logrus.Fields{
			"from" : this.Addr.Ip,
			"to" : this.predecessor.Ip,
		}).Error("Diag Failed. " + err.Error())*/
	} else {
		err = client.Call("ReceiverType.SuccListUpdate", &this.succList[SuccListLen-1], nil)
		client.Close()
	}

	this.succListFlush()

	client, err = Diag(this.succList[0].Ip)
	if err != nil {
		/*Log.WithFields(logrus.Fields{
			"from" : this.Addr.Ip,
			"to" : this.succList[0].Ip,
		}).Error("Diag Failed. " + err.Error())*/
		return
	}
	defer client.Close()

	err = client.Call("ReceiverType.PredecessorUpdate", &this.predecessor, nil)

	for key := range this.data.hashMap {
		var ret bool
		client.Call("ReceiverType.DirectlyPut", StrPair{key, this.data.hashMap[key]}, &ret)
	}

	this.data.Init() //clear
}

func (this *NodeType) succListUpdate(tail *AddrType) {
	this.mux.Lock()
	defer this.mux.Unlock()

	Log.WithFields(logrus.Fields{
		"ip" : this.Addr.Ip,
	}).Info("Tracing... succListupd")

	for i := 1; i < SuccListLen; i++ {
		this.succList[i-1] = this.succList[i]
	}

	this.succList[SuccListLen-1] = *tail

	fmt.Println(this.succList)
}

//Put K-V Pair
func (this *NodeType) Put(key string, value string) bool {
	keyId := Hash(key)

	Log.WithFields(logrus.Fields{
		"ip" : this.Addr.Ip,
		"id" : this.Addr.Id,
		"key" : keyId,
	}).Info("Put Tracing...")

	if IsIn(&this.predecessor.Id, &this.Addr.Id, &keyId, false, true) {
		founded, dat := this.data.Load(key)
		if founded {
			Log.WithFields(logrus.Fields{
				"ip" : this.Addr.Ip,
				"key" : keyId,
				"value" : value,
				"dat" : dat,
			}).Error("Put Failed. ")
			return false
		}
		this.data.Store(key, value)
		Log.WithFields(logrus.Fields{
			"ip" : this.Addr.Ip,
			"key" : keyId,
			"value" : value,
		}).Info("Put Success. ")
		return true
	}

	tar := this.findSuccessor(&keyId)
	client, err := Diag(tar.Ip)

	if err != nil {
		Log.WithFields(logrus.Fields{
			"from" : this.Addr.Ip,
			"to" : tar.Ip,
		}).Error("Diag Failed. " + err.Error())
		return false
	}
	defer client.Close()

	var ret bool
	err = client.Call("ReceiverType.Put", StrPair{key, value}, &ret)

	if err != nil {
		Log.WithFields(logrus.Fields{
			"from" : this.Addr.Ip,
			"to" : tar.Ip,
		}).Error("Call Failed. " + err.Error())
		return false
	}

	return ret
}

//Get V
func (this *NodeType) Get(key string) (founded bool, value string) {
	keyId := Hash(key)

	Log.WithFields(logrus.Fields{
		"ip" : this.Addr.Ip,
		"id" : this.Addr.Id,
		"key" : keyId,
	}).Info("Get Tracing...")

	if IsIn(&this.predecessor.Id, &this.Addr.Id, &keyId, false, true) {
		return this.data.Load(key)
	}

	tar := this.findSuccessor(&keyId)
	client, err := Diag(tar.Ip)

	if err != nil {
		Log.WithFields(logrus.Fields{
			"from" : this.Addr.Ip,
			"to" : tar.Ip,
		}).Error("Diag Failed. " + err.Error())
		return
	}
	defer client.Close()

	var ret BoolStrPair
	err = client.Call("ReceiverType.Get", key, &ret)
	return ret.First, ret.Second
}

//Delete Key
func (this *NodeType) Delete(key string) bool {
	keyId := Hash(key)

	if IsIn(&this.predecessor.Id, &this.Addr.Id, &keyId, false, true) {
		return this.data.Delete(key)
	}

	tar := this.findSuccessor(&keyId)
	client, err := Diag(tar.Ip)

	if err != nil {
		Log.WithFields(logrus.Fields{
			"from" : this.Addr.Ip,
			"to" : tar.Ip,
		}).Error("Diag Failed. " + err.Error())
		return false
	}

	var ret bool
	err = client.Call("ReceiverType.Delete", key, &ret)
	return ret
}

//Updating 3 stages, When a Node is Running, it will Updating its data
func (this *NodeType) Updating() {
	next := 0
	for this.Running {
		this.checkPredecessor()
		this.stabilize()
		this.fixFingers(&next)
		time.Sleep(UpdateInterval)
	}
}

//Ping Predecessor, if failed set this.pre nil
func (this *NodeType) checkPredecessor()  {
	this.mux.Lock()
	defer this.mux.Unlock()
	err := Ping(this.predecessor.Ip)
	if err != nil {
		/*
		Log.WithFields(logrus.Fields{
			"from" : this.Addr.Ip,
			"to" : this.predecessor.Ip,
		}).Error("Ping Failed. " + err.Error())
		*/
		this.predecessor = AddrType{"", *big.NewInt(0)}
	}
}

func (this *NodeType) succListFlush() {
	this.mux.Lock()
	defer this.mux.Unlock()
	var j int
	for j = 0; j < SuccListLen; j++ {
		if Ping(this.succList[j].Ip) == nil {
			break
		}
		this.succList[j] = AddrType{} //clear useless
	}
	for i := 0; j < SuccListLen; j++ {
		this.succList[i] = this.succList[j]
		i++
	}
}

func (this *NodeType) stabilize()  {
	this.succListFlush()

	this.mux.Lock()
	succ := this.succList[0]
	this.mux.Unlock()

	var succPre AddrType
	client, err := Diag(succ.Ip)
	if err != nil {
		/*Log.WithFields(logrus.Fields{
			"from" : this.Addr.Ip,
			"to" : succ.Ip,
		}).Error("Diag Failed. " + err.Error())*/
		return
	}

	err = client.Call("ReceiverType.GetPredecessor", 0, &succPre)
	if succPre.Ip != "" && IsIn(&this.Addr.Id, &succ.Id, &succPre.Id, false, false) {
		this.mux.Lock()
		this.succList[0] = succPre
		this.mux.Unlock()
	}

	var succSuccList SuccListType
	err = client.Call("ReceiverType.GetSuccList", 0, &succSuccList)

	this.mux.Lock()
	for i := 1; i < SuccListLen; i++ {
		this.succList[i] = succSuccList[i-1]
	}
	this.mux.Unlock()

	var empty int
	err = client.Call("ReceiverType.Notify", &this.Addr, &empty)
	/*
	Log.WithFields(logrus.Fields{
		"ip" : this.Addr.Ip,
	}).Info("Stabilize Success. ")
	*/
}

//notify this, check whether to modify the pre
func (this *NodeType) notify(addr *AddrType) {
	this.mux.Lock()
	defer this.mux.Unlock()

	if this.predecessor.Ip == "" || IsIn(&this.predecessor.Id, &this.Addr.Id, &addr.Id, false, false) {
		this.predecessor = *addr
	}
}

//fixPos is the index to fix, every time fixPos++ to fix the next pos
func (this *NodeType) fixFingers(fixPos *int)  {
	var tar big.Int
	tar.Add(&this.Addr.Id, big.NewInt(0).Exp(big.NewInt(2), big.NewInt(int64(*fixPos)), Mod))
	tar.Mod(&tar, Mod)

	fixFinger := this.findSuccessor(&tar)

	this.mux.Lock()
	this.finger[*fixPos] = fixFinger
	this.mux.Unlock()
	/*
	Log.WithFields(logrus.Fields{
		"ip" : this.Addr.Ip,
		"fixpos" : *fixPos,
	}).Info("Fix Finger Success")
	*/
	*fixPos++
	if *fixPos >= M {
		*fixPos = 0
	}
}