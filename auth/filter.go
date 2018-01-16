package auth

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/ricardocunha/golang-react-app/users"
	"gopkg.in/gin-gonic/gin.v1"
	"gopkg.in/jose.v1/crypto"
	"gopkg.in/jose.v1/jws"
)

//AuthMiddleware You can custom middlewares yourself as the doc: https://github.com/gin-gonic/gin#custom-middleware
func AuthMiddleware(auto401 bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		users.UpdateContextUserModel(c, 0)
		bytes, _ := ioutil.ReadFile("./sample_key.pub")
		rsaPublic, _ := crypto.ParseRSAPublicKeyFromPEM(bytes)

		auth := c.Request.Header.Get("Authorization")
		if len(auth) == 0 {
			c.AbortWithError(http.StatusUnauthorized, errors.New("Authorization not found"))
		}
		accessToken := auth[7:len(auth)]
		jwt, err := jws.ParseJWT([]byte(accessToken))
		if err != nil {
			if auto401 {
				c.AbortWithError(http.StatusUnauthorized, err)
			}
		}
		// Validate token
		if err = jwt.Validate(rsaPublic, crypto.SigningMethodRS256); err != nil {
			c.AbortWithError(http.StatusUnauthorized, err)
		}
		//j, erro := jws.ParseFromHeader(c.Request, jws.General)
		j, erro := jws.ParseFromRequest(c.Request, jws.General)
		if erro != nil {
			if auto401 {
				c.AbortWithError(http.StatusUnauthorized, erro)
			}
			return
		}
		fmt.Println(j)
		//fmt.Println(auth)
		//fmt.Println(token)
		/*
				token, err := request.ParseFromRequest(c.Request, request.OAuth2Extractor, func(token *jwt.Token) (interface{}, error) {
					b := ([]byte(common.NBSecretPassword))
					return b, nil
				})
			if err != nil {
				if auto401 {
					c.AbortWithError(http.StatusUnauthorized, err)
				}
				return
			}
		*/
		//_, err := session.FindSession(&session.SessionModel{Token: token.Raw})
		/*
			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				myUserID := uint(claims["id"].(float64))
				users.UpdateContextUserModel(c, myUserID)
			}*/
	}
}
