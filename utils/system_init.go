package utils

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB      *gorm.DB
	RedisDB *redis.Client
	ctx     = context.Background()
)

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

	fmt.Println("MySQL inited")
}

func InitRedis() {
	redisAddr := viper.GetString("redis.addr")
	redisPwd := viper.GetString("redis.password")
	redisDB := viper.GetInt("redis.DB")
	redisPoolSize := viper.GetInt("redis.poolSize")
	redisMinIdleConn := viper.GetInt("redis.minIdleConn")

	RedisDB = redis.NewClient(&redis.Options{
		Addr:         redisAddr,
		Password:     redisPwd,
		DB:           redisDB,
		PoolSize:     redisPoolSize,
		MinIdleConns: redisMinIdleConn,
	})

	pong, err := RedisDB.Ping(ctx).Result()

	if err != nil {
		fmt.Printf("Ping failed, err: %v\n", err)
	} else {
		fmt.Printf("Ping success, pong: %v\n", pong)
	}
}

const (
	PublishKey = "websocket"
)

// Publish发布消息到Redis
func Publish(ctx context.Context, channel string, msg string) error {
	fmt.Println("Publish:", msg)
	err := RedisDB.Publish(ctx, channel, msg).Err()
	if err != nil {
		fmt.Println(err)
	}
	return err
}

// Subscribe 订阅消息
func Subscribe(ctx context.Context, channel string) (string, error) {
	sub := RedisDB.Subscribe(ctx, channel)
	fmt.Println("Subscribe收到消息", ctx)
	msg, err := sub.ReceiveMessage(ctx)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	fmt.Println("Subscribe", msg.Payload)
	return msg.Payload, err
}
