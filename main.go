package main

import (
	"gopkg.in/gin-gonic/gin.v1"

	"github.com/jinzhu/gorm"
	"github.com/ricardocunha/golang-react-app/auth"
	"github.com/ricardocunha/golang-react-app/common"
	"github.com/ricardocunha/golang-react-app/users"
)

func Migrate(db *gorm.DB) {
	users.AutoMigrate()
}

func main() {

	db := common.Init()
	Migrate(db)
	defer db.Close()

	r := gin.Default()

	v1 := r.Group("/api")
	users.RegisterNoAuth(v1.Group("/users"))
	//v1.Use(users.AuthMiddleware(false))

	v1.Use(auth.AuthMiddleware(true))
	users.Register(v1.Group("/user"))
	users.ProfileRegister(v1.Group("/profiles"))

	testAuth := r.Group("/api/ping")

	testAuth.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// test 1 to 1
	/*
		tx1 := db.Begin()
		_, err := users.FindOneUser(&users.UserModel{Username: "ricardocunha"})
		if err != nil {
			bytePassword := []byte("Ab123456")
			// Make sure the second param `bcrypt generator cost` between [4, 32)
			passwordHash, _ := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)
			userA := users.UserModel{
				Username:     "ricardocunha",
				Email:        "ricardocunha@gmail.com",
				PasswordHash: string(passwordHash),
			}
			tx1.Save(&userA)
			tx1.Commit()
			fmt.Println(userA)
		} else {
			fmt.Println("User exists ")
		}
	*/
	r.Run() // listen and serve on 0.0.0.0:8080
}
