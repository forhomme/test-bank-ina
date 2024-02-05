package config

import (
	"context"
	"embed"
	"gopkg.in/yaml.v3"
	"os"
	"sync"
	"test_ina_bank/pkg/baselogger"
	"time"
)

var Config *Cfg

type Cfg struct {
	Server             Server                        `yaml:"server"`
	Database           Database                      `yaml:"database"`
	General            General                       `yaml:"general"`
	ApplicationMessage map[string]ApplicationMessage `yaml:"application_message"`

	mutex sync.RWMutex
}

type Server struct {
	Port     int  `yaml:"port"`
	CORS     bool `yaml:"cors"`
	BodyDump bool `yaml:"body_dump"`
}

type Database struct {
	Host     string `yaml:"host"`
	DbName   string `yaml:"db_name"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	MaxConn  int    `yaml:"max_connection"`
	MaxIdle  int    `yaml:"max_idle"`
}

type General struct {
	SecretKey            string        `yaml:"secret_key"`
	TokenDuration        time.Duration `yaml:"token_duration"`
	RefreshTokenDuration time.Duration `yaml:"refresh_token_duration"`
	CurrentLanguage      string        `yaml:"current_language"`
	UptraceDSN           string        `yaml:"uptrace_dsn"`
	AppName              string        `yaml:"app_name"`
	AppVersion           string        `yaml:"app_version"`
	Env                  string
}

type ApplicationMessage struct {
	InvalidStatus MessageFormat `yaml:"invalid_status"`
}

type MessageFormat struct {
	Key     string `yaml:"key"`
	Message string `yaml:"message"`
}

//go:embed *
var files embed.FS

func InitConfig(ctx context.Context, log *baselogger.Logger) {
	var data []byte

	bytes, err := files.ReadFile("config.yaml")
	if err != nil {
		log.Fatal("error when load config: ", err)
	}
	data = bytes

	err = yaml.Unmarshal(data, &Config)
	if err != nil {
		log.Fatal("error when unmarshal config: ", err)
	}
	Config.SetDefault()
}

func (c *Cfg) Reload(ctx context.Context, log *baselogger.Logger) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	InitConfig(ctx, log)
}

func (c *Cfg) SetDefault() {
	if os.Getenv("DATABASE_HOST") != "" {
		c.Database.Host = os.Getenv("DATABASE_HOST")
	}
	if os.Getenv("DATABASE_NAME") != "" {
		c.Database.DbName = os.Getenv("DATABASE_NAME")
	}
	if os.Getenv("DATABASE_USER") != "" {
		c.Database.User = os.Getenv("DATABASE_USER")
	}
	if os.Getenv("DATABASE_PASSWORD") != "" {
		c.Database.Password = os.Getenv("DATABASE_PASSWORD")
	}
}
