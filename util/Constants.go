package util

const VideoFavoriteKeyPrefix = "VIDEO_FAVORITE_KEY:"   //用户角度，用处：该用户点赞视频
const AuthorFollowedKeyPrefix = "AUTHOR_FOLLOWED_KEY:" //用户角度，用处：该用户点关注的人
const StaticRooterPrefix = "http://127.0.0.1:8080/"    //
const AuthorBeLikedNum = "AUTHOR_BE_LIKED_NUM_KEY:"

const UserHashKeyPrefix = "USER_HASH_KEY:" //用户角度：用户的各个字段使用hash存储到redis中，其中Hset对应的key前缀
