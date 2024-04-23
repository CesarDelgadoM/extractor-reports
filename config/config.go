package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Server   *ServerConfig
	Worker   *WorkerConfig
	RabbitMQ *RabbitMQ
	Mongo    *MongoConfig
	Producer *Producer
}

type ServerConfig struct {
	Port string
}

type WorkerConfig struct {
	Pool int
}

type RabbitMQ struct {
	URI      string
	User     string
	Password string
}

type MongoConfig struct {
	URI      string
	User     string
	Password string
	DBName   string
}

type Producer struct {
	DataBus *DataBus
	Branch  *Branch
}

type DataBus struct {
	ExchangeName string
	ExchangeType string
	ContentType  string
}

type Branch struct {
	ExchangeName string
	ExchangeType string
	ContentType  string
}

func LoadConfig(filename string) *viper.Viper {
	v := viper.New()

	v.SetConfigName(filename)
	v.SetConfigType("yaml")
	v.AddConfigPath("../config")
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		log.Panic(err)
	}

	return v
}

func ParseConfig(v *viper.Viper) *Config {
	var c Config

	if err := v.Unmarshal(&c); err != nil {
		log.Panic(err)
	}

	return &c
}
