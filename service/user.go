package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/goForward/tictok_simple_version/dao"
	"github.com/goForward/tictok_simple_version/domain"
	"github.com/goForward/tictok_simple_version/util"
	"golang.org/x/crypto/bcrypt"
)

// TODO: LoginLimit 中间件服务，限制注册登录操作过于频繁。
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

func Register(username, password string) (id int64, tokenString string, err error) {
	// 校验数据合法性
	if len(username) > 32 {
		return 0, "", errors.New("用户名过长，不可超过32位")
	}
	if len(password) > 32 {
		return 0, "", errors.New("密码过长，不可超过32位")
	}

	// 判断用户是否存在
	user := domain.User{}
	dao.DB.Model(&domain.User{}).Where("name = ?", username).Find(&user)
	if user.Id != 0 {
		return 0, "", errors.New("用户已存在")
	}

	// 加盐加密存储用户密码
	user.Name = username
	user.Salt = randSalt()
	buf := bytes.Buffer{}
	buf.WriteString(username)
	buf.WriteString(password)
	buf.WriteString(user.Salt)
	pwd, err := bcrypt.GenerateFromPassword(buf.Bytes(), bcrypt.MinCost)
	if err != nil {
		return 0, "", fmt.Errorf("bcrypt加密错误: %w", err)
	}
	user.Pwd = string(pwd)

	// 创建用户
	dao.DB.Model(&domain.User{}).Create(&user)

	// 生成jwt
	tokenString, err = GenerateJWT(user.Id, util.JWTSecret())
	if err != nil {
		return 0, "", errors.New("生成jwt错误")
	}

	// 缓存jwt
	dao.RedisClient.Set(context.Background(), tokenString, user.Id, 0)

	return user.Id, tokenString, nil
}

func Login(username, password string) (id int64, token string, err error) {
	// 查询用户
	user := domain.User{}
	dao.DB.Model(&domain.User{}).Where("name = ?", username).Find(&user)
	if user.Id == 0 {
		err = errors.New("用户不存在！")
		return
	}

	// 核对密码
	buf := bytes.Buffer{}
	buf.WriteString(username)
	buf.WriteString(password)
	buf.WriteString(user.Salt)
	if err = bcrypt.CompareHashAndPassword([]byte(user.Pwd), buf.Bytes()); err != nil {
		err = errors.New("密码错误")
		return
	}

	// 生成jwt
	tokenString, err := GenerateJWT(user.Id, util.JWTSecret())
	if err != nil {
		err = errors.New("生成jwt错误")
		return
	}

	// 缓存jwt
	err = dao.RedisClient.Set(context.Background(), util.Key(util.TokenRefreshPrefix, strconv.FormatInt(user.Id, 10)), tokenString, time.Hour*24).Err()
	if err != nil {
		return
	}

	return user.Id, tokenString, nil
}

func User(userId int64) (user domain.User, err error) {
	if userId <= 0 {
		return user, errors.New("不合法的用户id")
	}
	err = dao.DB.Model(&domain.User{}).Where("id = ?", userId).Find(&user).Error
	if err != nil {
		return domain.User{}, errors.New("用户不存在")
	}
	userFollowNum := dao.RedisClient.HLen(context.Background(), util.Key(util.UserFollowHashPrefix, userId)).Val()
	userFollowerNum := dao.RedisClient.HLen(context.Background(), util.Key(util.UserFollowersHashPrefix, userId)).Val()
	user.FollowCount = userFollowNum
	user.FollowerCount = userFollowerNum
	return user, nil
}

func GenerateJWT(userId int64, secret string) (tokenString string, err error) {
	// 指定用HS256算法
	token := jwt.New(jwt.SigningMethodHS256)

	// token声明：sub (subject), iat (issued at time)
	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = userId
	claims["iat"] = time.Now().Unix()

	// token签名
	tokenString, err = token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func VerifyJWT(tokenString, secret string) (userId int64, err error) {
	// 解析和验证token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		// 验证算法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("不支持的加密算法: %v", token.Header["alg"])
		}

		// 返回密钥
		return []byte(secret), nil
	})
	if err != nil {
		return 0, fmt.Errorf("验证失败：%w", err)
	}

	// 获取声明
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		sub := int64(claims["sub"].(float64))
		cachedTokenString, err := dao.RedisClient.Get(context.Background(), util.Key(util.TokenRefreshPrefix, sub)).Result()
		if err != nil {
			return 0, err
		}
		if tokenString != cachedTokenString {
			return 0, errors.New("验证失败，请重新登录")
		}
		// iat := claims["iat"].(float64)
		return sub, nil
	}
	return 0, errors.New("验证失败，无法获取声明")
}

func RefreshJWT(userId int64) (err error) {
	return dao.RedisClient.Expire(context.Background(), util.Key(util.TokenRefreshPrefix, userId), time.Hour*24).Err()
}

// 随机盐长度固定为4
func randSalt() (salt string) {
	buf := strings.Builder{}
	for i := 0; i < 4; i++ {
		// 如果写byte会无法兼容mysql编码
		buf.WriteRune(rune(rand.Intn(256)))
	}
	return buf.String()
}
