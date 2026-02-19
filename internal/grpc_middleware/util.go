package grpc_middleware

import (
	"crypto/rand"
	"time"
)

var (
	charSet  = "0123456789abcdefghijklmnopqrstuvwxyz"
	charLen = len(charSet)
)

// generateID generates a random ID
func generateID() string {
	var b [8]byte
	_, err := rand.Read(b[:])
	if err != nil {
		// Fallback to time-based ID
		return time.Now().Format("20060102150405")
	}

	id := make([]byte, 16)
	for i := 0; i < 8; i++ {
		id[i*2] = charSet[int(b[i])%charLen]
		id[i*2+1] = charSet[int(b[i]/16)%charLen]
	}
	return string(id[:8])
}
