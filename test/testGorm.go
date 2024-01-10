package main

import (
	"ginchat/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(mysql.Open("root:huzeqin@tcp(172.23.176.1:3306)/ginchat?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// 迁移 schema
	db.AutoMigrate(&models.Message{})
	db.AutoMigrate(&models.GroupBasic{})
	db.AutoMigrate(&models.Contact{})

	// Create
	// user := &models.UserBasic{}
	// user.Name = "申专"
	// user.LoginTime = time.Now()
	// user.LoginOutTime = time.Now()
	// user.HeartbeatTime = time.Now()
	// db.Create(user)

	// Read
	// fmt.Println(db.First(user, 2)) // 根据整型主键查找
	// db.First(&product, "code = ?", "D42") // 查找 code 字段值为 D42 的记录

	// Update - 将 product 的 price 更新为 200
	// db.Model(user).Update("Password", "123478")
	// Update - 更新多个字段
	// db.Model(&product).Updates(Product{Price: 200, Code: "F42"}) // 仅更新非零值字段
	// db.Model(&product).Updates(map[string]interface{}{"Price": 200, "Code": "F42"})

	// Delete - 删除 product
	// db.Delete(&product, 1)
}
