package strings

import (
	"crypto/aes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

var (
	SESSION_SECRET string
)

func init() {
	SESSION_SECRET = MustEnv("SESSION_SECRET", "FF2nkf91HWiCkY+xGMhtaQ==")
}

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func String(length int) string {
	return StringWithCharset(length, charset)
}

func RandomHex(n int) string {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return ""
	}
	return hex.EncodeToString(bytes)
}

func MustEnv(env string, target string) string {
	newEnv, ok := os.LookupEnv(env)
	if ok {
		target = newEnv
	}

	return target
}

func Split(s string, sep string) []string {
	return strings.Split(s, sep)
}

func TrimSpace(s string) string {
	return strings.TrimSpace(s)
}

func AesEncrypt(plainText string) string {
	b, err := aes.NewCipher([]byte(SESSION_SECRET))
	if err != nil {
		panic(err)
	}

	chiperText := make([]byte, len(plainText))
	b.Encrypt(chiperText, []byte(plainText))

	return string(chiperText)
}

func Hash256(v string) string {
	h := sha256.New()
	h.Write([]byte(v))
	bs := h.Sum(nil)
	return fmt.Sprintf("%x", bs)
}

func NewGravatar(email string) string {
	h := sha256.New()
	hasher := h.Sum([]byte(strings.TrimSpace(email)))
	hash := hex.EncodeToString(hasher[:])
	return hash
}
