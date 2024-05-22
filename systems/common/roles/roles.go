package roles

type RoleType uint8

const (
	TYPE_INVALID RoleType = 0
	TYPE_OWNER   RoleType = 1
	TYPE_ADMIN   RoleType = 2
	TYPE_VENDOR  RoleType = 3
	TYPE_USERS   RoleType = 4
)

func (e *RoleType) Scan(value interface{}) error {
	*e = RoleType(uint8(value.(int64)))

	return nil
}

func (e RoleType) Value() (uint8, error) {
	return uint8(e), nil
}
