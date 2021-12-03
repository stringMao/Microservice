package db

import (
	"Common/log"
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
)

//github.com/gomodule/redigo/redis库文档: https://pkg.go.dev/github.com/gomodule/redigo/redis#pkg-overview

type RedisData struct {
	Host         string
	Port         int
	UserName     string
	PassWord     string
	DataBase     int //连接几号db
	MaxIdleConns int //连接池的空闲数大小
	MaxOpenConns int //连接池最大打开连接数
	Pool         *redis.Pool
}

func (r *RedisData) Open() bool {
	//setdb := redis.DialDatabase(1)
	//setPasswd := redis.DialPassword(r.PassWord)
	r.Pool = &redis.Pool{
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", fmt.Sprintf("%s:%d", r.Host, r.Port))
			if err != nil {
				log.Logger.Fatalln("redis Dial err:", err)
				return nil, err
			}
			//使用 Dial 函数通过 AUTH 命令验证连接或使用 SELECT 命令选择数据库：
			if _, err := c.Do("AUTH", r.PassWord); err != nil {
				c.Close()
				log.Logger.Fatalln("redis DO AUTH err:", err)
				return nil, err
			}
			if _, err := c.Do("SELECT", r.DataBase); err != nil {
				c.Close()
				log.Logger.Fatalln("redis SELECT err:", err)
				return nil, err
			}
			return c, nil
		},
		//在将连接返回给应用程序之前，使用 TestOnBorrow 函数检查空闲连接的运行状况。此示例 PING 已空闲超过一分钟的连接：
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
		MaxIdle:   r.MaxIdleConns,
		MaxActive: r.MaxOpenConns,
	}

	// r.Pool = redis.NewPool(func() (redis.Conn, error) {
	// 	c, err := redis.Dial("tcp", fmt.Sprintf("%s@%s:%d", r.PassWord, r.Host, r.Port))
	// 	if err != nil {
	// 		log.Logger.Fatalln("redis connet is err:", err)
	// 		return nil, err
	// 	}
	// 	return c, nil
	// }, r.MaxOpenConns)

	//test redis connect

	//log.Logger.Info("redis init success")
	// pool := &redis.Pool{
	// 	MaxActive:   100,                              //  最大连接数，即最多的tcp连接数，一般建议往大的配置，但不要超过操作系统文件句柄个数（centos下可以ulimit -n查看）
	// 	MaxIdle:     10,                               // 最大空闲连接数，即会有这么多个连接提前等待着，但过了超时时间也会关闭。
	// 	IdleTimeout: time.Duration(100) * time.Second, // 空闲连接超时时间，但应该设置比redis服务器超时时间短。否则服务端超时了，客户端保持着连接也没用
	// 	Wait:        true,                             // 当超过最大连接数 是报错还是等待， true 等待 false 报错
	// 	Dial: func() (redis.Conn, error) {
	// 		conn, err := redis.Dial("tcp", fmt.Sprintf("%s:%s", host, port))
	// 		if err != nil {
	// 			log.Logger.Fatalln("redis connet is err:", err)
	// 			return nil, err
	// 		}
	// 		return conn, nil
	// 	},
	// }
	return true
}

//GetRedis ..
func (r *RedisData) GetRedis() redis.Conn {
	return r.Pool.Get()
}
