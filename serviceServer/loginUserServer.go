package serviceServer

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/DualVectorFoil/Zelda/conf"
	"github.com/DualVectorFoil/Zelda/db"
	pb "github.com/DualVectorFoil/Zelda/protobuf"
	"github.com/DualVectorFoil/Zelda/utils/jwt"
	"github.com/DualVectorFoil/Zelda/utils/ptr"
	"github.com/sirupsen/logrus"
	"sync"
)

type LoginUserServer struct {
	DBInstance *db.DB
}

var LUS *LoginUserServer
var LUSOnce sync.Once

func NewLoginUserServer(instance *db.DB) *LoginUserServer {
	LUSOnce.Do(func() {
		LUS = &LoginUserServer{DBInstance: instance}
	})
	return LUS
}

func (l *LoginUserServer) LoginUserWithPWD(ctx context.Context, loginInfo *pb.LoginInfo) (*pb.LoginResp, error) {
	userNameInfo := loginInfo.GetUserNameInfo()
	pwd := loginInfo.GetPwd()
	if userNameInfo == "" || pwd == "" {
		logrus.WithFields(logrus.Fields{
			"userNameInfo": userNameInfo,
			"pwd":          pwd,
		}).Error("Uncorrected login info.")
		return &pb.LoginResp{
			Status:      ptr.BoolPtr(false),
			ProfileInfo: &pb.ProfileInfo{},
		}, errors.New("Uncorrected login info.")
	}

	profileInfoModel, err := l.DBInstance.LoginWithPWD(userNameInfo, pwd)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"userNameInfo": userNameInfo,
			"pwd":          pwd,
		}).Error(err.Error())
		return &pb.LoginResp{
			Status:      ptr.BoolPtr(false),
			ProfileInfo: &pb.ProfileInfo{},
		}, err
	}

	token, err := jwt.GetToken()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"userNameInfo": userNameInfo,
			"pwd":          pwd,
		}).Error("Get token failed, login failed.")
		return &pb.LoginResp{
			Status:      ptr.BoolPtr(false),
			ProfileInfo: &pb.ProfileInfo{},
		}, err
	}

	profileInfoModelBytes, err := json.Marshal(profileInfoModel)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"userNameInfo": userNameInfo,
			"pwd":          pwd,
		}).Error("Marshal profile info failed, login failed.")
		return &pb.LoginResp{
			Status:      ptr.BoolPtr(false),
			ProfileInfo: &pb.ProfileInfo{},
		}, err
	}

	err = l.DBInstance.SetCacheKV(token, string(profileInfoModelBytes), conf.USER_INFO_TTL)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"userNameInfo": userNameInfo,
			"pwd":          pwd,
			"token":        token,
		}).Error("Save token failed, login failed.")
		return &pb.LoginResp{
			Status:      ptr.BoolPtr(false),
			ProfileInfo: &pb.ProfileInfo{},
		}, err
	}

	logrus.WithFields(logrus.Fields{
		"userNameInfo": userNameInfo,
	}).Info("Login success.")
	return &pb.LoginResp{
		Status: ptr.BoolPtr(true),
		ProfileInfo: &pb.ProfileInfo{
			PhoneNum:     ptr.StringPtr(profileInfoModel.PhoneNum),
			AvatarUrl:    ptr.StringPtr(profileInfoModel.AvatarUrl),
			UserName:     ptr.StringPtr(profileInfoModel.UserName),
			Locale:       ptr.StringPtr(profileInfoModel.Locale),
			Bio:          ptr.StringPtr(profileInfoModel.Bio),
			Followers:    ptr.Int32Ptr(profileInfoModel.Followers),
			Following:    ptr.Int32Ptr(profileInfoModel.Following),
			ArtworkCount: ptr.Int32Ptr(profileInfoModel.ArtworkCount),
			Pwd:          ptr.StringPtr(profileInfoModel.PWD),
			RegisteredAt: ptr.Int64Ptr(profileInfoModel.RegisteredAt),
			LastLoginAt:  ptr.Int64Ptr(profileInfoModel.LastLoginAt),
			Token:        ptr.StringPtr(token),
		},
	}, nil
}

func (l *LoginUserServer) LoginUserWithToken(ctx context.Context, loginInfo *pb.LoginInfo) (*pb.LoginResp, error) {
	token := loginInfo.GetToken()
	if token == "" {
		logrus.WithFields(logrus.Fields{
			"token": token,
		}).Error("Uncorrected login token info, login failed.")
		return &pb.LoginResp{
			Status:      ptr.BoolPtr(false),
			ProfileInfo: &pb.ProfileInfo{},
		}, errors.New("Uncorrected login token info, login failed.")
	}

	profileInfoModel, err := l.DBInstance.LoginWithToken(token)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"token": token,
		}).Error(err.Error())
		return &pb.LoginResp{
			Status:      ptr.BoolPtr(false),
			ProfileInfo: &pb.ProfileInfo{},
		}, err
	}

	logrus.WithFields(logrus.Fields{
		"token": token,
	}).Info("Login success.")
	return &pb.LoginResp{
		Status: ptr.BoolPtr(true),
		ProfileInfo: &pb.ProfileInfo{
			PhoneNum:     ptr.StringPtr(profileInfoModel.PhoneNum),
			AvatarUrl:    ptr.StringPtr(profileInfoModel.AvatarUrl),
			UserName:     ptr.StringPtr(profileInfoModel.UserName),
			Locale:       ptr.StringPtr(profileInfoModel.Locale),
			Bio:          ptr.StringPtr(profileInfoModel.Bio),
			Followers:    ptr.Int32Ptr(profileInfoModel.Followers),
			Following:    ptr.Int32Ptr(profileInfoModel.Following),
			ArtworkCount: ptr.Int32Ptr(profileInfoModel.ArtworkCount),
			Pwd:          ptr.StringPtr(profileInfoModel.PWD),
			RegisteredAt: ptr.Int64Ptr(profileInfoModel.RegisteredAt),
			LastLoginAt:  ptr.Int64Ptr(profileInfoModel.LastLoginAt),
			Token:        ptr.StringPtr(token),
		},
	}, nil
}
