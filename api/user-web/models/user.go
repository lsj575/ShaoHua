package models

type UserInfo struct {
	Id               uint64 `json:"id"`
	Username         string `json:"username" form:"username"`                 // 用户名
	Email            string `json:"email" form:"email"`                       // 邮箱
	EmailVerified    bool   `json:"emailVerified" form:"emailVerified"`       // 邮箱是否验证
	Avatar           string `json:"avatar" form:"avatar"`                     // 头像
	BackgroundImage  string `json:"backgroundImage" form:"backgroundImage"`   // 个人中心背景图片
	Description      string `json:"description" form:"description"`           // 个人描述
	Score            uint64 `json:"score" form:"score"`                       // 积分
	Status           int32  `json:"status" form:"status"`                     // 状态
	ArticleCount     uint64 `json:"topicCount" form:"topicCount"`             // 帖子数量
	CommentCount     uint64 `json:"commentCount" form:"commentCount"`         // 评论数量
	Roles            int32  `json:"roles" form:"roles"`                       // 角色
	Type             int32  `json:"type" form:"type"`                         // 用户类型
	Gender           uint32 `json:"gender" form:"gender"`                     // 用户性别
	ForbiddenEndTime uint64 `json:"forbiddenEndTime" form:"forbiddenEndTime"` // 禁言结束时间
	CreateTime       uint64 `json:"createTime" form:"createTime"`             // 创建时间
}
