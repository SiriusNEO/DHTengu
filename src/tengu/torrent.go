package tengu

import (
	"bytes"
	"crypto/sha1"
	"github.com/jackpal/bencode-go"
	"os"
	"strconv"
)

//credit to https://github.com/veggiedefender/torrent-client/blob/master/torrentfile/torrentfile.go

type bencodeInfo struct {
	Pieces      string `bencode:"pieces"`
	PieceLength int    `bencode:"piece length"`
	Length      int    `bencode:"length"`
	Name        string `bencode:"name"`
}

type bencodeTorrent struct {
	Announce string      `bencode:"announce"`
	Info     bencodeInfo `bencode:"info"`
}

func (this *bencodeInfo) Display() {
	yellow.Println("\n", "* FileName: ", this.Name, " Size:", this.Length, " bytes", "\n")
}

// Open parses a torrent file
func Open(path string) (*bencodeTorrent, error) {
	file, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	bto := bencodeTorrent{}
	err = bencode.Unmarshal(file, &bto)
	if err != nil {
		return nil, err
	}

	return &bto, nil
}

//InfoHash hash the bencodeInfo
func (i *bencodeInfo) InfoHash() ([SHA1Len]byte, error) {
	var buf bytes.Buffer
	err := bencode.Marshal(&buf, *i)
	if err != nil {
		return [20]byte{}, err
	}
	h := sha1.Sum(buf.Bytes())
	return h, nil
}

//PiecesHash hash one piece, file & index unique
func PiecesHash(piece DataPiece, index int) ([SHA1Len]byte, error) {
	piece = append(piece, []byte(strconv.Itoa(index))...)
	pieceHash := sha1.Sum(piece)
	return pieceHash, nil
}

//Save writes a torrent file into target path
func (this *bencodeTorrent) Save(path string) error {
	fp, err := os.Create(path)
	if err != nil {
		return err
	}
	torrent := bencodeTorrent{
		Announce: "",
		Info: this.Info,
	}
	err = bencode.Marshal(fp, torrent)
	if err != nil {
		return err
	}
	return nil
}

