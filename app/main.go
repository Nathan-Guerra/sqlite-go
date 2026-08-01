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

func varint(b []byte) (int64, int) {
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
	var highOrderBit byte = 0x80
	length := 0
	for _, v := range b {
		if length == 9 {
			break
		}
		length++
		if v&highOrderBit == 0 {
			break
		}
	}

	var value int64 = 0x0
	var mask byte = 0x7f
	for i, v := range b[:length] {
		if v&highOrderBit != 0 && i < 8 {
			value = (value << 7) | int64(v&mask)
			// fmt.Fprintf(os.Stderr, "value = (value << 7) | int64(v&mask) %b %b\n", value, v&mask)
		} else {
			value = (value << 8) | int64(v)
			// fmt.Fprintf(os.Stderr, "value = (value << 8) | int64(v) %b %b\n", value, v)

			break
		}
	}

	return value, length
}

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

		dbHeaderBuffer := make([]byte, 100)
		n, err := databaseFile.Read(dbHeaderBuffer)
		if err != nil {
			log.Fatal(err, n)
		}

		dbHeader := DatabaseHeader{}
		dbHeader.Fill(dbHeaderBuffer)

		pageHeaderBuffer := make([]byte, 12)
		n, err = databaseFile.ReadAt(pageHeaderBuffer, 100)
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

		dbHeaderBuffer := make([]byte, 100)
		n, err := databaseFile.Read(dbHeaderBuffer)
		if err != nil {
			log.Fatal(err, n)
		}

		dbHeader := DatabaseHeader{}
		dbHeader.Fill(dbHeaderBuffer)

		pageHeaderBuffer := make([]byte, 12)
		n, err = databaseFile.ReadAt(pageHeaderBuffer, 100)
		if err != nil {
			log.Fatal(err, n)
		}

		pageHeader := PageHeader{}
		pageHeader.Fill(pageHeaderBuffer)

		cellPos := make([]uint16, pageHeader.NumberOfCells)
		buff := make([]byte, pageHeader.NumberOfCells*2)
		_, err = databaseFile.ReadAt(buff, int64(100+pageHeader.Size))
		if err != nil {
			log.Fatal(err)
		}

		for i := 0; i < len(buff); i += 2 {
			if err = binary.Read(bytes.NewReader(buff[i:i+2]), binary.BigEndian, &cellPos[i/2]); err != nil {
				log.Fatal(err)
			}
		}

		fmt.Fprintf(os.Stderr, "%#v\n", cellPos)

		// tableContent := make([]byte, pageHeader.Size-byte(pageHeader.StartOfCellContentArea))
		// _, err = databaseFile.ReadAt(tableContent, int64(cellsOffset))
		// if err != nil {
		// 	log.Fatal(err)
		// }

		// totOffset := 0
		// recordSize, n := varint(tableContent)
		// totOffset += n

		// rowID, n := varint(tableContent[totOffset:])
		// totOffset += n

		// headerSize, n := varint(tableContent[totOffset:])
		// columns := make([]int64, 0)
		// header := tableContent[totOffset+n : totOffset+int(headerSize)]

		// i := 0
		// for i < len(header) {
		// 	serialType, o := varint(header[i:])
		// 	columns = append(columns, serialType)
		// 	i += o
		// }
		// totOffset += int(headerSize)

		// fmt.Fprintf(os.Stderr, "record size: %d\n", recordSize)
		// fmt.Fprintf(os.Stderr, "rowID: %d\n", rowID)
		// fmt.Fprintf(os.Stderr, "header size: %d\n", headerSize)
		// fmt.Fprintf(os.Stderr, "columns: %v\n", columns)

		// body := tableContent[totOffset:]
		// bodyOffset := 0
		// for _, col := range columns {
		// 	switch {
		// 	case col == 0:
		// 		fmt.Fprintf(os.Stderr, "Value is a NULL.\n")
		// 	case col == 1:
		// 		fmt.Fprintf(os.Stderr, "Value is an 8-bit twos-complement integer. %v\n", body[bodyOffset:bodyOffset+1])
		// 		bodyOffset++
		// 	case col == 2:
		// 		fmt.Fprintf(os.Stderr, "Value is an 8-bit twos-complement integer. %v\n", body[bodyOffset:bodyOffset+2])
		// 		bodyOffset += 2
		// 	case col == 3:
		// 		fmt.Fprintf(os.Stderr, "Value is a big-endian 24-bit twos-complement integer. %v\n", body[bodyOffset:bodyOffset+3])
		// 		bodyOffset += 3
		// 	case col == 4:
		// 		fmt.Fprintf(os.Stderr, "Value is a big-endian 32-bit twos-complement integer. %v\n", body[bodyOffset:bodyOffset+4])
		// 		bodyOffset += 4
		// 	case col == 5:
		// 		fmt.Fprintf(os.Stderr, "Value is a big-endian 48-bit twos-complement integer. %v\n", body[bodyOffset:bodyOffset+6])
		// 		bodyOffset += 6
		// 	case col == 6:
		// 		fmt.Fprintf(os.Stderr, "Value is a big-endian 64-bit twos-complement integer. %v\n", body[bodyOffset:bodyOffset+8])
		// 		bodyOffset += 8
		// 	case col == 7:
		// 		fmt.Fprintf(os.Stderr, "Value is a big-endian IEEE 754-2008 64-bit floating point number. %v\n", body[bodyOffset:bodyOffset+8])
		// 		bodyOffset += 8
		// 	case col == 8:
		// 		fmt.Fprintf(os.Stderr, "Value is the integer 0. (Only available for schema format 4 and higher.) %v\n", 0)
		// 	case col == 9:
		// 		fmt.Fprintf(os.Stderr, "Value is the integer 1. (Only available for schema format 4 and higher.) %v\n", 1)
		// 	case (col >= 12 && col%2 == 0):
		// 		length := (int(col) - 12) / 2
		// 		fmt.Fprintf(os.Stderr, "Value is a BLOB that is (N-12)/2 bytes in length. '%s'\n", body[bodyOffset:bodyOffset+length])
		// 		bodyOffset += length
		// 	case (col >= 13 && col%2 == 1):
		// 		length := (int(col) - 13) / 2
		// 		fmt.Fprintf(os.Stderr, "Value is a string in the text encoding and (N-13)/2 bytes in length. The nul terminator is not stored. '%s'\n", body[bodyOffset:bodyOffset+length])
		// 		bodyOffset += length
		// 	}
		// }
	default:
		fmt.Println("Unknown command", command)
		os.Exit(1)
	}
}
