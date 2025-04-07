package opo_uid

import (
	crand "crypto/rand"
	"encoding/binary"
	"fmt"
	"math/big"
	"math/rand"
	"strings"
	"sync"
	"time"
)

type criptoRandIDGenerator struct {
	sync.Mutex
	randSource *rand.Rand
}

var (
	customBase62Alphabet = "23456789abcdefghijkmnopqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ" // No '0', 'O', 'l', 'I', '1' to avoid confusion
	base                 = big.NewInt(int64(len(customBase62Alphabet)))
	criptoRand           = &criptoRandIDGenerator{}
)

// The UIDs are 11 characters long
func New() string {
	uid := getEncodedDate2() + getRandomNumber()
	return encodeHexNumberString(uid)
}

func DecodeUid(base62 string) (string, error) {
	decimalValue := base62ToDecimal(base62)
	if decimalValue == nil {
		return "", fmt.Errorf("Characters not in the alphabet %s", base62)
	}
	return decimalToHex(decimalValue), nil
}

func getRandomNumber() string {
	criptoRand.Lock()
	defer criptoRand.Unlock()

	if criptoRand.randSource == nil { // first time only
		var rngSeed int64
		_ = binary.Read(crand.Reader, binary.LittleEndian, &rngSeed)
		criptoRand.randSource = rand.New(rand.NewSource(rngSeed))
	}

	buf := make([]byte, 3) // Number between 0 and 0xFFFFFF (16.777.215)
	for {
		_, _ = criptoRand.randSource.Read(buf)
		if buf[0] != 0 {
			break
		}
	}

	hexString := fmt.Sprintf("%X", buf)
	return hexString
}

func packetNumbers(day int, hours int, minutes int) string {
	var ret uint16 = (uint16(day) << 11) | (uint16(hours) << 6) | uint16(minutes)
	return fmt.Sprintf("%04X", ret)
}

func getEncodedDate2() string {
	now := time.Now()

	year := now.Year() - 2000
	month := now.Month()
	day := now.Day()        // 5 bits
	hours := now.Hour()     // 5 bits
	minutes := now.Minute() // 6 bits
	seconds := now.Second()
	milliseconds := now.Nanosecond() / 1e6

	// change scale to 0..99 instead of 0..59)
	centiSeconds := float64(seconds) + float64(milliseconds)/1000
	centiSeconds = centiSeconds * (100.0 / 60.0)
	seconds = int(centiSeconds)

	packeted := packetNumbers(day, hours, minutes)

	return fmt.Sprintf("%02d%s%s%d", seconds, packeted, fmt.Sprintf("%X", int(month)), year)
}

func getEncodedDate() string {
	now := time.Now()

	year := now.Year() - 2000
	month := now.Month()
	day := now.Day()        // 5 bits
	hours := now.Hour()     // 5 bits
	minutes := now.Minute() // 6 bits
	seconds := now.Second()
	milliseconds := now.Nanosecond() / 1e6

	centiSeconds := float64(seconds) + float64(milliseconds)/1000
	centiSeconds = centiSeconds * (100.0 / 60.0)
	seconds = int(centiSeconds)

	return fmt.Sprintf("%02d%02d%02d%02d%s%d",
		seconds, minutes, hours, day, strings.ToUpper(fmt.Sprintf("%X", int(month))), year)
}

func encodeHexNumberString(input string) string {
	// Convert hexadecimal string to big.Int
	num := new(big.Int)
	num.SetString(input, 16)

	var base62Str strings.Builder
	zero := big.NewInt(0)
	mod := new(big.Int)

	// Convert big.Int a Base62
	for num.Cmp(zero) > 0 {
		num.DivMod(num, base, mod)
		base62Str.WriteString(string(customBase62Alphabet[mod.Int64()]))
	}

	runes := []rune(base62Str.String())
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	result := string(runes)
	result = padLeftWithZero(result, 11)
	return result
}

func padLeftWithZero(input string, length int) string {
	if len(input) >= length {
		return input
	}

	zeroRune := string(customBase62Alphabet[0])
	paddingSize := length - len(input)
	padding := strings.Repeat(zeroRune, paddingSize)

	return padding + input
}

func base62ToDecimal(base62 string) *big.Int {
	decimal := big.NewInt(0)

	for _, char := range base62 {
		index := big.NewInt(int64(strings.IndexRune(customBase62Alphabet, char)))
		if index.Int64() == -1 {
			return nil
		}
		decimal.Mul(decimal, base)  // decimal *= BASE62
		decimal.Add(decimal, index) // decimal += index
	}

	return decimal
}

func decimalToHex(decimal *big.Int) string {
	return fmt.Sprintf("%015X", decimal)
}
