package main

import "encoding/binary"
import "bytes"

func (in *MyStruct) Unpack(data []byte) error {
	r := bytes.NewReader(data)

	// browsers
	var browsersLenRaw uint32
	binary.Read(r, binary.LittleEndian, &browsersLenRaw)
	browsersRaw := make([]byte, browsersLenRaw)
	binary.Read(r, binary.LittleEndian, &browsersRaw)
	in.browsers = string(browsersRaw)

	// email
	var emailLenRaw uint32
	binary.Read(r, binary.LittleEndian, &emailLenRaw)
	emailRaw := make([]byte, emailLenRaw)
	binary.Read(r, binary.LittleEndian, &emailRaw)
	in.email = string(emailRaw)
	return nil
}

