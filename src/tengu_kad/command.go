package tengu

import (
	"bytes"
	"fmt"
	"github.com/jackpal/bencode-go"
	"os"
	"time"
)

var self Peer
var port int
var bootstrapAddr string

func Welcome() {
	hiBlue.Println("Hello, this is Tengu, Welcome.")

	hiBlue.Println("* Please Input your port and bootstrap address")
	fmt.Scanln(&port, &bootstrapAddr)
	self.Login(port, bootstrapAddr)

	time.Sleep(AfterLoginSleep)

	for {
		hiBlue.Println("Tengu beta 1.0")
		hiBlue.Println("Type \"help\" for more infomation.")

		var cmd, fp, sp, fn, sn, mg string
		var args [3]string
		fmt.Scanln(&cmd, &args[0], &args[1], &args[2])

		for i := 0; i < 3; i++ {
			if args[i] == "" {
				continue
			}
			switch args[i][1:3] {
				case "fp": fp = args[i][4:]
				case "sp": sp = args[i][4:]
				case "sn": sn = args[i][4:]
				case "fn": fn = args[i][4:]
				case "mg": mg = args[i][4:]
			}
		}

		switch cmd {
			case "help": {
				yellow.Println("[Tengu Commands]")
				yellow.Println("upload -fp=<file-path> -sp=<seed-path> -fn=<fileName>                            #to upload a file")
				yellow.Println("download -fp=<file-path> -sp=<seed-path> -sn=<seedFileName> -mg=<magnet>          #to download a file")
				yellow.Println("help                                                                             #show help")
				yellow.Println("quit                                                                             #quit from tengu")
				yellow.Println("\n[Tengu Environment Setting]")
				yellow.Println("Default Torrent Path: ", TorrentPath)
				yellow.Println("Default Upload Path: ", UploadPath)
				yellow.Println("Default Download Path: ", DownloadPath)
			}
			case "quit": {
				self.Quit()
			}
			case "upload": {
				if fp == "" {
					fp = UploadPath
				}
				if sp == "" {
					sp = TorrentPath
				}
				if fn == "" {
					fn = "file"
				}
				keyPackage, dataPackage, magnet, torrentStr := UploadFileProcessing(fp + fn, fn, sp)
				ok := self.Upload(&keyPackage, &dataPackage)
				if ok {
					yellow.Println("Finish Upload and create torrent to: ", sp+fn+".torrent")
					self.node.Put(magnet, torrentStr)
					yellow.Println("Magnet URL: ", magnet, " saved to: ", sp+fn+"-magnet.txt")
				}
			}
			case "download": {
				if fp == "" {
					fp = DownloadPath
				}
				if sp == "" {
					sp = TorrentPath
				}
				if sn == "" {
					sn = "file.torrent"
				}
				if mg != "" {
					ok, torrentStr := self.node.Get(mg)
					if ok {
						reader := bytes.NewBufferString(torrentStr)
						torrent := bencodeTorrent{}
						err := bencode.Unmarshal(reader, &torrent)
						if err != nil {
							red.Println("Failed to analysis Magnet URL: torrent broken.")
						} else {
							torrent.Save(sp + sn)
							green.Println("Magnet to Torrent Success! saved to: ", sp + sn)
						}
					} else {
						red.Println("Failed to analysis Magnet URL: torrent not founded.")
					}
				}
				keyPackage, fileName := DownloadFileProcessing(sp + sn)

				fileIO, err := os.Create(fp + fileName)

				if err != nil {
					red.Println("File Path invalid in ", fp)
					continue
				}
				ok, data := self.DownLoad(&keyPackage)
				if ok {
					fileIO.Write(data)
					time.Sleep(DownloadWriteInterval)
					yellow.Println("Finish Download to: ", fp+fileName)
				}
			}
		}
	}
}
