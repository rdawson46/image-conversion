package storage

type SampleDB struct {
    db map[string]string
}

func NewSampleDB() *SampleDB {
    return &SampleDB{
        db: make(map[string]string),
    }
}

func (s *SampleDB) CheckForImage(hashId string) (string, error) {
    if value, ok := s.db[hashId]; ok {
        return value, nil
    }
    return "", NotStoredError
}

func (s *SampleDB) StoreImage(hashId, ansi string) error {
    s.db[hashId] = ansi
    return nil
}
