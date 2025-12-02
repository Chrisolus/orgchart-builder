package handlers

import (
	"log"
	"net/http"
	"org_chart/scylla"

	"github.com/gin-gonic/gin"
)

func GetConversation(c *gin.Context) {
	key := c.Query("key")
	if key == "" {
		c.JSON(http.StatusBadRequest, gin.H{"messages": nil, "error": "conversation key missing in query parameter"})
		return
	}
	messages, err := scylla.FetchMessages(key)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"messages": nil, "error": err.Error()})
		return
	}
	log.Println(messages)
	c.JSON(http.StatusOK, gin.H{"messages": messages, "error": nil})
}
