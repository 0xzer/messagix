package methods

import (
	"encoding/base64"
	"encoding/hex"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)


var (
	epochMutex    sync.Mutex
	lastTimestamp int64
	counter       int64
)
var Charset = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890")

/*
	Counter + Mutex logic to ensure unique epoch id for all calls
*/
func GenerateEpochId() int64 {
	epochMutex.Lock()
	defer epochMutex.Unlock()

	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	if timestamp == lastTimestamp {
		counter++
	} else {
		lastTimestamp = timestamp
		counter = 0
	}
	id := (timestamp << 22) | (counter << 12) | 42
	return id
}

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
	str := RandStr(6) + ":" + RandStr(6) + ":" + RandStr(6)
	return strings.ToLower(str)
}

func GenerateTraceId() string {
	uuidHex := strings.ReplaceAll(uuid.NewString(), "-", "")
	decodedHex, err := hex.DecodeString(uuidHex)
	if err != nil {
		log.Fatalf("failed to decode traceId string: %s", err)
	}
	return "#" + base64.RawURLEncoding.EncodeToString(decodedHex)
}