package utils

import (
	"log"
	"time"
)

func Retry(attempts int, sleep time.Duration, returnErr error) (err error) {
	for i := 0; i < attempts; i++ {
		if i > 0 {
			log.Println("retrying after error ===>", err)
			time.Sleep(sleep)
		}
		err = returnErr
		if err == nil {
			return nil
		}
	}
	return err
}
