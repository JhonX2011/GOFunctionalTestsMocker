package mock

import "sync"

type Repository interface {
	Save(info mock) error
	GetAll() []mock
}

type inMemoryRepository struct {
	storage *sync.Map
}

func newRepository() Repository {
	return &inMemoryRepository{
		storage: &sync.Map{},
	}
}

func (repo *inMemoryRepository) Save(info mock) error {
	LogInfo("storing aggregate &v", info)
	repo.storage.Store(info.ID, info)
	LogInfo("aggregate stored")
	return nil
}

func (repo *inMemoryRepository) GetAll() []mock {
	LogInfo("getting aggregates")
	var results []mock
	repo.storage.Range(func(_, value any) bool {
		results = append(results, value.(mock))
		return true
	})
	LogInfo("the aggregates &v is returned", results)
	return results
}
