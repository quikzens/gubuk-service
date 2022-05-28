package user

type UserRegisterRequest struct {
	Fullname    string `form:"fullname" binding:"required"`
	Username    string `form:"username" binding:"required,min=3"`
	Email       string `form:"email" binding:"required,email"`
	Role        string `form:"role" binding:"required,oneof=tenant owner"`
	Gender      string `form:"gender" binding:"required,oneof=male female"`
	PhoneNumber string `form:"phone_number" binding:"required"`
	Password    string `form:"password" binding:"required,min=8"`
	Address     string `form:"address" binding:"required"`
}

type UserLoginRequest struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

type UserUpdatePasswordRequest struct {
	Password    string `form:"password" binding:"required"`
	NewPassword string `form:"new_password" binding:"required,min=8"`
}

type UserUpdateProfileRequest struct {
	Fullname    string `form:"fullname" binding:"required"`
	Email       string `form:"email" binding:"required,email"`
	Gender      string `form:"gender" binding:"required,oneof=male female"`
	PhoneNumber string `form:"phone_number" binding:"required"`
	Address     string `form:"address" binding:"required"`
}
