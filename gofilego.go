package gofilego

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// VERSION represents the version of the library
const VERSION = "v0.2"

// Connection represents a basic Gofile upload struct
type Connection struct {
	// Upload server name
	Server string
	// Map of files, represented by [filename]io.Reader
	FilesUploaded map[string]io.Reader
	// Upload email
	Email string
	// Upload description
	Description string
	// Upload Password
	Password string
	// Upload tags provided as string array
	Tags []string
	// Expiration date timestamp
	Expire int64
}

type bestServerResponse struct {
	Status string
	Data   map[string]string
}

// UploadResponse is a struct representing the json response of the Gofile server
type UploadResponse struct {
	Status string
	Data   map[string]interface{}
}

// GetNewServer refreshes best server via contacting gofile api
func (conn *Connection) GetNewServer() error {
	const getServerUrl = "https://apiv2.gofile.io/getServer"
	request, err := http.Get(getServerUrl)
	if err != nil {
		return err
	}
	defer request.Body.Close()
	response, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return err
	}
	// Parse response
	bestServer := new(bestServerResponse)
	err = json.Unmarshal(response, bestServer)
	if err != nil {
		return err
	}
	if bestServer.Status == "ok" {
		conn.Server = bestServer.Data["server"]
		return nil
	}
	return errors.New("Couldn't obtain new server")
}

func (conn *Connection) noRepeat(name string) string {
	if _, exists := conn.FilesUploaded[name]; exists {
		name = "(Copy)- " + name
		conn.noRepeat(name)
	}
	return name
}

// Setters

func (conn *Connection) AddFile(name string, file io.Reader) {
	key := conn.noRepeat(name)
	conn.FilesUploaded[key] = file
}

func (conn *Connection) SetEmail(email string) {
	conn.Email = email
}

func (conn *Connection) SetDescription(desc string) {
	conn.Description = desc
}

func (conn *Connection) SetPassword(pass string) {
	conn.Password = pass
}

func (conn *Connection) SetTags(tags ...string) {
	conn.Tags = tags
}

func (conn *Connection) SetExpire(timeStamp time.Time) {
	conn.Expire = timeStamp.Unix()
}

// Creates a new default connection. Returns a basic connection struct.
func NewConnection() (*Connection, error) {
	conn := new(Connection)
	err := conn.build()
	if err != nil {
		return conn, err
	}
	return conn, nil
}

func (conn *Connection) build() error {
	err := conn.GetNewServer()
	conn.FilesUploaded = make(map[string]io.Reader)
	if err != nil {
		return err
	}
	return nil
}

// Acts as a constructor. Requires a Connection struct.
func (conn *Connection) Construct(providedStruct *Connection) {
	*conn = *providedStruct
	conn.build()
}

// Uploads files to gofile, using the Connection struct parameters
func (conn *Connection) Upload() (*UploadResponse, error) {
	url := fmt.Sprintf("https://%s.gofile.io/upload", conn.Server)
	buffer := new(bytes.Buffer)
	multiWriter := multipart.NewWriter(buffer)
	for _, reader := range conn.FilesUploaded {
		var fileWriter io.Writer
		var err error
		if x, ok := reader.(io.Closer); ok {
			defer x.Close()
		}
		if x, ok := reader.(*os.File); ok {
			if fileWriter, err = multiWriter.CreateFormFile("filesUploaded", x.Name()); err != nil {
				return nil, err
			}
		} else {
			if fileWriter, err = multiWriter.CreateFormField("filesUploaded"); err != nil {
				return nil, err
			}
		}
		if _, err := io.Copy(fileWriter, reader); err != nil {
			return nil, err
		}
	}
	conn.generateFormFields(multiWriter)
	multiWriter.Close()

	request, err := http.NewRequest("POST", url, buffer)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", multiWriter.FormDataContentType())
	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	responseJson := new(UploadResponse)
	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(content, responseJson)

	return responseJson, nil
}

func (conn *Connection) generateFormFields(multiWriter *multipart.Writer) {
	if conn.Email != "" {
		multiWriter.WriteField("email", conn.Email)
	}
	if conn.Description != "" {
		multiWriter.WriteField("description", conn.Description)
	}
	if conn.Password != "" {
		multiWriter.WriteField("password", conn.Password)
	}
	if conn.Tags != nil {
		var tags string
		for _, tag := range conn.Tags {
			// TODO sanitize inputs
			if strings.Contains(tag, ",") {
				continue
			}
			tags = tags + "," + tag
		}
		tags = strings.TrimPrefix(tags, ",")
		multiWriter.WriteField("tags", tags)
	}
	if conn.Expire != 0 {
		multiWriter.WriteField("expire", strconv.Itoa(int(conn.Expire)))
	}
}
