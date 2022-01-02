package userapp

import (
	"net/http"

	"github.com/nns7/userapp/models"

	"github.com/gin-gonic/gin"
)

func GetUsers(c *gin.Context) {
	users, err := scanUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error()})
		return
	}
	var resp []models.User
	for _, u := range users {
		u := u
		resp = append(resp, models.User{
			UserID: u.UserID,
			Name:   u.UserName,
		})
	}
	c.JSON(http.StatusOK, resp)
}

func PostUsers(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}
	u := User{
		UserID:   user.UserID,
		UserName: user.Name,
	}

	putUser(u)
	c.JSON(http.StatusOK, user)
}

func DeleteUser(c *gin.Context) {
	user_id := c.Param("user_id")
	if err := deleteUser(user_id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "user has been deleted"})
}

func PutUser(c *gin.Context) {
	user_id := c.Param("user_id")
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}
	u := User{
		UserID:   user_id,
		UserName: user.Name,
	}

	updateUser(u)
	c.JSON(http.StatusOK, user)
}

func SearchUser(c *gin.Context) {
	id := c.Query("user_id")
	user, err := getUser(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "user not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}
