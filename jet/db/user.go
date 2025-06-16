package db

import (
	"errors"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        string     `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	Email     string     `gorm:"unique" json:"email"`
	Password  string     `json:"-"`
	IsAdmin   bool       `gorm:"default:false" json:"isAdmin"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
}

func (u *User) CreateAdmin() error {
	db := GetDB()
	if db == nil {
		log.Println("DBConn is not initialized")
		return errors.New("DBConn is not initialized")
	}

	log.Println("DBConn is initialized")

	user := User{
		Email:    "your email",
		Password: "your password",
		IsAdmin:  true,
	}

	log.Println("Creating user with email:", user.Email)

	password, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		log.Println("Error creating password hash:", err)
		return errors.New("error creating password")
	}
	user.Password = string(password)

	if err := db.Create(&user).Error; err != nil {
		log.Println("Error creating user:", err)
		return errors.New("error creating user")
	}

	log.Println("User created successfully")
	return nil
}

func (u *User) LoginAsAdmin(email string, password string) (*User, error) {
	log.Println("Initializing login process")
	if DBConn == nil {
		log.Println("DBConn is nil")
		return nil, errors.New("DBConn is not initialized")
	}
	log.Println("DBConn is initialized")
	if err := DBConn.Where("email = ? AND is_admin = ? ", email, true).First(&u).Error; err != nil {
		log.Println("User not found:", err)
		return nil, errors.New("user not found")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		log.Println("Invalid password:", err)
		return nil, errors.New("invalid password")
	}
	log.Println("User logged in successfully")
	return u, nil
}
