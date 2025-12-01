package handlers

import (
	"log"
	"net/http"
	mdl "org_chart/middleware"
	"org_chart/models"
	"org_chart/psql"

	"github.com/gin-gonic/gin"
)

func RegisterUser(c *gin.Context) {

	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}
	err := psql.RegisterController(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}
	tokens, err := mdl.GenAuthAndRefreshToken(&mdl.Claims{UserID: user.ID, Email: user.Email})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"token": tokens,
		"error": nil,
	})
}

func LoginUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Email/Password invalid"})
		return
	}
	ResUser, err := psql.AuthController(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}
	tokens, err := mdl.GenAuthAndRefreshToken(&mdl.Claims{UserID: ResUser.ID, Email: ResUser.Email})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"token": tokens,
		"error": nil,
	})
}

func FetchAllUsers(c *gin.Context) {
	if summary, err := psql.GetBasicUserInfo(); err != nil {
		c.JSON(http.StatusOK, gin.H{"users": nil, "error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"users": summary, "error": nil})
	}
}

func RefreshToken(c *gin.Context) {
	var jwt models.JWT
	if err := c.ShouldBindJSON(&jwt); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"isvalid": false, "error": "Token Check Failed"})
		return
	}
	oldClaims, err := mdl.ValidateToken(jwt.Token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"isvalid": false, "error": err.Error()})
		return
	}
	if !oldClaims.IsRefresh {
		c.JSON(http.StatusBadRequest, gin.H{"isvalid": false, "error": "The token is not valid refresh token"})
		return
	}
	tokens, err := mdl.GenAuthAndRefreshToken(oldClaims)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"isvalid": false,
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"isvalid": true,
		"error":   nil,
		"token":   tokens,
	})
}

func VerifyToken(c *gin.Context) {
	var jwt models.JWT
	if err := c.ShouldBindJSON(&jwt); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"isvalid": false, "error": "Token Check Failed"})
		return
	}
	_, err := mdl.ValidateToken(jwt.Token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"isvalid": false, "status": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"isvalid": true, "token": jwt.Token})
}
