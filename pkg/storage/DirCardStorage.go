package storage

import (
	"github.com/LostPetInitiative/poiskzoo-ru-crawler/pkg/types"
)

type DirCardStorage struct {
	baseDir string
}

func NewDirCardStorage(path string) DirCardStorage {
	return DirCardStorage{
		baseDir: path,
	}
}

func (s *DirCardStorage) getLatestKnownCards() []types.CardID {
	panic("not implemented")
}
