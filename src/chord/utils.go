package chord

import (
	"crypto/sha1"
	"errors"
	"github.com/sasha-s/go-deadlock"
	"math/big"
	"net/rpc"
	"time"
)

const (
	M = 160

	SuccListLen = 10

	UpdateInterval = 50 * time.Millisecond

	RemoteTryTime = 3
	RemoteTryInterval = 25 * time.Millisecond
)

//can not declared as const

var Mod = big.NewInt(0).Exp(big.NewInt(2), big.NewInt(M), nil)

//A IP address store as a pair due to the hash value is frequently called

type AddrType struct {
	Ip    string
	Id    big.Int
}

type SuccListType [SuccListLen]AddrType

type StrPair struct {
	First   string
	Second  string
}

type BoolStrPair struct {
	First   bool
	Second  string
}

//a SHA-1 hash, hash % mod, mod = 2^M

func Hash(str string) big.Int{
	hasher := sha1.New()
	hasher.Write([]byte(str))
	var ret big.Int
	ret.SetBytes(hasher.Sum(nil))
	ret.Mod(&ret, Mod)
	return ret
}

//val in [ll, rr], if ll >= rr, rr += Mod

func IsIn(ll *big.Int, rr *big.Int, val *big.Int, lClosed bool, rClosed bool) bool {
	var l big.Int
	l.Add(ll, big.NewInt(0))
	var r big.Int
	if ll.Cmp(rr) >= 0 {r.Add(rr, Mod)} else {r.Add(rr, big.NewInt(0))}
	lCmp := -1
	rCmp := 1
	if lClosed {lCmp = 0}
	if rClosed {rCmp = 0}

	if l.Cmp(val) <= lCmp && r.Cmp(val) >= rCmp {
		return true
	}

	var val1 big.Int
	val1.Add(val, Mod)
	return l.Cmp(&val1) <= lCmp && r.Cmp(&val1) >= rCmp
}

//RPC Service: Diag, Ping

func Diag(addr string) (*rpc.Client, error) {
	var ret *rpc.Client
	var err error
	if addr == "" {
		return nil, errors.New("ERROR: empty IP addr")
	}
	for i := 0; i < RemoteTryTime; i++ {
		ret, err = rpc.Dial("tcp", addr)
		if err == nil {
			return ret, err
		}
		time.Sleep(RemoteTryInterval)
	}
	return nil, err
}

func Ping(addr string) error {
	var ret *rpc.Client
	var err error
	if addr == "" {
		return errors.New("ERROR: empty IP addr")
	}
	for i := 0; i < RemoteTryTime; i++ {
		ret, err = rpc.Dial("tcp", addr)
		if err == nil {
			ret.Close()
			return err
		}
		time.Sleep(RemoteTryInterval)
	}
	return err
}

//A map with lock, for safe in para env
//not use sync.Map due to the need of other operations

type LockMap struct {
	hashMap     map[string]string
	lock     	deadlock.RWMutex
}

func (this *LockMap) Init() {
	this.hashMap = make(map[string]string)
}

func (this *LockMap) Store(key string, val string) {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.hashMap[key] = val
}

func (this LockMap) Load(key string) (founded bool, ret string) {
	this.lock.RLock()
	defer this.lock.RUnlock()
	ret, founded = this.hashMap[key]
	return
}

func (this *LockMap) Delete(key string) (founded bool) {
	this.lock.Lock()
	defer this.lock.Unlock()
	_, founded = this.hashMap[key]
	delete(this.hashMap, key)
	return founded
}

func (this *LockMap) Copy() map[string]string {
	ret := make(map[string]string)

	this.lock.Lock()
	defer this.lock.Unlock()

	for key, value := range this.hashMap {
		ret[key] = value
	}

	return ret
}