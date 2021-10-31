package util

import (
    "crypto/rand"
    "encoding/base64"
)

// GenerateRandomBytes returns securely generated random bytes.
// It will return an errors if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomBytes(n int) ([]byte, error) {
    b := make([]byte, n)
    _, err := rand.Read(b)
    // Note that err == nil only if we read len(b) bytes.
    if err != nil {
        return nil, err
    }

    return b, nil
}

// GenerateRandomString returns a URL-safe, base64 encoded
// securely generated random string.
func GenerateRandomString(s int) (string, error) {
    b, err := GenerateRandomBytes(s)
    return base64.URLEncoding.EncodeToString(b), err
}

func RandomUnsignedInt() uint {
    b := make([]byte, 4)
    _, _ = rand.Read(b)
    return (uint(b[0]) << 24) | (uint(b[1]) << 16) | (uint(b[2]) << 8) | uint(b[3])
}

func RandomInt() int {
    res := int(RandomUnsignedInt())
    if res < 0 {
        res = -res
    }
    return res
}
