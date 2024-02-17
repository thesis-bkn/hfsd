package entity

import "github.com/golang-jwt/jwt/v5"

type User struct {
	ID        string `json:"id"        db:"id"`
	Password  string `json:"password"  db:"password"`
	Email     string `json:"email"     db:"email"`
	Activated bool   `json:"activated" db:"activated"`
	CreatedAt int64  `json:"createdAt" db:"created_at"`
}

type ProfileClaim struct {
	*jwt.RegisteredClaims
	UserID    string `json:"userID"`
	Activated bool   `json:"act"`
}
