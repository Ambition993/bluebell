package mysql

import (
	"bluebell/models"
	"github.com/jmoiron/sqlx"
	"strings"
)

// CreatePost 创建一个post记录
func CreatePost(p *models.Post) (err error) {
	sqlStr := `insert into post
	(post_id,title,content,author_id, community_id)
	values(?,?,?,?,?)`
	_, err = db.Exec(sqlStr, p.ID, p.Title, p.Content, p.AuthorID, p.CommunityID)
	return err
}

// GetPostById 根据一个ID查找一个post记录
func GetPostById(pid int64) (post *models.Post, err error) {
	post = new(models.Post)
	sqlStr := `select post_id, title , content, author_id , create_time , community_id
               from post where post_id = ?`
	err = db.Get(post, sqlStr, pid)
	return post, err
}

// GetPostList 查询post列表
func GetPostList(page int64, size int64) (posts []*models.Post, err error) {
	sqlStr := `select post_id, title , content, author_id , create_time , community_id
               from post
               order by create_time
               desc 
               limit ?, ?`
	posts = make([]*models.Post, 0, 2)
	err = db.Select(&posts, sqlStr, (page-1)*size, size)
	return
}

// GetPostListByIDs 根据ID列表查询post
func GetPostListByIDs(IDs []string) (postList []*models.Post, err error) {
	sqlStr := `select 
	post_id, title , content, author_id, community_id , create_time 
	from post 
	where post_id in (?)
	order by  FIND_IN_SET(post_id, ?)`
	query, args, err := sqlx.In(sqlStr, IDs, strings.Join(IDs, ","))
	if err != nil {
		return nil, err
	}
	query = db.Rebind(query)
	err = db.Select(&postList, query, args...) // args...
	return
}
