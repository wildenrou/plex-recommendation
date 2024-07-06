package pg

import (
	b64 "encoding/base64"
	"fmt"
)

func buildStringFromSlice[T any](slice []T) string {
	return fmt.Sprintf("%+v", slice)
}

func toBase64(i string) string {
	return b64.StdEncoding.EncodeToString([]byte(i))
}
