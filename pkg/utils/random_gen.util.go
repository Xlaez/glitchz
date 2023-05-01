package utils

import (
	"math/rand"
	"strings"
	"time"
)

const alp = "abcdefghijklmnopqrstuvwxyz"
const num = "1234567890"

func init() {
	rand.NewSource(time.Now().Unix())
}

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func SixDigitsCode() string {
	var sb strings.Builder

	k := len(num)

	for i := 0; i < 6; i++ {
		c := num[rand.Intn(k)]
		sb.WriteByte(c)
	}
	return sb.String()

}

func RandomStr(n int) string {
	var sb strings.Builder

	k := len(alp)

	for i := 0; i < n; i++ {
		c := alp[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}
