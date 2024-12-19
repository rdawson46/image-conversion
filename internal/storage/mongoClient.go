package storage

import (
	"errors"
	"image"
)

type MongoClient struct {
}

func NewMongo() *MongoClient {
    return &MongoClient{}
}

// error handling not really needed for sample
func (m *MongoClient) GetImage(img image.Image, width int) (string, bool, error) {
    return "", false, errors.New("Not implemented")
}
