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

	//localAddress := GetLocalAddress()
	localAddress := "localhost"
	this.addr = portToAddr(localAddress, port)
	this.node = NewNode(port)
	go this.node.Run()
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
	green.Println("Client Login Success in port: ", port)
}

func (this *Peer) Quit() {
	this.node.Quit()
	time.Sleep(AfterQuitSleep)
	green.Println("Client Quit Success.")
}

type Worker struct {
	index     int
	retry     int
	success      bool
	result    DataPiece
}

func (this *Peer) uploadPieceWork(keyPackage *KeyPackage, dataPackage *DataPackage, index int, retry int, workQueue *chan Worker) {
	key := keyPackage.getKey(index)
	//fmt.Println(key, string(dataPackage.data[i]))
	ok := this.node.Put(key, string(dataPackage.data[index]))
	if ok {
		*workQueue <- Worker{index: index, success: true, retry: retry}
		time.Sleep(UploadInterval)
	} else {
		*workQueue <- Worker{index: index, success: false, retry: retry}
	}
}

func (this *Peer) Upload(keyPackage *KeyPackage, dataPackage *DataPackage) bool {
	if keyPackage.size != dataPackage.size {
		fmt.Println("Package Error, Key & Data don't match.")
		return false
	}

	workQueue := make(chan Worker, WorkQueueBuffer)

	for i := 0; i < keyPackage.size; i++ {
		go this.uploadPieceWork(keyPackage, dataPackage, i, 0, &workQueue)
	}

	var cnt int

	for cnt != keyPackage.size{
		select {
		case worker := <-workQueue:
			{
				if worker.success {
					cnt++
					yellow.Printf("Piece #%v Upload Finish. (%.2f", worker.index+1, float64(cnt*100)/float64(keyPackage.size))
					yellow.Println("%)")
					time.Sleep(UploadInterval)
				} else {
					red.Printf("Piece #%v Upload Failed. (%.2f", worker.index+1, float64(cnt*100)/float64(keyPackage.size))
					red.Println("%)")
					time.Sleep(UploadInterval)
					if worker.retry < RetryTime {
						go this.uploadPieceWork(keyPackage, dataPackage, worker.index, worker.retry+1, &workQueue)
					} else {
						red.Println("Upload Failed: Retry Too Much, Killed")
						return false
					}
				}
			}
		case <-time.After(time.Duration(int64(UploadTimeout) * int64(keyPackage.size))):
			{
				red.Println("Upload Failed: TimeOut!")
				return false
			}
		}
	}

	green.Println("Uploading to network Success")
	return true
}

func (this *Peer) downloadPieceWork(keyPackage *KeyPackage, index int, retry int, workQueue *chan Worker) {
	key := keyPackage.getKey(index)
	ok, piece := this.node.Get(key)
	//fmt.Println(key, piece)
	if ok {
		//ret = append(ret, []byte(piece)[:]...)
		*workQueue <- Worker{index: index, success: true, result: []byte(piece), retry: retry}
	} else {
		*workQueue <- Worker{index: index, success: false, retry: retry}
	}
}

func (this *Peer) DownLoad(keyPackage *KeyPackage) (bool, []byte) {
	ret := make([]byte, keyPackage.length)
	workQueue := make(chan Worker, WorkQueueBuffer)
	for i := 0; i < keyPackage.size; i++ {
		go this.downloadPieceWork(keyPackage, i, 0, &workQueue)
	}

	var cnt int

	checker := make([]KeyPiece, keyPackage.size)

	for cnt != keyPackage.size{
		select {
		case worker := <-workQueue:
			{
				if worker.success {
					cnt++
					bound := (worker.index+1) * PieceSize
					if bound > keyPackage.length {
						bound = keyPackage.length
					}
					copy(ret[worker.index*PieceSize:bound], worker.result)
					pieceHash, _ := PiecesHash(worker.result, worker.index)
					checker[worker.index] = pieceHash
					yellow.Printf("Piece #%v Download Finish. (%.2f", worker.index+1, float64(cnt*100)/float64(keyPackage.size))
					yellow.Println("%)")
					time.Sleep(DownloadInterval)
				} else {
					red.Printf("Piece #%v Download Failed. Retry Times: %d (%.2f", worker.index+1, worker.retry, float64(cnt*100)/float64(keyPackage.size))
					red.Println("%)")
					time.Sleep(DownloadInterval)
					if worker.retry < RetryTime {
						go this.downloadPieceWork(keyPackage, worker.index, worker.retry+1, &workQueue)
					} else {
						red.Println("Download Failed: Retry Too Much, Killed")
						return false, []byte{}
					}
				}
			}
		case <-time.After(time.Duration(int64(DownloadTimeout) * int64(keyPackage.size))):
			{
				red.Println("Download Failed: TimeOut!")
				return false, []byte{}
			}
		}
	}

	green.Println("Download Data Success")

	yellow.Println("Check Integrity...")

	var checkerPieces string
	for i := 0; i < keyPackage.size; i++ {
		checkerPieces += fmt.Sprintf("%x", checker[i])
	}
	buf := []byte(checkerPieces)
	for i := 0; i < keyPackage.size; i++ {
		copy(checker[i][:], buf[i*SHA1StrLen:(i+1)*SHA1StrLen])
		if checker[i] != keyPackage.key[i] {
			red.Println("Check Failed in Piece #", i+1)
			return false, []byte{}
		}
	}

	green.Println("Check Success")
	return true, ret
}
