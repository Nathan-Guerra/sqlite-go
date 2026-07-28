package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	// Available if you need it!
	// "github.com/xwb1989/sqlparser"
)

// Usage: your_program.sh sample.db .dbinfo
func main() {
	databaseFilePath := os.Args[1]
	command := os.Args[2]

	switch command {
	case ".dbinfo":
		databaseFile, err := os.Open(databaseFilePath)
		if err != nil {
			log.Fatal(err)
		}

		header := make([]byte, 100)

		_, err = databaseFile.Read(header)
		if err != nil {
			log.Fatal(err)
		}

		var pageSize uint16
		if err := binary.Read(bytes.NewReader(header[16:18]), binary.BigEndian, &pageSize); err != nil {
			fmt.Println("Failed to read integer:", err)
			return
		}
		// You can use print statements as follows for debugging, they'll be visible when running tests.
		fmt.Fprintln(os.Stderr, "Logs from your program will appear here!")

		/**
				The one-byte flag at offset 0 indicating the b-tree page type.
		A value of 2 (0x02) means the page is an interior index b-tree page.
		A value of 5 (0x05) means the page is an interior table b-tree page.
		A value of 10 (0x0a) means the page is a leaf index b-tree page.
		A value of 13 (0x0d) means the page is a leaf table b-tree page.
		Any other value for the b-tree page type is an error.
		*/
		pageHeader := make([]byte, 12)
		_, err = databaseFile.ReadAt(pageHeader, 100)
		if err != nil {
			panic(err)
		}

		// var pageType string
		// switch pageHeader[0] {
		// case 0x02:
		// 	pageType = "interior index b-tree"
		// case 0x05:
		// 	pageType = "interior table b-tree"
		// case 0x0a:
		// 	pageType = "leaf index b-tree"
		// case 0x0d:
		// 	pageType = "leaf table b-tree"
		// default:
		// 	pageType = "error"
		// }

		var startOfFirstFreeblock uint16
		if err := binary.Read(bytes.NewReader(pageHeader[1:3]), binary.BigEndian, &startOfFirstFreeblock); err != nil {
			panic(err)
		}

		var numCells uint16
		if err := binary.Read(bytes.NewReader(pageHeader[3:5]), binary.BigEndian, &numCells); err != nil {
			panic(err)
		}

		var startOfCellContentArea uint16
		if err := binary.Read(bytes.NewReader(pageHeader[5:7]), binary.BigEndian, &startOfCellContentArea); err != nil {
			panic(err)
		}

		var fragmentedFreeBytes uint8
		if err := binary.Read(bytes.NewReader(pageHeader[7:8]), binary.BigEndian, &fragmentedFreeBytes); err != nil {
			panic(err)
		}

		fmt.Printf("database page size: %v\n", pageSize)
		fmt.Printf("number of tables: %v\n", numCells)
		// fmt.Printf("typeof page: %v\n", pageType)
		// fmt.Printf("start of first freeblock on the page: 0x%x\n", startOfFirstFreeblock)
		// fmt.Printf("number of cells on the page: %v\n", numCells)
		// fmt.Printf("start of cell content area: %v\n", startOfCellContentArea)
		// fmt.Printf("fragmented free bytes: %v\n", fragmentedFreeBytes)
	default:
		fmt.Println("Unknown command", command)
		os.Exit(1)
	}
}
