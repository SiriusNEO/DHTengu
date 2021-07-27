package tengu

import (
	"fmt"
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

		var cmd string
		var arg1, arg2, arg3 string
		fmt.Scanln(&cmd, &arg1, &arg2, &arg3)

		switch cmd {
			case "help": {
				yellow.Println("[Tengu Commands]")
				yellow.Println("upload <file-path> <seed-path> <fileName>                #to upload a file")
				yellow.Println("download <file-path> <seed-path> <seedFileName>          #to download a file")
				yellow.Println("help                                                     #show help")
				yellow.Println("quit                                                     #quit from tengu")
				yellow.Println("\n[Tengu Environment Setting]")
				yellow.Println("Default Torrent Path: ", TorrentPath)
				yellow.Println("Default Upload Path: ", UploadPath)
				yellow.Println("Default Download Path: ", DownloadPath)
				yellow.Println("\n[Tengu Hints]")
				yellow.Println("<", DefaultSymbol ,"> represents default name & path")
			}
			case "quit": {
				self.Quit()
			}
			case "upload": {
				if arg1 == DefaultSymbol {
					arg1 = UploadPath
				}
				if arg2 == DefaultSymbol {
					arg2 = TorrentPath
				}
				if arg3 == DefaultSymbol {
					arg3 = "file"
				}
				keyPackage, dataPackage := UploadFileProcessing(arg1 + arg3, arg3, arg2)
				self.Upload(&keyPackage, &dataPackage)
				yellow.Println("Finish Upload and create torrent to: ", arg2 + arg3 + ".torrent")
			}
			case "download": {
				if arg1 == DefaultSymbol {
					arg1 = DownloadPath
				}
				if arg2 == DefaultSymbol {
					arg2 = TorrentPath
				}
				if arg3 == DefaultSymbol {
					arg3 = "file.torrent"
				}
				keyPackage, fileName := DownloadFileProcessing(arg2 + arg3)

				fp, err := os.Create(arg1 + fileName)

				if err != nil {
					red.Println("File Path invalid in ", arg1)
					continue
				}
				fp.Write(self.DownLoad(&keyPackage))
				time.Sleep(DownloadWriteInterval)
				yellow.Println("Finish Download to: ", arg1 + fileName)
			}
		}
	}
}
