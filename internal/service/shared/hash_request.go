package shared

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strings"
)

func HashRequest(req *http.Request) string {
	hash := sha256.New()
	hash.Write([]byte(req.Method))        // Метод запроса
	hash.Write([]byte(req.URL.String()))  // URL
	for key, values := range req.Header { // Заголовки
		for _, value := range values {
			hash.Write([]byte(strings.ToLower(key)))
			hash.Write([]byte(":"))
			hash.Write([]byte(value))
		}
	}
	return hex.EncodeToString(hash.Sum(nil))
}
