package url
import (
	"time"
)

type URL struct {
	ID string
	CreatedAt time.time
	Destination string
}