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
		ServerAddress: "localhost:8080",
		BaseURL:       "http://localhost:8080/",
	}
}

func (o *Options) OptionsInit() {
	serverAddress := flag.String("a", "localhost:8080", "адрес HTTP-сервера")
	baseURL := flag.String("b", "http://localhost:8080/", "базовый адрес URL")

	flag.Parse()

	o.ServerAddress = *serverAddress
	o.BaseURL = *baseURL
}
