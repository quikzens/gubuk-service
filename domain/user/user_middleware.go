package user

import (
	"errors"

	"gubuk-service/util"

	"github.com/gin-gonic/gin"
)

func VerifyAuth(c *gin.Context) {
	token, err := c.Cookie("token")
	if err != nil {
		util.SendUnauthorized(c, err)
		return
	}

	payload, err := util.VerifyToken(token)
	if err != nil {
		if errors.Is(err, util.ErrExpiredToken) {
			c.SetCookie("token", "", 0, "", "", true, true)
		}

		util.SendUnauthorized(c, err)
		return
	}

	c.Set("user", payload)
	c.Next()
}

func VerifyRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		payload, _ := c.Get("user")
		userPayload, _ := payload.(*util.UserPayload)
		userRole := userPayload.UserRole

		if userRole != role {
			util.SendUnauthorized(c, errors.New("your role could not access this api"))
			return
		}

		c.Next()
	}
}
