package gl

import (
	"bufio"
	"bytes"
	"context"
	"encoding/binary"
	"glock/src/model"
	"glock/src/treaty"
	"log"
	"net"
	"strconv"
	"time"
)

func tcpHandler(conn *net.Conn) {
	reader := bufio.NewReader(*conn)
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	for {
		d, err := treaty.Decode(reader)
		if err != nil {
			log.Println("accept failed, err:%V", err)
			break
		}
		ttx, _ := context.WithTimeout(ctx, time.Second*5)
		go processMessage(ttx, conn, d)
	}
}

func processMessage(ctx context.Context, conn *net.Conn, d []byte) {
	// Check data integrity
	if len(d) == 12 {
		// get lock structure
		// gLock := model.GlobalLockSet[d[0]]

		if d[1] == 0 {
			var id uint64
			var rd uint16
			_ = binary.Read(bytes.NewReader(d[4:]), binary.LittleEndian, &id)
			_ = binary.Read(bytes.NewReader(d[2:4]), binary.LittleEndian, &rd)

			v := strconv.Itoa(int(id)) + strconv.Itoa(int(rd))

			val, ok := model.GlobalLockMap.Load(d[0])
			if ok {
				s, _ := val.(string)
				if s == v {
					d[1] = 1
					(*conn).Write(d)
					return
				} else {
					<-model.GlobalLockChan
					model.GlobalLockMap.Store(d[0], v)
					d[1] = 1
					(*conn).Write(d)
					return
				}
			} else {
				model.GlobalLockMap.Store(d[0], v)
				d[1] = 1
				(*conn).Write(d)
				return
			}
		} else {
			model.GlobalLockChan <- true
			return
		}
	} else {
		log.Println("wrong message body")
		(*conn).Write(treaty.Encode(d[0], 0))
	}
}
