package tengu

import (
	"fmt"
	"io"
	"os"
)

type KeyPiece [SHA1Len]byte
type DataPiece []byte

type DataPackage struct {
	size int
	data []DataPiece
}

type KeyPackage struct {
	size int
	length int
	infoHash [SHA1Len]byte
	key []KeyPiece
}

func (this *KeyPackage) getKey(index int) string {
	var ret KeyPiece
	for i := 0; i < SHA1Len; i++ {
		ret[i] = this.key[index][i] ^ this.infoHash[i]
	}
	return fmt.Sprintf("%x", ret)
}

func UploadFileProcessing(filePath string, fileName string, seedPath string) (KeyPackage, DataPackage) {
	dataPackage, length := makeDataPackage(filePath)

	green.Println("Data Packaged Finish. Total Piece: ", dataPackage.size)

	var pieces string

	for i := 0; i < dataPackage.size; i++ {
		piece, _ := PiecesHash(dataPackage.data[i], i)
		pieces += fmt.Sprintf("%x", piece)
	}

	torrent := bencodeTorrent{
		Announce: "",
		Info: bencodeInfo{
			Length: length,
			Pieces: pieces,
			PieceLength: PieceSize,
			Name: fileName,
		},
	}

	err := torrent.Save(seedPath + fileName + ".torrent")

	if err != nil {
		red.Println("Failed to Make Torrent File ", err.Error())
	}

	keyPackage := torrent.makeKeyPackage()

	green.Println("Torrent Resolved Finish: ")
	torrent.Info.Display()

	return keyPackage, dataPackage
}

func DownloadFileProcessing(seedPath string) (KeyPackage, string) {
	torrent, err := Open(seedPath)

	if err != nil {
		red.Println("Torrent Open Failed in path: ", seedPath, err.Error())
		return KeyPackage{}, ""
	}

	keyPackage := torrent.makeKeyPackage()

	green.Println("Torrent Resolved Finish: ")
	torrent.Info.Display()

	return keyPackage, torrent.Info.Name
}

func makeDataPackage(path string) (DataPackage, int) {
	fp, err := os.Open(path)
	length := 0
	if err != nil {
		red.Println("File Open Failed in path: ", path, err.Error())
		return DataPackage{}, 0
	}

	var ret DataPackage

	for {
		buf := make([]byte, PieceSize)
		bufSize, err := fp.Read(buf)

		if err != nil && err != io.EOF {
			red.Println("File Read Error.")
			return DataPackage{}, 0
		}

		if bufSize == 0 {
			break //finish read
		}

		ret.size++
		length += bufSize
		ret.data = append(ret.data, buf[:bufSize][:])
	}

	return ret, length
}

func (this *bencodeTorrent) makeKeyPackage() KeyPackage {
	var ret KeyPackage
	buf := []byte(this.Info.Pieces)

	ret.size = len(buf) / SHA1StrLen
	ret.infoHash, _ = this.Info.InfoHash()
	ret.length = this.Info.Length
	ret.key = make([]KeyPiece, ret.size)

	for i := 0; i < ret.size; i++ {
		copy(ret.key[i][:], buf[i*SHA1StrLen:(i+1)*SHA1StrLen])
	}
	return ret
}