package url

type memoryRepository struct {
	urls map[string]*Url
}

func NewMemoryRepository() *memoryRepository {
	return &memoryRepository{make(map[string]*Url)}
}

func (r *memoryRepository) IdExists(id string) bool {
	_, exists := r.urls[id]
	return exists
}

func (r *memoryRepository) SearchById(id string) *Url {
	return r.urls[id]
}

func (r *memoryRepository) SearchByUrl(url string) *Url {
	for _, u := range r.urls {
		if u.Destination == url {
			return u
		}
	}
	return nil
}

func (r *memoryRepository) Save(url Url) error {
	r.urls[url.ID] = &url
	return nil
}