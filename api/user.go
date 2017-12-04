package api

import (
	"fmt"
	"log"
	"net/http"
	"tx_base_api/models"

	"golang.org/x/crypto/bcrypt"

	"github.com/labstack/echo"
)

// Logon 登录行为
func Logon(c echo.Context) error {
	fmt.Printf("===========>In Logon\n")
	jsonBody := c.Get("JSONBody").(map[string]interface{})
	username := jsonBody["username"].(string)
	password := jsonBody["password"].(string)
	var user models.User
	models.Mdb.First(&user, "name = ?", username)
	log.Print(username)
	log.Print(user)
	if user.ID == 0 {
		return c.JSON(403, map[string]interface{}{
			"success": false,
			"error":   "用户不存在",
		})
	}
	fmt.Printf("===========>hashpassword is %s\n", user.PasswordDigest)
	fmt.Printf("===========>password is %s\n", password)
	result := bcrypt.CompareHashAndPassword([]byte(user.PasswordDigest), []byte(password))
	if result != nil {
		fmt.Printf("===========>password not right,%s\n", result.Error())
		return c.JSON(403, map[string]interface{}{
			"success": false,
			"error":   "密码错误",
		})
	}
	log.Print("user_id is ", user.ID)
	var apiKey models.APIKey
	models.Mdb.First(&apiKey, "login_type = 'User' and login_id = ? expires_at <= now()", user.ID)
	if apiKey.ID == 0 {
		apiKey = models.APIKey{
			LoginType: "User",
			LoginID:   user.ID,
		}
		models.Mdb.Create(&apiKey)
	} else {
		models.Mdb.Save(&apiKey)
	}
	return c.JSON(200, map[string]interface{}{
		"success": true,
		"user":    user,
	})
}

// GetAllUser 获得全部用户
// @Title Get User List
// @Description 获得全部用户
// @Success 200 {object[]} models.User[]
// @Failure 500 get users common error
// @router /get_all_user
func GetAllUser(c echo.Context) error {
	var users []models.User
	models.Mdb.Preload("Role").Preload("UserRoles").Preload("Role").Find(&users)
	result := struct {
		Success bool        `json:"success"`
		UserArr interface{} `json:"user_arr"`
	}{}
	result.Success = true
	result.UserArr = users
	return c.JSON(http.StatusOK, result)
}
