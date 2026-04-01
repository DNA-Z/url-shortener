package config

import (
	"flag"
	"os"
)

type Options struct {
	ServerAddress   string
	BaseURL         string
	FileStoragePath string
}

func NewOptions() *Options {
	return &Options{
		ServerAddress:   "localhost:8080",
		BaseURL:         "http://localhost:8080/",
		FileStoragePath: "./short_url",
	}
}

func (o *Options) OptionsInit() {
	defaultServerAddress := o.ServerAddress
	defaultBaseURL := o.BaseURL
	defaultFileStoragePath := o.FileStoragePath

	serverAddressFlag := flag.String("a", defaultServerAddress, "адрес HTTP-сервера")
	baseURLFlag := flag.String("b", defaultBaseURL, "базовый адрес URL")
	fileStoragePath := flag.String("f", defaultFileStoragePath, "файл в корне проекта")

	flag.Parse()

	o.ServerAddressSet(serverAddressFlag)
	o.BaseURLSet(baseURLFlag)
	o.PathToFile(fileStoragePath)
}

func (o *Options) ServerAddressSet(serverAddressFlag *string) {
	switch {
	case os.Getenv("SERVER_ADDRESS") != "":
		o.ServerAddress = os.Getenv("SERVER_ADDRESS")
	case *serverAddressFlag != o.ServerAddress:
		o.ServerAddress = *serverAddressFlag
	}
}

func (o *Options) BaseURLSet(baseURLFlag *string) {
	switch {
	case os.Getenv("BASE_URL") != "":
		o.BaseURL = os.Getenv("BASE_URL")
	case *baseURLFlag != o.BaseURL:
		o.BaseURL = *baseURLFlag
	}
}

func (o *Options) PathToFile(fileStoragePath *string) {
	switch {
	case os.Getenv("FILE_STORAGE_PATH") != "":
		o.FileStoragePath = os.Getenv("FILE_STORAGE_PATH")
	case *fileStoragePath != o.FileStoragePath:
		o.FileStoragePath = *fileStoragePath
	}
}
