package web

import "time"

func Worker() {
	for conn := range ConnChan {
		conn.Do()
	}
}

type TempWorker struct {
	maxTempWorkerNum uint
	maxIdleTime time.Duration
}

func (w *TempWorker) Work(conn *Connection) {
	//start the temp worker
	conn.Do()
	for {
		select {
		case conn = <- ConnChan:
			conn.Do()
		case <- time.After(tempWorker.maxIdleTime):
			return
		}
	}
}
