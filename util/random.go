package util

import (
	"math/rand"
	"strings"
	"time"
	"unicode"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"
const nums = "0123456789"
const specialCharacters = "!@#$%^&*()_-"

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

func RandomFullName() (firstName string, lastName string) {
	firstName = RandomOwner()
	lastName = RandomOwner()
	return firstName, lastName

}

func RandomMoney() int64 {
	return rand.Int63n(10000000)
}

func RandomCurrency() string {
	currencies := []string{EUR, USD, CAD}
	return currencies[rand.Intn(len(currencies))]
}

func RandomHashPassword() string {
	var sb strings.Builder
	var ch byte
	for i := 0; i < 7; i++ {
		n := rand.Intn(2)
		if n == 0 {
			ch = alphabet[rand.Intn(len(alphabet))]
		} else if n == 1 {
			ch = nums[rand.Intn(len(nums))]
		} else {
			ch = specialCharacters[rand.Intn(len(specialCharacters))]
		}
		sb.WriteByte(ch)
	}
	return sb.String()
}

func RandomEmail() string {
	var sb strings.Builder
	var ch byte
	for i := 0; i < 7; i++ {
		ch = alphabet[rand.Intn(len(alphabet))]
		sb.WriteByte(ch)
	}
	sb.WriteString("@gmail.com")
	return sb.String()
}
