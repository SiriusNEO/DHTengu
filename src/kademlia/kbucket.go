package kadmelia

type KBucketType struct {
	size   		int
	addr 		[K]AddrType
}

func (this *KBucketType) Update(addr AddrType) {
	if addr.Ip == "" {
		return
	}

	founded := -1
	for i := 0; i < this.size; i++ {
		if this.addr[i].Ip == addr.Ip {
			founded = i
			break
		}
	}

	if founded == -1 {
		if this.size < K {
			this.addr[this.size] = addr
			this.size++
			return
		} else {
			if Ping(this.addr[0].Ip) != nil {
				for i := 1; i < K; i++ {
					this.addr[i-1] = this.addr[i]
				}
				this.addr[K-1] = addr
				return
			} else {
				head := this.addr[0]
				for i := 1; i < K; i++ {
					this.addr[i-1] = this.addr[i]
				}
				this.addr[K-1] = head
				return
			}
		}
	} else {
		for i := founded+1; i < this.size; i++ {
			this.addr[i-1] = this.addr[i]
		}
		this.addr[this.size-1] = addr
	}
}
