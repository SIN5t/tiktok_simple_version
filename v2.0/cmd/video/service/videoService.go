package service

import (
	"context"
	"github.com/prometheus/common/log"
	"tiktok_v2/dal/dao"
	video "tiktok_v2/kitex_gen/video"
)

// VideoServiceImpl implements the last service interface defined in the IDL.
type VideoServiceImpl struct{}

// Feed implements the VideoServiceImpl interface.
// 要求按照投稿时间倒叙输出
func (s *VideoServiceImpl) Feed(ctx context.Context, req *video.FeedRequest) (resp *video.FeedResponse, err error) {
	latestTime := req.LatestTime
	//数据库获取视频
	videoListRaw, err := dao.GetVideosByLastTime(ctx, latestTime, 30)
	if err != nil {
		log.Error(err)
		response := &video.FeedResponse{
			StatusCode: 1,
			StatusMsg:  "视频查询失败",
		}
		return response, nil
	}

	flag := false
	if req.Token != "" {
		//TODO 校验 jwt
		flag = true //已登入
	}

	videoListRes := make([]*video.Video, 0)
	for _, value := range videoListRaw {
		curVideo := video.Video{
			Id:            value.ID,
			AuthorId:      value.AuthorID,
			PlayUrl:       value.PlayUrl,
			CoverUrl:      value.CoverUrl,
			FavoriteCount: 0,
			CommentCount:  0,
			IsFavorite:    false,
			Title:         "",
		}
		//如果登入
		if flag {
			//TODO 当前视频是否点赞

			//当前用户是否关注视频作者

		}
		videoListRes = append(videoListRes, &curVideo)
	}
	return &video.FeedResponse{
		StatusCode: 0,
		StatusMsg:  "刷新成功",
		VideoList:  videoListRes,
		NextTime:   videoListRaw[len(videoListRaw)].CreatedAt.UnixMilli() - 1, //注意下次的时间是这次的最后
	}, nil
}

// PublishAction implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) PublishAction(ctx context.Context, req *video.PublishActionRequest) (resp *video.PublishActionResponse, err error) {
	// TODO: Your code here...
	return
}

// PublishList implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) PublishList(ctx context.Context, req *video.PublishListRequest) (resp *video.PublishListResponse, err error) {
	// TODO: Your code here...
	return
}
