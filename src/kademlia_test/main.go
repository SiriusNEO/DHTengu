package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

const (
	myselfTestNodeSize = 20
	kvPairSize = 100
)

func main() {
	testKey := [] string {"THU", "PKU", "FDU", "SJTU", "ZJU", "NJU", "USTC"}
	testVal := [] string {"Beijing", "Beijing", "Shanghai", "Shanghai", "Zhejiang", "Nanjing", "Anhui"}

	firstPort := 20000

	nodes := new([myselfTestNodeSize + 1]PubNodeType)
	nodeAddresses := new([myselfTestNodeSize + 1]string)
	nodesInNetwork := make([]int, 0, myselfTestNodeSize+1)

	for i := 0; i <= myselfTestNodeSize; i++ {
		nodes[i] = *NewPubNode("localhost:" + strconv.Itoa(firstPort + i))
		nodeAddresses[i] = "localhost:" + strconv.Itoa(firstPort + i)
		go nodes[i].Run()
	}

	time.Sleep(time.Second)

	nodes[0].Create()
	nodesInNetwork = append(nodesInNetwork, 0)

	for i := 1; i <= myselfTestNodeSize; i++ {
		fmt.Println("Join Round ", i)
		addr := nodeAddresses[nodesInNetwork[rand.Intn(len(nodesInNetwork))]]
		fmt.Println(nodes[i].Join(addr))
		nodesInNetwork = append(nodesInNetwork, i)
		fmt.Println("Join Finish ", i)
		time.Sleep(time.Second)
	}

	time.Sleep(time.Second)
	for i := 0; i <= myselfTestNodeSize; i++ {
		nodes[i].receiver.Node.Display()
	}

	for i := 0; i < 7; i++ {
		fmt.Println(nodes[nodesInNetwork[rand.Intn(len(nodesInNetwork))]].Put(testKey[i], testVal[i]))
		time.Sleep(time.Second)
	}


	for i := 0; i < 7; i++ {
		fmt.Println(nodes[nodesInNetwork[rand.Intn(len(nodesInNetwork))]].Get(testKey[i]))
		time.Sleep(time.Second)
	}

	for i := 0; i <= kvPairSize; i++ {
		fmt.Println("Put ", i)
		nodes[nodesInNetwork[rand.Intn(len(nodesInNetwork))]].Put(strconv.Itoa(i), strconv.Itoa(i))
	}

	time.Sleep(time.Second)

	for i := 0; i <= kvPairSize; i++ {
		fmt.Println(nodes[nodesInNetwork[rand.Intn(len(nodesInNetwork))]].Get(strconv.Itoa(i)))
	}

	for i := myselfTestNodeSize / 2 + 1; i <= myselfTestNodeSize; i++ {
		nodes[i].Quit()
		fmt.Println("Node ", i, " Quit.")
		time.Sleep(10 * time.Second)
	}

	for i := 0; i <= kvPairSize; i++ {
		fmt.Println(nodes[nodesInNetwork[rand.Intn(10)]].Get(strconv.Itoa(i)))
	}

	time.Sleep(5 * time.Second)

	for i := 0; i <= myselfTestNodeSize; i++ {
		nodes[i].receiver.Node.Display()
	}
}
