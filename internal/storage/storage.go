package storage

/*

NEEDS:
caching
    TTL
hashing images

*/

import (
    "crypto/sha256"
    "encoding/hex"
    "errors"
    "image"
    "io"
)

var NotStoredError = errors.New("Hash not found in storage")

type Client interface {
    // TODO: create error for not in cache
    CheckForImage(hashId string) (string, error)
    StoreImage(hashId, ansi string) error
}

// TODO: make public
func CalculateImageHash(img image.Image) string {
    hash := sha256.New()
    io.WriteString(hash, img.Bounds().String())
    return hex.EncodeToString(hash.Sum(nil))
}
