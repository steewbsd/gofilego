package main

import (
	"filego"
	"flag"
	"fmt"
	"os"
)

func main() {
	inputFile := flag.String("file", "", "Path of the file to be uploaded")
	email := flag.String("email", "", "Email of the upload")
	desc := flag.String("description", "", "File description")
	pass := flag.String("password","", "File upload password")

	flag.Parse()

	conn,_ := filego.NewConnection()
	upload := &filego.Connection{Email: *email, Description: *desc, Password: *pass}
	conn.Construct(upload)
	b, err := os.Open(*inputFile)
	if err != nil {
		fmt.Println("Specified file is not valid or couldn't be read")
		return
	}
	conn.AddFile(*inputFile, b)
	rs, _ := conn.Upload()
	fmt.Println("Status: ", rs.Status)
	fmt.Println("Data: ")
	fmt.Println("\t-Code: ", rs.Data["code"])
	fmt.Println("\t-removalCode: ", rs.Data["removalCode"])
}
