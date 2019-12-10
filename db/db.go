package db

import (
	"errors"
	"github.com/DualVectorFoil/Zelda/conf"
	"github.com/DualVectorFoil/Zelda/model"
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

type DB struct {
	Lock  sync.Locker
	Mysql *gorm.DB
	Redis *redis.Client
}

var dbInstance *DB
var dbInstanceOnce sync.Once

func InitDB(mySqlInfo string) {
	dbInstanceOnce.Do(func() {
		mysqlDB, err := gorm.Open("mysql", mySqlInfo)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"mySqlInfo": mySqlInfo,
				"err":       err.Error(),
			}).Fatal("Mysql init failed.")
		}
		redisCli := redis.NewClient(&redis.Options{
			Addr:        conf.REDIS_ADDR,
			Password:    conf.REDIS_PASSWORD,
			DB:          conf.REDIS_DB_NUM,
			DialTimeout: conf.REDIS_TIMEOUT,
		})
		dbInstance = &DB{
			Lock:  &sync.Mutex{},
			Mysql: mysqlDB,
			Redis: redisCli,
		}
	})
}

func GetDB() *DB {
	return dbInstance
}

func (instance *DB) Close() {
	instance.Lock.Lock()
	defer instance.Lock.Unlock()
	instance.Mysql.Close()
	instance.Redis.Close()
}

func (instance *DB) GetCacheKV(key string) (string, error) {
	instance.Lock.Lock()
	defer instance.Lock.Unlock()
	value, err := instance.Redis.Get(key).Result()
	return value, err
}

func (instance *DB) SetCacheKV(key string, value interface{}, expiration time.Duration) error {
	instance.Lock.Lock()
	defer instance.Lock.Unlock()
	err := instance.Redis.Set(key, value, expiration).Err()
	return err
}

// TODO decouple bussiness and basice db operation
func (instance *DB) IsVerifyCodeAvailable(phoneNum string, verifyCode string) bool {
	if phoneNum == "" || verifyCode == "" {
		return false
	}

	storageVerifyCode, err := instance.GetCacheKV(phoneNum)
	if err != nil {
		return false
	}

	if storageVerifyCode != verifyCode {
		return false
	}

	return true
}

func (instance *DB) SaveRegisterUserInfo(phoneNum string, userName string, pwdEncoded string, verifyCode string) error {
	if phoneNum == "" || userName == "" || pwdEncoded == "" || verifyCode == "" {
		return errors.New("Register failed, uncorrected register info.")
	}

	if !instance.IsVerifyCodeAvailable(phoneNum, verifyCode) {
		return errors.New("Register failed, uncorrected verify code.")
	}

	instance.Lock.Lock()
	defer instance.Lock.Unlock()
	rows, err := instance.Mysql.Table(conf.PROFILE_TABLE_NAME).Select([]string{"phone_num", "user_name"}).Rows()
	if err != nil {
		return errors.New("Register failed, server error, err: " + err.Error())
	}

	for rows.Next() {
		var phoneNumTmp string
		var userNameTmp string
		if err := rows.Scan(&phoneNumTmp, &userNameTmp); err != nil {
			return errors.New("Register failed, server error, err: " + err.Error())
		}

		if phoneNumTmp == phoneNum {
			return errors.New("phoneNum has registered.")
		} else if userNameTmp == userName {
			return errors.New("userName has registered.")
		}
	}

	profileInfo := &model.ProfileInfo{
		PhoneNum:     phoneNum,
		UserName:     userName,
		PWDEncoded:   pwdEncoded,
		RegisteredAt: time.Now().Unix(),
		LastLoginAt:  time.Now().Unix(),
	}

	errs := instance.Mysql.Create(profileInfo).GetErrors()
	if len(errs) > 0 {
		return errors.New("Register failed.")
	}

	return nil
}
