package permissions

import (
	"encoding/json"

	"github.com/itsmontoya/turtle"
)

func marshal(val turtle.Value) (b []byte, err error) {
	var (
		rm resourceMap
		ok bool
	)

	if rm, ok = val.(resourceMap); !ok {
		err = ErrInvalidType
		return
	}

	return json.Marshal(rm)
}

func unmarshal(b []byte) (val turtle.Value, err error) {
	var (
		rm resourceMap
	)

	if err = json.Unmarshal(b, &rm); err != nil {
		return
	}

	val = rm
	return
}

func isValidPermissions(permissions uint8) bool {
	switch permissions {
	case PermissionNone:
	case PermissionRead:
	case PermissionWrite:
	case PermissionReadWrite:
	default:
		return false
	}

	return true
}
