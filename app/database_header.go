package main

import (
	"bytes"
	"encoding/binary"
	"log"
)

/*
Offset	Size	Description
0		16		The header string: "SQLite format 3\000"
16		2		The database page size in bytes. Must be a power of two between 512 and 32768 inclusive, or the value 1 representing a page size of 65536.
18		1		File format write version. 1 for legacy; 2 for WAL.
19		1		File format read version. 1 for legacy; 2 for WAL.
20		1		Bytes of unused "reserved" space at the end of each page. Usually 0.
21		1		Maximum embedded payload fraction. Must be 64.
22		1		Minimum embedded payload fraction. Must be 32.
23		1		Leaf payload fraction. Must be 32.
24		4		File change counter.
28		4		Size of the database file in pages. The "in-header database size".
32		4		Page number of the first freelist trunk page.
36		4		Total number of freelist pages.
40		4		The schema cookie.
44		4		The schema format number. Supported schema formats are 1, 2, 3, and 4.
48		4		Default page cache size.
52		4		The page number of the largest root b-tree page when in auto-vacuum or incremental-vacuum modes, or zero otherwise.
56		4		The database text encoding. A value of 1 means UTF-8. A value of 2 means UTF-16le. A value of 3 means UTF-16be.
60		4		The "user version" as read and set by the user_version pragma.
64		4		True (non-zero) for incremental-vacuum mode. False (zero) otherwise.
68		4		The "Application ID" set by PRAGMA application_id.
72		20		Reserved for expansion. Must be zero.
92		4		The version-valid-for number.
96		4		SQLITE_VERSION_NUMBER
*/

type DatabaseHeader struct {
	HeaderString               string // The header string: "SQLite format 3\000"
	PageSize                   int16  // The database page size in bytes. Must be a power of two between 512 and 32768 inclusive, or the value 1 representing a page size of 65536.
	ReadVersion                byte   // File format write version. 1 for legacy; 2 for WAL.
	WriteVersion               byte   // File format read version. 1 for legacy; 2 for WAL.
	ReservedBytes              byte   // Bytes of unused "reserved" space at the end of each page. Usually 0.
	MaxEmbeddedPayloadFraction byte   // Maximum embedded payload fraction. Must be 64.
	MinEmbeddedPayloadFraction byte   // Minimum embedded payload fraction. Must be 32.
	LeafPayloadFraction        byte   // Leaf payload fraction. Must be 32.
	FileChangeCounter          int32  // File change counter.
	SizeOfDatabaseInPages      int32  // Size of the database file in pages. The "in-header database size".
	FirstFreelistTrunkPage     int32  // Page number of the first freelist trunk page.
	TotalFreelistPages         int32  // Total number of freelist pages.
	SchemaCookie               int32  // The schema cookie.
	SchemaFormatNumber         int32  // The schema format number. Supported schema formats are 1, 2, 3, and 4.
	PageCacheSize              int32  // Default page cache size.
	LargestRootBTree           int32  // The page number of the largest root b-tree page when in auto-vacuum or incremental-vacuum modes, or zero otherwise.
	TextEncoding               int32  // The database text encoding. A value of 1 means UTF-8. A value of 2 means UTF-16le. A value of 3 means UTF-16be.
	UserVersino                int32  // The "user version" as read and set by the user_version pragma.
	IncrementalVacuumMode      int32  // True (non-zero) for incremental-vacuum mode. False (zero) otherwise.
	ApplicationID              int32  // The "Application ID" set by PRAGMA application_id.
	VersionValidFor            int32  // The version-valid-for number.
	SqliteVersionNumber        int32  // SQLITE_VERSION_NUMBER
}

func (d *DatabaseHeader) Fill(header []byte) {
	d.HeaderString = string(header[:16])

	if err := binary.Read(bytes.NewReader(header[16:18]), binary.BigEndian, &d.PageSize); err != nil {
		log.Fatal(err, "d.PageSize")
	}

	d.ReadVersion = header[18]
	d.WriteVersion = header[19]
	d.ReservedBytes = header[20]
	d.MaxEmbeddedPayloadFraction = header[21]
	d.MinEmbeddedPayloadFraction = header[22]
	d.LeafPayloadFraction = header[23]

	if err := binary.Read(bytes.NewReader(header[24:28]), binary.BigEndian, &d.FileChangeCounter); err != nil {
		log.Fatal(err, "d.FileChangeCounter")
	}

	if err := binary.Read(bytes.NewReader(header[28:32]), binary.BigEndian, &d.SizeOfDatabaseInPages); err != nil {
		log.Fatal(err, "d.SizeOfDatabaseInPages")
	}

	if err := binary.Read(bytes.NewReader(header[32:36]), binary.BigEndian, &d.FirstFreelistTrunkPage); err != nil {
		log.Fatal(err, "d.FirstFreelistTrunkPage")
	}

	if err := binary.Read(bytes.NewReader(header[36:40]), binary.BigEndian, &d.TotalFreelistPages); err != nil {
		log.Fatal(err, "d.TotalFreelistPages")
	}

	if err := binary.Read(bytes.NewReader(header[40:44]), binary.BigEndian, &d.SchemaCookie); err != nil {
		log.Fatal(err, "d.SchemaCookie")
	}

	if err := binary.Read(bytes.NewReader(header[44:48]), binary.BigEndian, &d.SchemaFormatNumber); err != nil {
		log.Fatal(err, "d.SchemaFormatNumber")
	}

	if err := binary.Read(bytes.NewReader(header[48:52]), binary.BigEndian, &d.PageCacheSize); err != nil {
		log.Fatal(err, "d.PageCacheSize")
	}

	if err := binary.Read(bytes.NewReader(header[52:56]), binary.BigEndian, &d.LargestRootBTree); err != nil {
		log.Fatal(err, "d.LargestRootBTree")
	}

	if err := binary.Read(bytes.NewReader(header[56:60]), binary.BigEndian, &d.TextEncoding); err != nil {
		log.Fatal(err, "d.TextEncoding")
	}

	if err := binary.Read(bytes.NewReader(header[60:64]), binary.BigEndian, &d.UserVersino); err != nil {
		log.Fatal(err, "d.UserVersino")
	}

	if err := binary.Read(bytes.NewReader(header[64:68]), binary.BigEndian, &d.IncrementalVacuumMode); err != nil {
		log.Fatal(err, "d.IncrementalVacuumMode")
	}

	if err := binary.Read(bytes.NewReader(header[68:72]), binary.BigEndian, &d.ApplicationID); err != nil {
		log.Fatal(err, "d.ApplicationID")
	}

	if err := binary.Read(bytes.NewReader(header[92:96]), binary.BigEndian, &d.VersionValidFor); err != nil {
		log.Fatal(err, "d.VersionValidFor")
	}

	if err := binary.Read(bytes.NewReader(header[96:]), binary.BigEndian, &d.SqliteVersionNumber); err != nil {
		log.Fatal(err, "d.SqliteVersionNumber")
	}

}
