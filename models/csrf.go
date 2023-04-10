package models

import (
	"crypto/rand"
	"fmt"
)

func GenerateRandByteSlice() []byte {
	byteSlice := make([]byte, 32)
	if _, err := rand.Read(byteSlice); err != nil {
		fmt.Errorf("error generating byteSlice: %v", err)
	}

	return byteSlice
}
