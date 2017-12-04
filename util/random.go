package util

import (
	"math/rand"
	"strings"
	"time"
)

var r *rand.Rand // Rand for this package.

func init() {
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
}

// RandomString 随机字符串
func RandomString(strlen int) string {
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := []string{}
	for i := 0; i < strlen; i++ {
		index := r.Intn(len(chars))
		result[i] = chars[index : index+1]
	}
	return strings.Join(result, "")
}
