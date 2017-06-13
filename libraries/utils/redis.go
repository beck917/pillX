package utils

import (
	"fmt"

	"github.com/beck917/pillX/libraries/toml"

	"github.com/garyburd/redigo/redis"
)

var _instanceRedis map[string]*Redis
var GlobalRedis *Redis

func init() {
	_instanceRedis = make(map[string]*Redis)
	//GlobalRedis = InstanceRedis(toml.GlobalTomlConfig.Redis0)
}

type Redis struct {
	Conn   redis.Conn
	Config toml.DBConfig
	Pool   *redis.Pool
}

func InstanceRedis(config toml.DBConfig) *Redis {
	name := config.DBname
	if _instanceRedis[name] == nil {
		_instanceRedis[name] = new(Redis)
		_instanceRedis[name].Config = config
		_instanceRedis[name].Conn = _instanceRedis[name].Dail()
		_instanceRedis[name].Pool = _instanceRedis[name].NewPool()
	}

	return _instanceRedis[name]
}

func (this *Redis) NewPool() *redis.Pool {
	redisconfig := this.Config

	redisConnStr := fmt.Sprintf("%s:%d", redisconfig.Host, redisconfig.Port)
	fmt.Println(redisConnStr)
	redisPool := redis.NewPool(func() (redis.Conn, error) {
		redisClient, err := redis.Dial("tcp", redisConnStr)

		if err != nil {
			return nil, err
		}
		redisClient.Do("AUTH", redisconfig.Password)

		return redisClient, err
	}, 500)

	return redisPool
}

func (this *Redis) Dail() redis.Conn {
	redisconfig := this.Config
	fmt.Println(fmt.Sprintf("%s:%d", redisconfig.Host, redisconfig.Port))
	redisClient, err := redis.Dial("tcp", fmt.Sprintf("%s:%d", redisconfig.Host, redisconfig.Port))
	if err != nil {
		panic(err)
	}
	redisClient.Do("AUTH", redisconfig.Password)

	return redisClient
}

func (this *Redis) HIncrBy(key string, field string, by int) (int64, error) {
	conn := this.Pool.Get()
	defer conn.Close()
	num, err := conn.Do("HINCRBY", key, field, by)

	if err != nil {
		fmt.Println(num)
		panic(err)
	}

	retnum := num.(int64)
	return retnum, err
}

func (this *Redis) HGet(key string, field interface{}) (string, error) {
	conn := this.Pool.Get()
	defer conn.Close()

	res, err := conn.Do("HGET", key, field)

	if err != nil {
		return "", err
	}

	if res == nil {
		return "", nil
	}

	getRes, err := redis.String(res, err)

	return getRes, err
}

func (this *Redis) HGetOrgin(key string, field interface{}) (interface{}, error) {
	conn := this.Pool.Get()
	defer conn.Close()

	res, err := conn.Do("HGET", key, field)

	return res, err
}

func (this *Redis) HSet(key string, field interface{}, value interface{}) error {
	conn := this.Pool.Get()
	defer conn.Close()

	_, err := conn.Do("HSET", key, field, value)

	return err
}
