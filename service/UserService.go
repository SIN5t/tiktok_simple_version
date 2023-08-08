package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bwmarrin/snowflake"
	"github.com/goTouch/TicTok_SimpleVersion/dao"
	"github.com/goTouch/TicTok_SimpleVersion/domain"
	"golang.org/x/crypto/bcrypt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

// LoginLimit 中间件服务，限制注册登录操作过于频繁。
func LoginLimit(ipAddress string) bool {
	// 错误可忽略
	times, _ := dao.RdbToken.Get(context.Background(), ipAddress).Int64()
	if times > 10 {
		return false
	} else {
		dao.RdbToken.Set(context.Background(), ipAddress, times+1, time.Minute)
	}
	return true
}
func SaveTokenToRedis(userIdInt64 int64) (token int64) {
	//雪花算法生成token
	node, err := snowflake.NewNode(userIdInt64) //这里的userIdInt64就是 User.Id(主键)
	if err != nil {
		log.Println("雪花算法生成id错误!")
		log.Println(err)
	}
	token = node.Generate().Int64()
	// 检查键是否存在
	key := strconv.FormatInt(token, 10)
	exists, err := dao.RedisClient.Exists(context.Background(), key).Result()
	if exists == 1 {
		fmt.Println("token已经存在")
		return
	}
	return token
}

// 随机盐长度固定为4
func randSalt() string {
	buf := strings.Builder{}
	for i := 0; i < 4; i++ {
		// 如果写byte会无法兼容mysql编码
		buf.WriteRune(rune(rand.Intn(256)))
	}
	return buf.String()
}
func SaveToRedis(user domain.User) (err error, jsonUser string) {
	// 将结构体序列化为JSON字符串
	jsonData, err := json.Marshal(user)
	if err != nil {
		fmt.Println("序列化失败:", err)
		return
	}
	// 返回序列化后的结构体
	jsonUser1 := string(jsonData)
	return nil, jsonUser1
}
func Register(username, password string) (id int64, err error) {
	if len(username) > 32 {
		return 0, errors.New("用户名过长，不可超过32位")
	}
	if len(password) > 32 {
		return 0, errors.New("密码过长，不可超过32位")
	}

	user := domain.User{}
	dao.DB.Model(&domain.User{}).Where("name = ?", username).Find(&user)
	if user.Id != 0 {
		return 0, errors.New("用户已存在")
	}
	user.Name = username
	// 加密存储用户密码
	user.Salt = randSalt()
	buf := bytes.Buffer{}
	buf.WriteString(username)
	buf.WriteString(password)
	buf.WriteString(user.Salt)
	pwd, err1 := bcrypt.GenerateFromPassword(buf.Bytes(), bcrypt.MinCost)
	if err1 != nil {
		return 0, err
	}
	user.Pwd = string(pwd)
	dao.DB.Model(&domain.User{}).Create(&user)
	//为每个用户生成一个token作为唯一标识存储在redis
	token := SaveTokenToRedis(user.Id)
	fmt.Println(token)
	//将没有加密的信息用户信息存在redis中去减轻mysql查询的压力
	_, userdata := SaveToRedis(user)
	fmt.Println(userdata)
	//假设永久有效
	dao.RedisClient.Set(context.Background(), strconv.FormatInt(token, 10), userdata, 0)
	return user.Id, nil
}

func Login(username, password string) (id int64, err error) {
	user := domain.User{}
	dao.DB.Model(&domain.User{}).Where("name = ?", username).Find(&user)
	if user.Id != 0 {
		user.Name = username
		// 加密存储用户密码
		user.Salt = randSalt()
		buf := bytes.Buffer{}
		buf.WriteString(username)
		buf.WriteString(password)
		buf.WriteString(user.Salt)
		pwd, err1 := bcrypt.GenerateFromPassword(buf.Bytes(), bcrypt.MinCost)
		if err1 != nil {
			return 0, err
		}
		user.Pwd = string(pwd)
		dao.DB.Model(&domain.User{}).Where("name = ? AND pwd = ?", username, password).Find(&user)
		if user.Id != 0 {
			return 1, errors.New("用户登陆成功！")
		} else {
			return 0, errors.New("密码错误！")
		}
	} else {
		return 0, errors.New("用户不存在！")
	}
}
