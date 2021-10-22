package gob

import (
	"bytes"
	"encoding/gob"
	"io"
)

func encode(data interface{}) []byte{
	m := new(bytes.Buffer)
	enc := gob.NewEncoder(m)//Create a encoder
	err := enc.Encode(data) //encode
	if err!=nil{
		panic(err)
	}
	return m.Bytes()
}

func load(e interface{}, buf io.ReadCloser) {
	dec := gob.NewDecoder(buf)

	err := dec.Decode(e)
	if err != nil {
		panic(err)
	}
}
