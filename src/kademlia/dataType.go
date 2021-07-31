package kademlia

import (
	"fmt"
	"sync"
	"time"
)

//kademlia DataType, a LockMap with validTime

const (
	NeedRePublishTime = 120 * time.Second
	ExpiredTime = 600 * time.Second
)

type DataType struct {
	hashMap     	map[string]string
	validTime   	map[string]time.Time
	republishTime   map[string]time.Time
	lock     sync.RWMutex
}

func (this *DataType) Init() {
	this.hashMap = make(map[string]string)
	this.validTime = make(map[string]time.Time)
	this.republishTime = make(map[string]time.Time)
}

func (this *DataType) RePublishList() (republishList []string) {
	this.lock.RLock()
	for key, tim := range this.republishTime {
		if time.Now().After(tim) {
			republishList = append(republishList, key)
		}
	}
	this.lock.RUnlock()
	return
}

func (this *DataType) Expire() {
	var expiredKeys []string

	this.lock.RLock()
	for key, tim := range this.validTime {
		if time.Now().After(tim) {
			expiredKeys = append(expiredKeys, key)
		}
	}
	this.lock.RUnlock()

	this.lock.Lock()
	for _, key := range expiredKeys {
		fmt.Println("Data Expired: ", key, this.hashMap[key])
		delete(this.hashMap, key)
		delete(this.validTime, key)
	}
	this.lock.Unlock()
}

func (this *DataType) Store(key string, val string) {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.hashMap[key] = val
	this.validTime[key] = time.Now().Add(ExpiredTime)
	this.republishTime[key] = time.Now().Add(NeedRePublishTime)
}

func (this DataType) Load(key string) (founded bool, ret string) {
	this.lock.RLock()
	defer this.lock.RUnlock()

	ret, founded = this.hashMap[key]

	return
}

func (this *DataType) Delete(key string) (founded bool) {
	this.lock.Lock()
	defer this.lock.Unlock()
	_, founded = this.hashMap[key]
	delete(this.hashMap, key)
	return founded
}

//deep copy

func (this *DataType) Copy() map[string]string {
	ret := make(map[string]string)

	this.lock.Lock()
	defer this.lock.Unlock()

	for key, value := range this.hashMap {
		ret[key] = value
	}

	return ret
}
