package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	pathname := make([]string, 4)
	pathname[0] = "1.ico"
	pathname[1] = "2.png"
	pathname[2] = "Tone Damli - Stupid.mp3"
	pathname[3] = "Drake-Hotline Bling (2016年AMA全美音乐奖最受欢迎说唱-嘻哈歌曲获奖歌曲-2016年VMA奖最佳说唱录影带获奖歌曲)(高清).mp4"
	for i := 0; i <= 3; i++ {
		_, err := os.Stat(pathname[i])
		if os.IsNotExist(err) {
			fmt.Printf("The file %s doesn't exist\n", pathname[i])
			return
		}
		inputFile, err := os.Open(pathname[i])
		if err != nil {
			fmt.Fprintf(os.Stderr, "An error occurred on opening the inputfile: %s\nDoes the file exist?\n", pathname)
		}
		defer inputFile.Close()
		outputFile, err := os.OpenFile(strings.Split(pathname[i], ".")[0]+"_out."+strings.Split(pathname[i], ".")[1], os.O_CREATE|os.O_TRUNC, 0600)
		if err != nil {
			fmt.Println(err)
		}
		defer outputFile.Close()
		//io.Copy(outputFile, inputFile)
		for {
			data := make([]byte, 200)
			_, err1 := inputFile.Read(data)
			_, err2 := outputFile.Write(data)
			if err2 != nil {
				fmt.Println(err2)
				break
			}
			if err1 == io.EOF {
				fmt.Println("Read to end")
				break
			}
		}
	}
	return
}
