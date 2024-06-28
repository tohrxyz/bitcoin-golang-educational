package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"time"

	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
)

func main() {
	start_time := time.Now()
	i := 1
	for i < 100000 {
		data := RandStringRunes(i)
		add := gen_add(data)
		matched, _ := regexp.Compile("^1(P(P(K(E)?)?)?)?$")
		if matched.MatchString(add) {
			fmt.Println("-----> address:", add, "seed:", data)
			break
		}

		if i == 99999 {
			i = 1
		} else {
			i++
		}
	}

	fmt.Println("Elapsed:", time.Since(start_time))
}

func gen_string(length int) string {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZÅÄÖ" +
		"abcdefghijklmnopqrstuvwxyzåäö" +
		"0123456789")
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}

	return b.String()
}

func gen_add(input string) string {
	sha256_pubkey := sha256.Sum256(bytes.NewBufferString(input).Bytes())

	// fmt.Println("sha256 pubkey:", hex.EncodeToString(sha256_pubkey[:]))

	h := ripemd160.New()
	h.Write(sha256_pubkey[:])
	ripemd_sha256_pubkey := h.Sum(nil)

	// fmt.Println("ripemnd-160 of sha256:", hex.EncodeToString(ripemd_sha256_pubkey))

	var buffer2 bytes.Buffer
	buffer2.WriteByte(0x00)
	buffer2.Write(ripemd_sha256_pubkey)
	var result2 []byte = buffer2.Bytes()

	// fmt.Println("version byte + ripemd-160:", hex.EncodeToString(result2))

	sha256_result2 := sha256.Sum256(result2)
	// fmt.Println("1st sha256:", hex.EncodeToString(sha256_result2[:]))
	sha256_sha256_result2 := sha256.Sum256(sha256_result2[:])
	// fmt.Println("2nd sha256:", hex.EncodeToString(sha256_sha256_result2[:]))

	address_checksum := sha256_sha256_result2[:4]

	// fmt.Println("checksum (last 4B of 2nd sha256): ", hex.EncodeToString(address_checksum))

	var buffer3 bytes.Buffer
	buffer3.Write(result2)
	buffer3.Write(address_checksum)
	var result3 []byte = buffer3.Bytes()

	// fmt.Println("version B + ripemd + checksum:", hex.EncodeToString(result3))

	base58encoded := base58.Encode(result3)
	// fmt.Println("(base58) bitcoin address: ", base58encoded)
	// if base58encoded[:1] == "1P" || base58encoded[:2] == "1PP" || base58encoded[:3] == "1PPK" || base58encoded[:4] == "1PPKE" {
	// 	fmt.Println("address:", base58encoded, "seed:", input)
	// }

	// matched, _ := regexp.MatchString("^1(P(P(K(E)?)?)?)?$", base58encoded)
	// if matched {
	// 	fmt.Println("-----> address:", base58encoded, "seed:", input)
	// }

	// fmt.Println("NOT", base58encoded)
	return base58encoded
}

func RandStringRunes(n int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789/?.,-=]:!@#$%^&*()±§")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
