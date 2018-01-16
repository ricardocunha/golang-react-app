package users

import (
	"github.com/ricardocunha/golang-react-app/common"
	"gopkg.in/gin-gonic/gin.v1"
)

//UpdateContextUserModel A helper to write user_id and user_model to the context
func UpdateContextUserModel(c *gin.Context, myUserID uint) {
	var myUserModel UserModel
	if myUserID != 0 {
		db := common.GetDB()
		db.First(&myUserModel, myUserID)
	}
	c.Set("my_user_id", myUserID)
	c.Set("my_user_model", myUserModel)
}
