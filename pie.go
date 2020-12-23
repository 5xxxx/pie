package pie

import (
	"github.com/NSObjects/pie/driver"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewClient(db string, options ...*options.ClientOptions) (Client, error) {
	return driver.NewClient(db, options...)
}
