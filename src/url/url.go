package url
import (
	"time"
	"math/rand"
	"net/url"
	"fmt"
)

const (
	size = 5
	tokens = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890_-+"
)

var repository Repository

type Url struct {
	ID string
	CreatedAt time.Time
	Destination string
}

type Repository interface {
	IdExists(id string) bool
	SearchById(id string) *Url
	SearchByUrl(url string) *Url
	Save(url Url) error
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
func CreateIfNotExists(destination string) (u *Url, isNew bool, err error) {
	if u = repository.SearchByUrl(destination); u != nil {
		return u, false, nil
	}
	if _, err = url.ParseRequestURI(destination); err != nil {
		return nil, false, err
	}
	url := Url{generateID(), time.Now(), destination}
	repository.Save(url)
	return &url, true, nil
}

func Search(id string) *Url {
	return repository.SearchById(id)
}
 
func SetUpRepository(r Repository) {
	repository = r
}

func generateID() string {
	newId := func() string {
		id := make([]byte, size, size)
		for i := range id {
			id[i] = tokens[rand.Intn(len(tokens))]
		}
		return string(id)
	}
	for {
		if id := newId(); !repository.IdExists(id) {
			return id
		}
	}
}