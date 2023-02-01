package usersdto

type UserResponse struct {
	ID    int    `json:"id"`
	Name  string `json:"name" form:"name"`
	Email string `json:"email" form:"email"`
	Image string `json:"image" form:"image"`
}
