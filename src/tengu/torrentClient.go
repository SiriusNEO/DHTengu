package tengu

import (
	"fmt"
	"os"
	"time"
)

type Peer struct {
	addr string
	node dhtNode
}

func (this *Peer) Login(port int, bootstrapAddr string) {
	os.Mkdir("upload", os.ModePerm)
	os.Mkdir("download", os.ModePerm)
	os.Mkdir("torrent", os.ModePerm)

	localAddress := GetLocalAddress()
	this.addr = portToAddr(localAddress, port)
	this.node = NewNode(port)
	if bootstrapAddr == "" {
		this.node.Create()
		green.Println("Finish Create the network.")
	} else {
		ok := this.node.Join(bootstrapAddr)
		if ok {
			green.Println("Join the network guided by ", bootstrapAddr)
		} else {
			red.Println("Join Failed.")
			return
		}
	}
	go this.node.Run()
	green.Println("Client Login Success in port: ", port)
}

func (this *Peer) Quit() {
	this.node.Quit()
	time.Sleep(AfterQuitSleep)
	green.Println("Client Quit Success.")
}

type Worker struct {
	index     int //index = -1: failed.
	success      bool
	result    DataPiece
}

func (this *Peer) uploadPieceWork(keyPackage *KeyPackage, dataPackage *DataPackage, index int, workQueue *chan Worker) {
	key := keyPackage.getKey(index)
	//fmt.Println(key, string(dataPackage.data[i]))
	ok := this.node.Put(key, string(dataPackage.data[index]))
	if ok {
		*workQueue <- Worker{index: index+1, success: true}
		time.Sleep(UploadInterval)
	} else {
		*workQueue <- Worker{index: index+1, success: false}
	}
}

func (this *Peer) Upload(keyPackage *KeyPackage, dataPackage *DataPackage) {
	if keyPackage.size != dataPackage.size {
		fmt.Println("Package Error, Key & Data don't match.")
		return
	}

	workQueue := make(chan Worker, WorkQueueBuffer)

	for i := 0; i < keyPackage.size; i++ {
		go this.uploadPieceWork(keyPackage, dataPackage, i, &workQueue)
	}

	var cnt, successCnt int

	for cnt != keyPackage.size{
		worker := <- workQueue
		cnt++
		if worker.success {
			successCnt++
			yellow.Printf("Piece #%v Upload Finish. (%.2f",worker.index, float64(successCnt*100)/float64(keyPackage.size))
			yellow.Println("%)")
			time.Sleep(UploadInterval)
		} else {
			red.Printf("Piece #%v Upload Failed. (%.2f",worker.index, float64(successCnt*100)/float64(keyPackage.size))
			red.Println("%)")
		}
	}

	green.Println("Uploading to network Success")
}

func (this *Peer) downloadPieceWork(keyPackage *KeyPackage, index int, workQueue *chan Worker) {
	key := keyPackage.getKey(index)
	ok, piece := this.node.Get(key)
	//fmt.Println(key, piece)
	if ok {
		//ret = append(ret, []byte(piece)[:]...)
		*workQueue <- Worker{index: index+1, success: true, result: []byte(piece)}
	} else {
		*workQueue <- Worker{index: index+1, success: false}
	}
}

func (this *Peer) DownLoad(keyPackage *KeyPackage) []byte {
	ret := make([]byte, keyPackage.length)
	workQueue := make(chan Worker, WorkQueueBuffer)
	for i := 0; i < keyPackage.size; i++ {
		go this.downloadPieceWork(keyPackage, i, &workQueue)
	}

	var cnt, successCnt int

	for cnt != keyPackage.size{
		worker := <- workQueue
		cnt++
		if worker.success {
			bound := worker.index*PieceSize
			if bound > keyPackage.length {
				bound = keyPackage.length
			}
			copy(ret[(worker.index-1)*PieceSize:bound], worker.result)
			successCnt++
			yellow.Printf("Piece #%v Download Finish. (%.2f",worker.index, float64(successCnt*100)/float64(keyPackage.size))
			yellow.Println("%)")
			time.Sleep(DownloadInterval)
		} else {
			red.Printf("Piece #%v Download Failed. (%.2f",worker.index, float64(successCnt*100)/float64(keyPackage.size))
			red.Println("%)")
		}
	}

	green.Println("Download Data Success, Start to Write")
	return ret
}
