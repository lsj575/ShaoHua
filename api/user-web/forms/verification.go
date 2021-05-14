package forms

type EmailVerificationForm struct {
	Email    string `json:"email" binding:"required,email"`
	Username string `json:"username" binding:"loginUsername"`
	//Type     string `json:"type" binding:"required,oneof=register login"`
}
