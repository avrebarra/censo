package censo

import "fmt"

// ErrNotCensorable is returned when censo is used to censor unsupported or
// unknown data types.
var ErrNotCensorable = fmt.Errorf("type not censorable")
