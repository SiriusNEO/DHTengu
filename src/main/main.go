package main

import "tengu"

func main() {
	tengu.Welcome()
	//tengu.Play("music/Melty_Land_Nightmare.mp3", "MLN", "V")
	/*fp, _ := os.Open("upload/Monet.pptx")

	buf := make([]byte, 1048576)

	var data []byte
	for {
		bufSize, err := fp.Read(buf)

		if err != nil && err != io.EOF {
			return
		}

		if bufSize == 0 {
			fmt.Println(len(data))
			fw, _ := os.Create("upload/Monet1.pptx")
			tmp := string(data)
			fw.Write([]byte(tmp))
			return
		}

		pieceStr := string(buf[:bufSize])

		data = append(data, []byte(pieceStr)...)
	}*/
}

