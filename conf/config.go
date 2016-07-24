package conf

import (
	"flag"
	"github.com/coreos/etcd/pkg/fileutil"
	"github.com/uber-go/zap"
	"github.com/spf13/viper"
)

var confLogger = zap.NewJSON()

func init() {
	confPath := flag.String("conf", "conf/apm.toml", "conf file")
	flag.Parse()
	if !fileutil.Exist(*confPath) {
		confLogger.Fatal("Config file not exists", zap.String("currentPath", *confPath))
	}
	viper.SetConfigFile(*confPath)
	err := viper.ReadInConfig()
	if err != nil {
		confLogger.Fatal("Read Config failure", zap.Error(err))
	}
}