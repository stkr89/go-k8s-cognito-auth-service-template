package types

type SignInRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type SignInResponse struct {
	AccessToken string `json:"accessToken"`
}

type SignUpRequest struct {
	FirstName string `json:"firstName" validate:"required" conform:"name"`
	LastName  string `json:"lastName" validate:"required" conform:"name"`
	Email     string `json:"email" validate:"required" conform:"email"`
	Password  string `json:"password" validate:"required,min=8"`
}

type SignUpResponse struct {
	ID        string `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}
