package main

import (
	"chord"
	"fmt"
	"math/rand"
	"net"
	"runtime"
	"time"
)

const (
	myselfTestNodeSize = 5
	kvPairSize = 100
)

func GetLocalAddress() string {
	var localaddress string

	ifaces, err := net.Interfaces()
	if err != nil {
		panic("init: failed to find network interfaces")
	}

	// find the first non-loopback interface with an IP address
	for _, elt := range ifaces {
		if elt.Flags&net.FlagLoopback == 0 && elt.Flags&net.FlagUp != 0 {
			addrs, err := elt.Addrs()
			if err != nil {
				panic("init: failed to get addresses for network interface")
			}

			for _, addr := range addrs {
				ipnet, ok := addr.(*net.IPNet)
				if ok {
					if ip4 := ipnet.IP.To4(); len(ip4) == net.IPv4len {
						localaddress = ip4.String()
						break
					}
				}
			}
		}
	}
	if localaddress == "" {
		panic("init: failed to find non-loopback interface with valid address on this node")
	}

	return localaddress
}

func portToAddr(ip string, port int) string {
	return fmt.Sprintf("%s:%d", ip, port)
}

func main() {
	/*
	a := big.NewInt(23)
	b := big.NewInt(23)
	c := big.NewInt(23)

	fmt.Println(chord.IsIn(a, b, c, false, true))
	fmt.Println(a, b, c)
	*/
	chord.LogInit()
	runtime.GOMAXPROCS(runtime.NumCPU())

	localAddress := GetLocalAddress()
	firstPort := 20000

	nodes := new([myselfTestNodeSize + 1]PubNodeType)
	nodeAddresses := new([myselfTestNodeSize + 1]string)
	nodesInNetwork := make([]int, 0, myselfTestNodeSize+1)

	for i := 0; i <= myselfTestNodeSize; i++ {
		nodes[i] = *NewPubNode(portToAddr(localAddress, firstPort + i))
		nodeAddresses[i] = portToAddr(localAddress, firstPort + i)
		go nodes[i].Run()
	}

	time.Sleep(200 * time.Millisecond)

	fmt.Println("Run Finish")

	fmt.Println(chord.Mod)

	nodes[0].Create()

	nodesInNetwork = append(nodesInNetwork, 0)

	testKey := [] string {"THU", "PKU", "FDU", "SJTU", "ZJU"}
	testVal := [] string {"Beijing", "Beijing", "Shanghai", "Shanghai", "Zhejiang"}

	for i := 0; i < 5; i++ {
		fmt.Println(nodes[nodesInNetwork[rand.Intn(len(nodesInNetwork))]].Put(testKey[i], testVal[i]))
	}

	time.Sleep(time.Second)

	for i := 1; i <= myselfTestNodeSize - 2; i++ {
		fmt.Println("Join Round ", i)
		addr := nodeAddresses[nodesInNetwork[rand.Intn(len(nodesInNetwork))]]
		fmt.Println(nodes[i].Join(addr))
		nodesInNetwork = append(nodesInNetwork, i)
		fmt.Println("Join Finish ", i)
		time.Sleep(time.Second)
	}

	time.Sleep(time.Second * 5)

	fmt.Println("Join Finish.")

	nodes[3].Delete("FDU")
	nodes[3].Delete("ZJU")
	nodes[3].Quit()

	nodes[4].Join(nodeAddresses[0])

	time.Sleep(time.Second * 5)

	fmt.Println("Join 4")

	for i := 0; i <= myselfTestNodeSize; i++ {
		nodes[i].receiver.Node.Display()
	}


	time.Sleep(time.Second * 10)

	/*
	nodes[nodesInNetwork[rand.Intn(len(nodesInNetwork)-1)]].Put("FDU", "Shanghai")
	nodes[nodesInNetwork[rand.Intn(len(nodesInNetwork)-1)]].Put("ZJU", "Zhejiang")

	for i := 0; i <= myselfTestNodeSize; i++ {
		nodes[i].receiver.Node.Display()
	}


	for i := 0; i < 5; i++ {
		fmt.Printf("%v %v ", i, testKey[i])
		fmt.Println(nodes[nodesInNetwork[rand.Intn(len(nodesInNetwork)-1)]].Get(testKey[i]))
	}*/


	/*nodes[0].Quit()

	time.Sleep(time.Second * 10)

	for i := 0; i <= myselfTestNodeSize; i++ {
		nodes[i].receiver.Node.Display()
	}

	for i := 0; i < 5; i++ {
		fmt.Println(nodes[nodesInNetwork[rand.Intn(len(nodesInNetwork)-1)+1]].Get(testKey[i]))
	}*/

	/*for i := 0; i <= myselfTestNodeSize; i++ {
		nodes[i].receiver.Node.Display()
	}

	for i := 0; i <= kvPairSize; i++ {
		nodes[nodesInNetwork[rand.Intn(len(nodesInNetwork))]].Put(strconv.Itoa(i), strconv.Itoa(i))
	}

	for i := 0; i <= kvPairSize; i++ {
		fmt.Println(nodes[nodesInNetwork[rand.Intn(len(nodesInNetwork))]].Get(strconv.Itoa(i)))
	}

	for i := 0; i <= kvPairSize; i+=2 {
		nodes[nodesInNetwork[rand.Intn(len(nodesInNetwork))]].Delete(strconv.Itoa(i))
	}

	for i := 0; i <= kvPairSize; i++ {
		fmt.Println(nodes[nodesInNetwork[rand.Intn(len(nodesInNetwork))]].Get(strconv.Itoa(i)))
	}

	for i := 0; i < 10; i++ {
		nodes[i].Quit()
		time.Sleep(time.Second)
	}

	time.Sleep(10 * time.Second)

	for i := 0; i <= myselfTestNodeSize; i++ {
		nodes[i].receiver.Node.Display()
	}

	for i := 0; i <= kvPairSize; i++ {
		fmt.Println(nodes[nodesInNetwork[rand.Intn(len(nodesInNetwork)-10)+10]].Get(strconv.Itoa(i)))
	}

	for i := 0; i <= kvPairSize; i+=2 {
		nodes[nodesInNetwork[rand.Intn(len(nodesInNetwork)-10)+10]].Put(strconv.Itoa(i), strconv.Itoa(i))
	}

	for i := 0; i <= kvPairSize; i++ {
		fmt.Println(nodes[nodesInNetwork[rand.Intn(len(nodesInNetwork)-10)+10]].Get(strconv.Itoa(i)))
	}*/

	/*

	for i := 0; i <= myselfTestNodeSize; i++ {
		nodes[i].receiver.Node.Display()
	}

	for i := 0; i < 5; i++ {
		_, city := nodes[nodesInNetwork[rand.Intn(len(nodesInNetwork))]].Get(testKey[i])
		fmt.Println(testKey[i], "@", city)
	}

	nodes[nodesInNetwork[rand.Intn(len(nodesInNetwork))]].Delete("SJTU")

	fmt.Println("SJTU Deleted.")

	founded, city := nodes[nodesInNetwork[rand.Intn(len(nodesInNetwork))]].Get("SJTU")

	if !founded {
		fmt.Println("什么破学校 找不到捏")
	} else {
		fmt.Println(city)
	}
	*/

	/*nodes[0].Quit()

	time.Sleep(5 * time.Second)

	for i := 0; i <= myselfTestNodeSize; i++ {
		nodes[i].receiver.Node.Display()
	}*/

	//time.Sleep(100000 * time.Millisecond)
}
