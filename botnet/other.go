package botnet

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"math/rand"
	"time"
)

var StaticPrefix, _ = hex.DecodeString("42")
var NdcSignatureKey, _ = hex.DecodeString("F8E7A61AC3F725941E3AC7CAE2D688BE97F30B93")
var NdcDeviceKey, _ = hex.DecodeString("02B258C63559D8804321C5D5065AF320358D366F")

func RandInt(min int, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}

func RandBytes(length int) []byte {
	randBytes := make([]byte, length)
	rand.Seed(time.Now().UnixNano())
	rand.Read(randBytes)
	return randBytes
}

func NdcDevice() string {
	randToken := RandBytes(15)
	mac := hmac.New(sha1.New, NdcDeviceKey)
	mac.Write(StaticPrefix)
	mac.Write(randToken)
	return hex.EncodeToString(StaticPrefix) + hex.EncodeToString(randToken) + hex.EncodeToString(mac.Sum(nil))
}

func NdcSignature(data []byte) string {
	mac := hmac.New(sha1.New, NdcSignatureKey)
	mac.Write(data)
	return base64.StdEncoding.EncodeToString(append(StaticPrefix, mac.Sum(nil)...))
}

func TimeInMillis() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
