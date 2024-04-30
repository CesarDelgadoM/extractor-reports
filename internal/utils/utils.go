package utils

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/CesarDelgadoM/extractor-reports/pkg/logger/zap"
)

func ToBytes(data any) []byte {
	bytes, err := json.Marshal(&data)
	if err != nil {
		zap.Log.Error("Failed to make marshal: ", err)
	}
	return bytes
}

func TimestampID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func MD5(key string) string {
	hash := md5.New()
	hash.Write([]byte(key))

	return hex.EncodeToString(hash.Sum(nil))
}
