package serviceServer

import (
	"context"
	"errors"
	"fmt"
	"github.com/DualVectorFoil/Zelda/db"
	pb "github.com/DualVectorFoil/Zelda/protobuf"
	ptr2 "github.com/DualVectorFoil/Zelda/utils/ptr"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

const (
	VERIFYCODE_EXPIRATION = time.Minute * 5
)

type VerifyCodeServer struct {
	DBInstance *db.DB
}

var vcs *VerifyCodeServer
var vcsOnce sync.Once

func NewVerifyCodeServer(dbInstance *db.DB) *VerifyCodeServer {
	vcsOnce.Do(func() {
		vcs = &VerifyCodeServer{
			DBInstance: dbInstance,
		}
	})
	return vcs
}

func (s *VerifyCodeServer) SetVerifyCodeInfo(ctx context.Context, verifyCodeInfo *pb.VerifyCodeInfo) (*pb.VerifyCodeRespStatus, error) {
	phoneNum := verifyCodeInfo.GetPhoneNum()
	verifyCode := verifyCodeInfo.GetVerifyCode()
	if phoneNum == "" || verifyCode == "" {
		logrus.WithFields(logrus.Fields{
			"phoneNum":   phoneNum,
			"verifyCode": verifyCode,
		}).Error("SetVerifyCodeInfo failed, uncorrect verifyCodeInfo.")
		return nil, errors.New("SetVerifyCodeInfo failed, uncorrect verifyCodeInfo.")
	}

	err := s.DBInstance.SetCacheKV(phoneNum, verifyCode, VERIFYCODE_EXPIRATION)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"phoneNum":   phoneNum,
			"verifyCode": verifyCode,
		}).Error("SetVerifyCodeInfo failed, err: " + err.Error())
		return &pb.VerifyCodeRespStatus{
			Status: ptr2.BoolPtr(false),
		}, err
	}
	fmt.Println("verifyCode setted: " + verifyCode)
	return &pb.VerifyCodeRespStatus{
		Status: ptr2.BoolPtr(true),
	}, nil
}

func (s *VerifyCodeServer) IsVerifyCodeAvailable(ctx context.Context, verifyCodeInfo *pb.VerifyCodeInfo) (*pb.VerifyCodeRespStatus, error) {
	phoneNum := verifyCodeInfo.GetPhoneNum()
	verifyCode := verifyCodeInfo.GetVerifyCode()

	if s.DBInstance.IsVerifyCodeAvailable(phoneNum, verifyCode) {
		return &pb.VerifyCodeRespStatus{
			Status: ptr2.BoolPtr(true),
		}, nil
	} else {
		return nil, errors.New("SetVerifyCodeInfo failed.")
	}
}
