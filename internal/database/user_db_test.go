package database

import (
	"apis/internal/entity"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestCreateUser(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}
	db.AutoMigrate(&entity.User{})
	user, _ := entity.NewUser("Jean Jacques", "j@example.com", "password123")
	userDB := NewUserDB(db)
	err = userDB.Create(user)
	assert.Nil(t, err)
	var userFound entity.User
	err = db.First(&userFound, "email = ?", "j@example.com").Error
	assert.Nil(t, err)
	assert.Equal(t, "Jean Jacques", userFound.Name)
	assert.Equal(t, "j@example.com", userFound.Email)
}
