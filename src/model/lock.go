package model

import "sync/atomic"

type GLock struct {
	d uint32
	C chan bool
}

func (gl *GLock) GetLock() {
	if atomic.LoadUint32(&gl.d) == 0 {
		atomic.StoreUint32(&gl.d, 1)
	}
}
