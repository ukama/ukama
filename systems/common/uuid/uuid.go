// This package wraps satori's uuid package "github.com/satori/go.uuid"
// inside ukama's common utilities.

// The wrapper keeps the same interface as the wrappee in order to provide
// perfect substition everywhere, between the two types.

package uuid

import (
	"database/sql/driver"

	suuid "github.com/satori/go.uuid"
)

// Nil represents the nil UUID
var Nil = UUID{}

// Generator provides interface for generating UUIDs.
type Generator interface {
	NewV1() UUID
	NewV2(domain byte) UUID
	NewV3(ns UUID, name string) UUID
	NewV4() UUID
	NewV5(ns UUID, name string) UUID
}

// Equal wraps suuid.Equal().
func Equal(u1 UUID, u2 UUID) bool {
	return suuid.Equal(suuid.UUID(u1), suuid.UUID(u2))
}

// FromBytes wraps suuid.FromBytes().
func FromBytes(input []byte) (u UUID, err error) {
	val, err := suuid.FromBytes(input)

	return UUID(val), err
}

// FromBytesOrNil wraps suuid.FromBytesOrNil().
func FromBytesOrNil(input []byte) UUID {
	return UUID(suuid.FromBytesOrNil(input))
}

// FromString wraps suuid.FromString().
func FromString(input string) (u UUID, err error) {
	val, err := suuid.FromString(input)

	return UUID(val), err
}

// FromStringOrNil wraps suuid.FromStringOrNil().
func FromStringOrNil(input string) UUID {
	return UUID(suuid.FromStringOrNil(input))
}

// Must wraps suuid.Must().
func Must(u UUID, err error) UUID {
	return UUID(suuid.Must(suuid.UUID(u), err))
}

// NewV1 wraps suuid.NewV1().
func NewV1() UUID {
	return UUID(suuid.NewV1())
}

// NewV2 wraps suuid.NewV2().
func NewV2(domain byte) UUID {
	return UUID(suuid.NewV2(domain))
}

// NewV3 wraps suuid.NewV3().
func NewV3(ns UUID, name string) UUID {
	return UUID(suuid.NewV3(suuid.UUID(ns), name))
}

// NewV4 wraps suuid.NewV4().
func NewV4() UUID {
	return UUID(suuid.NewV4())
}

// NewV5 wraps suuid.NewV5().
func NewV5(ns UUID, name string) UUID {
	return UUID(suuid.NewV5(suuid.UUID(ns), name))
}

// NullUUID wraps suuid.NullUUID.
type NullUUID struct {
	UUID  UUID
	Valid bool
}

// Scan wraps suuid.NullUUID.Scan().
func (u *NullUUID) Scan(src interface{}) error {
	if src == nil {
		u.UUID, u.Valid = Nil, false
		return nil
	}

	// Delegate to UUID Scan function
	u.Valid = true
	return u.UUID.Scan(src)
}

// Value wraps suuid.NullUUID.Value().
func (u NullUUID) Value() (driver.Value, error) {
	if !u.Valid {
		return nil, nil
	}
	// Delegate to UUID Value function
	return u.UUID.Value()
}

// UUID wraps suuid.UUID.
type UUID suuid.UUID

// Bytes wraps suuid.UUID.Bytes().
func (u UUID) Bytes() []byte {
	return suuid.UUID(u).Bytes()
}

// MarshalBinary wraps suuid.UUID.MarshalBinary().
func (u UUID) MarshalBinary() (data []byte, err error) {
	return suuid.UUID(u).MarshalBinary()
}

// MarshalText wraps suuid.UUID.MarshalText().
func (u UUID) MarshalText() (text []byte, err error) {
	return suuid.UUID(u).MarshalText()
}

// Scan wraps suuid.UUID.Scan().
func (u *UUID) Scan(src interface{}) error {
	return (*suuid.UUID)(u).Scan(src)
}

// SetVariant wraps suuid.UUID.SetVariant().
func (u *UUID) SetVariant(v byte) {
	(*suuid.UUID)(u).SetVariant(v)
}

// SetVersion wraps suuid.UUID.SetVersion().
func (u *UUID) SetVersion(v byte) {
	(*suuid.UUID)(u).SetVersion(v)
}

// String wraps suuid.UUID.String().
func (u UUID) String() string {
	return suuid.UUID(u).String()
}

// UnmarshalBinary wraps suuid.UUID.UnmarshalBinary().
func (u *UUID) UnmarshalBinary(data []byte) (err error) {
	return (*suuid.UUID)(u).UnmarshalBinary(data)
}

// UnmarshalText wraps suuid.UUID.UnmarshalText().
func (u *UUID) UnmarshalText(text []byte) (err error) {
	return (*suuid.UUID)(u).UnmarshalText(text)
}

// Value wraps suuid.UUID.Value().
func (u UUID) Value() (driver.Value, error) {
	return suuid.UUID(u).Value()
}

// Variant wraps suuid.UUID.Variant().
func (u UUID) Variant() byte {
	return suuid.UUID(u).Variant()
}

// Version wraps suuid.UUID.Version().
func (u UUID) Version() byte {
	return suuid.UUID(u).Version()
}
