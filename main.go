package main

import (
	"comic_converter/copymanga"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("url mismatch")
		return
	}

	copymanga.DownloadCopymanga(os.Args[1])
}
