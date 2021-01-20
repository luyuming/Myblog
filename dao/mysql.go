package dao

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

var (
	DB *gorm.DB
)

type mysqlcfg struct {
	User     string `yaml:user`
	Password string `yaml:password`
	Host     string `yaml:host`
	Port     string `yaml:port`
	Db       string `yaml:db`
}

func getMysqlcfg(icfg *mysqlcfg) {
	pwd, _ := os.Getwd()
	filepath := fmt.Sprintf("%s/dao/mysql.yaml", pwd)
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

/* root:luyuming@/GoTest?charset=utf8&parseTime=True&loc=Local */
func InitMysql() (err error) {
	mysqlcfg := mysqlcfg{}
	getMysqlcfg(&mysqlcfg)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", mysqlcfg.User, mysqlcfg.Password,
		mysqlcfg.Host, mysqlcfg.Port, mysqlcfg.Db)
	//fmt.Println("dsn: ", dsn)
	DB, err = gorm.Open("mysql", dsn)
	return
}

func Clsoe() {
	DB.Close()
}
