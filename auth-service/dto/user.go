package dto

type CreateAuthRequest struct {
	FullName string `json:"full_name" binding:"required"`
	Login    string `json:"login" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
}
