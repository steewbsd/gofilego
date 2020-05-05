package main

import (
	"filego"
	"flag"
	"fmt"
	"os"
)

func main() {
	inputFile := flag.String("input", "", "Path of the file to be uploaded")

	flag.Parse()

	conn,_ := filego.NewConnection()
	upload := &filego.Connection{Email: "steew0x8@protonmail.com"}
	conn.Construct(upload)
	b, _ := os.Open(*inputFile)
	conn.AddFile(*inputFile, b)
	rs, _ := conn.Upload()
	fmt.Println(rs)
}
