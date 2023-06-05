package client

import (
	"errors"
	"fmt"
	"net/http"
	"encoding/json"
	"io"
	"io/ioutil"
	"github.com/jivid/passman/passman"
)

type PassmanClient struct {
	ServerHost string
	ServerPort string
}

func (c *PassmanClient) ServerUrlBase() string {
	return fmt.Sprintf("http://%s:%s/passwords/", c.ServerHost, c.ServerPort)
}

func (c *PassmanClient) GetAll() ([]passman.PassmanEntry, error) {
	resp, err := http.Get(c.ServerUrlBase())
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusOK {
		return c.readManyJson(resp.Body)
	} else {
		msg, err := c.readErrorJson(resp.Body)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(msg)
	}
}

func (c *PassmanClient) Get(site string) ([]passman.PassmanEntry, error) {
	resp, err := http.Get(c.ServerUrlBase() + site)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusOK {
		return c.readManyJson(resp.Body)
	} else {
		msg, err := c.readErrorJson(resp.Body)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(msg)
	}
}

func (c *PassmanClient) readErrorJson(bytes io.ReadCloser) (string, error) {
	defer bytes.Close()
	data, err := ioutil.ReadAll(bytes)
	if err != nil {
		return "", err
	}
	var resp map[string]string
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return "", err
	}
	return resp["message"], nil
}

func (c *PassmanClient) readManyJson(bytes io.ReadCloser) ([]passman.PassmanEntry, error) {
	defer bytes.Close()
	data, err := ioutil.ReadAll(bytes)
	if err != nil {
		return nil, err
	}
	var entries []passman.PassmanEntry
	err = json.Unmarshal(data, &entries)
	if err != nil {
		return nil, err
	}
	return entries, nil
}
