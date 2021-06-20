/*
Johan Calderon
Simulator send router google api service Direction to UDP Scoket
go run main.go <IDUnit> <Origin> <Destination> <Interval to send message> <Server> <Port>
*/
package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/nahhoj/HanaCarTrack/def"
	"github.com/nahhoj/HanaCarTrack/utils"
	"github.com/subosito/gotenv"
)

func main() {
	gotenv.Load()
	url := os.Getenv("google_api")
	key := os.Getenv("google_key")
	if len(os.Args) < 4 {
		fmt.Println("There are no arguments")
		return
	}
	idunit := os.Args[1]
	origin := os.Args[2]
	destination := os.Args[3]
	delay, _ := strconv.ParseInt(os.Args[4], 10, 16)
	host := os.Args[5]
	port := os.Args[6]
	url += "?origin=" + origin + "&destination=" + destination + "&key=" + key
	res := utils.CallService(url, "GET", nil, "")
	if res.StatusCode != 200 {
		fmt.Println(string(res.Response))
		return
	}
	var data def.Directions
	json.Unmarshal(res.Response, &data)
	ip := "n/a"
	if resIP := utils.CallService("https://api.ipify.org/", "GET", nil, ""); resIP.StatusCode == 200 {
		ip = string(resIP.Response)
	}
	server, err := net.ResolveUDPAddr("udp4", host+":"+port)
	if err != nil {
		fmt.Println(err)
		return
	}
	connection, err := net.DialUDP("udp4", nil, server)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer connection.Close()
	for _, steps := range data.Routes[0].Legs[0].Steps {
		send := fmt.Sprintln(idunit, ip, steps.StartLocation.Lat, steps.StartLocation.Lng, nil, rand.Intn(100))
		_, err = connection.Write([]byte(send))
		if err != nil {
			fmt.Println(err)
			return
		}
		time.Sleep(time.Duration(delay) * time.Second)
	}
}
