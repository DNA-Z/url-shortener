package infrastructure

import (
	"bufio"
	"encoding/json"
	"os"

	"github.com/DNA-Z/url-shortener/internal/model"
)

type URLConsumer struct {
	file    *os.File
	scanner *bufio.Scanner
}

func NewConsumer(filename string) (*URLConsumer, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	return &URLConsumer{
		file:    file,
		scanner: bufio.NewScanner(file),
	}, nil
}

func (c *URLConsumer) ReadURL() ([]*model.URLFile, error) {
	var urls []*model.URLFile

	for c.scanner.Scan() {
		data := c.scanner.Bytes()

		url := model.URLFile{}
		err := json.Unmarshal(data, &url)
		if err != nil {
			return nil, err
		}

		urls = append(urls, &url)
	}

	if err := c.scanner.Err(); err != nil {
		return nil, err
	}

	return urls, nil
}

func (c *URLConsumer) Close() error {
	return c.file.Close()
}
