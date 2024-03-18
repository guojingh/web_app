package mysql

import (
	"github.com/jmoiron/sqlx"
	"strings"
	"web_app/models"
)

// CreatePost 创建帖子
func CreatePost(p *models.Post) (err error) {
	sqlStr := `insert into post (post_id, title, content, author_id, community_id) values (?,?,?,?,?)`

	_, err = db.Exec(sqlStr, p.ID, p.Title, p.Content, p.AuthorID, p.CommunityID)
	return
}

// GetPostByID 根据ID查询单个贴子数据
func GetPostByID(pid int64) (post *models.Post, err error) {
	post = new(models.Post)
	sqlStr := "select post_id, title, content, author_id, status, community_id, create_time from post where post_id = ?"

	err = db.Get(post, sqlStr, pid)
	return
}

// GetPostList 查询帖子列表
func GetPostList(page, size int64) (post []*models.Post, err error) {
	sqlStr := `select
		post_id, title, content, author_id, status, community_id, create_time 
		from post 
		ORDER BY create_time
		DESC 
		limit ?,?`

	post = make([]*models.Post, 0, 2)
	err = db.Select(&post, sqlStr, (page-1)*size, size)
	return
}

// GetPostListByIDs 根据给定的ID列表查询帖子
func GetPostListByIDs(ids []string) (postList []*models.Post, err error) {
	sqlStr := `select post_id, title, content, author_id, community_id, create_time
				from post
				where post_id in (?)
				order by FIND_IN_SET(post_id, ?)`

	query, args, err := sqlx.In(sqlStr, ids, strings.Join(ids, ","))
	if err != nil {
		return nil, err
	}

	query = db.Rebind(query)
	err = db.Select(&postList, query, args...)
	return
}