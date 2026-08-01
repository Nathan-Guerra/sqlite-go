package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
)

/*
Offset	Size	Description
0		1		The one-byte flag at offset 0 indicating the b-tree page type.
				A value of 2 (0x02) means the page is an interior index b-tree page.
				A value of 5 (0x05) means the page is an interior table b-tree page.
				A value of 10 (0x0a) means the page is a leaf index b-tree page.
				A value of 13 (0x0d) means the page is a leaf table b-tree page.
				Any other value for the b-tree page type is an error.
1		2		The two-byte integer at offset 1 gives the start of the first freeblock on the page, or is zero if there are no freeblocks.
3		2		The two-byte integer at offset 3 gives the number of cells on the page.
5		2		The two-byte integer at offset 5 designates the start of the cell content area. A zero value for this integer is interpreted as 65536.
7		1		The one-byte integer at offset 7 gives the number of fragmented free bytes within the cell content area.
8		4		The four-byte page number at offset 8 is the right-most pointer. This value appears in the header of interior b-tree pages only and is omitted from all other pages.
*/

type PageHeader struct {
	PageType               byte
	FirstFreeblock         int16
	NumberOfCells          int16
	StartOfCellContentArea int16
	FragmentedFreeBytes    byte
	RightmostPointer       int32
	Size                   byte
}

func (p *PageHeader) Dump() map[string]any {
	response := make(map[string]any, 7)
	response["PageType"] = p.PageType
	response["FirstFreeblock"] = p.FirstFreeblock
	response["NumberOfCells"] = p.NumberOfCells
	response["StartOfCellContentArea"] = p.StartOfCellContentArea
	response["FragmentedFreeBytes"] = p.FragmentedFreeBytes
	response["RightmostPointer"] = p.RightmostPointer
	response["Size"] = p.Size
	return response
}

func (p *PageHeader) Fill(header []byte) {
	p.Size = 8
	p.PageType = header[0]

	if err := binary.Read(bytes.NewReader(header[1:3]), binary.BigEndian, &p.FirstFreeblock); err != nil {
		log.Fatal(err, "p.FirstFreeblock")
	}
	if err := binary.Read(bytes.NewReader(header[3:5]), binary.BigEndian, &p.NumberOfCells); err != nil {
		log.Fatal(err, "p.NumberOfCells")
	}
	if err := binary.Read(bytes.NewReader(header[5:7]), binary.BigEndian, &p.StartOfCellContentArea); err != nil {
		log.Fatal(err, "p.StartOfCellContentArea")
	}

	p.FragmentedFreeBytes = header[7]

	if p.PageType == 0x02 || p.PageType == 0x05 {
		if err := binary.Read(bytes.NewReader(header[8:12]), binary.BigEndian, &p.RightmostPointer); err != nil {
			log.Fatal(err, "p.RightmostPointer")
		}

		p.Size += 4
	}
}

func (p *PageHeader) Type() string {
	switch p.PageType {
	case 0x02:
		return "interior index b-tree page"
	case 0x05:
		return "interior table b-tree page"
	case 0x0a:
		return "leaf index b-tree page"
	case 0x0d:
		return "leaf table b-tree page"
	default:
		return fmt.Sprintf("unknown page type (0x%x)", p.PageType)
	}
}
