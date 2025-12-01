package handlers

import (
	"log"
	"net/http"
	"org_chart/models"
	"org_chart/psql"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetRoles(c *gin.Context) {

	if roles, err := psql.FetchAllRoles(); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
			"roles": nil,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"error": nil,
			"roles": *roles,
		})
	}

}

func GetRoleById(c *gin.Context) {
	if id := c.Param("id"); id != "" {
		rid, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"role": nil, "error": err.Error()})
			return
		}
		role, err := psql.FetchRoleById(uint(rid))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"role": nil, "error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"error": nil, "role": role})

	} else {
		err := "Parameter id is required"
		c.IndentedJSON(http.StatusBadRequest, gin.H{"role": nil, "error": err})
	}
}

func PostRoles(c *gin.Context) {
	var role models.Role

	if err := c.ShouldBindJSON(&role); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"role": nil, "error": err.Error()})
	} else {
		if err = psql.CreateRole(&role); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"role": nil, "error": err.Error()})
		} else {
			c.IndentedJSON(http.StatusOK, gin.H{"role": role, "error": nil})
		}
	}
}

func PatchRole(c *gin.Context) {
	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"role": nil, "error": err.Error()})
		return
	}

	if id := c.Param("id"); id != "" {
		rid, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"role": nil, "error": "Parameter should be an integer"})
			return
		}

		if role, err := psql.UpdateRole(int(rid), input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"role": nil, "error": err.Error()})
		} else {
			c.IndentedJSON(http.StatusOK, gin.H{"role": *role, "error": nil})
		}
	} else {
		err := "Paameter id is required"
		c.IndentedJSON(http.StatusBadRequest, gin.H{"role": nil, "error": err})
	}
}

func DeleteRole(c *gin.Context) {
	if id := c.Param("id"); id != "" {
		rid, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"role": nil, "error": "Parameter is not an integer"})
		}
		err = psql.RemoveRole(uint(rid))
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"role": nil, "error": err.Error()})
		} else {
			c.IndentedJSON(http.StatusOK, gin.H{"role": map[string]int64{"id": rid}, "error": nil})
		}
	} else {
		err := "Parameter id is required"
		c.IndentedJSON(http.StatusBadRequest, gin.H{"role": nil, "error": err})
	}

}
