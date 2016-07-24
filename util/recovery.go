package util

var GoroutinePanicHandler func(r interface{})

func WithRecover(fn func()) {
	defer func() {
		handler := GoroutinePanicHandler
		if handler != nil {
			if r := recover(); r != nil {
				handler(r)
			}
		}
	}()
	fn()
}
