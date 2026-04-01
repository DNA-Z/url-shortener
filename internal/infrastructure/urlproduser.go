package infrastructure

import (
	"bufio"
	"encoding/json"
	"os"

	"github.com/DNA-Z/url-shortener/internal/model"
)

type URLProducer struct {
	file   *os.File
	writer *bufio.Writer
}

func NewURLProducer(filename string) (*URLProducer, error) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return &URLProducer{
		file:   file,
		writer: bufio.NewWriter(file),
	}, nil
}

func (p *URLProducer) WriteURL(file *model.URLFile) error {
	data, err := json.Marshal(&file)
	if err != nil {
		return err
	}

	if _, err := p.writer.Write(data); err != nil {
		return err
	}

	if err := p.writer.WriteByte('\n'); err != nil {
		return err
	}

	return p.writer.Flush()
}
