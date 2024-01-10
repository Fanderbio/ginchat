package service

import (
	"fmt"
	"ginchat/models"
	"ginchat/utils"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// GetUserList
// @Summary 所有用户
// @Tags 用户模块
// @Success 200 {string} json{"code", "message"}
// @Router /user/getUserList [get]
func GetUserList(c *gin.Context) {
	data := models.GetUserList()

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "获取列表成功",
		"data":    data,
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
			"code":    -1, // 0：成功；-1失败
			"message": "用户名已经注册",
			"data":    user.Name,
		})
		return
	}

	if password != repassword {
		c.JSON(-1, gin.H{
			"code":    -1, // 0：成功；-1失败
			"message": "密码不一致",
			"data":    gin.H{"password": password, "repassword": repassword}})
		return
	}
	user.Password = utils.MakePassword(password, salt)
	user.Salt = salt

	models.CreateUser(user)
	c.JSON(http.StatusOK, gin.H{
		"code":    0, // 0：成功；-1失败
		"message": "新增用户成功！",
		"data":    user,
	})
}

// GetUserList
// @Summary 用户登录
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
			"code":    -1, // 0：成功；-1失败
			"message": "该用户不存在",
			"data":    nil,
		})
		return
	}

	fmt.Println(queryUser)
	flag := utils.ValidPassword(password, queryUser.Salt, queryUser.Password)
	if !flag {
		c.JSON(http.StatusOK, gin.H{
			"code":    -1, // 0：成功；-1失败
			"message": "密码错误",
			"data":    nil,
		})
		return
	}

	data := models.FindUserByNameAndPwd(name, queryUser.Password)

	c.JSON(http.StatusOK, gin.H{
		"code":    0, // 0：成功；-1失败
		"message": "登录成功",
		"data":    data,
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
		"code":    0, // 0：成功；-1失败
		"message": "删除用户成功！",
		"data":    user,
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
			"code":    -1, // 0：成功；-1失败
			"message": "修改参数不匹配",
			"data":    user,
		})
		return
	}

	models.UpdateUser(user)
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "修改用户成功！",
		"data":    user,
	})
}

// 防止跨域站点伪造请求
var upGrade = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func SendMsg(c *gin.Context) {
	// 以下为发送消息的准备工作
	ws, err := upGrade.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer func(ws *websocket.Conn) {
		err = ws.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(ws)

	// 发送流程
	MsgHeader(ws, c)
}

func MsgHeader(ws *websocket.Conn, c *gin.Context) {
	for {
		msg, err := utils.Subscribe(c, utils.PublishKey)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("发送消息", msg)
		tm := time.Now().Format("2006-01-02 15:04:05")
		m := fmt.Sprintf("[ws][%s]:%s\n", tm, msg)
		err = ws.WriteMessage(1, []byte(m))
		if err != nil {
			fmt.Println(err)
		}
	}
}

func SendUserMsg(c *gin.Context) {
	models.Chat(c.Writer, c.Request)
}
