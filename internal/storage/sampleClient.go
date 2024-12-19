package storage

import (
    "github.com/rdawson46/pic-conversion/internal/conversion"
	"image"
)

type SampleDB struct {
    db map[string]string
}

func NewSampleDB() *SampleDB {
    return &SampleDB{
        db: make(map[string]string),
    }
}

// error handling not really needed for sample
func (s *SampleDB) GetImage(img image.Image, width int) (string, error) {
    // create hash
    hash := calculateImageHash(img)

    // check in db
    if value, ok := s.db[hash]; ok {
        return value, nil
    }

    // else create everything
    ansi := conversion.ConvertImage(img, width)

    // and store
    s.db[hash] = ansi

    return ansi, nil
}
