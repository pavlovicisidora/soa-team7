package model

type User struct {
	//Id       string `json:"id" gorm:"not null;type:string"`     $Nisam siguran jel treba id, s obzirom da mi Mongo stalno svoj izbacuje
	Username string  `json:"username" gorm:"not null;type:string"`
	Password string  `json:"password" gorm:"not null;type:string"`
	Mail     string  `json:"mail"`
	Role     string  `json:"role" gorm:"not null;type:string"`
	Blocked  bool    `json:"blocked"`
	Profile  Profile `json:"profile"`
}
