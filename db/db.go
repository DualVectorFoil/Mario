package db

import (
	"github.com/DualVectorFoil/Zelda/conf"
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
	instance.Mysql.Close()
	instance.Redis.Close()
	instance.Lock.Unlock()
}

func (instance *DB) RedisGetKV(key string) (string, error) {
	instance.Lock.Lock()
	value, err := instance.Redis.Get(key).Result()
	instance.Lock.Unlock()
	return value, err
}

func (instance *DB) RedisSetKV(key string, value interface{}, expiration time.Duration) error {
	instance.Lock.Lock()
	err := instance.Redis.Set(key, value, expiration).Err()
	instance.Lock.Unlock()
	return err
}
