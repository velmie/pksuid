package pksuid

import (
	"github.com/segmentio/ksuid"

	"bytes"
	"database/sql/driver"
	"fmt"
	"time"
)

const (
	// KSUIDs are 20 bytes when binary encoded
	ksuidByteLength = 20
	// The length of a KSUID when string (base62) encoded
	ksuidStringEncodedLength = 27
	// The Prefix maximum bytes length
	prefixByteLength = 16
	// PKSUIDs are conjunction of the Prefix and encoded KSUID
	pksuidByteLength = ksuidByteLength + prefixByteLength
)

var (
	errMinSize = fmt.Errorf("pksuid: valid PKSUIDs cannot be less than %v bytes", ksuidByteLength)
	errMaxSize = fmt.Errorf("pksuid: valid PKSUIDs cannot be greater than %v bytes", pksuidByteLength)
	errStrSize = fmt.Errorf(
		"pksuid: valid encoded PKSUIDs cannot be less than %v characters",
		ksuidStringEncodedLength,
	)
	errPrefixSize = fmt.Errorf("pksuid: prefix must be less than %v bytes", prefixByteLength)

	Nil       = PKSUID{}
	NilPrefix = Prefix{}
	nilKSUID  = make([]byte, ksuidByteLength)
)

// Prefix is the extension to the KSUID
type Prefix [prefixByteLength]byte

func (p Prefix) String() string {
	return string(trimAfterNullChar(p[:]))
}

// PKSUIDs are 36 bytes:
//  00-15 byte: arbitrary string prefix
//  KSUID:
//  16-19 byte: uint32 BE UTC timestamp with custom epoch
//  20-35 byte: random "payload"
type PKSUID [pksuidByteLength]byte

// New generates a new PKSUID with the given prefix
func New(prefix Prefix) PKSUID {
	puid := PKSUID{}
	puid.SetPrefix(prefix)
	uid := ksuid.New()
	copy(puid[prefixByteLength:], uid[:])
	return puid
}

// Parse decodes a string-encoded representation of a PKSUID object
func Parse(s string) (PKSUID, error) {
	if len(s) < ksuidStringEncodedLength {
		return Nil, errStrSize
	}

	ksuidPart := s[len(s)-ksuidStringEncodedLength:]

	uid, err := ksuid.Parse(ksuidPart)
	if err != nil {
		return Nil, fmt.Errorf("pksuid: failed to parse encoded KSUID: %s", err.Error())
	}

	prefixPart := s[0 : len(s)-ksuidStringEncodedLength]
	if len(prefixPart) > prefixByteLength {
		return Nil, errPrefixSize
	}

	pksuid := PKSUID{}
	copy(pksuid[prefixByteLength:], uid[:])

	if len(prefixPart) > 0 {
		copy(pksuid[:prefixByteLength], prefixPart)
	}
	return pksuid, nil
}

// FromBytes constructs a PKSUID from a binary representation
func FromBytes(b []byte) (PKSUID, error) {
	var pksuid PKSUID

	if len(b) < ksuidByteLength {
		return Nil, errMinSize
	}
	if len(b) > pksuidByteLength {
		return Nil, errMaxSize
	}
	copy(pksuid[prefixByteLength:], b[len(b)-ksuidByteLength:])
	if len(b) > ksuidByteLength {
		copy(pksuid[:prefixByteLength], b[:len(b)-ksuidByteLength])
	}
	return pksuid, nil
}

// SetPrefix set prefix bytes
func (p *PKSUID) SetPrefix(prefix Prefix) {
	copy(p[:prefixByteLength], prefix[:])
}

// IsNil returns true if the KSUID part is "nil"
// prefix is ignored only the KSUID part matter
func (p PKSUID) IsNil() bool {
	return bytes.Equal(p[prefixByteLength:], nilKSUID)
}

// IsNilWithPrefix returns nil if prefix and the KSUID parts nil
func (p PKSUID) IsNilWithPrefix() bool {
	return p == Nil
}

// String returns string representation of the PKSUID
func (p PKSUID) String() string {
	return string(trimAfterNullChar(p.PrefixBytes())) + p.KSUID().String()
}

// Bytes returns raw byte representation of PKSUID
func (p PKSUID) Bytes() []byte {
	return p[:]
}

//
func (p PKSUID) MarshalText() ([]byte, error) {
	return []byte(p.String()), nil
}

func (p PKSUID) MarshalBinary() ([]byte, error) {
	return p.Bytes(), nil
}

func (p *PKSUID) UnmarshalText(b []byte) error {
	id, err := Parse(string(b))
	if err != nil {
		return err
	}
	*p = id
	return nil
}

func (p *PKSUID) UnmarshalBinary(b []byte) error {
	id, err := FromBytes(b)
	if err != nil {
		return err
	}
	*p = id
	return nil
}

// Value converts the PKSUID into a SQL driver value which can be used to
// directly use the PKSUID as parameter to a SQL query.
func (p PKSUID) Value() (driver.Value, error) {
	if p.IsNil() {
		return nil, nil
	}
	return p.String(), nil
}

// Scan implements the sql.Scanner interface. It supports converting from
// string, []byte, or nil into a PKSUID value. Attempting to convert from
// another type will return an error.
func (p *PKSUID) Scan(src interface{}) error {
	switch v := src.(type) {
	case nil:
		*p = Nil
		return nil
	case []byte:
		if len(v) > pksuidByteLength {
			return p.UnmarshalText(v)
		}
		if len(v) >= ksuidStringEncodedLength && isBase62Bytes(v[len(v)-ksuidStringEncodedLength:]) {
			return p.UnmarshalText(v)
		}
		return p.UnmarshalBinary(v)
	case string:
		return p.UnmarshalText([]byte(v))
	default:
		return fmt.Errorf("scan: unable to scan type %T into PKSUID", v)
	}
}

// Time retrieves the timestamp portion of the ID as a Time object
func (p PKSUID) Time() time.Time {
	return p.KSUID().Time()
}

// Timestamp retrieves the timestamp portion of the ID as a bare integer which is uncorrected
// for KSUID's special epoch.
func (p PKSUID) Timestamp() uint32 {
	return p.KSUID().Timestamp()
}

// Payload retrieves the 16-byte random payload of the KSUID without the timestamp
func (p PKSUID) Payload() []byte {
	return p.KSUID().Payload()
}

// Prefix fetches a prefix part of the PKSUID
func (p PKSUID) Prefix() Prefix {
	prefx := Prefix{}
	copy(prefx[:], p.PrefixBytes())
	return prefx
}

// ID returns encoded to string KSUID
func (p PKSUID) ID() string {
	return p.KSUID().String()
}

// IDBytes returns raw KSUID bytes
func (p PKSUID) IDBytes() []byte {
	return p.ksuidBytes()
}

// PrefixBytes raw byte representation of the prefix part
func (p PKSUID) PrefixBytes() []byte {
	return p[:prefixByteLength]
}

// KSUID retrieves a native KSUID object
func (p PKSUID) KSUID() ksuid.KSUID {
	uid := ksuid.KSUID{}
	copy(uid[:], p.ksuidBytes())
	return uid
}

func (p PKSUID) ksuidBytes() []byte {
	return p[prefixByteLength:]
}

// trimAfterNullChar trims everything from the first occurrence of the null character
// including the character
func trimAfterNullChar(b []byte) []byte {
	i := bytes.Index(b, []byte("\x00"))
	if i != -1 {
		return b[:i]
	}
	return b
}
