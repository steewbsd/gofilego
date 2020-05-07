# Gofilego!

Gofilego is a simple Golang wrapper around Gofile.io's API. It provides a basic library for uploading files inside your code and an executable to quickly upload files to Gofile.

# Usage

To install the executable run:

```bash
go install github.com/0xSteeW/gofilego/gofilego-upload
```
To run it:
```bash
gofilego-upload -f test.file
```
Make sure you've set your GOPATH properly.

To get a basic help menu:

```bash
gofilego-upload -h
```
### Library
To just install the library:
```bash
go get -u github.com/0xSteeW/gofilego
```
Then, import it in your code as follows:
```go
import "github.com/0xSteeW/gofilego"
```

# Documentation

Here is a basic example for uploading a file with email and password:

```go
package main

import (
	"fmt"
	"github.com/0xSteeW/gofilego"
	"os"
)

func main() {
	conn, err := gofilego.NewConnection()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	uploadStruct := &gofilego.Connection{Email: "example@email.com", Password: "rockyou"}
	conn.Construct(uploadStruct)
	// Open a file for uploading
	reader, _ := os.Open("test.file")
	conn.AddFile("testfilename", reader)
	response, err := conn.Upload()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(response.Data["code"])
	
}
```

## Index

```go
*Connection* - Struct
	Server	string
	FilesUploaded map[string]io.Reader
	Email	string
	Description	string
	Password	string
	Expire	time.Time
	Tags	[]string

*UploadResponse* - Struct 
	Status	string
	Data map[string]interface{}

func (conn *Connection):
	GetNewServer() error
	NewConnection() (*Connection, error)
	Construct(providedStruct *Connection)
	Upload() (*UploadResponse, error)
	AddFile(fileName string, reader io.Reader)
	AddEmail(email string)	
	AddDescription(desc string)
	AddPassword(pass string)
	AddExpire(tstamp time.Time)
	AddTags(tags ...string)
```
