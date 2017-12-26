package pb

import (
	"encoding/gob"
)

func init() {
	gob.Register(Position{})
	// gob.Register(BroadCast_Content{})
	// gob.Register(isBroadCast_Data{})
}
