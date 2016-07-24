package http

import (
	"errors"
	"fmt"
	"github.com/go-apm/go-apm/port/pprof"
	"github.com/go-apm/go-apm/util/xhttp"
	"github.com/uber-go/zap"
	"golang.org/x/net/context"
	"io"
	"net/http"
	"os"
)

type httpPProf struct {
	hc *http.Client
}

func (h *httpPProf) Fetch(c context.Context, typ pprof.PProfType, host string, file string) error {
	logger := xhttp.CurrentLogger(c)
	url := fmt.Sprintf("http://%s/debug/pprof/heap?debug=1", host)
	logger.Info("Start request heap pprof", zap.String("url", url))
	resp, err := h.hc.Get(url)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		logger.Error("Fetch heap pprof failure", zap.Int("code", resp.StatusCode))
		return errors.New("Fetch heap ret " + resp.Status)
	}
	_, err = os.Stat(file)
	if os.IsNotExist(err) {
		_, err = os.Create(file)
		if err != nil {
			logger.Error("Create heap file failure", zap.Error(err))
			return err
		}
	}
	if err != nil {
		logger.Error("Open heap file failure", zap.Error(err))
		return err
	}
	dumpFile, err := os.OpenFile(file, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		logger.Error("Open heap file failure", zap.Error(err))
		return err
	}
	size, err := io.Copy(dumpFile, resp.Body)
	if err != nil {
		logger.Error("Fetch save heap failure", zap.Error(err))
		return err
	}
	logger.Info("Finsh fetch heap pprof", zap.String("file", dumpFile.Name()), zap.Int64("size", size))
	return nil
}

func init() {
	pprof.PProfetch = &httpPProf{hc: http.DefaultClient}
}
