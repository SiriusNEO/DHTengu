package kademlia

import (
	"crypto/sha1"
	"errors"
	"math/big"
	"net/rpc"
	"time"
)

const (
	M = 16
	K = 6
	Alpha = 3

	UpdateInterval = 25 * time.Millisecond

	RemoteTryTime = 3
	RemoteTryInterval = 25 * time.Millisecond

	RePublishInterval = time.Second
)

//can not declared as const

var Mod = big.NewInt(0).Exp(big.NewInt(2), big.NewInt(M), nil)

//A IP address store as a pair due to the hash value is frequently called
type AddrType struct {
	Ip    string
	Id    big.Int
}

type StrPair struct {
	First   string
	Second  string
}

type BoolStrPair struct {
	First   bool
	Second  string
}

//Arg and Ret Type
type FindNodeArg struct {
	TarID big.Int
	Sender AddrType
}

type FindValueArg struct {
	Key string
	Hash big.Int
	Sender AddrType
}

type StoreArg struct {
	Key string
	Value string
	Sender AddrType
}

type FindValueRet struct {
	First     ClosestList
	Second    string
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

//dis(a, b) = a xor b
func dis(obj1 *big.Int, obj2 *big.Int) big.Int {
	var ret big.Int
	ret.Xor(obj1, obj2)
	return ret
}

//common prefix length between two ID
func cpl(obj1 *big.Int, obj2 *big.Int) int {
	xorDis := dis(obj1, obj2)
	return xorDis.BitLen() - 1
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