package common

import "github.com/Benyam-S/asseri/entity"

// ICommonRepository is an interface that defines all the common repository methods
type ICommonRepository interface {
	IsUnique(columnName string, columnValue interface{}, tableName string) bool
	CreateJobAttribute(newAttribute *entity.JobAttribute, tableName string) error
	FindJobAttribute(identifier, tableName string) (*entity.JobAttribute, error)
	AllJobAttributes(tableName string) []*entity.JobAttribute
	UpdateJobAttribute(attribute *entity.JobAttribute, tableName string) error
	DeleteJobAttribute(identifier, tableName string) (*entity.JobAttribute, error)
}
