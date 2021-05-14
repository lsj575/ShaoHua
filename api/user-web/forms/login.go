package forms

type PasswordLoginForm struct {
	Username string `json:"username" binding:"loginUsername"`
	Email    string `json:"email" binding:"email"`
	Password string `json:"password" binding:"required,min=7,max=30"`
}

type NonPasswordLoginForm struct {
	Email string `json:"email" binding:"required,email"`
	Code  string `json:"code" binding:"required,min=6,max=6"`
}

type RegisterForm struct {
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password"` // 暂时未采用注册时输入密码的方案，若采用可以用如下验证规则 binding:"required,min=7,max=30"
	RePassword string `json:"rePassword"`
	Code  string `json:"code" binding:"required,min=6,max=6"`
}
