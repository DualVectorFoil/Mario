package serviceServer

import (
	"context"
	"github.com/DualVectorFoil/Zelda/db"
	pb "github.com/DualVectorFoil/Zelda/protobuf"
	ptr2 "github.com/DualVectorFoil/Zelda/utils/ptr"
	"github.com/sirupsen/logrus"
	"sync"
)

type RegisterUserServer struct {
	DBInstance *db.DB
}

var RUS *RegisterUserServer
var RUSOnce sync.Once

func NewRegisterUserServer(dbInstance *db.DB) *RegisterUserServer {
	RUSOnce.Do(func() {
		RUS = &RegisterUserServer{
			DBInstance: dbInstance,
		}
	})
	return RUS
}

func (s *RegisterUserServer) RegisterUser(ctx context.Context, registerInfo *pb.RegisterInfo) (*pb.RegisterRespStatus, error) {
	phoneNum := registerInfo.GetPhoneNum()
	userName := registerInfo.GetUserName()
	pwd := registerInfo.GetPwd()
	verifyCode := registerInfo.GetVerifyCode()

	err := s.DBInstance.SaveRegisterUserInfo(phoneNum, userName, pwd, verifyCode)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"phoneNum":   phoneNum,
			"userName":   userName,
			"pwd":        pwd,
			"verifyCode": verifyCode,
		}).Error(err.Error())
		return &pb.RegisterRespStatus{
			Status: ptr2.BoolPtr(false),
			Err:    ptr2.StringPtr(err.Error()),
		}, err
	}

	logrus.WithFields(logrus.Fields{
		"phoneNum":   phoneNum,
		"userName":   userName,
		"verifyCode": verifyCode,
	}).Info("Register success.")
	return &pb.RegisterRespStatus{
		Status: ptr2.BoolPtr(true),
		Err:    ptr2.StringPtr(""),
	}, nil
}
