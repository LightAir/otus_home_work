package config

type Config struct {
	Logger     LoggerConf
	DB         DB
	Server     Server
	GRPCServer GRPCServer
}

type LoggerConf struct {
	Level string
}

type DB struct {
	Type string // "mem", "sql"
	SQL  SQLDatabase
}

type SQLDatabase struct {
	Driver   string
	Name     string
	User     string
	Password string
	Host     string
	Port     string
}

type Server struct {
	Port string
	Host string
}

type GRPCServer struct {
	Port string
	Host string
}

func NewConfig() Config {
	return Config{}
}
