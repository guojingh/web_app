package logic

import (
	"go.uber.org/zap"
	"web_app/dao/mysql"
	"web_app/dao/redis"
	"web_app/models"
	"web_app/pkg/snowflake"
)

// CreatePost 创建帖子
func CreatePost(p *models.Post) (err error) {
	//1.生成postID
	p.ID = snowflake.GetID()
	//2.保存到数据库)
	if err = mysql.CreatePost(p); err != nil {
		return
	}
	err = redis.CreatePost(p.ID, p.CommunityID)
	return

}

// GetPostByID 获取帖子详情
func GetPostByID(pid int64) (data *models.ApiPostDetail, err error) {
	//查询并组合我们接口想用的数据
	post, err := mysql.GetPostByID(pid)
	if err != nil {
		zap.L().Error("mysql.GetPostByID(pid) failed", zap.Int64("pid", pid), zap.Error(err))
		return
	}
	//根据作者id查询作者信息
	user, err := mysql.GetUserByID(post.AuthorID)
	if err != nil {
		zap.L().Error("mysql.GetUserByID(pid) failed", zap.Int64("authorID", post.AuthorID), zap.Error(err))
		return
	}
	//根据社区ID查询社区ID
	communityDetail, err := mysql.GetCommunityDetailByID(post.CommunityID)
	if err != nil {
		zap.L().Error("mysql.GetCommunityDetailByID(pid) failed", zap.Int64("communityID", post.CommunityID), zap.Error(err))
		return
	}
	data = &models.ApiPostDetail{
		Post:            post,
		AuthorName:      user.Username,
		CommunityDetail: communityDetail,
	}
	return
}

func GetPostList(page, size int64) (data []*models.ApiPostDetail, err error) {
	posts, err := mysql.GetPostList(page, size)
	if err != nil {
		return nil, err
	}
	data = make([]*models.ApiPostDetail, 0, len(posts))
	for _, post := range posts {
		//根据作者id查询作者信息
		user, err := mysql.GetUserByID(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserByID(pid) failed", zap.Int64("authorID", post.AuthorID), zap.Error(err))
			continue
		}
		//根据社区ID查询社区ID
		communityDetail, err := mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetailByID(pid) failed", zap.Int64("communityID", post.CommunityID), zap.Error(err))
			continue
		}
		postDetail := &models.ApiPostDetail{
			Post:            post,
			AuthorName:      user.Username,
			CommunityDetail: communityDetail,
		}
		data = append(data, postDetail)
	}
	return
}

func GetPostList2(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	// 去redis查询id列表
	ids, err := redis.GetPostIDsInOrder(p)
	if err != nil {
		return
	}

	if ids == nil {
		zap.L().Warn("redis.GetPostList2(p) return 0 data")
		return
	}
	//根据id去mysql数据库查询详细信息
	//返回的时候根据给定的数据返回
	posts, err := mysql.GetPostListByIDs(ids)
	if err != nil {
		return
	}
	//提前查询好每篇帖子的投票数
	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		return
	}

	//将帖子的作者及分区信息查询出来填充到帖子中
	for idx, post := range posts {
		//根据作者id查询作者信息
		user, err := mysql.GetUserByID(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserByID(pid) failed", zap.Int64("authorID", post.AuthorID), zap.Error(err))
			continue
		}
		//根据社区ID查询社区ID
		communityDetail, err := mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetailByID(pid) failed", zap.Int64("communityID", post.CommunityID), zap.Error(err))
			continue
		}
		postDetail := &models.ApiPostDetail{
			Post:            post,
			VoteNum:         voteData[idx],
			AuthorName:      user.Username,
			CommunityDetail: communityDetail,
		}
		data = append(data, postDetail)
	}
	return
}

func GetCommunityPostList(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	// 去redis查询id列表
	ids, err := redis.GetCommunityPostIDsInOrder(p)
	if err != nil {
		return
	}

	if ids == nil {
		zap.L().Warn("redis.GetPostList2(p) return 0 data")
		return
	}
	//根据id去mysql数据库查询详细信息
	//返回的时候根据给定的数据返回
	posts, err := mysql.GetPostListByIDs(ids)
	if err != nil {
		return
	}
	//提前查询好每篇帖子的投票数
	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		return
	}

	//将帖子的作者及分区信息查询出来填充到帖子中
	for idx, post := range posts {
		//根据作者id查询作者信息
		user, err := mysql.GetUserByID(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserByID(pid) failed", zap.Int64("authorID", post.AuthorID), zap.Error(err))
			continue
		}
		//根据社区ID查询社区ID
		communityDetail, err := mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetailByID(pid) failed", zap.Int64("communityID", post.CommunityID), zap.Error(err))
			continue
		}
		postDetail := &models.ApiPostDetail{
			Post:            post,
			VoteNum:         voteData[idx],
			AuthorName:      user.Username,
			CommunityDetail: communityDetail,
		}
		data = append(data, postDetail)
	}
	return
}

// GetPostListNew 将两个查询逻辑合二为一
func GetPostListNew(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	if p.CommunityID == 0 {
		//查所有
		data, err = GetPostList2(p)
	} else {
		//根据社区id查询
		data, err = GetCommunityPostList(p)
	}
	if err != nil {
		zap.L().Error("logic GetPostListNew failed", zap.Error(err))
		return nil, err
	}
	return
}
