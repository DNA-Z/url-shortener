// Package config предоставляет конфигурацию для сервиса сокращения URL.
//
// Конфигурация загружается из:
//   - флагов командной строки
//   - переменных окружения
//   - значений по умолчанию
//
// Приоритет: переменные окружения > флаги > значения по умолчанию.
package config

import (
	"bytes"
	"encoding/json"
	"flag"
	"os"
)

// JSONConfig представляет структуру конфигурационного файла.
type JSONConfig struct {
	ServerAddress   string `json:"server_address"`
	BaseURL         string `json:"base_url"`
	FileStoragePath string `json:"file_storage_path"`
	DatabaseDSN     string `json:"database_dsn"`
	EnableHTTPS     bool   `json:"enable_https"`
}

// Options содержит все настройки сервиса.
type Options struct {
	ServerAddress    string
	BaseURL          string
	FileStoragePath  string
	ConnectionString string
	SecretKey        string
	AuditFile        string
	AuditURL         string
	EnableHTTPS      bool
	ConfigFile       string
}

// NewOptions создает новый экземпляр Options со значениями по умолчанию.
func NewOptions() *Options {
	return &Options{
		ServerAddress:    "localhost:8080",
		BaseURL:          "http://localhost:8080/",
		FileStoragePath:  "short_url",
		ConnectionString: "",
		SecretKey:        "superSecretKey",
		AuditFile:        "",
		AuditURL:         "",
		EnableHTTPS:      false,
		ConfigFile:       "",
	}
}

// OptionsInit инициализирует конфигурацию из флагов и переменных окружения.
func (o *Options) OptionsInit() {
	defaultServerAddress := o.ServerAddress
	defaultBaseURL := o.BaseURL
	defaultFileStoragePath := o.FileStoragePath
	defaultConnectionStr := o.ConnectionString
	defaultSecretKey := o.SecretKey
	defaultAuditFile := o.AuditFile
	defaultAuditURL := o.AuditURL
	defaultEnableHTTPS := o.EnableHTTPS
	defaultConfigFile := o.ConfigFile

	var jsonConfig *JSONConfig

	if flag.Lookup("a") == nil {
		serverAddressFlag := flag.String("a", defaultServerAddress, "адрес HTTP-сервера")
		baseURLFlag := flag.String("b", defaultBaseURL, "базовый адрес URL")
		fileStoragePath := flag.String("f", defaultFileStoragePath, "файл в корне проекта")
		connectionStringFlag := flag.String("d", defaultConnectionStr, "строка подключения к БД")
		secretKeyFlag := flag.String("k", defaultSecretKey, "секретный ключ для подписи JWT")
		auditFileFlag := flag.String("audit-file", defaultAuditFile, "аудит запросов с записью логов в файл")
		auditURLFlag := flag.String("audit-url", defaultAuditURL, "URL сервера для отправки логов аудита")
		enableHTTPSFlag := flag.Bool("s", defaultEnableHTTPS, "включить HTTPS")
		configFile := flag.String("c", defaultConfigFile, "конфигурационный файл")

		flag.Parse()

		jsonConfig = o.readConfigFile(configFile)

		o.ServerAddressSet(serverAddressFlag, jsonConfig)
		o.BaseURLSet(baseURLFlag, jsonConfig)
		o.PathToFile(fileStoragePath, jsonConfig)
		o.ConnectionStringSet(connectionStringFlag, jsonConfig)
		o.SecretKeySet(secretKeyFlag)
		o.AuditFileSet(auditFileFlag)
		o.AuditURLSet(auditURLFlag)
		o.EnableHTTPSSet(enableHTTPSFlag, jsonConfig)
		o.ConfigFileSet(configFile)
	} else {
		// Флаги уже проинициализированы
		o.ServerAddressSet(&o.ServerAddress, jsonConfig)
		o.BaseURLSet(&o.BaseURL, jsonConfig)
		o.PathToFile(&o.FileStoragePath, jsonConfig)
		o.ConnectionStringSet(&o.ConnectionString, jsonConfig)
		o.SecretKeySet(&o.SecretKey)
		o.AuditFileSet(&o.AuditFile)
		o.AuditURLSet(&o.AuditURL)
		o.EnableHTTPSSet(&o.EnableHTTPS, jsonConfig)
		o.ConfigFileSet(&o.ConfigFile)
	}
}

