package util

import (
	"crypto/sha1"
	"encoding/hex"
	"io"

	"go.uber.org/zap"
)

func GetFileHash(f io.Reader) (string, error) {
	h := sha1.New()
	if _, err := io.Copy(h, f); err != nil {
		zap.L().Error("compute hash wrong", zap.Error(err))
		return "", err
	}
	hashString := hex.EncodeToString(h.Sum(nil))
	return hashString, nil
}
