package psql

import (
	"errors"
	"fmt"
	"log"
	"org_chart/initializers"
	"org_chart/models"
)

func RegisterController(user *models.User) error {
	errortxt := ""
	if err := user.Validate(); err != nil {
		errortxt = "email/password invalid"
		log.Println(errortxt, err.Error())
		return errors.New(errortxt)
	}
	if err := user.HashPassword(); err != nil {
		errortxt = "password hashing failed"
		log.Println(errortxt, err.Error())
		return errors.New(errortxt)
	}
	result := initializers.DB.Create(user)
	if result.Error != nil {
		errortxt = "user already exist. try logging in"
		log.Println(result.Error)
		return errors.New(errortxt)
	}
	return nil
}

// Validating credentials while user login
func AuthController(reqUser *models.User) (*models.User, error) {
	if err := reqUser.Validate(); err != nil {
		errortxt := fmt.Sprint("Invaid User: ", reqUser)
		log.Println(errortxt, err.Error())
		return nil, errors.New(errortxt)
	}
	var user models.User
	err := initializers.DB.
		Model(&models.User{}).
		Where("email = ?", reqUser.Email).
		Scan(&user).Error
	if err != nil {
		errortxt := "selection error: " + err.Error()
		log.Println(errortxt)
		return nil, errors.New(errortxt)
	}
	if user.ID == 0 {
		log.Println("email doesn't exist")
		return nil, errors.New("email doesn't exist")
	}
	err = user.ValidateHash(reqUser.Password)
	if err != nil {
		errortxt := "password invalid"
		log.Print(errortxt)
		return nil, errors.New(errortxt)
	}
	return &user, nil
}

func IsValidUserId(id uint) bool {
	var user models.User
	if err := initializers.DB.Model(&models.User{}).Where("id = ?", id).Scan(&user).Error; err != nil {
		log.Println("db level validation: ", err.Error())
		return false
	}
	if user.ID == 0 {
		return false
	}
	return true
}

func GetUsersByIDs(userIDs []uint) ([]models.UserSummary, error) {
	if len(userIDs) == 0 {
		return []models.UserSummary{}, nil
	}

	var summaries []models.UserSummary
	if err := initializers.DB.
		Model(&models.User{}).
		Select("id", "username").
		Where("id IN ?", userIDs).
		Find(&summaries).Error; err != nil {
		return nil, err
	}

	return summaries, nil
}

func GetBasicUserInfo() ([]models.UserSummary, error) {
	var summaries []models.UserSummary
	if err := initializers.DB.
		Model(&models.User{}).
		Select("id", "username").
		Find(&summaries).Error; err != nil {
		return nil, err
	}
	return summaries, nil

}
