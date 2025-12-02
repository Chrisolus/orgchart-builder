package handlers

import (
	"net/http"
	"org_chart/initializers"
	"org_chart/models"
	"org_chart/psql"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetEmployeeSummary(c *gin.Context) {
	if summaries, err := psql.FetchEmployeeList(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"employees": nil, "error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"employees": summaries, "error": nil})
	}
}

func GetEmployees(c *gin.Context) {
	var employees []models.Employee
	query := initializers.DB.Preload("Role").Preload("Manager")

	if id := c.Query("id"); id != "" {
		query = query.Where("employees.id = ?", id)
	}
	if fname := c.Query("first_name"); fname != "" {
		query = query.Where("employees.first_name ILIKE ?", "%"+fname+"%")
	}
	if lname := c.Query("last_name"); lname != "" {
		query = query.Where("employees.last_name ILIKE ?", "%"+lname+"%")
	}
	if roleID := c.Query("role_id"); roleID != "" {
		query = query.Where("employees.role_id = ?", roleID)
	}

	if err := query.Find(&employees).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"employees": nil, "error": err.Error()})
		return
	}
	if len(employees) == 0 {
		c.JSON(http.StatusOK, gin.H{"employees": nil, "error": "No records found"})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"employees": employees, "error": nil})
}

func PostEmployee(c *gin.Context) {
	var req models.RequestBody
	if err := c.ShouldBindJSON(&req); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"employee": nil, "error": err.Error()})
	} else {
		if err = req.Validate(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"employee": nil, "error": err.Error()})
		} else {
			emp, err := psql.CreateEmployee(&req)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"employee": nil, "error": err.Error()})
			} else {
				c.JSON(http.StatusOK, gin.H{"employee": emp, "error": nil})
			}
		}
	}
}

func PatchEmployee(c *gin.Context) {
	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"employee": nil, "error": err.Error()})
		return
	}
	if id := c.Param("id"); id != "" {
		eid, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"employee": nil, "error": "Paramerte ID should be an integer"})
			return
		}

		if employee, err := psql.UpdateEmployee(int(eid), input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"employee": nil, "error": err.Error()})
		} else {
			c.IndentedJSON(http.StatusOK, gin.H{"employee": *employee, "error": nil})
		}
	} else {
		err := "Parameter id is required"
		c.IndentedJSON(http.StatusBadRequest, gin.H{"employee": nil, "error": err})
	}
}

func DeleteEmployee(c *gin.Context) {
	if id := c.Param("id"); id != "" {
		eid, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"employee": nil, "error": "Parameter is not an integer"})
		}
		err = psql.RemoveEmployee(uint(eid))
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"employee": nil, "error": err.Error()})
		} else {
			c.IndentedJSON(http.StatusOK, gin.H{"employee": map[string]int64{"id": eid}, "error": nil})
		}
	} else {
		err := "Parameter id is required"
		c.IndentedJSON(http.StatusBadRequest, gin.H{"employee": nil, "error": err})
	}
}
