package pb

import (
	"encoding/gob"
)

func init() {
	gob.Register(BroadCast{})
	gob.Register(BroadCast_Content{})
	gob.Register(isBroadCast_Data{})
}
