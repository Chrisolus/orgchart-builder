package psql

import (
	"errors"
	"log"
	"org_chart/initializers"
	"org_chart/models"
)

func CreateRole(role *models.Role) error {
	res := initializers.DB.Create(role)
	if res.Error != nil {
		log.Println("CreateRole:", res.Error)
		return errors.New("error while creating role")
	}
	log.Println("Role Created: ", role.Name)
	return nil
}

func FetchAllRoles() (*[]models.Role, error) {
	var (
		roles []models.Role
		err   error
	)
	if err = initializers.DB.Find(&roles).Error; err != nil {
		log.Println("FetchAllRoles:", err.Error())
		return nil, errors.New("cannot fetch roles")
	}
	log.Println("Fetch Successfull: ", len(roles))
	return &roles, nil
}

func FetchRoleById(id uint) (*models.Role, error) {
	var (
		role models.Role
		err  error
	)
	if err = initializers.DB.Find(&role, id).Error; err != nil {
		log.Println("FetchRoleById:", err.Error())
		return nil, errors.New("error while finding role")
	}
	if role.ID == 0 {
		log.Println("FetchRoleByID: role doesn't exist")
		return nil, errors.New("role doesn't exist")
	}
	log.Println("Fetch Successfull: ", role.Name)
	return &role, nil
}

func UpdateRole(id int, input map[string]interface{}) (*models.Role, error) {
	var role models.Role
	if res := initializers.DB.First(&role, id); res.Error != nil {
		log.Println("UpdateRole:", res.Error.Error())
		return nil, errors.New("error while finding the role to be updated")
	}

	if err := initializers.DB.Model(&role).Updates(input).Error; err != nil {
		log.Println("UpdateRole:", err.Error())
		return nil, errors.New("errir while updating the role")
	}
	log.Println("Role Updated: ", role.ID)
	return &role, nil
}

func RemoveRole(rid uint) error {
	res := initializers.DB.Delete(&models.Role{}, rid)
	if res.RowsAffected < 1 {
		log.Println("RemoveRole: No rows affected")
		return errors.New("the role doesn't seem to exist")
	}
	if res.Error != nil {
		log.Println("RemoveRole:", res.Error.Error())
		return errors.New("error while deleting the role")
	}
	log.Println("Deletetion Successfull | RoleID:", rid)
	return nil
}
