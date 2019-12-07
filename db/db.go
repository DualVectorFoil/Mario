package db

import "sync"

type DB struct {
	Lock sync.Locker
}
