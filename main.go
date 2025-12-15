package main

import (
	"fmt"
	"github.com/CoupDeGrace92/gator/internal/config"
)

func main() {
	conf := config.Read()
	conf.SetUser("nickoboy1992")
	conf = config.Read()
	fmt.Println(conf)
}