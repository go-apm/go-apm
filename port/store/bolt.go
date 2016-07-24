package store

import (
	"github.com/boltdb/bolt"
	_ "github.com/go-apm/go-apm/conf"
	"github.com/spf13/viper"
	"github.com/uber-go/zap"
	"os"
	"time"
)

var dbLogger = zap.NewJSON()

var DefaultDB *bolt.DB

func init() {
	db, err := bolt.Open(viper.GetString("db.dbFile"),
		os.ModePerm, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		dbLogger.Fatal("Init DB failure", zap.Error(err))
	}
	DefaultDB = db
}
