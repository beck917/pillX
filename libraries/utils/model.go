package utils

import (
	"libraries/toml"
)

type Model struct {
	DB        *Database
	Redis     *Redis
	TableName string
}

func (this *Model) Construct(tableName string) {
	//this.DB = InstanceDatabase(toml.GlobalTomlConfig.MysqlLottery)

	this.Redis = InstanceRedis(toml.GlobalTomlConfig.Redis0)
	this.TableName = tableName
}

func (this *Model) Update(id interface{}, entity interface{}, updateEntity interface{}) (err error) {
	//效率微微有点低,这里做一下读取操作
	data, _ := this.Redis.HGet("$:"+this.TableName, id)

	UnPack([]byte(data), entity)

	packed, _ := Pack(entity)
	//更新redis
	err = this.Redis.HSet("$:"+this.TableName, id, packed)

	if err != nil {
		return err
	}

	//更新表
	_, err = this.DB.XORM.Id(id).Update(updateEntity)

	return err
}
