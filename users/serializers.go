package users

import (
	"gopkg.in/gin-gonic/gin.v1"
)

type ProfileSerializer struct {
	C *gin.Context
	UserModel
}

// Declare your response schema here
type ProfileResponse struct {
	ID        uint    `json:"-"`
	Username  string  `json:"username"`
}

// Put your response logic including wrap the userModel here.
func (self *ProfileSerializer) Response() ProfileResponse {
	//myUserModel := self.C.MustGet("my_user_model").(UserModel)
	profile := ProfileResponse{
		ID:        self.ID,
		Username:  self.Username,
	}
	return profile
}

type UserSerializer struct {
	c *gin.Context
}

type UserResponse struct {
	Username string  `json:"username"`
	Email    string  `json:"email"`
}

func (self *UserSerializer) Response() UserResponse {
	myUserModel := self.c.MustGet("my_user_model").(UserModel)
	user := UserResponse{
		Username: myUserModel.Username,
		Email:    myUserModel.Email,
	}
	return user
}
