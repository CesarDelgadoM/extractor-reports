package utils

import (
	"encoding/json"

	"github.com/CesarDelgadoM/extractor-reports/pkg/logger/zap"
)

func ToBytes(data any) []byte {
	bytes, err := json.Marshal(&data)
	if err != nil {
		zap.Log.Error("Failed to make marshal: ", err)
	}
	return bytes
}
