package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	if len(os.Args) != 2 {
		printUsage()
		os.Exit(1)
	}

	cidrInput := os.Args[1]
	cidrs := strings.Split(cidrInput, ",")

	if len(cidrs) == 0 {
		printUsage()
		os.Exit(1)
	}

	dbPath := "assets/cnc.db"
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	for _, cidr := range cidrs {
		err := insertCIDR(db, cidr)
		if err != nil {
			log.Printf("Error inserting CIDR %s: %v\n", cidr, err)
		} else {
			fmt.Printf("CIDR %s has been added to the whitelist.\n", cidr)
		}
	}
}

func insertCIDR(db *sql.DB, cidr string) error {
	_, err := db.Exec("INSERT INTO whitelist (prefix, netmask) VALUES (?, ?)", cidr, 24)
	return err
}

func printUsage() {
	fmt.Println("Usage: ./whitelist <CIDR1,CIDR2,CIDR3,...>")
	fmt.Println("Example: ./whitelist 1.1.1.0/24,8.8.8.0/24")
	fmt.Println("Specific IP: ./whitelist 1.3.3.7/32")
}
