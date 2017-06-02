package utils

import (
	"entities"
	"libraries/toml"
	"time"
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

	switch entity.(type) {
	case *entities.Plan:
		if updateEntity.(*entities.Plan).Status != 0 {
			entity.(*entities.Plan).Status = updateEntity.(*entities.Plan).Status
		}
		if updateEntity.(*entities.Plan).ResultStatus != 0 {
			entity.(*entities.Plan).ResultStatus = updateEntity.(*entities.Plan).ResultStatus
		}
		if updateEntity.(*entities.Plan).Prize != 0 {
			entity.(*entities.Plan).Prize = updateEntity.(*entities.Plan).Prize
		}
		if updateEntity.(*entities.Plan).PrizeTax != 0 {
			entity.(*entities.Plan).PrizeTax = updateEntity.(*entities.Plan).PrizeTax
		}
		if updateEntity.(*entities.Plan).SendTicketFinishTime != 0 {
			entity.(*entities.Plan).SendTicketFinishTime = updateEntity.(*entities.Plan).SendTicketFinishTime
		}
		entity.(*entities.Plan).UpdateTime = entities.JsonTime(time.Now())
	case *entities.Phase:
		entity.(*entities.Phase).Status = updateEntity.(*entities.Phase).Status
		entity.(*entities.Phase).UpdateTime = entities.JsonTime(time.Now())
	case *entities.Order:
		entity.(*entities.Order).RewardStatus = updateEntity.(*entities.Order).RewardStatus
		entity.(*entities.Order).Amount = updateEntity.(*entities.Order).Amount
		entity.(*entities.Order).AfterTax = updateEntity.(*entities.Order).AfterTax
		entity.(*entities.Order).UpdateTime = entities.JsonTime(time.Now())
	default:
	}

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
