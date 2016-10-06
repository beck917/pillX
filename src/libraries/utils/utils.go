package utils

import (
	"fmt"
	"io/ioutil"
	"strconv"

	"log"
	"time"

	"libraries/toml"

	"github.com/coreos/etcd/clientv3"
	"github.com/garyburd/redigo/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"gopkg.in/mgo.v2"
)

func AnyTypeInt(sint interface{}) (ret int) {
	switch sint.(type) {
	case string:
		ret, _ = strconv.Atoi(sint.(string))
		break
	case int:
		ret = sint.(int)
		break
	case int64:
		ret = int(sint.(int64))
		break
	}
	return
}

func LoadFile(filename string) (filedata string, err error) {
	fileDataByte, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("Read failed", err)
		return
	}

	filedata = string(fileDataByte)
	return
}

func MongoDail(mongoconfig toml.DBConfig) *mgo.Session {
	mgoUrl := fmt.Sprintf("%s:%s@%s", mongoconfig.User, mongoconfig.Password, mongoconfig.Host)
	session, err := mgo.Dial(mgoUrl)
	if err != nil {
		panic(err)
	}
	//defer session.Close()
	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
	return session
}

func RedisDail(redisconfig toml.DBConfig) redis.Conn {
	redisClient, err := redis.Dial("tcp", fmt.Sprintf("%s:%d", redisconfig.Host, redisconfig.Port))
	if err != nil {
		panic(err)
	}
	redisClient.Do("AUTH", redisconfig.Password)
	return redisClient
}

func MysqlDail(mysqlconfig toml.DBConfig) *xorm.Engine {
	db, err := xorm.NewEngine("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", mysqlconfig.User, mysqlconfig.Password, mysqlconfig.Host, mysqlconfig.Port, mysqlconfig.DBname))
	if err != nil {
		panic(err)
	}
	return db
}

func MysqlDailName(mysqlconfig toml.DBConfig, dbName string) *xorm.Engine {
	db, err := xorm.NewEngine("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", mysqlconfig.User, mysqlconfig.Password, mysqlconfig.Host, mysqlconfig.Port, dbName))
	if err != nil {
		panic(err)
	}
	return db
}

func EtcdDail(etcdconfig toml.DBConfig) *clientv3.Client {
	cfg := clientv3.Config{
		Endpoints: []string{etcdconfig.Host},
		//Transport:   client.DefaultTransport,
		DialTimeout: 5 * time.Second,
	}
	c, err := clientv3.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	return c
}
