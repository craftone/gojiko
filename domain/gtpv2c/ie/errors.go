package ie

import (
	"fmt"
)

type UnknownIEError struct {
	ieTypeValue ieTypeNum
	instance    byte
}

func (e *UnknownIEError) Error() string {
	return fmt.Sprintf("Unknown IE type : TypeValue:%d, Instance:%d", e.ieTypeValue, e.instance)
}
