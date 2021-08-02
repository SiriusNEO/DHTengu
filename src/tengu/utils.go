package tengu

import (
	"fmt"
	"github.com/fatih/color"
	"net"
	"time"
)

var (
	green  = color.New(color.FgGreen)
	red    = color.New(color.FgRed)
	yellow = color.New(color.FgYellow)
	cyan   = color.New(color.FgCyan)
	blue   = color.New(color.FgBlue)
	hiBlue = color.New(color.FgHiBlue)
	magenta = color.New(color.FgHiMagenta)
)

const (
	 SHA1Len = 20
	 SHA1StrLen = 40

	 DefaultTorrentPath = "torrent/"
	 DefaultUploadPath = "upload/"
	 DefaultDownloadPath = "download/"
	 DefaultMusicPath = "music/"

	 DefaultFileName = "file"
	 DefaultMusicName = "music"
	 DefaultAlbumName = "Default Album"

	 SongDelim = '$'

	 PieceSize = 1048576 //1MB
	 WorkQueueBuffer = 1024

	 AfterLoginSleep = time.Second
	 AfterQuitSleep = time.Second

	 UploadTimeout = time.Second
	 DownloadTimeout = time.Second

	 RetryTime = 2

	 UploadInterval = 100 * time.Millisecond
	 DownloadInterval = 100 * time.Millisecond
	 DownloadWriteInterval = time.Second
)

func MakeMagnet(infoHash string) string {
	return "magnet:?xt=urn:btih:" + infoHash
}

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
