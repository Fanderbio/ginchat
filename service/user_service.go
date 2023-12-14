package service

import (
	"fmt"
	"ginchat/models"
	"ginchat/utils"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
)

// GetUserList
// @Summary 所有用户
// @Tags 用户模块
// @Success 200 {string} json{"code", "message"}
// @Router /user/getUserList [get]
func GetUserList(c *gin.Context) {
	data := models.GetUserList()

	c.JSON(http.StatusOK, gin.H{
		"message": data,
	})
}

// CreateUser
// @Summary 新增用户
// @Tags 用户模块
// @param name query string false "用户名"
// @param password query string false "密码"
// @param repassword query string false "确认密码"
// @Success 200 {string} json{"code", "message"}
// @Router /user/createUser [get]
func CreateUser(c *gin.Context) {
	user := models.UserBasic{}
	user.Name = c.Query("name")
	password := c.Query("password")
	repassword := c.Query("repassword")

	salt := fmt.Sprintf("%06d", rand.Int31())

	userFind := models.FindUserByName(user.Name)
	if userFind.Name != "" {
		c.JSON(-1, gin.H{
			"message": "用户名已经注册",
		})
		return
	}

	if password != repassword {
		c.JSON(-1, gin.H{
			"message": "密码不一致",
		})
		return
	}
	user.Password = utils.MakePassword(password, salt)
	user.Salt = salt

	models.CreateUser(user)
	c.JSON(http.StatusOK, gin.H{
		"message": "新增用户成功！",
	})
}

// GetUserList
// @Summary 所有用户
// @Tags 用户模块
// @param name query string false "用户名"
// @param password query string false "密码"
// @Success 200 {string} json{"code", "message"}
// @Router /user/findUserByNameAndPwd [post]
func FindUserByNameAndPwd(c *gin.Context) {
	name := c.Query("name")
	password := c.Query("password")

	queryUser := models.FindUserByName(name)
	if queryUser.Name == "" {
		c.JSON(http.StatusOK, gin.H{
			"message": "该用户不存在",
		})
		return
	}

	fmt.Println(queryUser)
	flag := utils.ValidPassword(password, queryUser.Salt, queryUser.Password)
	if !flag {
		c.JSON(http.StatusOK, gin.H{
			"message": "密码错误",
		})
		return
	}

	data := models.FindUserByNameAndPwd(name, queryUser.Password)

	c.JSON(http.StatusOK, gin.H{
		"message": data,
	})
}

// DeleteUser
// @Summary 删除用户
// @Tags 用户模块
// @param id query string false "id"
// @Success 200 {string} json{"code", "message"}
// @Router /user/deleteUser [get]
func DeleteUser(c *gin.Context) {
	user := models.UserBasic{}
	id, _ := strconv.Atoi(c.Query("id"))
	user.ID = uint(id)
	models.DeleteUser(user)
	c.JSON(http.StatusOK, gin.H{
		"message": "删除用户成功！",
	})
}

// UpdateUser
// @Summary 修改用户
// @Tags 用户模块
// @param id formData string false "id"
// @param name formData string false "name"
// @param password formData string false "password"
// @param phone formData string false "phone"
// @param email formData string false "email"
// @Success 200 {string} json{"code", "message"}
// @Router /user/updateUser [post]
func UpdateUser(c *gin.Context) {
	user := models.UserBasic{}
	id, _ := strconv.Atoi(c.PostForm("id"))
	user.ID = uint(id)

	user.Name = c.PostForm("name")
	user.Password = c.PostForm("password")
	user.Phone = c.PostForm("phone")
	user.Email = c.PostForm("email")

	fmt.Println("update:", user)

	_, err := govalidator.ValidateStruct(user)
	if err != nil {
		fmt.Println(err)
		c.JSON(200, gin.H{
			"message": "修改参数不匹配",
		})
		return
	}

	models.UpdateUser(user)
	c.JSON(http.StatusOK, gin.H{
		"message": "修改用户成功！",
	})
}
