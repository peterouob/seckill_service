package model

type UserLoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserRegisterReq struct {
	Username      string `json:"username"`
	Password      string `json:"password"`
	CheckPassword string `json:"check_password"`
}

type User struct {
	Id       uint   `json:id gorm:"primary_key"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (u User) TableName() string {
	return "user"
}
