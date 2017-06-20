package permissions

import (
	"github.com/itsmontoya/turtle"
	"github.com/missionMeteora/toolkit/errors"
)

const (
	// ErrInvalidType is returned when an invalid type is stored for an id
	ErrInvalidType = errors.Error("invalid type")
	// ErrInvalidPermissions is returned when an invalid permissions value is attempted to be set
	ErrInvalidPermissions = errors.Error("invalid permissions, please see constant block for reference")
	// ErrPermissionsUnchanged is returned when matching permissions are set for a resource
	ErrPermissionsUnchanged = errors.Error("permissions match, unchanged")
)

const (
	// PermissionNone represents a zero value, no permissions
	PermissionNone uint8 = iota
	// PermissionRead represents read-only permissions
	PermissionRead
	// PermissionWrite represents write-only permissions
	PermissionWrite
	// PermissionReadWrite represents read/write permissions
	PermissionReadWrite
)

const (
	// ActionNone represents a zero value, no action
	ActionNone uint8 = iota
	// ActionRead represents a reading action
	ActionRead
	// ActionWrite represents a writing action
	ActionWrite
)

// New will return a new instance of Permissions
func New(path string) (pp *Permissions, err error) {
	var p Permissions
	if p.db, err = turtle.New("permissions", path, marshal, unmarshal); err != nil {
		return
	}

	pp = &p
	return
}

// Permissions manages permissions
type Permissions struct {
	db *turtle.Turtle
}

func (p *Permissions) get(txn turtle.Txn, id string) (rm resourceMap, err error) {
	var (
		val turtle.Value
		ok  bool
	)

	if val, err = txn.Get(id); err != nil {
		return
	}

	if rm, ok = val.(resourceMap); !ok {
		err = ErrInvalidType
		return
	}

	return
}

// Get will get the permissions for a given group for a resource id
func (p *Permissions) Get(id, group string) (permissions uint8) {
	var rm resourceMap

	p.db.Read(func(txn turtle.Txn) (err error) {
		if rm, err = p.get(txn, id); err != nil {
			return
		}

		permissions, _ = rm.Get(group)
		return
	})

	return
}

// Set will set the permissions for a given group for a resource id
func (p *Permissions) Set(id, group string, permissions uint8) (err error) {
	var rm resourceMap
	if !isValidPermissions(permissions) {
		return ErrInvalidPermissions
	}

	return p.db.Update(func(txn turtle.Txn) (err error) {
		if rm, err = p.get(txn, id); err != nil {
			rm = make(resourceMap)
			err = nil
		}

		if !rm.Set(group, permissions) {
			return ErrPermissionsUnchanged
		}

		txn.Put(id, rm)
		return
	})
}

// Can will return if a set of groups can perform a given action
func (p *Permissions) Can(id string, action uint8, groups ...string) (can bool) {
	var (
		rm  resourceMap
		err error
	)

	if err = p.db.Read(func(txn turtle.Txn) (err error) {
		if rm, err = p.get(txn, id); err != nil {
			return
		}

		switch action {
		case ActionNone:
			can = true

		case ActionRead:
			for _, group := range groups {
				if perm, _ := rm.Get(group); perm == PermissionRead || perm == PermissionReadWrite {
					can = true
					break
				}
			}

		case ActionWrite:
			for _, group := range groups {
				perm, _ := rm.Get(group)

				if perm == PermissionWrite || perm == PermissionReadWrite {
					can = true
					break
				}
			}
		}

		return
	}); err != nil {
		return
	}

	return
}

// Close will close permissions
func (p *Permissions) Close() (err error) {
	return p.db.Close()
}
