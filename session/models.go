package session

import (
	"github.com/ricardocunha/golang-react-app/common"
)

//SessionModel TODO
type SessionModel struct {
	ID      uint   `gorm:"primary_key"`
	Email   string `gorm:"column:email;unique_index"`
	Token   string `gorm:"column:token;not null"`
	Expires string `gorm:"column:expires;not null"`
}

//AutoMigrate Migrate the schema of database if needed
func AutoMigrate() {
	session := common.GetSession()

	session.AutoMigrate(&SessionModel{})
}

//FindSession TODO
func FindSession(condition interface{}) (SessionModel, error) {
	session := common.GetSession()
	var model SessionModel
	err := session.Where(condition).First(&model).Error
	return model, err
}

//SaveOne TODO
func SaveOne(data interface{}) error {
	session := common.GetSession()
	err := session.Save(data).Error
	return err
}
