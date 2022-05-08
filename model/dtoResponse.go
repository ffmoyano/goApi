package model

type UserResponse struct {
	ID       int            `json:"id"`
	Name     string         `json:"name"`
	Username string         `json:"username"`
	Email    string         `json:"email"`
	Roles    []RoleResponse `json:"roles"`
	Active   bool           `json:"active"`
}

type RoleResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type TokenResponse struct {
	Jwt          string `json:"jwt"`
	RefreshToken string `json:"refreshToken"`
}
