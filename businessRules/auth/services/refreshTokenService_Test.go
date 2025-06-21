package services

import (
	"fmt"
	"time"
)

type TestChan struct {
	Data string
}

var testChan = make(chan TestChan, 1000)

func produce(data string) {
	testChan <- TestChan{Data: data}
}
func consume() {
	for {
		select {
		case data := <-testChan:
			fmt.Println(data.Data)
		default:
			fmt.Println("No data available")
			time.Sleep(1 * time.Second)
		}
	}
}
func init() {
	go consume()
}
