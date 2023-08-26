package methods

import (
	"math/rand"
	"strconv"
	"time"
)

var Charset = []rune("abcdefghijklmnopqrstuvwxyz1234567890")

func GenerateTimestampString() string {
	return strconv.Itoa(int(time.Now().UnixMilli()))
}

func GenerateSessionId() int64 {
	min := int64(2171078810009599)
	max := int64(4613554604867583)
	return rand.Int63n(max-min+1) + min
}

func RandStr(length int) string {
	b := make([]rune, length)
    for i := range b {
        b[i] = Charset[rand.Intn(len(Charset))]
    }
    return string(b)
}

func GenerateWebsessionID() string {
	return RandStr(6) + ":" + RandStr(6) + ":" + RandStr(6)
}

func GenerateEpochId() int64 {
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	id := (timestamp << 22) | (42 << 12)
	return id
}