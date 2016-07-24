package model

type PProfStatus int

const (
	Created PProfStatus = 1 << iota
	PProfFetched
	CodeFetched
	BinaryFetched
	Cached
)

