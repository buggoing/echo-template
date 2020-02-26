package handler

type ReqUserLogin struct {
	Username string `minLength:"5" maxLength:"32" json:"username"`
	Password string `minLength:"5" maxLength:"32" json:"password"`
}

type RspUserLogin struct {
	Token string `json:"token"`
}
