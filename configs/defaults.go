package configs

import (
	"time"
)

type MessengerConfig struct {
	Url       string
	Durable   bool
	Attempts  int
	Delay     time.Duration
	Threshold int
}

//func GetConfig(url string) MessengerConfig {
//	return MessengerConfig{
//		"localhost",
//		true,
//		5,
//		time.Second * 2,
//		3,
//	}
//}
