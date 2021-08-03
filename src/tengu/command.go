package tengu

import (
	"bytes"
	"chord"
	"fmt"
	"github.com/jackpal/bencode-go"
	"math/rand"
	"os"
	"runtime"
	"time"
)

var self Peer
var port int
var bootstrapAddr string

func Welcome() {
	hiBlue.Println("Hello, this is Tengu, Welcome.")

	//Init
	chord.LogInit()
	rand.Seed(time.Now().UnixNano())
	runtime.GOMAXPROCS(runtime.NumCPU())

	hiBlue.Println("* Please Input your port and bootstrap address")
	fmt.Scanln(&port, &bootstrapAddr)
	self.Login(port, bootstrapAddr)

	time.Sleep(AfterLoginSleep)

	for {
		hiBlue.Println("Tengu beta 1.0")
		hiBlue.Println("Type \"help\" for more infomation.")

		var cmd, fp, sp, fn, sn, mg, so, mp, al string
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
				case "so": so = args[i][4:]
				case "mp":mp = args[i][4:]
				case "al":al = args[i][4:]
			}
		}

		switch cmd {
			case "help": {
				yellow.Println("Tengu is a File Sharing System based on DHT")
				yellow.Println("\nTengu Commands:\n")
				yellow.Println("	upload -fp=<file-path> -sp=<seed-path> -fn=<fileName>                            #to upload a file")
				yellow.Println("	download -fp=<file-path> -sp=<seed-path> -sn=<seedFileName> -mg=<magnet>         #to download a file")
				yellow.Println("	music-upload -mp=<music-path> -sp=<seed-path> -so=<song-name> -al=<album>        #to upload a song")
				yellow.Println("	music-play -al=<album> -so=<song-name> -sp=<seed-path>                           #to play a song.")
				yellow.Println("	help                                                                             #show help")
				yellow.Println("	quit                                                                             #quit from tengu")
				yellow.Println("\nTengu Environment Setting:\n")
				yellow.Println("	Default Torrent Path: ", DefaultTorrentPath)
				yellow.Println("	Default Upload Path: ", DefaultUploadPath)
				yellow.Println("	Default Download Path: ", DefaultDownloadPath)
				yellow.Println("	Default Music Path: ", DefaultMusicPath)
				yellow.Println("\nHints:\n")
				yellow.Println("	You can omit \"-fn\" \"-so\" arguments to upload the whole directory.")
			}
			case "quit": {
				self.Quit()
			}
			case "upload": {
				if fp == "" {
					fp = DefaultUploadPath
				}
				if sp == "" {
					sp = DefaultTorrentPath
				}
				var fileList []string
				if fn == "" {
					dir, err := os.Open(fp)
					if err != nil {
						red.Println("Open File Directory Error!")
					} else {
						list, _ := dir.Readdir(-1)
						for _, file := range list {
							fileList = append(fileList, file.Name())
						}
					}
				} else {
					fileList = append(fileList, fn)
				}
				for _, fileName := range fileList {
					yellow.Println("\nStart Upload File: ", fileName)
					keyPackage, dataPackage, magnet, torrentStr := UploadFileProcessing(fp+fileName, fileName, sp)
					ok := self.Upload(&keyPackage, &dataPackage)
					if ok {
						yellow.Println("Finish Upload and create torrent to: ", sp+fileName+".torrent")
						self.node.Put(magnet, torrentStr)
						yellow.Println("Magnet URL: ", magnet, " saved to: ", sp+fileName+"-magnet.txt")
					}
					time.Sleep(UploadFileInterval)
				}
			}
			case "download": {
				if fp == "" {
					fp = DefaultDownloadPath
				}
				if sp == "" {
					sp = DefaultTorrentPath
				}
				if sn == "" {
					sn = DefaultFileName + ".torrent"
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
			case "music-upload": {
				if mp == "" {
					mp = DefaultMusicPath
				}
				if sp == "" {
					sp = DefaultTorrentPath
				}
				var songList []string
				if so == "" {
					dir, err := os.Open(mp)
					if err != nil {
						red.Println("Open Music Directory Error!")
					} else {
						list, _ := dir.Readdir(-1)
						for _, musicFile := range list {
							songList = append(songList, musicFile.Name())
						}
					}
				} else {
					songList = append(songList, so)
				}
				if al == "" {
					al = DefaultAlbumName
				}
				for _, songName := range songList {
					magenta.Println("\nStart Upload Song: ", songName)
					keyPackage, dataPackage, magnet, torrentStr := UploadFileProcessing(mp+songName, songName, sp)
					ok := self.Upload(&keyPackage, &dataPackage)
					if ok {
						yellow.Println("Finish Upload and create torrent to: ", sp+songName+".torrent")
						self.node.Put(magnet, torrentStr)
						yellow.Println("Magnet URL: ", magnet, " saved to: ", sp+songName+"-magnet.txt")
						self.node.Put(al+"/"+songName, torrentStr)
						founded, listStr := self.node.Get(al)
						if !founded {
							magenta.Println("Create Album: ", al)
							listStr = ""
						}
						listStr += songName + string(SongDelim)
						self.node.Delete(al)
						self.node.Put(al, listStr)
						magenta.Println("Song: ", songName, " has been collected to Album: ", al)
					}
					time.Sleep(UploadFileInterval)
				}
			}
			case "music-play": {
				if al == "" {
					al = DefaultAlbumName
				}
				if sp == "" {
					sp = DefaultTorrentPath
				}
				if mp == "" {
					mp = DefaultMusicPath
				}
				founded, listStr := self.node.Get(al)
				if !founded {
					red.Println("Album Not Founded!")
				} else {
					var songList []string
					lastPos := 0
					for i := 0; i < len(listStr); i++ {
						if listStr[i] == SongDelim {
							songList = append(songList, listStr[lastPos : i])
							lastPos = i + 1
						}
					}
					magenta.Println(al)
					for _, song := range songList {
						magenta.Println("* ", song)
					}
				}
				if so != "" {
					ok, torrentStr := self.node.Get(al + "/" + so)
					if ok {
						reader := bytes.NewBufferString(torrentStr)
						torrent := bencodeTorrent{}
						err := bencode.Unmarshal(reader, &torrent)
						if err != nil {
							red.Println("Failed to create Torrent: torrent broken.")
						} else {
							sn = so + ".torrent"
							torrent.Save(sp + sn)
							green.Println("Song Torrent saved to: ", sp + sn)
						}
					} else {
						red.Println("Failed to find target song")
					}

					magenta.Println("Start to Loading Song...")
					keyPackage, fileName := DownloadFileProcessing(sp + sn)

					fileIO, err := os.Create(mp + fileName)

					if err != nil {
						red.Println("Music Path invalid in ", mp)
						continue
					}
					ok, data := self.DownLoad(&keyPackage)
					if ok {
						fileIO.Write(data)
						time.Sleep(DownloadWriteInterval)
						magenta.Println("Temporary Song File Download to: ", mp + fileName)
					}

					Play(mp + fileName, so, al)
				}
			}
		}
	}
}
