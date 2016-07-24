package endpoint

import (
	"github.com/go-apm/go-apm/model"
	"github.com/go-apm/go-apm/port/pprof"
	"github.com/go-apm/go-apm/port/validate"
	"github.com/go-apm/go-apm/util/xhttp"
	"github.com/labstack/echo"
	"github.com/uber-go/zap"
	"net/http"
	"strconv"
)

func ListHeaps(c echo.Context) error {
	return c.JSON(http.StatusOK, struct{}{})
}

func NewHeap(c echo.Context) error {
	var req model.NewHeapRequest
	err := c.Bind(&req)
	if err != nil {
		return err
	}
	err = validate.Validator.Struct(&req)
	if err != nil {
		return err
	}
	heap := model.PProf{
		Typ: pprof.Heap, Host: req.Host, Port: req.Port, Git: req.Git, RemoteBin: req.Binary,
	}
	err = heap.Create(c)
	if err != nil {
		return err
	}
	err = heap.FetchPProf(c)
	if err != nil {
		return err
	}

	err = heap.FetchBinary(c)
	if err != nil {
		return err
	}

	xhttp.CurrentLogger(c).Info("Now", zap.String("dump", heap.DumpPath), zap.String("bin", heap.LocalBin))

	return c.JSON(http.StatusOK, struct {
		ID uint64 `json:"heapID"`
	}{ID: heap.ID})
}

func ViewHeap(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return err
	}
	heap, err := model.GetPProf(c, id)
	if err != nil {
		return err
	}
	heap.AllocObjects(c)
	return c.JSON(http.StatusOK, heap)
}
