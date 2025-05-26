package model

type SignUpRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type SignInRequest struct {
	IDToken string `json:"id_token"`
}

type ChangePasswordRequest struct {
	Email       string `json:"email"`
	NewPassword string `json:"new_password"`
}
