package core

import (
	"fmt"
	"testing"
)

func Test_LoadScheduleConfig(t *testing.T) {
	conf, err := LoadScheduleConfig("../conf/config.yaml")
	if err != nil {
		panic(err)
	}
	fmt.Println(conf)
}
