package service

import (
	"context"
	"errors"
	"strconv"

	"github.com/goTouch/TicTok_SimpleVersion/dao"
	"github.com/goTouch/TicTok_SimpleVersion/domain"
	"github.com/goTouch/TicTok_SimpleVersion/util"
)

// Action 进行关注和取消关注，关键维护两个hash：当前用户的关注列表，当前用户的粉丝列表
func FollowAction(userIdInt64 int64, toUserIdInt64 int64, actionType int) (err error) {

	userIdStr := strconv.FormatInt(userIdInt64, 10)
	toUserIdStr := strconv.FormatInt(toUserIdInt64, 10)

	//1. 确认关注对象存在
	toUserNameStr, _ := getUserInfoName(toUserIdInt64)
	if toUserNameStr == "" {
		return errors.New("关注对象不存在")
	}
	userNameStr, _ := getUserInfoName(userIdInt64)

	//2.关注
	if actionType == 1 {
		/*//2.1 先判断是否已经关注  问题：什么时候创建Hash表结构？---不存在自动创建！
		followedStatus := dao.RedisClient.
			HExists(
			....
		HSetNX解决以上所有问题：当哈希表的键key不存在时，会自动创建。当字段field不存在时才会插入，字段存在返回false
		*/

		//关注列表加字段
		followRes := dao.RedisClient.HSetNX(
			context.Background(),
			util.UserFollowHashPrefix+userIdStr,
			toUserIdStr,
			toUserNameStr,
		).Val()
		if !followRes {
			return errors.New("已关注，请勿重复操作")
		}
		//toUser那边，粉丝列表维护起来
		fanInrRes := dao.RedisClient.HSetNX(
			context.Background(),
			util.UserFollowersHashPrefix+toUserIdStr,
			userIdStr,
			userNameStr,
		).Val()
		if !fanInrRes {
			return errors.New("添加粉丝列表失败！")
		}
	} else if actionType == 2 {
		//取关 ：HDel方法删除的时候，如果删除的字段不存在，会返回0

		//当前用户关注列表减少
		unFollowRes := dao.RedisClient.HDel(
			context.Background(),
			util.UserHashKeyPrefix+strconv.FormatInt(userIdInt64, 10),
			strconv.FormatInt(userIdInt64, 10),
		).Val()
		if unFollowRes == 0 {
			//说明本来就没关注
			return errors.New("已经取关，请勿重复操作")
		}
		//目标用户粉丝列表删除对应
		deleteFanRes := dao.RedisClient.HDel(
			context.Background(),
			util.UserFollowersHashPrefix+toUserIdStr,
			userIdStr,
		).Val()
		if deleteFanRes == 0 {
			return errors.New("粉丝删除失败！")
		}
	}
	return nil
}
func getUserInfoName(userIdInt64 int64) (res string, err error) {
	user, err := User(userIdInt64)
	if err != nil {
		return "", err
	}
	return user.Name, nil
}

// FollowList 返回当前用户关注的人。user中只分装name与id
func FollowList(userIdStr string) (userList []domain.User, err error) {

	//查询用户关注的Hash表. 注意HGetAll返回的是map[string]string类型
	userStringMap, err := dao.RedisClient.
		HGetAll(context.Background(), util.UserFollowHashPrefix+userIdStr).
		Result()
	if err != nil || userStringMap == nil {
		return nil, errors.New("redis searching for userFollowList error")
	}
	//遍历
	for userIdStr, username := range userStringMap {
		//一般不会解析错误,忽略错误
		userIdInt64, _ := strconv.ParseInt(userIdStr, 10, 64)
		user := domain.User{
			Id:   userIdInt64,
			Name: username,
		}
		userList = append(userList, user)
	}

	return userList, nil
}
func FollowerList(userIdStr string) (userList []domain.User, err error) {
	userMap := dao.RedisClient.HGetAll(context.Background(), util.UserFollowersHashPrefix+userIdStr).Val()

	for userId, userName := range userMap {
		if err != nil || userMap == nil {
			return nil, errors.New("redis searching for userFollowList error")
		}
		userIdInt64, _ := strconv.ParseInt(userId, 10, 64)
		user := domain.User{
			Id:   userIdInt64,
			Name: userName,
		}
		userList = append(userList, user)
	}
	return userList, nil
}
