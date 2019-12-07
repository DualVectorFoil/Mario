package serviceServer

import (
	"context"
	"errors"
	"github.com/DualVectorFoil/Zelda/db"
	pb "github.com/DualVectorFoil/Zelda/protobuf"
	ptr "github.com/DualVectorFoil/Zelda/utils"
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
		return nil, errors.New("SetVerifyCodeInfo failed, uncorrect verifyCodeInfo.")
	}

	err := s.DBInstance.RedisSetKV(phoneNum, verifyCode, VERIFYCODE_EXPIRATION)
	if err != nil {
		return &pb.VerifyCodeRespStatus{
			Status: ptr.BoolPtr(false),
		}, err
	}

	return &pb.VerifyCodeRespStatus{
		Status: ptr.BoolPtr(true),
	}, nil
}

func (s *VerifyCodeServer) IsVerifyCodeAvailable(ctx context.Context, verifyCodeInfo *pb.VerifyCodeInfo) (*pb.VerifyCodeRespStatus, error) {
	phoneNum := verifyCodeInfo.GetPhoneNum()
	verifyCode := verifyCodeInfo.GetVerifyCode()
	if phoneNum == "" || verifyCode == "" {
		return nil, errors.New("SetVerifyCodeInfo failed.")
	}

	storageVerifyCode, err := s.DBInstance.RedisGetKV(phoneNum)
	if err != nil {
		return nil, err
	}

	if storageVerifyCode != verifyCode {
		return &pb.VerifyCodeRespStatus{
			Status: ptr.BoolPtr(false),
		}, nil
	}

	return &pb.VerifyCodeRespStatus{
		Status: ptr.BoolPtr(true),
	}, nil
}
