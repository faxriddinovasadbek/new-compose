package models

type User struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	LastName     string `json:"last_name"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	UserName     string `json:"user_name"`
}

type UserByAccess struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	LastName     string `json:"last_name"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	UserName     string `json:"user_name"`
	AccessToken  string `json:"access_token"`
}

type RegisterUser struct {
	Name         string `json:"name"`
	LastName     string `json:"last_name"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	UserName     string `json:"user_name"`
}
