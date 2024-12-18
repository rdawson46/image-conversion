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
    "image"
    "io"
)

// TODO: add a client to the server
type Client interface {
    GetImage(img image.Image, width int) (string, error)
}

func CalculateImageHash(img image.Image) string {
    hash := sha256.New()
    io.WriteString(hash, img.Bounds().String())
    return hex.EncodeToString(hash.Sum(nil))
}
