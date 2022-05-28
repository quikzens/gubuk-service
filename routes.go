package main

import (
	"gubuk-service/domain/house"
	"gubuk-service/domain/transaction"
	"gubuk-service/domain/user"

	"github.com/gin-gonic/gin"
)

func SetRoutes(router *gin.Engine) {
	apiGroup := router.Group("/api")

	// User Authentication
	apiGroup.POST("/register", user.Register)
	apiGroup.POST("/login", user.Login)
	apiGroup.GET("/logout", user.VerifyAuth, user.Logout)
	apiGroup.GET("/auth", user.VerifyAuth, user.CheckAuth)

	// User
	apiGroup.PATCH("/user/avatar", user.VerifyAuth, user.UpdateUserAvatar)
	apiGroup.PATCH("/user/password", user.VerifyAuth, user.UpdateUserPassword)
	apiGroup.PATCH("/user", user.VerifyAuth, user.UpdateUserProfile)
	apiGroup.GET("/user", user.VerifyAuth, user.GetUserDetail)

	// House
	apiGroup.POST("/houses", user.VerifyAuth, user.VerifyRole("owner"), house.CreateHouse)
	apiGroup.PATCH("/houses/:id", user.VerifyAuth, user.VerifyRole("owner"), house.UpdateHouse)
	apiGroup.DELETE("/houses/:id", user.VerifyAuth, user.VerifyRole("owner"), house.DeleteHouse)
	apiGroup.GET("/houses", house.GetHouseList)
	apiGroup.GET("/houses/me", user.VerifyAuth, user.VerifyRole("owner"), house.GetMyHouseList)
	apiGroup.GET("/houses/:id", house.GetHouseDetail)
	apiGroup.GET("/houses/count", house.GetHouseCount)

	// Transaction
	apiGroup.POST("/transactions", user.VerifyAuth, user.VerifyRole("tenant"), transaction.CreateTransaction)
	apiGroup.GET("/transactions", user.VerifyAuth, transaction.ListTransaction)
	apiGroup.PATCH("/transactions/pay/:id", user.VerifyAuth, user.VerifyRole("tenant"), transaction.PayTransaction)
	apiGroup.PATCH("/transactions/status/:id", user.VerifyAuth, user.VerifyRole("owner"), transaction.UpdateTransactionStatus)
}
