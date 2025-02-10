package storage

type RedisClient struct {
}

func NewRedis() *RedisClient {
    return &RedisClient{}
}

func (m *RedisClient) CheckForImage(hashId string) (string, error) {
    return "", NotStoredError
}

func (m *RedisClient) StoreImage(hashId, ansi string) error {
    return nil
}

