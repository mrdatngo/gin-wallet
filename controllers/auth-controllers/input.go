package auth_controllers

type InputActivation struct {
	Email  string `json:"email" validate:"required,email"`
	Active bool   `json:"active"`
	Token  string `json:"token" validate:"required"`
}

type InputForgot struct {
	Email string `json:"email" validate:"required,email"`
}

type InputLogin struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type InputRegister struct {
	Fullname string `json:"fullname" validate:"required,lowercase"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,gte=8"`
}

type InputResend struct {
	Email string `json:"email" validate:"required,email"`
}

type InputReset struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,gte=8"`
	Cpassword string `json:"cpassword" validate:"required,gte=8"`
	Active    bool   `json:"active"`
}

type InputUpdateEmail struct {
	Email    string `json:"email" validate:"required,email"`
	NewEmail string `json:"new_email" validate:"required,email"`
}
