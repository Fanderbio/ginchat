package models

import (
	"fmt"
	"ginchat/utils"
	"time"

	"gorm.io/gorm"
)

type UserBasic struct {
	gorm.Model
	Identity   string
	Name       string
	Password   string
	Phone      string `valid:"matches(^1[3-9]{1}[0-9]{9}$)"`
	Email      string `valid:"email"`
	ClientIp   string
	ClientPort string

	Salt string

	LoginTime     time.Time
	HeartbeatTime time.Time
	LoginOutTime  time.Time `gorm:"column:login_out_time" json:"login_out_time"`

	IsLogout   bool
	DeviceInfo string
}

func (table *UserBasic) TableName() string {
	return "user_basic"
}

func GetUserList() []*UserBasic {
	data := make([]*UserBasic, 10)
	utils.DB.Find(&data)

	for _, v := range data {
		fmt.Println(v)
	}

	return data
}

func FindUserByNameAndPwd(name, password string) UserBasic {
	user := UserBasic{}
	utils.DB.Where("name = ? and password = ?", name, password).First(&user)

	// token加密
	str := fmt.Sprintf("%d", time.Now().Unix())
	temp := utils.MD5Encode(str)

	utils.DB.Model(&user).Where("id = ?", user.ID).Update("identity", temp)
	return user
}

func FindUserByName(name string) UserBasic {
	user := UserBasic{}
	utils.DB.Where("name = ?", name).First(&user)
	return user
}

func FindUserByPhone(phone string) UserBasic {
	user := UserBasic{}
	utils.DB.Where("phone = ?", phone).First(&user)
	return user
}

func FindUserByEmail(email string) UserBasic {
	user := UserBasic{}
	utils.DB.Where("email = ?", email).First(&user)
	return user
}

func CreateUser(user UserBasic) *gorm.DB {
	// 首先判断表是否存在
	migrator := utils.DB.Migrator()
	exist := migrator.HasTable(user.TableName())
	// exist := migrator.HasTable("user_basic")
	if !exist {
		fmt.Println("表不存在")
		utils.DB.AutoMigrate(&user)
	} else {
		fmt.Println("表已经存在")
	}

	fmt.Println("创建用户")
	return utils.DB.Create(&user)
}

func DeleteUser(user UserBasic) *gorm.DB {
	return utils.DB.Delete(&user)
}

func UpdateUser(user UserBasic) *gorm.DB {
	return utils.DB.Model(&user).Updates(UserBasic{
		Name:     user.Name,
		Password: user.Password,
		Phone:    user.Phone,
		Email:    user.Email,
	})
}
