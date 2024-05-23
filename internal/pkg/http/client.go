package httpinternal

import (
	"net/http"
	"time"
)

func NewClient() *http.Client {
	return &http.Client{
		Timeout: 30 * time.Second,
	}
}
