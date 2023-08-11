PROJECT_ID="tiktok"
function v() {
    echo "$PROJECT_ID:$1"
}

# mysql
export `v mysql_user`="root"
export `v mysql_pswd`="123456"
export `v mysql_addr`="localhost:3306"
export `v mysql_name`="tiktok"

# redis
export `v redis_addr`="localhost:6379"
export `v redis_pswd`="123456"
export `v redis_master_name`="mymaster"
export `v redis_sentinel_addrs`=":17000 :17001 :17002 "

# jwt
export `v jwt_secret`="f05ad7412aa192ddc121ba50a64e585943b3e6d8fca4a3a19a8eea26e76496a7"