package replice

import (
	"srvs/user_srv/models"
	"srvs/user_srv/proto"
)

func UserInfoResponse(user *models.User) *proto.UserInfoResponse {
	return &proto.UserInfoResponse{
		Id:               user.Model.Id,
		Username:         user.Username.String,
		Password:         user.Password,
		Email:            user.Email.String,
		EmailVerified:    user.EmailVerified,
		Avatar:           user.Avatar,
		BackgroundImage:  user.BackgroundImage,
		Description:      user.Description,
		Status:           user.Status,
		ArticleCount:     user.ArticleCount,
		CommentCount:     user.CommentCount,
		Roles:            user.Roles,
		Type:             user.Type,
		Gender:           user.Gender,
		ForbiddenEndTime: user.ForbiddenEndTime,
		CreateTime:       user.CreateTime,
		UpdateTime:       user.UpdateTime,
	}
}
