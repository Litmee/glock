package gl

import (
	"bufio"
	"context"
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

func processMessage(ctx context.Context, conn *net.Conn, b []byte) {
	// Check data integrity
	if len(b) == 8 {
		// get lock structure
		gLock := model.GlobalLockSet[b[0]]
		if gLock == nil {
			(*conn).Write(treaty.Encode(b[0], 1))
			return
		} else {
			// If the data body is 0, it means that the requester wants to acquire the lock
			if b[1] == 0 {
				select {
				case <-ctx.Done():
					log.Println("Failed to acquire lock, err:%V", ctx.Err())
					(*conn).Write(treaty.Encode(b[0], 0))
				case <-gLock.C:
					b[1] = 1
					(*conn).Write(treaty.Encode(b[0], 1))
				}
			}
		}
	} else {
		log.Println("wrong message body")
		(*conn).Write(treaty.Encode(b[0], 0))
	}
}
