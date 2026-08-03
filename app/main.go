package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"strings"
	// Available if you need it!
	// "github.com/xwb1989/sqlparser"
)

// from sqlite documentation
//A variable-length integer or "varint" is a static Huffman encoding
// of 64-bit twos-complement integers that uses less space for small
// positive values. A varint is between 1 and 9 bytes in length.
// The varint consists of either zero or more bytes which have the
// high-order bit set followed by a single byte with the high-order
// bit clear, or nine bytes, whichever is shorter. The lower seven
// bits of each of the first eight bytes and all 8 bits of the ninth
// byte are used to reconstruct the 64-bit twos-complement integer.
// Varints are big-endian: bits taken from the earlier byte of the
// varint are more significant than bits taken from the later bytes.

// source: https://github.com/go-sqlite/sqlite3/blob/53dd8e640ee7dd6005bd7199eed0c470ab43a16e/utils.go#L30-L45
func varint(data []byte) (int64, int) {
	var val uint64
	for i := range 8 {
		if i > len(data)-1 {
			return 0, 0
		}
		val = (val << 7) | uint64(data[i]&0x7f)
		if data[i] < 0x80 {
			return int64(val), i + 1
		}
	}
	if len(data) < 9 {
		return 0, 0
	}
	return int64((val << 8) | uint64(data[8])), 9
}

const DatabaseHeaderSize = 100

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

		dbHeaderBuffer := make([]byte, DatabaseHeaderSize)
		n, err := databaseFile.Read(dbHeaderBuffer)
		if err != nil {
			log.Fatal(err, n)
		}

		dbHeader := DatabaseHeader{}
		dbHeader.Fill(dbHeaderBuffer)

		pageHeaderBuffer := make([]byte, 12)
		n, err = databaseFile.ReadAt(pageHeaderBuffer, DatabaseHeaderSize)
		if err != nil {
			log.Fatal(err, n)
		}

		pageHeader := PageHeader{}
		pageHeader.Fill(pageHeaderBuffer)

		fmt.Printf("database page size: %v\n", dbHeader.PageSize)
		fmt.Printf("number of tables: %v\n", pageHeader.NumberOfCells)
	case ".tables":
		databaseFile, err := os.Open(databaseFilePath)
		if err != nil {
			log.Fatal(err)
		}

		dbHeaderBuffer := make([]byte, DatabaseHeaderSize)
		n, err := databaseFile.Read(dbHeaderBuffer)
		if err != nil {
			log.Fatal(err, n)
		}

		dbHeader := &DatabaseHeader{}
		dbHeader.Fill(dbHeaderBuffer)

		pageHeaderBuffer := make([]byte, 12)
		n, err = databaseFile.ReadAt(pageHeaderBuffer, DatabaseHeaderSize)
		if err != nil {
			log.Fatal(err, n)
		}

		pageHeader := &PageHeader{}
		pageHeader.Fill(pageHeaderBuffer)

		// The cell pointer array of a b-tree page immediately follows
		// the b-tree page header. Let K be the number of cells on the btree.
		// The cell pointer array consists of K 2-byte integer offsets to the
		// cell contents. The cell pointers are arranged in key order with
		// left-most cell (the cell with the smallest key) first and the
		// right-most cell (the cell with the largest key) last.
		cellPositions := make([]uint16, pageHeader.NumberOfCells)
		cellPositionBuffer := make([]byte, pageHeader.NumberOfCells*2)
		_, err = databaseFile.ReadAt(cellPositionBuffer, int64(DatabaseHeaderSize+pageHeader.Size))
		if err != nil {
			log.Fatal(err)
		}

		// fmt.Fprintf(os.Stderr, "dbHeader\t: %#v\n", dbHeader.Dump())
		// fmt.Fprintf(os.Stderr, "pageHeader\t: %#v\n", pageHeader.Dump())

		for i := 0; i < len(cellPositionBuffer); i += 2 {
			if err = binary.Read(bytes.NewReader(cellPositionBuffer[i:i+2]), binary.BigEndian, &cellPositions[i/2]); err != nil {
				log.Fatal(err)
			}
		}

		var endOfContent uint16 = dbHeader.PageSize - uint16(dbHeader.ReservedBytes)
		var records []Record
		for _, pos := range cellPositions {
			cellContentBuffer := make([]byte, endOfContent-pos)
			_, err = databaseFile.ReadAt(cellContentBuffer, int64(pos))
			if err != nil {
				log.Fatal(err)
			}
			record := Record{}
			record.Fill(cellContentBuffer)
			records = append(records, record)

			endOfContent = pos
		}

		output := make([]string, 0)
		for _, record := range records {
			if record.Body.Payload[0] == "table" {
				s, ok := record.Body.Payload[2].(string)
				if ok && !strings.HasPrefix(s, "sqlite_") {
					output = append(output, s)
				}
			}
		}
		fmt.Fprintf(os.Stdout, "%s\n", strings.Join(output, " "))
	default:
		fmt.Println("Unknown command", command)
		os.Exit(1)
	}
}
