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


func calculateImageHash(img image.Image) string {
    hash := sha256.New()
    io.WriteString(hash, img.Bounds().String())
    return hex.EncodeToString(hash.Sum(nil))
}

// TODO: will need a db connection
func ProcessImage(img image.Image) string {
    imageHash := calculateImageHash(img)

    // check for calculated art with hash

    // if so return 

    // else caclulate

    // store

    // then return

    return ""
}
