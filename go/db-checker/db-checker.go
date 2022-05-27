package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
)

var (
	user     string
	password string
	protocol = "tcp"
	address  string
	port     = "3306"
)

func dbConnect(ctx context.Context, user, password, address, protocol, port string) error {

	var (
		db  *sql.DB
		err error
	)

	// Specify connection properties.
	cfg := mysql.Config{
		User:                 user,
		Passwd:               password,
		Net:                  protocol,
		Addr:                 address + ":" + port,
		AllowNativePasswords: true,
	}

	// Get a database handle.
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return err
	}
	defer db.Close()

	// Get a database handle.
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return err
	}

	c, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	// Get a database handle.
	if err := db.PingContext(c); err != nil {
		// status = "Unable to connect to db: " + dbAddr
		return err
	}

	return nil
}

func init() {
	// Setting all the parameters and its values
	flag.StringVar(&user, "user", "", "The database user name")
	flag.StringVar(&password, "password", "", "The database password")
	flag.StringVar(&protocol, "protocol", protocol, "The Net connection protocol, defaults TCP")
	flag.StringVar(&address, "host", "", "The database host name or IP address to connect")
	flag.StringVar(&port, "port", port, "The database port to connect")
	flag.Parse()

	switch {
	case user == "":
		fmt.Fprint(os.Stderr, "Please, enter the user name\n")
		os.Exit(1)
	case password == "":
		fmt.Fprint(os.Stderr, "Please, enter the password\n")
		os.Exit(1)
	case address == "":
		fmt.Fprint(os.Stderr, "Please, enter the host\n")
		os.Exit(1)
	}
}

func main() {

	// Background context
	ctx := context.Background()

	// Trying to connect to db
	fmt.Println("Connecting to db: " + address + "...")
	err := dbConnect(ctx, user, password, address, protocol, port)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error %s\n", err)
		os.Exit(1)
	}

	fmt.Println("Connected!")

}
