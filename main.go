package main

import (
	"fmt"
	"github.com/DualVectorFoil/Zelda/conf"
	"github.com/DualVectorFoil/Zelda/db"
	"github.com/DualVectorFoil/Zelda/manager"
	"github.com/DualVectorFoil/Zelda/model"
	pb "github.com/DualVectorFoil/Zelda/protobuf"
	"github.com/DualVectorFoil/Zelda/serviceServer"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func main() {
	initDB()
	initGRPCService()
	defer db.GetDB().Close()
}

func initDB() {
	mysqlInfo := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true", conf.MYSQL_USERNAME, conf.MYSQL_PASSWORD, conf.MYSQL_IP, conf.MYSQL_PORT, conf.MYSQL_DBNAME)
	db.InitDB(mysqlInfo)
	dbInstance := db.GetDB()
	dbInstance.Mysql.AutoMigrate(&model.ProfileInfo{})
}

func initGRPCService() {
	grpcServer := grpc.NewServer()
	defer grpcServer.GracefulStop()

	pb.RegisterLoginUserServiceServer(grpcServer, serviceServer.NewLoginUserServer(db.GetDB()))
	pb.RegisterRegisterUserServiceServer(grpcServer, serviceServer.NewRegisterUserServer(db.GetDB()))
	pb.RegisterVerifyCodeServiceServer(grpcServer, serviceServer.NewVerifyCodeServer(db.GetDB()))
	err := manager.GetServiceManger(conf.ETCD_ADDRESS).Register(conf.SERVICE_NAME, conf.LISTEN_IP, conf.SERVICE_IP, conf.SERVICE_PORT, grpcServer, conf.GRPC_TTL)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"service_name": conf.SERVICE_NAME,
			"listen_ip":    conf.LISTEN_IP,
			"service_port": conf.SERVICE_PORT,
			"grpc_ttl":     conf.GRPC_TTL,
		}).Error(fmt.Sprintf("Register service to etcd failed, err: %s.", err.Error()))
		panic(err)
	}
}
