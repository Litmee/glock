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
	"time"
)

func tcpHandler(conn *net.Conn) {
	reader := bufio.NewReader(*conn)
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelFunc()
	for {
		d, err := treaty.Decode(reader)
		if err != nil {
			log.Println("accept failed, err:%V", err)
			break
		}
		go processMessage(ctx, conn, d)
	}
}

func processMessage(ctx context.Context, conn *net.Conn, d []byte) {
	// Check data integrity
	if len(d) == 12 {
		// get lock structure
		gLock := model.GlobalLockSet[d[0]]

		var id uint64
		var rd uint16
		_ = binary.Read(bytes.NewReader(d[4:]), binary.LittleEndian, &id)
		_ = binary.Read(bytes.NewReader(d[2:4]), binary.LittleEndian, &rd)

		if gLock == nil {
			d[1] = 1
			(*conn).Write(d)
			return
		} else {
			// If the data body is 0, it means that the requester wants to acquire the lock
			if d[1] == 0 {
				if gLock.Id == id && gLock.Rd == uint32(rd) {
					d[1] = 1
					(*conn).Write(d)
					return
				} else {
					select {
					case <-ctx.Done():
						log.Println("Failed to acquire lock, err:%V", ctx.Err())
						(*conn).Write(d)
					case <-gLock.C:
						gLock.Id = id
						gLock.Rd = uint32(rd)
						d[1] = 1
						(*conn).Write(d)
					}
					return
				}
			} else {
				gLock.C <- true
			}
		}
	} else {
		log.Println("wrong message body")
		(*conn).Write(treaty.Encode(d[0], 0))
	}
}
