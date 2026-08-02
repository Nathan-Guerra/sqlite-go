package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
)

type Header struct {
	Size        int64
	SerialTypes []int64
}

func (h *Header) FillSerialTypes(b []byte) int {
	var offset int64 = 0

	for offset < h.Size || int(offset) >= len(b) {
		st, n := varint(b[offset:])
		h.SerialTypes = append(h.SerialTypes, st)
		offset += int64(n)
	}

	return int(offset)
}

type Body struct {
	Payload map[string]any
}

func (b *Body) Dump() {
	fmt.Fprintf(os.Stderr, "%#v\n", b.Payload)
}

type Record struct {
	Size   int64
	Header *Header
	Body   *Body
}

func (r *Record) Fill(b []byte) {
	offset := 0
	_, n := varint(b) // payloadSize, n := varint(b)
	offset += n
	_, n = varint(b[offset:]) // rowID, n := varint(b)
	offset += n
	headerSize, headerSizeOffset := varint(b[offset:])
	offset += headerSizeOffset
	header := &Header{Size: headerSize - int64(headerSizeOffset)}
	n = header.FillSerialTypes(b[offset:])
	offset += n

	body := &Body{}
	body.Payload = make(map[string]any, 0)
	r.Header = header
	r.Body = body

	for i, st := range header.SerialTypes {
		if st == 0 {
			body.Payload[fmt.Sprintf("st-%v", i)] = "nil"
		} else if st == 0x1 {
			body.Payload[fmt.Sprintf("st-%v", i)] = int8(b[offset])
			offset += 1
		} else if st == 0x2 {
			var i int16
			if err := binary.Read(bytes.NewReader(b[offset:offset+2]), binary.BigEndian, &i); err != nil {
				log.Fatal(err, "0x2")
			}
			body.Payload[fmt.Sprintf("st-%v", i)] = i
			offset += 2
		} else if st == 0x3 {
			var i int32
			if err := binary.Read(bytes.NewReader(append([]byte{0}, b[offset:offset+3]...)), binary.BigEndian, &i); err != nil {
				log.Fatal(err, "0x3")
			}
			body.Payload[fmt.Sprintf("st-%v", i)] = i
			offset += 3
		} else if st == 0x4 {
			var i int32
			if err := binary.Read(bytes.NewReader(b[offset:offset+4]), binary.BigEndian, &i); err != nil {
				log.Fatal(err, "0x4")
			}
			body.Payload[fmt.Sprintf("st-%v", i)] = i
			offset += 4
		} else if st == 0x5 {
			var i int64
			if err := binary.Read(bytes.NewReader(append([]byte{0, 0}, b[offset:offset+6]...)), binary.BigEndian, &i); err != nil {
				log.Fatal(err, "0x5")
			}
			body.Payload[fmt.Sprintf("st-%v", i)] = i
			offset += 6
		} else if st == 0x6 {
			var i int64
			if err := binary.Read(bytes.NewReader(b[offset:offset+8]), binary.BigEndian, &i); err != nil {
				log.Fatal(err, "0x6")
			}
			body.Payload[fmt.Sprintf("st-%v", i)] = i
			offset += 8
		} else if st == 0x7 {
			var i float64
			if err := binary.Read(bytes.NewReader(b[offset:offset+8]), binary.BigEndian, &i); err != nil {
				log.Fatal(err, "0x7")
			}
			body.Payload[fmt.Sprintf("st-%v", i)] = i
			offset += 8
		} else if st == 0x8 {
			body.Payload[fmt.Sprintf("st-%v", i)] = int8(0)
		} else if st == 0x9 {
			body.Payload[fmt.Sprintf("st-%v", i)] = int8(1)
		} else if st >= 12 && st%2 == 0 {
			body.Payload[fmt.Sprintf("st-%v", i)] = string(b[offset : offset+(int((st-12)/2))])
			offset += int((st - 12) / 2)
		} else if st >= 13 && st%2 == 1 {
			body.Payload[fmt.Sprintf("st-%v", i)] = string(b[offset : offset+(int((st-13)/2))])
			offset += int((st - 13) / 2)
		}
	}

	body.Dump()
}
