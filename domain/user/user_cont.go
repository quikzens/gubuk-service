package user

import (
	"context"
	"errors"
	"time"

	db "gubuk-service/db"
	sqlc "gubuk-service/db/sqlc"
	"gubuk-service/media"
	"gubuk-service/util"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Register(c *gin.Context) {
	var req UserRegisterRequest
	err := c.Bind(&req)
	if err != nil {
		util.SendBadRequest(c, err)
		return
	}

	_, err = db.Queries.GetUserByUsername(context.TODO(), req.Username)
	if err == nil {
		util.SendBadRequest(c, errors.New("username is already used"))
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		util.SendServerError(c, err)
		return
	}

	createdUserID, err := db.Queries.CreateUser(context.TODO(), sqlc.CreateUserParams{
		ID:          uuid.New(),
		Fullname:    req.Fullname,
		Username:    req.Username,
		Email:       req.Email,
		Role:        req.Role,
		Gender:      req.Gender,
		PhoneNumber: req.PhoneNumber,
		Password:    hashedPassword,
		Address:     req.Address,
		Avatar:      "",
	})
	if err != nil {
		util.SendServerError(c, err)
		return
	}

	token, _, err := util.CreateToken(&util.UserPayload{
		ID:        uuid.NewString(),
		Username:  req.Username,
		UserID:    createdUserID.String(),
		UserRole:  req.Role,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(24 * time.Hour),
	})
	if err != nil {
		util.SendServerError(c, err)
		return
	}

	c.SetCookie("token", token, 86400, "", "", true, true)
	util.SendSuccess(c, gin.H{
		"created_id": createdUserID,
	})
}

func Login(c *gin.Context) {
	var req UserLoginRequest
	err := c.Bind(&req)
	if err != nil {
		util.SendBadRequest(c, err)
		return
	}

	user, err := db.Queries.GetUserByUsername(context.TODO(), req.Username)
	if err != nil {
		util.SendBadRequest(c, errors.New("wrong username or password"))
		return
	}

	err = util.CheckPassword(req.Password, user.Password)
	if err != nil {
		util.SendBadRequest(c, errors.New("wrong username or password"))
		return
	}

	token, _, err := util.CreateToken(&util.UserPayload{
		ID:        uuid.NewString(),
		Username:  user.Username,
		UserID:    user.ID.String(),
		UserRole:  user.Role,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(24 * time.Hour),
	})
	if err != nil {
		util.SendServerError(c, err)
		return
	}

	c.SetCookie("token", token, 86400, "", "", true, true)
	util.SendSuccess(c, nil)
}

func Logout(c *gin.Context) {
	c.SetCookie("token", "", 0, "", "", true, true)
	util.SendSuccess(c, nil)
}

// CheckAuth is validate a user session from it's token and return it's role & avatar
func CheckAuth(c *gin.Context) {
	payload, _ := c.Get("user")
	userPayload, _ := payload.(*util.UserPayload)
	userID := userPayload.UserID
	userRole := userPayload.UserRole

	id, err := uuid.Parse(userID)
	if err != nil {
		util.SendServerError(c, err)
		return
	}

	userAvatar, err := db.Queries.GetUserAvatarById(context.TODO(), id)
	if err != nil {
		util.SendServerError(c, err)
		return
	}

	util.SendSuccess(c, gin.H{
		"user_id":     userID,
		"user_role":   userRole,
		"user_avatar": userAvatar,
	})
}

// GetUserDetail return profile data of the currently logged in user
func GetUserDetail(c *gin.Context) {
	payload, _ := c.Get("user")
	userPayload, _ := payload.(*util.UserPayload)
	userID := userPayload.UserID

	id, err := uuid.Parse(userID)
	if err != nil {
		util.SendServerError(c, err)
		return
	}

	user, err := db.Queries.GetUserById(context.TODO(), id)
	if err != nil {
		util.SendServerError(c, err)
		return
	}

	util.SendSuccess(c, user)
}

func UpdateUserAvatar(c *gin.Context) {
	payload, _ := c.Get("user")
	userPayload, _ := payload.(*util.UserPayload)
	userID := userPayload.UserID

	id, err := uuid.Parse(userID)
	if err != nil {
		util.SendServerError(c, err)
		return
	}

	avatarImage, err := c.FormFile("avatar")
	if err != nil {
		util.SendBadRequest(c, err)
		return
	}

	err = media.ValidateImage(avatarImage)
	if err != nil {
		util.SendBadRequest(c, err)
		return
	}

	oldAvatarUrl, err := db.Queries.GetUserAvatarById(context.TODO(), id)
	if err != nil {
		util.SendBadRequest(c, err)
		return
	}

	var newAvatarUrl string
	if oldAvatarUrl != "" {
		newAvatarUrl, err = media.UpdateMedia("avatar", oldAvatarUrl, avatarImage)
		if err != nil {
			util.SendServerError(c, err)
			return
		}
	} else {
		newAvatarUrl, err = media.UploadMedia("avatar", avatarImage)
		if err != nil {
			util.SendServerError(c, err)
			return
		}
	}

	err = db.Queries.UpdateUserAvatarById(context.TODO(), sqlc.UpdateUserAvatarByIdParams{
		ID:        id,
		Avatar:    newAvatarUrl,
		UpdatedAt: time.Now(),
	})
	if err != nil {
		util.SendServerError(c, err)
	}

	util.SendSuccess(c, gin.H{
		"new_image": newAvatarUrl,
	})
}

func UpdateUserPassword(c *gin.Context) {
	payload, _ := c.Get("user")
	userPayload, _ := payload.(*util.UserPayload)
	username := userPayload.Username
	userID := userPayload.UserID

	id, err := uuid.Parse(userID)
	if err != nil {
		util.SendServerError(c, err)
		return
	}

	var req UserUpdatePasswordRequest
	err = c.Bind(&req)
	if err != nil {
		util.SendBadRequest(c, err)
		return
	}

	user, err := db.Queries.GetUserByUsername(context.TODO(), username)
	if err != nil {
		util.SendServerError(c, err)
		return
	}

	err = util.CheckPassword(req.Password, user.Password)
	if err != nil {
		util.SendBadRequest(c, errors.New("wrong password"))
		return
	}

	hashedPassword, err := util.HashPassword(req.NewPassword)
	if err != nil {
		util.SendServerError(c, err)
		return
	}

	err = db.Queries.UpdateUserPasswordById(context.TODO(), sqlc.UpdateUserPasswordByIdParams{
		ID:        id,
		Password:  hashedPassword,
		UpdatedAt: time.Now(),
	})
	if err != nil {
		util.SendServerError(c, err)
		return
	}

	c.SetCookie("token", "", 0, "", "", true, true)
	util.SendSuccess(c, nil)
}

func UpdateUserProfile(c *gin.Context) {
	payload, _ := c.Get("user")
	userPayload, _ := payload.(*util.UserPayload)
	userID := userPayload.UserID

	id, err := uuid.Parse(userID)
	if err != nil {
		util.SendServerError(c, err)
		return
	}

	var req UserUpdateProfileRequest
	err = c.Bind(&req)
	if err != nil {
		util.SendBadRequest(c, err)
		return
	}

	err = db.Queries.UpdateUserById(context.TODO(), sqlc.UpdateUserByIdParams{
		ID:          id,
		Fullname:    req.Fullname,
		Email:       req.Email,
		Gender:      req.Gender,
		PhoneNumber: req.PhoneNumber,
		Address:     req.Address,
		UpdatedAt:   time.Now(),
	})
	if err != nil {
		util.SendServerError(c, err)
		return
	}

	util.SendSuccess(c, nil)
}
