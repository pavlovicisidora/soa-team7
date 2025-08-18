package model

type Profile struct {
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	ProfilePic string `json:"picture"`
	Bio        string `json:"bio"`
	Motto      string `json:"motto"`
}
