package model

import (
	"encoding/binary"
	"encoding/json"
	"github.com/boltdb/bolt"
	"github.com/go-apm/go-apm/port/pprof"
	"github.com/go-apm/go-apm/port/ssh"
	"github.com/go-apm/go-apm/port/store"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
	"path/filepath"
	"strconv"
	"time"
)

type PProf struct {
	ID         uint64
	Typ        pprof.PProfType
	Host       string
	Port       string
	Git        string
	Code       string
	RemoteBin  string
	AppName    string
	DumpPath   string
	LocalBin   string
	Status     PProfStatus
	CreateTime time.Time
	UpdateTime time.Time
}

func (h *PProf) Create(c context.Context) error {
	h.Status = Created
	h.CreateTime = time.Now()
	h.AppName = h.genAppName(c)
	h.DumpPath = h.genDumpPath(c)
	h.LocalBin = h.genLocalBin(c)
	return h.save(c)
}

func (h *PProf) genAppName(c context.Context) string {
	return filepath.Base(h.RemoteBin)
}

func (h *PProf) genDumpPath(c context.Context) string {
	return filepath.Join(viper.GetString("dump.dumpFolder"),
		strconv.Itoa(int(h.Typ))+"_"+h.AppName+"_"+h.Host+"_"+h.Port+"_"+h.CreateTime.Format("20060102150405")+".out")
}

func (h *PProf) genLocalBin(c context.Context) string {
	return filepath.Join(viper.GetString("bin.binFolder"),
		strconv.Itoa(int(h.Typ))+"_"+h.AppName+"_"+h.Host+"_"+h.Port+"_"+h.CreateTime.Format("20060102150405"))
}

func (h *PProf) FetchPProf(c context.Context) error {
	err := pprof.PProfetch.Fetch(c, h.Typ, h.Host+":"+h.Port, h.DumpPath)
	if err != nil {
		return err
	}
	h.Status = h.Status | PProfFetched
	h.UpdateTime = time.Now()
	return h.save(c)
}

func (h *PProf) FetchBinary(c context.Context) error {
	err := ssh.FetchBinary(c, h.Host, h.RemoteBin, h.LocalBin)
	if err != nil {
		return err
	}
	h.Status = h.Status | BinaryFetched
	h.UpdateTime = time.Now()
	return h.save(c)
}

func (h *PProf) AllocObjects(c context.Context) error {

	return nil
}

func (h *PProf) save(c context.Context) error {
	return store.DefaultDB.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("pprofs"))
		if err != nil {
			return err
		}
		b := tx.Bucket([]byte("pprofs"))
		if h.ID == 0 {
			h.ID, _ = b.NextSequence()
		}
		buf, err := json.Marshal(h)
		if err != nil {
			return err
		}
		return b.Put(itob(h.ID), buf)
	})
}

func GetPProf(c context.Context, id uint64) (*PProf, error) {
	var p PProf
	err := store.DefaultDB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("pprofs"))
		buf := b.Get(itob(id))
		err := json.Unmarshal(buf, &p)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func itob(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, v)
	return b
}
