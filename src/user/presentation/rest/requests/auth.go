package requests

type Register struct {
	Username string `json:"username" validate:"required,alphanum"`
	Password string `json:"password" validate:"required,min=3"`
	Role     string `json:"role" validate:"required,oneof=seller buyer"`
}

type Login struct {
	Username string `json:"username" validate:"required,alphanum"`
	Password string `json:"password" validate:"required,min=3"`
}
