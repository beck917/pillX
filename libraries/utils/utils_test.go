package utils

import (
	"entities"
	"fmt"
	"models"
	"testing"
)

func Test_UpdateStatus(t *testing.T) {
	planModel := models.NewPlan()
	updateEntity := &entities.Plan{}
	updateEntity.Status = 5
	planModel.Model.UpdateStatus(317, new(entities.Plan), updateEntity)
}
