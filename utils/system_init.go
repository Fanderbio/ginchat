package utils

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitConfig() {
	viper.SetConfigName("app")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("config/")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("config app inited")
}

func InitMySQL() {
	// 自定义日志模板，打印SQL语句
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second, // 慢SQL阈值
			LogLevel:      logger.Info, // 日志级别
			Colorful:      true,        // 彩色
		},
	)

	DB, _ = gorm.Open(mysql.Open(viper.GetString("mysql.dns")),
		&gorm.Config{Logger: newLogger})

	// 初始化时确保数据库中可以存入时间为'0000-00-00 00:00:00'的值
	DB.Exec("SET @@sql_mode = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION';")

	fmt.Println("MySQL inited:")
}
