package model

import "sync"

var GlobalLockSet = [256]*GLock{}

var GlobalLockMap sync.Map

var GlobalLockChan = make(chan interface{})
