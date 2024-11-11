package helpers

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash/crc32"
	"strings"
	"time"
)

// parseToken splits the token into ID and actual token string
func ParseToken(fullToken string) (uint, string, error) {
	parts := strings.Split(fullToken, "|")
	if len(parts) != 2 {
		return 0, "", fmt.Errorf("invalid token format")
	}

	var id uint
	_, err := fmt.Sscanf(parts[0], "%d", &id)
	if err != nil {
		return 0, "", fmt.Errorf("invalid token ID")
	}

	return id, parts[1], nil
}

func HashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

func GenerateToken() (string, error) {
	// Generate random bytes for token entropy (40 characters)
	bytes := make([]byte, 20) // 20 bytes will give us 40 hex chars
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	// Convert to hex string for the tokenEntropy
	tokenEntropy := hex.EncodeToString(bytes)

	// Calculate CRC32B hash of the tokenEntropy
	crc32q := crc32.MakeTable(crc32.IEEE)
	hash := fmt.Sprintf("%x", crc32.Checksum([]byte(tokenEntropy), crc32q))

	// Combine tokenEntropy and hash (prefix + tokenEntropy + crc32b hash)
	// Note: prefix is empty string by default in Laravel
	return fmt.Sprintf("%s%s%s", "", tokenEntropy, hash), nil
}

func GetAccessTokenExpiry() time.Time {
	now := time.Now()
	atExpiry := now.Add(time.Minute * 60)

	return atExpiry
}

func GetRefreshTokenExpiry() time.Time {
	now := time.Now()
	atExpiry := now.Add(time.Hour * 24 * 7)

	return atExpiry
}
