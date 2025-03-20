package netw

import (
	"context"
)

// Validator is an object that can be validated.
//
// The Valid method takes a context (which is optional) and returns a map.
// If there is a problem with a field, its name is used as the key, and a human-readable explanation of the issue is set as the value.
//
// The method can do whatever it needs to validate the fields of the struct. For example, it can check to make sure:
//     Required fields are not empty
//     Strings with a specific format (like email) are correct
//     Numbers are within an acceptable range
//
// If you need to do anything more complicated, like check the field in a database, that should happen elsewhere;
//

// Ver ejemplo de uso en DecodeValid() en serialize.go
type Validator interface {
	// Valid checks the object and returns any problems. If len(problems) == 0 then the object is valid.
	Valid(ctx context.Context) (problems map[string]string)
}
