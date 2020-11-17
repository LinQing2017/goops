package sys

import (
	"github.com/sirupsen/logrus"
	"time"
)

func WaitAllThreadFinish(threadNum int, tChan chan int, timeOutInSec int) {

	isFinish := make(chan bool)
	isTimeout := make(chan bool)
	go wait(threadNum, tChan, isFinish)
	go timeOut(timeOutInSec, isTimeout)
	select {
	case <-isFinish:
	case <-isTimeout:
		logrus.Warn("操作超时终止")
	}

}

func wait(threadNum int, tChan chan int, isFinish chan bool) {
	for i := 0; i < threadNum; i++ {
		<-tChan
	}
	isFinish <- true
	close(isFinish)
}

func timeOut(timeOutInSec int, isTimeOut chan bool) {
	time.Sleep(time.Duration(timeOutInSec * 1e9))
	isTimeOut <- true
	close(isTimeOut)
}
