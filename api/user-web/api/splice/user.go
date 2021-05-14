package splice

import (
	"api/user-web/models"
	"api/user-web/proto"
)

func UserInfo(response *proto.UserInfoResponse) *models.UserInfo {
	return &models.UserInfo{
		Id:               response.Id,
		Username:         response.Username,
		Email:            response.Email,
		EmailVerified:    response.EmailVerified,
		Avatar:           response.Avatar,
		BackgroundImage:  response.BackgroundImage,
		Description:      response.Description,
		Score:            response.Score,
		Status:           response.Status,
		ArticleCount:     response.ArticleCount,
		CommentCount:     response.CommentCount,
		Roles:            response.Roles,
		Type:             response.Type,
		Gender:           response.Gender,
		ForbiddenEndTime: response.ForbiddenEndTime,
		CreateTime:       response.CreateTime,
	}
}
