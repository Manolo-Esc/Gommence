package opo_uid

import (
	"log"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type packetDateSample struct {
	day, hours, minutes int
	expected            string
}

var packetDateSamples = []packetDateSample{
	{1, 0, 0, "0800"},    // Día 1, 00:00
	{15, 23, 59, "7DFB"}, // Día 15, 23:59
	{31, 23, 59, "FDFB"}, // Día 31, 23:59
	{31, 12, 30, "FB1E"}, // Día 31, 12:30
	{5, 6, 15, "298F"},   // Día 5, 06:15
	{7, 7, 7, "39C7"},    // Día 7, 07:07
	{2, 2, 2, "1082"},    // Día 2, 02:02
	{28, 10, 5, "E285"},  // Día 28, 10:05
	{3, 23, 45, "1DED"},  // Día 3, 23:45
}

func TestPacketDate(t *testing.T) {
	for _, tc := range packetDateSamples {
		result := packetNumbers(tc.day, tc.hours, tc.minutes)
		assert.Equal(t, tc.expected, result)
	}
}

type uidPairs struct {
	uid, uidRaw string
}

var uidSamples = []uidPairs{
	{"5cwFANhpap8", "FFFFFFFFFFFFFFF"}, // bigger theorical value: it takes 11 digits
	{"32222222222", "50633659656D971"}, // first value that needs 11 digits
	{"ZZZZZZZZZZ", "50633659656D970"},  // last value with 10 digits 0x6336 -> day 12, hours 12, minutes 54 so we can exced this numeber easily
	{"3gZH3h2fV", "000800125000000"},   // smallest theorical value: it takes 9 digits
	{"32quX3yrp2b", "50F800125000000"}, // seconds 50, day 31 -> 11 digits
	{"TRgQviFLxt", "47B59F2256CA2D4"},
	{"32f2BrLGJzx", "50B59F2258D40B9"},
	{"35M7nR9pbMP", "55B59F225368D59"},
	{"37UmCuXva5J", "58B59F225516270"},
	{"37UmCuXc3eX", "58B59F2251E2E45"},
	{"3h8LpXwzxVA", "65B59F2251D9C9D"},
	{"3Fwwf5AkAY4", "86B59F22563CAA6"},
	{"3Q3wgDTzkPn", "92B59F225B273F6"},
	{"3VvXDyQmYs", "02B5A02254232C4"},
	{"64LeioBSqK", "05B5A022528E296"},
	{"fhbZKY3XZm", "12B5A0225CDD7F8"},
	{"iPgLbmzocL", "17B5A0225FEAEB5"},
	{"2NUdMuUztLY", "40BB172259FE721"},
	{"2VhXyrnR6qq", "49BB1722548D2EE"},
	{"33EHknRxSUp", "52BB17225D72CAD"},
	{"3grdSfLTXgU", "64BB17225098A44"},
	{"3h9CC95op5L", "65BB17225901418"},
}

func TestDecodeUIDs(t *testing.T) {
	for _, tc := range uidSamples {
		decimalValue := base62ToDecimal(tc.uid)
		hexStr := decimalToHex(decimalValue)
		//log.Printf("{\"%s\", \"%s\", %X}", tc.uid, hexStr, decimalValue)
		assert.Equal(t, hexStr, tc.uidRaw)
	}
	uidRaw, err := DecodeUid("223gZH3h2fV")
	assert.Nil(t, err)
	assert.Equal(t, "000800125000000", uidRaw)
}

func TestEncodingPadWithZeros(t *testing.T) {
	var uid, uidRaw string
	uidRaw = "000800125000000"
	uid = encodeHexNumberString(uidRaw)
	assert.Equal(t, "223gZH3h2fV", uid) // added 2 zeros ("2")
	uidRaw = "47B59F2256CA2D4"
	uid = encodeHexNumberString(uidRaw)
	assert.Equal(t, "2TRgQviFLxt", uid) // added 1 zeros ("2")
}

func TestDecodingPadWithZeros(t *testing.T) {
	// we dont have 0 or 1 in our alphabet so the number 0 is represented by "2"
	uidRawExpected := "000800125000000"
	uidRaw, err := DecodeUid("3gZH3h2fV")
	assert.Equal(t, uidRaw, uidRawExpected)
	assert.Nil(t, err)
	uidRaw, err = DecodeUid("23gZH3h2fV")
	assert.Equal(t, uidRaw, uidRawExpected)
	assert.Nil(t, err)
	uidRaw, err = DecodeUid("223gZH3h2fV")
	assert.Equal(t, uidRaw, uidRawExpected)
	assert.Nil(t, err)
	uidRaw, err = DecodeUid("2222223gZH3h2fV")
	assert.Equal(t, uidRaw, uidRawExpected)
	assert.Nil(t, err)
}

func TestUidsHave11Digits(t *testing.T) {
	var uid string
	start := time.Now()
	numIterations := 0
	for {
		if time.Since(start) >= 1000*time.Millisecond {
			break
		}
		uid = New()
		assert.Equal(t, 11, len(uid))
		//time.Sleep(time.Millisecond)
		numIterations++
	}
	log.Printf("Number of iterations: %d", numIterations)
	uidRawMax := "99FDFBC45FFFFFF" // FDFB -> day 31, 23:59, December/2045, rand: FFFFFF
	uid = encodeHexNumberString(uidRawMax)
	assert.Equal(t, 11, len(uid))
	uidRawMin := "000800125000000" // 0800 -> day 1, 00:00, January/2025, rand: 000000
	uid = encodeHexNumberString(uidRawMin)
	assert.Equal(t, 11, len(uid))
}

func TestGenerateSomeUIDs(t *testing.T) {
	t.SkipNow()
	min := 10
	max := 50
	for i := 0; i < 20; i++ {
		//uid := New()
		encDate := getEncodedDate2()
		encRand := getRandomNumber()
		uidRaw := encDate + encRand
		uid := encodeHexNumberString(uidRaw)
		log.Printf("{\"%s\", \"%s\"}", uid, uidRaw)
		randomNum := time.Duration(80*rand.Intn(max-min+1) + min)
		time.Sleep(randomNum * time.Millisecond)
	}
}

func TestGenerateWellKnownUIDs(t *testing.T) {
	t.SkipNow()
	var uid, uidRaw string
	uidRaw = "FFFFFFFFFFFFFFF"
	uid = encodeHexNumberString(uidRaw)
	log.Printf("{\"%s\", \"%s\"}", uid, uidRaw)
	uidRaw = "50633659656D970"
	uid = encodeHexNumberString(uidRaw)
	log.Printf("{\"%s\", \"%s\"}", uid, uidRaw)
	uidRaw = "50633659656D971"
	uid = encodeHexNumberString(uidRaw)
	log.Printf("{\"%s\", \"%s\"}", uid, uidRaw)
	uidRaw = "000800125000000"
	uid = encodeHexNumberString(uidRaw)
	log.Printf("{\"%s\", \"%s\"}", uid, uidRaw)
	uidRaw = "50F800125000000"
	uid = encodeHexNumberString(uidRaw)
	log.Printf("{\"%s\", \"%s\"}", uid, uidRaw)

}
