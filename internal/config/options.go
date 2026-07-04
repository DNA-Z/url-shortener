package config

import (
	"flag"
	"os"
)

type Options struct {
	ServerAddress    string
	BaseURL          string
	FileStoragePath  string
	ConnectionString string
	SecretKey        string
	AuditFile        string
	AuditURL         string
}

func NewOptions() *Options {
	return &Options{
		ServerAddress:    "localhost:8080",
		BaseURL:          "http://localhost:8080/",
		FileStoragePath:  "short_url",
		ConnectionString: "",
		SecretKey:        "superSecretKey",
		AuditFile:        "",
		AuditURL:         "",
	}
}

func (o *Options) OptionsInit() {
	defaultServerAddress := o.ServerAddress
	defaultBaseURL := o.BaseURL
	defaultFileStoragePath := o.FileStoragePath
	defaultConnectionStr := o.ConnectionString
	defaultSecretKey := o.SecretKey
	defaultAuditFile := o.AuditFile
	defaultAuditURL := o.AuditURL

	if flag.Lookup("a") == nil {
		serverAddressFlag := flag.String("a", defaultServerAddress, "адрес HTTP-сервера")
		baseURLFlag := flag.String("b", defaultBaseURL, "базовый адрес URL")
		fileStoragePath := flag.String("f", defaultFileStoragePath, "файл в корне проекта")
		connectionStringFlag := flag.String("d", defaultConnectionStr, "строка подключения к БД")
		secretKeyFlag := flag.String("k", defaultSecretKey, "секретный ключ для подписи JWT")
		auditFileFlag := flag.String("audit-file", defaultAuditFile, "аудит запросов с записью логов в файл")
		auditURLFlag := flag.String("audit-url", defaultAuditURL, "URL сервера для отправки логов аудита")

		flag.Parse()

		o.ServerAddressSet(serverAddressFlag)
		o.BaseURLSet(baseURLFlag)
		o.PathToFile(fileStoragePath)
		o.ConnectionStringSet(connectionStringFlag)
		o.SecretKeySet(secretKeyFlag)
		o.AuditFileSet(auditFileFlag)
		o.AuditURLSet(auditURLFlag)
	} else {
		// Флаги уже проинициализированы — просто используем текущие значения
		o.ServerAddressSet(&o.ServerAddress)
		o.BaseURLSet(&o.BaseURL)
		o.PathToFile(&o.FileStoragePath)
		o.ConnectionStringSet(&o.ConnectionString)
		o.SecretKeySet(&o.SecretKey)
		o.AuditFileSet(&o.AuditFile)
		o.AuditURLSet(&o.AuditURL)
	}
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

func (o *Options) ConnectionStringSet(connectionStringFlag *string) {
	switch {
	case os.Getenv("DATABASE_DSN") != "":
		o.ConnectionString = os.Getenv("DATABASE_DSN")
	case *connectionStringFlag != o.ConnectionString:
		o.ConnectionString = *connectionStringFlag
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
