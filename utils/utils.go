package utils

import (
	"crypto/md5"
	"fmt"
	rndm "math/rand"
	"mime/multipart"
	"net/http"
	"slices"

	"github.com/julienschmidt/httprouter"
)

// --- CSRF Token Generation ---

func CSRF(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	csrf := GenerateRandomString(12)
	RespondWithJSON(w, http.StatusOK, map[string]any{
		"ok":         true,
		"csrf_token": csrf,
	})
}

// --- Random String and ID Generators ---

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz0123456789_ABCDEFGHIJKLMNOPQRSTUVWXYZ")
var digitRunes = []rune("0123456789")

// GenerateRandomString creates a random alphanumeric string of length n.
func GenerateRandomString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rndm.Intn(len(letterRunes))]
	}
	return string(b)
}

// GenerateRandomDigitString creates a random numeric string of length n.
func GenerateRandomDigitString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = digitRunes[rndm.Intn(len(digitRunes))]
	}
	return string(b)
}

// --- Hashing ---

func EncrypIt(strToHash string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(strToHash)))
}

// --- HTTP Response Helpers ---

func SendResponse(w http.ResponseWriter, status int, data any, message string, err error) {
	resp := map[string]any{
		"status":  status,
		"message": message,
		"data":    data,
	}
	if err != nil {
		resp["error"] = err.Error()
	}
	RespondWithJSON(w, status, resp)
}

// func RespondWithJSON(w http.ResponseWriter, status int, data any) {
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(status)
// 	if err := json.NewEncoder(w).Encode(data); err != nil {
// 		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
// 	}
// }

func SendJSONResponse(w http.ResponseWriter, status int, response any) {
	RespondWithJSON(w, status, response)
}

// --- Slice Helpers ---

func Contains(slice []string, value string) bool {
	return slices.Contains(slice, value)
}

// --- Image Validation ---

var SupportedImageTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/webp": true,
	"image/gif":  true,
	"image/bmp":  true,
	"image/tiff": true,
}

func ValidateImageFileType(w http.ResponseWriter, header *multipart.FileHeader) bool {
	mimeType := header.Header.Get("Content-Type")
	if !SupportedImageTypes[mimeType] {
		http.Error(w, "Invalid file type. Supported formats: JPEG, PNG, WebP, GIF, BMP, TIFF.", http.StatusBadRequest)
		return false
	}
	return true
}
