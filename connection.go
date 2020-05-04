package filego

import (
	"encoding/json"
	"errors"
	"fmt"
	//"io"
	"io/ioutil"
	"multipart"
	"net/http"
	"time"
)

type Connection struct {
	// Attribute for FormFiles
	Server      string
	Files       []*File
	Email       string
	Description string
	Password    string
	Tags        []string
	Expire      int64
}

type File struct {
}

type bestServerResponse struct {
	Status string
	Data   map[string]string
}

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

// Setters
/*func (conn *Connection) AddFile(reader *io.Reader) {
	conn.Files = append(conn.Files, file)
}
*/
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

func (conn *Connection) Send() error {
	return nil
}

// Returns a new empty connection struct
func NewConnection() (*Connection, error) {
	conn := new(Connection)
	err := conn.GetNewServer()
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func NewDefaultConnection() {}
