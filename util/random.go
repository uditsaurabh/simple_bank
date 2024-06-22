package util

import (
	"math/rand"
	"strings"
	"time"
	"unicode"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandInt(max, min int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func RandomeString(n int) string {
	var sb strings.Builder
	for i := 0; i < n; i++ {
		ch := alphabet[rand.Intn(len(alphabet))]
		sb.WriteByte(ch)
	}
	return sb.String()

}

func RandomOwner() string {
	var sb strings.Builder
	var ch byte
	ch = alphabet[rand.Intn(len(alphabet))]
	upper_ch := unicode.ToUpper(rune(ch))
	sb.WriteRune(upper_ch)
	for i := 0; i < 7; i++ {
		ch = alphabet[rand.Intn(len(alphabet))]
		sb.WriteByte(ch)
	}
	return sb.String()
}

func RandomMoney() int64 {
	return rand.Int63n(10000000)
}
func RandomCurrency() string {
	currencies := []string{"EUR", "RUPEES", "USD"}
	return currencies[rand.Intn(len(currencies))]
}
