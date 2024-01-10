package service

import (
	"fmt"
	"html/template"

	"github.com/gin-gonic/gin"
)

// GetIndex
// @Tags 首页
// @Success 200 {string} welcome
// @Router /index [get]
func GetIndex(c *gin.Context) {
	ind, err := template.ParseFiles("index.html", "views/chat/head.html")
	if err != nil {
		panic(err)
	}
	err = ind.Execute(c.Writer, "index")
	if err != nil {
		fmt.Println(err)
	}
	// c.JSON(http.StatusOK, gin.H{
	// 	"message": "welcome!",
	// })
}
