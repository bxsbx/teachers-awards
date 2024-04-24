package util

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/google/uuid"
	"io"
)

func Getuid() string {
	b := make([]byte, 48)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return fmt.Sprintf("%x", md5.Sum([]byte(base64.URLEncoding.EncodeToString(b))))
}

func GetRandUUID() string {
	random, _ := uuid.NewRandom()
	return random.String()
}
