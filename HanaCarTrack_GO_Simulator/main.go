package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	gotenv.Load()
	if len(os.Args) < 4 {
		fmt.Println("There are no arguments")
		return
	}
	origin := os.Args[1]
	destination := os.Args[2]
	speed, _ := strconv.ParseInt(os.Args[3], 10, 16)
	fmt.Println(origin, destination, speed)
}
