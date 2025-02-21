package fileio

import (
	"fmt"
	"log"
	"os"
)

func WriteNewText(filename, text string) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err := file.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	fmt.Fprintln(file, text)
}
