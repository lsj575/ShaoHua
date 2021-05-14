package services

import (
	"context"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"srvs/user_srv/db"
	"srvs/user_srv/models"
	"srvs/user_srv/models/constants"
	"srvs/user_srv/proto"
	"srvs/user_srv/repositories"
	"srvs/user_srv/services/replice"
	"srvs/user_srv/validate"
	"strings"
	"time"
)

var UserService = newUserService()

//var mysqlDB = mysql.GetMysqlInstance().GetMysqlDB()

func newUserService() *userService {
	return &userService{}
}

type userService struct {
}

func (us *userService) GetUserById(ctx context.Context, request *proto.IDRequest) (*proto.UserInfoResponse, error) {
	//mysqlDB := mysql.GetMysqlInstance().GetMysqlDB()
	userInfo := repositories.UserRepository.GetUserById(db.GetMysqlInstance().GetMysqlDB(), request.Id)

	if userInfo != nil {
		return replice.UserInfoResponse(userInfo), nil
	} else {
		return nil, status.Errorf(codes.NotFound, "user not exist")
	}
}

func (us *userService) GetUserByUsername(username string) (*models.User, error) {
	userInfo := repositories.UserRepository.GetByUsername(db.GetMysqlInstance().GetMysqlDB(), username)

	if userInfo != nil {
		return userInfo, nil
	} else {
		return nil, errors.New("用户不存在")
	}
}

func (us *userService) GetUserList(ctx context.Context, request *proto.PageInfo) (*proto.UserListResponse, error) {
	return &proto.UserListResponse{}, nil
}

func (us *userService) GetUserByEmail(ctx context.Context, request *proto.EmailRequest) (*proto.UserInfoResponse, error) {
	userInfo := repositories.UserRepository.GetByEmail(db.GetMysqlInstance().GetMysqlDB(), request.Email)
	if userInfo != nil {
		return replice.UserInfoResponse(userInfo), nil
	}
	return nil, status.Errorf(codes.NotFound, "user not exist")
}

func (us *userService) CreateUser(ctx context.Context, request *proto.CreateUserInfo) (*proto.UserInfoResponse, error) {
	username := strings.TrimSpace(request.Username)
	email := strings.TrimSpace(request.Email)

	// 验证用户名
	if err := validate.IsUsername(username); err != nil {
		return nil, status.Error(codes.InvalidArgument, "用户名不合法")
	}
	if _, err := us.GetUserByUsername(username); err == nil {
		return nil, status.Error(codes.AlreadyExists, "用户名已存在")
	}

	// 验证邮箱
	if len(email) > 0 {
		if err := validate.IsEmail(email); err != nil {
			return nil, err
		}
		if _, err := us.GetUserByEmail(context.Background(), &proto.EmailRequest{
			Email: email,
		}); err == nil {
			return nil, status.Error(codes.AlreadyExists, "邮箱：" + email + " 已被占用")
		}
	} else {
		return nil, status.Error(codes.InvalidArgument, "邮箱不能为空")
	}

	//// 验证密码
	//if err := validate.IsPassword(request.Password, request.RePassword); err != nil {
	//	return nil, err
	//}
	//// 密码加密
	//h := md5.New()
	//h.Write([]byte(request.Password))
	//encryptPassword := hex.EncodeToString(h.Sum(nil))

	user := &models.User{
		Username:        sql.NullString{String:username, Valid: len(username) > 0},
		Email:           sql.NullString{String:email, Valid: len(email) > 0},
		EmailVerified:   true,
		Avatar:          "",
		BackgroundImage: "",
		Password:        "",
		Description:     "这个人很懒～什么都没有留下",
		Score:           0,
		Status:          constants.StatusOk,
		ArticleCount:    0,
		CommentCount:    0,
		Roles:           constants.RoleUser,
		Type:            constants.UserTypeNormal,
		CreateTime:      uint64(time.Now().UnixNano()),
		UpdateTime:      uint64(time.Now().UnixNano()),
	}
	err := repositories.UserRepository.Create(db.GetMysqlInstance().GetMysqlDB(), user)
	if err != nil {
		return nil, status.Error(codes.Internal, "创建用户失败")
	}

	return replice.UserInfoResponse(user), nil
}

func (us *userService) UpdateUser(ctx context.Context, request *proto.UpdateUserInfo) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (us *userService) CheckPassword(ctx context.Context, request *proto.PasswordCheckInfo) (*proto.CheckPasswordResponse, error) {
	res := false
	h := md5.New()
	h.Write([]byte(request.Password))
	if request.EncryptedPassword == hex.EncodeToString(h.Sum(nil)) {
		res = true
	}

	return &proto.CheckPasswordResponse{
		Success: res,
	}, nil
}
