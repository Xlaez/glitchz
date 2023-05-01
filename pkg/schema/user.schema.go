package schema

type AddUserReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,alphanum,min=6"`
	Username string `json:"username" binding:"required,alphanum,min=2"`
}

type VerifyEmailReq struct {
	Code string `json:"code" binding:"required,min=6,max=6"`
}
