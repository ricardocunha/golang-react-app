package users

import (
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/ricardocunha/golang-react-app/common"
	s "github.com/ricardocunha/golang-react-app/session"
	"gopkg.in/gin-gonic/gin.v1"
	"gopkg.in/jose.v1/crypto"
	"gopkg.in/jose.v1/jws"
)

//Register adding routers for /users
func Register(router *gin.RouterGroup) {
	router.GET("/", Retrieve)
	router.PUT("/", Update)
}

//Register adding routers for /users
func RegisterNoAuth(router *gin.RouterGroup) {
	router.POST("/signup", Signup)
	router.POST("/signin", Signin)
}

//ProfileRegister adding routers for /profile
func ProfileRegister(router *gin.RouterGroup) {
	router.GET("/:username", ProfileRetrieve)
}

//Signup User's registration
func Signup(c *gin.Context) {
	userModelValidator := NewUserModelValidator()
	if err := userModelValidator.Bind(c); err != nil {
		c.JSON(http.StatusUnprocessableEntity, common.NewValidatorError(err))
		return
	}

	if err := SaveOne(&userModelValidator.userModel); err != nil {
		c.JSON(http.StatusUnprocessableEntity, common.NewError("database", err))
		return
	}
	c.Set("my_user_model", userModelValidator.userModel)
	serializer := UserSerializer{c}
	c.JSON(http.StatusCreated, gin.H{"user": serializer.Response()})
}

//Signin is a function to handle with users login
func Signin(c *gin.Context) {
	session := common.InitSession()
	loginValidator := NewLoginValidator()
	if err := loginValidator.Bind(c); err != nil {
		c.JSON(http.StatusUnprocessableEntity, common.NewValidatorError(err))
		return
	}
	userModel, err := FindOneUser(&UserModel{Email: loginValidator.userModel.Email})

	if err != nil {
		c.JSON(http.StatusForbidden, common.NewError("login", errors.New("Not Registered email or invalid password")))
		return
	}

	if userModel.checkPassword(loginValidator.User.Password) != nil {
		c.JSON(http.StatusForbidden, common.NewError("login", errors.New("Not Registered email or invalid password")))
		return
	}
	UpdateContextUserModel(c, userModel.ID)
	serializer := UserSerializer{c}

	expire := time.Now().Add(common.ExpireTime) /*
			// Create the token
			token := jwt.New(jwt.SigningMethodHS256)
			// Set some claims
			claims := make(jwt.MapClaims)
			claims["id"] = userModel.ID
			claims["exp"] = expire.Unix()
			token.Claims = claims
			tokenString, err := token.SignedString([]byte(common.NBSecretPassword))
			if err != nil {
				c.JSON(http.StatusForbidden, common.NewError("login", errors.New("Create JWT Token failed")))
				return
			}
			reg, _ := regexp.Compile("[^a-zA-Z0-9]+")
			tokenString = reg.ReplaceAllString(tokenString, "")
			c.JSON(http.StatusOK, gin.H{
				"token":  tokenString,
				"expire": expire.Format(time.RFC3339),
				"user":   serializer.Response(),
			})
		claims := &jwt.StandardClaims{
			ExpiresAt: 60000000,
			Issuer:    "dev",
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		ss, _ := token.SignedString([]byte(common.NBSecretPassword))
		//fmt.Printf("%v %v", ss, err)
	*/
	bytes, _ := ioutil.ReadFile("./sample_key.priv")
	claims := jws.Claims{}
	claims.SetExpiration(time.Now().Add(time.Duration(10) * time.Minute))
	rsaPrivate, _ := crypto.ParseRSAPrivateKeyFromPEM(bytes)
	jwt := jws.NewJWT(claims, crypto.SigningMethodRS256)
	ss, _ := jwt.Serialize(rsaPrivate)

	c.JSON(http.StatusOK, gin.H{
		"token":  string(ss),
		"expire": expire.Format(time.RFC3339),
		"user":   serializer.Response(),
	})
	sessionData := s.SessionModel{

		Email:   loginValidator.userModel.Email,
		Token:   string(ss),
		Expires: expire.Format(time.RFC3339),
	}
	tx1 := session.Begin()
	tx1.Save(&sessionData)
	tx1.Commit()
}

//ProfileRetrieve retrieve the user's profile
func ProfileRetrieve(c *gin.Context) {
	username := c.Param("username")
	userModel, err := FindOneUser(&UserModel{Username: username})
	if err != nil {
		c.JSON(http.StatusNotFound, common.NewError("profile", errors.New("Invalid username")))
		return
	}
	profileSerializer := ProfileSerializer{c, userModel}
	c.JSON(http.StatusOK, gin.H{"profile": profileSerializer.Response()})
}

//Retrieve retrieve the user
func Retrieve(c *gin.Context) {
	serializer := UserSerializer{c}
	c.JSON(http.StatusOK, gin.H{"user": serializer.Response()})
}

//Update is a function to update the model
func Update(c *gin.Context) {
	myUserModel := c.MustGet("my_user_model").(UserModel)
	userModelValidator := NewUserModelValidatorFillWith(myUserModel)
	if err := userModelValidator.Bind(c); err != nil {
		c.JSON(http.StatusUnprocessableEntity, common.NewValidatorError(err))
		return
	}

	userModelValidator.userModel.ID = myUserModel.ID
	if err := myUserModel.Update(userModelValidator.userModel); err != nil {
		c.JSON(http.StatusUnprocessableEntity, common.NewError("database", err))
		return
	}
	UpdateContextUserModel(c, myUserModel.ID)
	serializer := UserSerializer{c}
	c.JSON(http.StatusOK, gin.H{"user": serializer.Response()})
}
