package model

import (
	"time"

	"gorm.io/datatypes"
)

type OperationType string

var (
	Create OperationType = "CREATE"
	Update OperationType = "UPDATE"
	DELETE OperationType = "DELETE"
)

type Audit struct {
	ID            string         `bson:"_id" gorm:"column:id;primaryKey" json:"id"`
	OperationType OperationType  `bson:"operation_type" gorm:"column:operation_type" json:"operation_type"`
	TableName     string         `bson:"table_name" gorm:"column:table_name" json:"table_name"`
	ObjectID      string         `bson:"object_id" gorm:"column:object_id" json:"object_id"`
	Data          datatypes.JSON `bson:"data" gorm:"column:data" json:"data"`
	UserID        string         `bson:"user_id" gorm:"column:user_id" json:"user_id"`
	CreatedAt     time.Time      `bson:"created_at" gorm:"column:created_at;default:now();not null" json:"created_at"`
}
