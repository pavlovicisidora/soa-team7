package model

type User struct {
	ID       int    `json:"id" gorm:"not null;type:int"`
	Username string `json:"username" gorm:"not null;type:string"`
	Password string `json:"password" gorm:"not null;type:string"`
	Mail     string `json:"mail"`
	Role     string `json:"role" gorm:"not null;type:string"`
	Blocked  bool   `json:"blocked"`
}
