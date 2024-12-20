package storage

type MongoClient struct {
}

func NewMongo() *MongoClient {
    return &MongoClient{}
}

func (m *MongoClient) CheckForImage(hashId string) (string, error) {
    return "", NotStoredError
}

func (m *MongoClient) StoreImage(hashId, ansi string) error {
    return nil
}
