package db

import (
	"encoding/json"
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

func (instance *DB) SaveRegisterUserInfo(phoneNum string, userName string, pwd string, verifyCode string) error {
	if phoneNum == "" || userName == "" || pwd == "" || verifyCode == "" {
		return errors.New("Register failed, uncorrected register info.")
	}

	if !instance.IsVerifyCodeAvailable(phoneNum, verifyCode) {
		return errors.New("Register failed, uncorrected verify code.")
	}

	instance.Lock.Lock()
	defer instance.Lock.Unlock()

	var profileInfo model.ProfileInfo
	instance.Mysql.Table(conf.PROFILE_TABLE_NAME).Where("phone_num = ? or user_name = ?", phoneNum, userName).Find(&profileInfo)
	if profileInfo.PhoneNum != "" || profileInfo.UserName != "" {
		return errors.New("phoneNum or userName has registered.")
	}

	profileInfo = model.ProfileInfo{
		PhoneNum:     phoneNum,
		UserName:     userName,
		PWD:          pwd,
		RegisteredAt: time.Now().Unix(),
		LastLoginAt:  time.Now().Unix(),
	}

	errs := instance.Mysql.Create(&profileInfo).GetErrors()
	if len(errs) > 0 {
		return errors.New("Register failed.")
	}

	return nil
}

func (instance *DB) LoginWithPWD(userNameInfo string, pwd string) (*model.ProfileInfo, error) {
	if userNameInfo == "" || pwd == "" {
		return nil, errors.New("Uncorrected login info, login failed.")
	}

	instance.Lock.Lock()
	defer instance.Lock.Unlock()

	var profileInfo model.ProfileInfo
	err := instance.Mysql.Where("(phone_num = ? OR user_name = ?) AND pwd = ?", userNameInfo, userNameInfo, pwd).Find(&profileInfo).Error
	if err != nil {
		return nil, err
	}
	if profileInfo.PhoneNum == "" || profileInfo.UserName == "" {
		return nil, errors.New("Uncorrected login info, login failed.")
	}

	return &profileInfo, err
}

func (instance *DB) LoginWithToken(token string) (*model.ProfileInfo, error) {
	if token == "" {
		return nil, errors.New("Uncorrected login token info, login failed.")
	}

	instance.Lock.Lock()
	defer instance.Lock.Unlock()

	profileInfoEncoded, err := instance.GetCacheKV(token)
	if err != nil {
		return nil, err
	}

	profileInfo := &model.ProfileInfo{}
	err = json.Unmarshal([]byte(profileInfoEncoded), profileInfoEncoded)
	if err != nil {
		return nil, err
	}
	if profileInfo.PhoneNum == "" || profileInfo.UserName == "" {
		return nil, errors.New("Uncorrected login info, login failed.")
	}

	return profileInfo, nil
}
