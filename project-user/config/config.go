package config

import (
	"log"
	"os"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"gwh.com/project-common/logs"
)

var AppConf = InitConfig()

type Config struct {
	viper *viper.Viper
	SC *ServerConfig
}

type ServerConfig struct {
	Name string
	Addr string
}



func InitConfig() *Config {
	v := viper.New()
	conf := &Config{viper: v}
	workDir, _ := os.Getwd()
	conf.viper.SetConfigName("app")
	conf.viper.SetConfigType("yaml")
	conf.viper.AddConfigPath(workDir + "/config")

	err := conf.viper.ReadInConfig()
	if err != nil {
		log.Fatalln(err)
		return nil
	}
	conf.ReadServerConfig()
	conf.InitZapLog()
	return conf
}


func (c *Config) InitZapLog() {
	//从配置中读取日志配置，初始化日志
	log.Println(c.viper.GetString("zap.DebugFileName"))
	lc := &logs.LogConfig{
		DebugFileName: c.viper.GetString("zap.DebugFileName"),
		InfoFileName:  c.viper.GetString("zap.InfoFileName"),
		WarnFileName:  c.viper.GetString("zap.WarnFileName"),
		MaxSize:       c.viper.GetInt("zap.MaxSize"),
		MaxAge:        c.viper.GetInt("zap.MaxAge"),
		MaxBackups:    c.viper.GetInt("zap.MaxBackups"),
	}
	err := logs.InitLogger(lc)
	if err != nil {
		log.Fatalln(err)
	}
}

func (c *Config) InitRedisOptions() *redis.Options {
	return &redis.Options{
		Addr:     c.viper.GetString("redis.host") + ":" + c.viper.GetString("redis.port"),
		Password: c.viper.GetString("redis.password"), // no password set
		DB:       c.viper.GetInt("redis.db"),                // use default DB
	}
}

func (c *Config) ReadServerConfig(){
	sc := &ServerConfig{
		Name: c.viper.GetString("server.name"),
		Addr: c.viper.GetString("server.addr"),
	}
	c.SC = sc
}
