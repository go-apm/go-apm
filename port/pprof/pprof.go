package pprof

import "golang.org/x/net/context"

type PProfType int

const (
	Heap PProfType = iota + 1
	CPU
)

var PProfetch PProfFetcher

type PProfFetcher interface {
	Fetch(c context.Context, typ PProfType, host string, file string) error
}
