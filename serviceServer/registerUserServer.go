package serviceServer

import (
	"context"
	"github.com/DualVectorFoil/Zelda/db"
	pb "github.com/DualVectorFoil/Zelda/protobuf"
	ptr "github.com/DualVectorFoil/Zelda/utils"
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
	pwdEncoded := registerInfo.GetPwd()
	verifyCode := registerInfo.GetVerifyCode()

	err := s.DBInstance.SaveRegisterUserInfo(phoneNum, userName, pwdEncoded, verifyCode)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"phoneNum":   phoneNum,
			"userName":   userName,
			"pwdEncoded": pwdEncoded,
			"verifyCode": verifyCode,
		}).Error(err.Error())
		return &pb.RegisterRespStatus{
			Status: ptr.BoolPtr(false),
			Err:    ptr.StringPtr(err.Error()),
		}, err
	}

	return &pb.RegisterRespStatus{
		Status: ptr.BoolPtr(true),
		Err:    ptr.StringPtr(""),
	}, nil
}
