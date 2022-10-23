package utils

import (
	"encoding/base64"
	"log"
)

type memoryWriter struct {
	data []byte
}

func (mw *memoryWriter) Write(p []byte) (n int, err error) {
	mw.data = append(mw.data, p...)
	return len(p), nil
}

func (mw *memoryWriter) GetBuffer() []byte {
	return mw.data
}

func newMemoryWriter() *memoryWriter {
	return &memoryWriter{}
}

func Base64Encode(s []byte) string {
	memWriter := newMemoryWriter()
	encoder := base64.NewEncoder(base64.StdEncoding, memWriter)
	data := []byte(s)
	_, err := encoder.Write(data)
	if err != nil {
		log.Panicln("Failed to encode data with base64")
	}
	encoder.Close()

	return string(memWriter.GetBuffer())
}
