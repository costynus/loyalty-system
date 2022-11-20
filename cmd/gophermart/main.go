package main

import (
	"fmt"

	"github.com/Konab/go-diplom-1/loyalty-system/config"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(cfg)
}
