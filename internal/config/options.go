package config

import (
	"flag"
)

type Options struct {
	ServerAddress string
	BaseURL       string
}

func NewOptions() *Options {
	return &Options{
		ServerAddress: "localhost:8888",
		BaseURL:       "http://localhost:8000/",
	}
}

func (o *Options) OptionsInit() {
	serverAddress := flag.String("a", "localhost:8888", "адрес HTTP-сервера")
	baseURL := flag.String("b", "http://localhost:8000/", "базовый адрес URL")

	flag.Parse()

	o.ServerAddress = *serverAddress
	o.BaseURL = *baseURL
}