func (o *Options) ServerAddressSet(serverAddressFlag *string, jsonConfig *JSONConfig) {
	switch {
	case os.Getenv("SERVER_ADDRESS") != "":
		o.ServerAddress = os.Getenv("SERVER_ADDRESS")
	case *serverAddressFlag != o.ServerAddress:
		o.ServerAddress = *serverAddressFlag
	case jsonConfig != nil && jsonConfig.ServerAddress != "":
		o.ServerAddress = jsonConfig.ServerAddress
	}
}

func (o *Options) BaseURLSet(baseURLFlag *string, jsonConfig *JSONConfig) {
	switch {
	case os.Getenv("BASE_URL") != "":
		o.BaseURL = os.Getenv("BASE_URL")
	case *baseURLFlag != o.BaseURL:
		o.BaseURL = *baseURLFlag
	case jsonConfig != nil && jsonConfig.BaseURL != "":
		o.BaseURL = jsonConfig.BaseURL
	}
}

func (o *Options) PathToFile(fileStoragePath *string, jsonConfig *JSONConfig) {
	switch {
	case os.Getenv("FILE_STORAGE_PATH") != "":
		o.FileStoragePath = os.Getenv("FILE_STORAGE_PATH")
	case *fileStoragePath != o.FileStoragePath:
		o.FileStoragePath = *fileStoragePath
	case jsonConfig != nil && jsonConfig.FileStoragePath != "":
		o.FileStoragePath = jsonConfig.FileStoragePath
	}
}

func (o *Options) ConnectionStringSet(connectionStringFlag *string, jsonConfig *JSONConfig) {
	switch {
	case os.Getenv("DATABASE_DSN") != "":
		o.ConnectionString = os.Getenv("DATABASE_DSN")
	case *connectionStringFlag != o.ConnectionString:
		o.ConnectionString = *connectionStringFlag
	case jsonConfig != nil && jsonConfig.DatabaseDSN != "":
		o.ConnectionString = jsonConfig.DatabaseDSN
	}
}

func (o *Options) SecretKeySet(secretKeyFlag *string) {
	switch {
	case os.Getenv("KEY") != "":
		o.SecretKey = os.Getenv("KEY")
	case *secretKeyFlag != o.SecretKey:
		o.SecretKey = *secretKeyFlag
	}
}

func (o *Options) AuditFileSet(auditFileFlag *string) {
	switch {
	case os.Getenv("AUDIT_FILE") != "":
		o.AuditFile = os.Getenv("AUDIT_FILE")
	case *auditFileFlag != o.AuditFile:
		o.AuditFile = *auditFileFlag
	}
}

func (o *Options) AuditURLSet(auditURLFlag *string) {
	switch {
	case os.Getenv("AUDIT_URL") != "":
		o.AuditURL = os.Getenv("AUDIT_URL")
	case *auditURLFlag != o.AuditURL:
		o.AuditURL = *auditURLFlag
	}
}

func (o *Options) EnableHTTPSSet(enableHTTPSFlag *bool, jsonConfig *JSONConfig) {
	switch {
	case os.Getenv("ENABLE_HTTPS") != "":
		o.EnableHTTPS = os.Getenv("ENABLE_HTTPS") == "true"
	case *enableHTTPSFlag != o.EnableHTTPS:
		o.EnableHTTPS = *enableHTTPSFlag
	case jsonConfig != nil && jsonConfig.EnableHTTPS == true:
		o.EnableHTTPS = jsonConfig.EnableHTTPS
	}
}

func (o *Options) ConfigFileSet(configFileFlag *string) {
	switch {
	case os.Getenv("CONFIG") != "":
		o.ConfigFile = os.Getenv("CONFIG")
	case *configFileFlag != o.ConfigFile:
		o.ConfigFile = *configFileFlag
	}
}

func (o *Options) readConfigFile(configFileFlag *string) *JSONConfig {
	configFile := *configFileFlag
	if configFile == "" {
		configFile = os.Getenv("CONFIG")
	}
	if configFile == "" {
		return nil
	}

	file, err := os.Open(configFile)
	if err != nil {
		return nil
	}
	defer file.Close()

	var buf bytes.Buffer
	_, err = buf.ReadFrom(file)
	if err != nil {
		return nil
	}

	data := buf.Bytes()

	var cfg JSONConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil
	}

	return &cfg
}
