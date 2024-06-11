package httpinternal

import (
	"fmt"
)

func buildStringFromSlice[T any](slice []T) string {
	return fmt.Sprintf("%+v", slice)
}
