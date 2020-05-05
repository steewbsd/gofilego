package main

import (
	"github.com/0xSteeW/filego/filego"
	"flag"
	"fmt"
	"os"
)

func main() {
	inputFile := flag.String("f", "", "Path of the file to be uploaded")
	help := flag.Bool("h", false, "Show help menu")
	email := flag.String("e", "", "Email of the upload")
	desc := flag.String("d", "", "File description")
	pass := flag.String("p", "", "File upload password")

	flag.Parse()

	if *help {
		fmt.Println("Filego uploader made by SteeW " + filego.VERSION)
		fmt.Println("Usage: (Flag)(Alt flag)(Description)")
		fmt.Println("\t-f file.txt\tRequired flag. Specifies the file path to be uploaded.")
		fmt.Println("\t-e\tOptional. Specifies the upload email.")
		fmt.Println("\t-d\tOptional. Specifies the upload description.")
		fmt.Println("\t-p\tOptional. Specifies the password.")
		return
	}
	conn, _ := filego.NewConnection()
	upload := &filego.Connection{Email: *email, Description: *desc, Password: *pass}
	conn.Construct(upload)
	b, err := os.Open(*inputFile)
	if err != nil {
		fmt.Println("Specified file is not valid or couldn't be read")
		return
	}
	conn.AddFile(*inputFile, b)
	rs, err := conn.Upload()
	if err != nil {
		fmt.Println("Error uploading file: " + err.Error())
	}
	fmt.Println("Status: ", rs.Status)
	fmt.Println("Data: ")
	fmt.Println("\t-Code: ", rs.Data["code"])
	fmt.Println("\t-removalCode: ", rs.Data["removalCode"])
}
