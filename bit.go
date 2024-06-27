package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
)

func main() {
	private_key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	public_key := private_key.PublicKey
	x_bytes := public_key.X.Bytes()[:]
	fmt.Println("privateKey:", "{ x:", hex.EncodeToString(private_key.X.Bytes()), "y: ", hex.EncodeToString(private_key.Y.Bytes()), "}")

	var buffer bytes.Buffer
	buffer.WriteByte(0x02)
	buffer.Write(x_bytes)
	var result []byte = buffer.Bytes()

	fmt.Println("combined public key:", hex.EncodeToString(result))

	sha256_pubkey := sha256.Sum256(result)

	fmt.Println("sha256 pubkey:", hex.EncodeToString(sha256_pubkey[:]))

	h := ripemd160.New()
	h.Write(sha256_pubkey[:])
	ripemd_sha256_pubkey := h.Sum(nil)

	fmt.Println("ripemnd-160 of sha256:", hex.EncodeToString(ripemd_sha256_pubkey))

	var buffer2 bytes.Buffer
	buffer2.WriteByte(0x00)
	buffer2.Write(ripemd_sha256_pubkey)
	var result2 []byte = buffer2.Bytes()

	fmt.Println("version byte + ripemd-160:", hex.EncodeToString(result2))

	sha256_result2 := sha256.Sum256(result2)
	fmt.Println("1st sha256:", hex.EncodeToString(sha256_result2[:]))
	sha256_sha256_result2 := sha256.Sum256(sha256_result2[:])
	fmt.Println("2nd sha256:", hex.EncodeToString(sha256_sha256_result2[:]))

	address_checksum := sha256_sha256_result2[:4]

	fmt.Println("checksum (last 4B of 2nd sha256): ", hex.EncodeToString(address_checksum))

	var buffer3 bytes.Buffer
	buffer3.Write(result2)
	buffer3.Write(address_checksum)
	var result3 []byte = buffer3.Bytes()

	fmt.Println("version B + ripemd + checksum:", hex.EncodeToString(result3))

	base58encoded := base58.Encode(result3)
	fmt.Println("(base58) bitcoin address: ", base58encoded)
}
