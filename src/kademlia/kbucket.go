package kademlia

import (
	"github.com/sasha-s/go-deadlock"
	"math/big"
)

type KBucketType struct {
	size   		int
	bucket 		[K]AddrType
	mux         deadlock.Mutex
}

func (this *KBucketType) Reflesh() {
	for i := 0; i < this.size; i++ {
		if Ping(this.bucket[i].Ip) != nil {
			for j := i+1; j < this.size; j++ {
				this.bucket[j-1] = this.bucket[j]
			}
			this.size--
			return
		}
	}
}

func (this *KBucketType) Update(addr AddrType) {
	this.mux.Lock()
	defer this.mux.Unlock()

	if addr.Ip == "" {
		return
	}

	founded := -1
	for i := 0; i < this.size; i++ {
		if this.bucket[i].Ip == addr.Ip {
			founded = i
			break
		}
	}

	if founded == -1 {
		if this.size < K {
			this.bucket[this.size] = addr
			this.size++
			return
		} else {
			if Ping(this.bucket[0].Ip) != nil {
				for i := 1; i < K; i++ {
					this.bucket[i-1] = this.bucket[i]
				}
				this.bucket[K-1] = addr
				return
			} else {
				head := this.bucket[0]
				for i := 1; i < K; i++ {
					this.bucket[i-1] = this.bucket[i]
				}
				this.bucket[K-1] = head
				return
			}
		}
	} else {
		for i := founded+1; i < this.size; i++ {
			this.bucket[i-1] = this.bucket[i]
		}
		this.bucket[this.size-1] = addr
	}
}

type ClosestList struct {
	Size     	int
	Standard    big.Int
	List	 	[K]AddrType
}

func (this *ClosestList) Insert(addr AddrType) (updated bool) { //promise in-order
	updated = false

	if Ping(addr.Ip) != nil {
		return false
	}

	for i := 0; i < this.Size; i++ { //founded
		if this.List[i].Ip == addr.Ip {
			return false
		}
	}

	newDis := dis(&addr.Id, &this.Standard)

	if this.Size < K {
		updated = true
		for i := 0; i < this.Size; i++ {
			nowDis := dis(&this.List[i].Id, &this.Standard)
			if newDis.Cmp(&nowDis) < 0 {
				for j := this.Size; j > i; j-- {
					this.List[j] = this.List[j-1]
				}
				this.Size++
				this.List[i] = addr
				return
			}
		}
		this.List[this.Size] = addr
		this.Size++
		return
	}


	for i := 0; i < K; i++ { //从小到大，找到第一个小的
		nowDis := dis(&this.List[i].Id, &this.Standard)
		if newDis.Cmp(&nowDis) < 0 {
			updated = true
			for j := K-1; j > i; j-- {
				this.List[j] = this.List[j-1]
			}
			this.List[i] = addr
			return
		}
	}

	return
}

func (this *ClosestList) Remove(addr AddrType) bool {
	for i := 0; i < this.Size; i++ {
		if this.List[i].Ip == addr.Ip {
			for j := i+1; j < this.Size; j++ {
				this.List[j-1] = this.List[j]
			}
			this.Size--
			return true
		}
	}
	return false
}