/*
Johan Calderon
Listen event for UDP Socket
go run main.go <Port> <database>
*/
package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/subosito/gotenv"

	"github.com/jackc/pgx/v4/pgxpool"
)

func main() {
	gotenv.Load()
	postgresUser := os.Getenv("postgres_user")
	postgresPasswd := os.Getenv("postgres_passwd")
	postgresServer := os.Getenv("postgres_server")
	postgresport := os.Getenv("postgres_port")
	postgresDatabase := os.Getenv("postgres_database")
	urlConnectionPostgres := "postgres://" + postgresUser + ":" + postgresPasswd + "@" + postgresServer + ":" + postgresport + "/" + postgresDatabase
	//urlConnectionHana := "hdb://" + os.Getenv("hana_user") + ":" + os.Getenv("hana_passwd") + "@" + os.Getenv("hana_host") + ":" + os.Getenv("hana_port") + "?DATABASENAME=" + os.Getenv("hana_database")
	pool, err := pgxpool.Connect(context.Background(), urlConnectionPostgres)
	defer pool.Close()
	//rows, err := conn.Query(context.Background(), "SELECT * FROM public.\"Events\";")
	if err != nil {
		panic(err)
	}
	port := ":" + os.Args[1]
	server, err := net.ResolveUDPAddr("udp4", port)
	if err != nil {
		panic(err)
	}
	connection, err := net.ListenUDP("udp4", server)
	if err != nil {
		panic(err)
	}
	defer connection.Close()
	buffer := make([]byte, 1024)
	for {
		fmt.Println("Waiting for connection...")
		n, _, err := connection.ReadFromUDP(buffer)
		if err != nil {
			panic(err)
		}
		message := strings.Split(string(buffer[0:n-1]), ",")
		fmt.Println(string(buffer[0:n]))
		query := "INSERT INTO public.\"Events\"" +
			"(\"Database\", \"IDUnit\", \"IP\", \"Date\", \"Time\", \"Latitude\", \"Longitude\", \"Altitude\", \"Speed\")" +
			"VALUES ('PG', '" + message[0] + "','" + message[1] + "','" + message[2] + "','" + message[3] + "'," + message[4] + "," + message[5] + "," + message[6] + "," + message[7] + ");"
		row, err := pool.Exec(context.Background(), query)
		fmt.Println(row, err)
	}
}

//[SIMU2 181.51.106.73 yyyy-MM-dd hh:mm:ss %!s(float64=35.9010822) %!s(float64=14.5150419)  %!s(int=59]
//INSERT INTO public."Events"("Database", "IDUnit", "IP", "Date", "Time", "Latitude", "Longitude", "Altitude", "Speed")VALUES ('PG', 'SIMU2','181.51.106.73','2021-06-28','18:58:06',35.899638,14.515699,,8);
