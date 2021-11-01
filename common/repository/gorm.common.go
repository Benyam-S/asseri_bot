package repository

import (
	"fmt"

	"github.com/Benyam-S/asseri/common"
	"github.com/Benyam-S/asseri/entity"
	"github.com/Benyam-S/asseri/tools"
	"github.com/jinzhu/gorm"
)

// CommonRepository is a type that defines a repository for common use
type CommonRepository struct {
	conn *gorm.DB
}

// NewCommonRepository is a function that returns a new common repository type
func NewCommonRepository(connection *gorm.DB) common.ICommonRepository {
	return &CommonRepository{conn: connection}
}

// IsUnique is a methods that checks if a given column value is unique in a certain table
func (repo *CommonRepository) IsUnique(columnName string, columnValue interface{}, tableName string) bool {
	return tools.IsUnique(columnName, columnValue, tableName, repo.conn)
}

// ======	======    ======   ======   ======   ======   ======   ======   ======   ======   ======   ======
//    ======	======    ======   ======   ======   ======   ======   ======   ======   ======   ======   ======
// ======	======    ======   ======   ======   ======   ======   ======   ======   ======   ======   ======

// CreateJobAttribute is a method that adds a new job attribute to the database
func (repo *CommonRepository) CreateJobAttribute(newAttribute *entity.JobAttribute, tableName string) error {

	var prefix string

	switch tableName {
	case "job_types":
		prefix = "TYPE"
	case "job_sectors":
		prefix = "SECTOR"
	case "education_levels":
		prefix = "LEVEL"
	}

	totalNumOfMembers := tools.CountMembers(tableName, repo.conn)
	newAttribute.ID = fmt.Sprintf(prefix+"-%s%d", tools.RandomStringGN(7), totalNumOfMembers+1)

	for !tools.IsUnique("id", newAttribute.ID, tableName, repo.conn) {
		totalNumOfMembers++
		newAttribute.ID = fmt.Sprintf(prefix+"-%s%d", tools.RandomStringGN(7), totalNumOfMembers+1)
	}

	err := repo.conn.Table(tableName).Create(newAttribute).Error
	if err != nil {
		return err
	}
	return nil
}

// FindJobAttribute is a method that finds a certain job attribute from the database using an identifier and table name.
// In FindJobAttribute() id and name are used as an key
func (repo *CommonRepository) FindJobAttribute(identifier, tableName string) (*entity.JobAttribute, error) {
	attribute := new(entity.JobAttribute)
	err := repo.conn.Table(tableName).
		Where("id = ? || name = ?", identifier, identifier).
		First(attribute).Error

	if err != nil {
		return nil, err
	}
	return attribute, nil
}

// AllJobAttributes is a method that returns all the job attributes of a single job attribute table in the database
func (repo *CommonRepository) AllJobAttributes(tableName string) []*entity.JobAttribute {
	var attributes []*entity.JobAttribute
	err := repo.conn.Table(tableName).Find(&attributes).Error

	if err != nil {
		return []*entity.JobAttribute{}
	}
	return attributes
}

// UpdateJobAttribute is a method that updates a certain job attribute value in the database
func (repo *CommonRepository) UpdateJobAttribute(attribute *entity.JobAttribute, tableName string) error {

	prevAttribute := new(entity.JobAttribute)
	err := repo.conn.Table(tableName).Where("id = ?", attribute.ID).First(prevAttribute).Error

	if err != nil {
		return err
	}

	err = repo.conn.Table(tableName).Save(attribute).Error
	if err != nil {
		return err
	}
	return nil
}

// DeleteJobAttribute is a method that deletes a certain job attribute from the database using an identifier.
// In DeleteJobAttribute() id and name are used as an key
func (repo *CommonRepository) DeleteJobAttribute(identifier, tableName string) (*entity.JobAttribute, error) {
	attribute := new(entity.JobAttribute)
	err := repo.conn.Table(tableName).Where("id = ? || name = ?", identifier, identifier).First(attribute).Error

	if err != nil {
		return nil, err
	}

	repo.conn.Table(tableName).Delete(attribute)
	return attribute, nil
}
