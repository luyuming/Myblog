package dao

import (
	"fmt"
	"github.com/go-redis/redis"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

var (
	RDB *redis.Client
)

type rediscfg struct {
	Host     string `yaml:host`
	Port     string `yaml:port`
	Password string `yaml:password`
	Db       int    `yaml:db`
}

func getRediscfg(icfg *rediscfg) {
	pwd, _ := os.Getwd()
	filepath := fmt.Sprintf("%s/dao/redis.yaml", pwd)
	file, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(data, icfg)
	if err != nil {
		panic(err)
	}
}

func InitRedis() (err error) {
	rediscfg := rediscfg{}
	getRediscfg(&rediscfg)
	RDB = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", rediscfg.Host, rediscfg.Port),
		Password: rediscfg.Password, // no password set
		DB:       rediscfg.Db,       // use default DB
	})

	_, err = RDB.Ping().Result()
	return
}

func RClose() {
	RDB.Close()
}
