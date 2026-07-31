package dto

type CreateUserGrpcRequest struct {
	ID       int
	FullName string
	Email    string
	Phone    string
}
