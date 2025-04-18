package models

type LoginInfo struct {
	App      string `json:"appid"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterInfo struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AcceptInfo struct {
	UserCode string `json:"usercode"`
}
type LoginResponse struct {
	AccessToken string `json:"access_token"`
}

type RegisterResponse struct {
	Id string `json:"id"`
}

type AcceptResponse struct {
	Status bool `json:"status"`
}
type ResponseError struct {
	Error string `json:"error"`
}
