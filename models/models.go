/*
* @Author: thepoy
* @Email:  thepoy@aliyun.com
* @File Name: models.go
* @Created:   2021-01-17 12:27:02
* @Modified:  2021-01-21 16:04:08
 */

package models

import (
	"time"

	"gorm.io/gorm"
)

// Blog blog模型
type Blog struct {
	gorm.Model
	Title          string `json:"title" binding:"required" gorm:"type:varchar(40);not null"`
	Content        string `json:"content" binding:"required" gorm:"type:text;not null"`
	Abstract       string `json:"abstract" binding:"required" gorm:"type:varchar(150);not null"`
	BlogTypeID     uint   `json:"type_id" binding:"required"`
	UserID         uint   `json:"user_id"`
	IsPrivate      bool   `json:"is_private" gorm:"type=boolean;not null;defult=false"`
	IsTop          bool   `json:"is_top" gorm:"type=boolean;not null;defult=false"`
	PublishIP      string `json:"publish_ip" gorm:"type:varchar(15);not null"`
	UpdateIP       string `json:"update_ip" gorm:"type:varchar(15)"`
	ClickCount     uint   `json:"click_count" gorm:"not null;default=0"`
	CoverImagePath string `json:"cover_img_path" gorm:"not null"`
}

// LikeBlog 记录用户对某文章的赞或踩
type LikeBlog struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	Status    bool      `json:"status"` // 0是踩，1是赞，不能为空，为空则删除记录
	CreatedAt time.Time `json:"create_at"`
	UserID    uint      `json:"user_id"`
	BlogID    uint      `json:"blog_id"`
	IsRead    bool      `json:"is_read"  gorm:"type=boolean;not null;defult=false"` // 对blog的赞或踩的状态，blog所有者是否已读
}

// User user模型
type User struct {
	gorm.Model
	Username      string     `json:"username" gorm:"type:varchar(20);not null;unique"`
	Password      string     `json:"password" gorm:"not null;type:varchar(256)"`
	Email         string     `json:"email" gorm:"type:varchar(60);not null;unique"`
	Phone         string     `json:"phone" gorm:"type:varchar(20);not null;unique"`
	IsAdmin       bool       `json:"is_admin" gorm:"type=boolean;not null;defult=false"`
	IsWriteOff    bool       `json:"is_write_off" gorm:"type=boolean;not null;defult=false"`
	IsVerified    bool       `json:"is_verify" gorm:"type=boolean;not null;defult=false"`
	RegisterIP    string     `json:"register_ip" gorm:"type:varchar(15);not null"`
	LastLoginTime *time.Time `json:"last_login_time"`
	LastLoginIP   string     `json:"last_login_ip" gorm:"type:varchar(15)"`
	Blogs         []Blog
}

// BlogType blog_type模型
type BlogType struct {
	ID    uint   `json:"id" gorm:"primary_key"`
	Type  string `json:"type" binding:"required" gorm:"not null;unique"`
	Blogs []Blog
}

// Comment 评论
type Comment struct {
	gorm.Model
	Content   string `json:"content" form:"content" binding:"required" gorm:"type:text;not null"`
	BlogID    uint   `json:"blog_id" binding:"required" gorm:"not null"`
	UserID    uint   `json:"user_id"`
	ReplyID   uint   `json:"reply_id"`
	PublishIP string `json:"publish_ip" gorm:"type:varchar(15);not null"`
	UpdateIP  string `json:"update_ip" gorm:"type:varchar(15)"`
	IsRead    bool   `json:"is_read" gorm:"type=boolean;not null;defult=false"` // 对blog的评论或对评论的回复，所有者是否已读
}

// LikeComment 记录用户对某评论的赞或踩
type LikeComment struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	CreatedAt time.Time `json:"create_at"`
	Status    bool      `json:"status"`
	UserID    uint      `json:"user_id"`
	CommentID uint      `json:"comment_id"`
	IsRead    bool      `json:"is_read" gorm:"type=boolean;not null;defult=false"` // 对评论的赞或踩的状态，评论所有者是否已读
}

// Favorite 收藏表
type Favorite struct {
	ID     uint
	BlogID uint
	UserID uint
}

// TableName LikeOrDislike的表名不能用复数
func (LikeBlog) TableName() string {
	return "like_blog"
}

// TableName LikeOrDislike的表名不能用复数
func (LikeComment) TableName() string {
	return "like_comment"
}
