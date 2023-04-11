package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/iamolegga/enviper"
	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/ukama/ukama/systems/common/sql"
)

// Common properties for all configs.
// Don't forget to use `mapstructure:",squash"`. See unittest for example
type BaseConfig struct {
	DebugMode bool
}

type Database struct {
	Host       string `default:"localhost"`
	Password   string `default:"Pass2020!"`
	DbName     string
	Username   string `default:"postgres"`
	SslEnabled bool   `default:"false"`
	Port       int    `default:"5432"`
}

func (p Database) GetConnString() string {
	sslMode := "disable"
	if p.SslEnabled {
		sslMode = "enable"
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s database=%s port=%d sslmode=%s",
		p.Host, p.Username, p.Password, p.DbName, p.Port, sslMode)
	return dsn
}

func (p Database) ChangeDbName(name string) sql.DbConfig {
	p.DbName = name
	return p
}

func (p Database) GetDbName() string {
	return p.DbName
}

/*
Route are key value pair
*/
type Route map[string]string

/*
Service would a have some finite number of routes
thorugh which it could be reached.
*/
type Pattern struct {
	Routes []Route
}

/*
Service would be listing on this
IP and Port for incoming messages
*/
type Forward struct {
	Ip   string `json:"ip"`
	Port int    `json:"port"`
	Path string `json:"default_path"`
}

/*
Service API interface is registered to Router service.
So that Router service know when and where to reach service.
*/
type ServiceApiIf struct {
	Name string  `json:"name"`
	P    []Route `json:"patterns"`
	F    Forward `json:"forward"`
}

/*
Messaging systems URI
*/
type Queue struct {
	Uri string `default:"amqp://guest:guest@localhost:5672"` // Env var name: QUEUE_URI or in file Queue: { Uri: "" }. Example: QUEUE_URI=amqp://guest:guest@localhost:5672/
}

// SafeString returns URI without password for logging purpose
func (q *Queue) SafeString() string {
	return q.Uri[strings.LastIndex(q.Uri, "@"):]
}

type Grpc struct {
	Port int `default:"9090"`
}

type GrpcService struct {
	Host    string        `default:"localhost:9090"`
	Timeout time.Duration `default:"3s"`
}

/*
Message Client for a system which talks to MsgBus
*/
type MsgClient struct {
	Host           string        `default:"localhost:9095"`
	Timeout        time.Duration `default:"3s"`
	RetryCount     int8          `default:"3"`
	Exchange       string        `default:"amq.topic"`
	ListenQueue    string        `default:""`
	PublishQueue   string        `default:""`
	ListenerRoutes []string
}

type Service struct {
	Host string `default:"localhost"`
	Port string `default:"9090"`
	Uri  string `default:"localhost:9090"`
}

type Metrics struct {
	Port    int  `default:"10250"`
	Enabled bool `default:"true"`
}

type Auth struct {
	AuthServerUrl string `default:"http://localhost:4434"`
	AuthAppUrl    string `default:"http://localhost:4455"`
	AuthAPIGW     string `default:"http://localhost:8080"`
}

// LoadConfig loads configuration into `config` object
// Pulls configuration from env vars and config file
// Config should be a yaml file with `configFileName` and 'yaml' extension
// for example: `.registry`
// Available paths: $HOME, same as binary
// ENV vars have precedence
// configFileName - name of config file without extension.
// Config file should have yaml format and property Names should start with lowercase latter
// Evn var should be uppercased
func LoadConfig(configFileName string, config interface{}) {

	e := enviper.New(viper.New())
	e.SetConfigType("yaml")

	// Find home directory.
	home, err := homedir.Dir()
	if err != nil {
		panic(err)
	}

	e.AddConfigPath(home)
	e.AddConfigPath("./")
	e.SetConfigName(configFileName + ".yaml")

	e.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	err = e.ReadInConfig()
	if err == nil {
		logrus.Info("Using config file:", e.ConfigFileUsed())
	} else {
		logrus.Infof("Config file was not loaded. Reason: %v\n", err)
	}

	err = e.Unmarshal(config)
	if err != nil {
		logrus.Fatalf("Unable to decode into struct, %v", err)
	}

	if e.GetBool("DebugMode") {
		logrus.Infoln("Debug mode is enabled")
		logrus.Infoln("vvvv Config file vvvv")
		logrus.Infof("%+v", config)
	}
}

func DefaultMetrics() *Metrics {
	return &Metrics{
		Enabled: true,
		Port:    10250,
	}
}

func DefaultDatabaseName(name string) Database {
	return Database{
		Host:       "localhost",
		Password:   "Pass2020!",
		DbName:     name,
		Username:   "postgres",
		Port:       5432,
		SslEnabled: false,
	}
}

func DefaultForwardConfig() Forward {
	return Forward{
		Ip:   "localhost",
		Port: 8080,
		Path: "/",
	}
}

func LoadServiceHostConfig(name string) *Service {
	s := &Service{}
	svcHost := "_SERVICE_HOST"
	svcPort := "_SERVICE_PORT"

	val, present := os.LookupEnv(strings.ToUpper(name + svcHost))
	if present {
		s.Host = val
	} else {
		logrus.Errorf("%s server host env not found", name)
	}

	val, present = os.LookupEnv(strings.ToUpper(name + svcPort))
	if present {
		s.Port = val
	} else {
		logrus.Errorf("%s server port env not found", name)
	}

	s.Uri = s.Host + ":" + s.Port

	return s
}

func LoadAuthHostConfig(name string) *Auth {
	s := &Auth{}
	serverUrl := "_SERVER_URL"
	appUrl := "_APP_URL"
	apigwUrl := "_API_GW_URL"

	val, present := os.LookupEnv(strings.ToUpper(name + serverUrl))
	if present {
		s.AuthServerUrl = val
	} else {
		logrus.Errorf("%s server url env not found", name)
	}

	val, present = os.LookupEnv(strings.ToUpper(name + appUrl))
	if present {
		s.AuthAppUrl = val
	} else {
		logrus.Errorf("%s app url env not found", name)
	}

	val, present = os.LookupEnv(strings.ToUpper(name + apigwUrl))
	if present {
		s.AuthAPIGW = val
	} else {
		logrus.Errorf("%s api gw url env not found", name)
	}

	return s
}
