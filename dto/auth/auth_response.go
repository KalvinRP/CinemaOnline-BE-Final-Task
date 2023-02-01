package authdto

type AuthResponse struct {
	Name  string `gorm:"type: varchar(255)" json:"name"`
	Token string `gorm:"type: varchar(255)" json:"token"`
	Role  string `gorm:"type: varchar(255)" json:"role"`
}
