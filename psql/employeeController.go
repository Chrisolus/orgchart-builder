package psql

import (
	"errors"
	"log"
	"org_chart/initializers"
	"org_chart/models"
)

func CreateEmployee(bodyReq *models.RequestBody) (*models.Employee, error) {
	var emp models.Employee
	if bodyReq != nil {
		emp = models.Employee{
			FirstName: bodyReq.FirstName,
			LastName:  bodyReq.LastName,
			RoleID:    bodyReq.RoleID,
			ManagerID: bodyReq.ManagerID,
		}
		res := initializers.DB.Create(&emp)
		if res.Error != nil {
			log.Println("CreateEmployee:", res.Error.Error())
			return nil, errors.New("error while creating employee")
		}
		log.Println("Employee Created: ", emp.ID)
		return &emp, nil
	} else {
		log.Println("CreateEmployee: Invalid Request")
		return nil, errors.New("required data not provided")
	}
}

func FetchAllEmployees() (*[]models.Employee, error) {
	var (
		employees []models.Employee
		err       error
	)
	if err = initializers.DB.Preload("Role").Preload("Manager").Find(&employees).Error; err != nil {
		log.Println("FetchAllEmployees:", err.Error())
		return nil, errors.New("cannot find employees")
	}
	log.Println("Fetch Successfull: ", len(employees))
	return &employees, nil
}

func FetchEmployeeList() (*[]models.EmployeeSummary, error) {
	var summaries []models.EmployeeSummary
	res := initializers.DB.Table("employees").
		Select("employees.id, CONCAT(employees.first_name, ' ', employees.last_name, ' [', roles.name, ']') AS name").
		Joins("JOIN roles ON employees.role_id = roles.id").
		Where("employees.deleted_at IS NULL").
		Scan(&summaries)
	if res.Error != nil {
		log.Println("FetchEmployeeList:", res.Error.Error())
		return nil, errors.New("cannot find employees/roles")
	}
	log.Println("List Fetched Successfully: ", len(summaries))
	return &summaries, nil
}

func UpdateEmployee(id int, input map[string]interface{}) (*models.Employee, error) {
	var employee models.Employee
	if res := initializers.DB.First(&employee, id); res.Error != nil {
		log.Println("UpdateEmployee:", res.Error.Error())
		return nil, errors.New("error while finding employee to be updated")
	}

	if err := initializers.DB.Model(&employee).Updates(input).Error; err != nil {
		log.Println("UpdateEmployee:", err.Error())
		return nil, errors.New("error while updating the employee")
	}
	log.Println("Employee Updated: ", employee.ID)
	return &employee, nil
}

func RemoveEmployee(id uint) error {
	res := initializers.DB.Delete(&models.Employee{}, id)
	if res.RowsAffected < 1 {
		log.Println("RemoveEmployee: No Rows Affected")
		return errors.New("employee doesn't seem to exist")
	}
	if res.Error != nil {
		log.Println("RemoveEmployee:", res.Error.Error())
		return errors.New("error while deleting the employee")
	}
	log.Println("Deletion Successfull | Employee ID:", id)
	return nil
}
