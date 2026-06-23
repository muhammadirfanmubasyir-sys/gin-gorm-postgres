package models

import "golang.org/x/crypto/bcrypt"

type User struct {
	Id       int    `json:"id" gorm:"primary_key"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"-"`
}

func (User) TableName() string {
	return "table_of_user"
}

func (u *User) HashPassword(plainPassword string) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashed)
	return nil
}

func (u *User) VerifyPassword(plainPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plainPassword))
}
