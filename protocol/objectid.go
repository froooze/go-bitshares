package protocol

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// ObjectID identifies BitShares protocol objects in the form space.type.instance.
type ObjectID struct {
	Space uint64
	Type  uint64
	ID    uint64
}

func (o ObjectID) String() string {
	return fmt.Sprintf("%d.%d.%d", o.Space, o.Type, o.ID)
}

func (o ObjectID) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.String())
}

func (o ObjectID) MarshalText() ([]byte, error) {
	return []byte(o.String()), nil
}

func (o *ObjectID) UnmarshalJSON(data []byte) error {
	value, err := unquote(string(data))
	if err != nil {
		return err
	}

	parsed, err := ParseObjectID(value)
	if err != nil {
		return err
	}

	*o = parsed
	return nil
}

func (o *ObjectID) UnmarshalText(text []byte) error {
	parsed, err := ParseObjectID(string(text))
	if err != nil {
		return err
	}
	*o = parsed
	return nil
}

// ParseObjectID parses a string in space.type.instance form.
func ParseObjectID(value string) (ObjectID, error) {
	var out ObjectID

	parts := strings.Split(value, ".")
	if len(parts) != 3 {
		return out, fmt.Errorf("invalid object id %q", value)
	}

	space, err := strconv.ParseUint(parts[0], 10, 64)
	if err != nil {
		return out, fmt.Errorf("invalid object id space %q: %w", value, err)
	}

	typ, err := strconv.ParseUint(parts[1], 10, 64)
	if err != nil {
		return out, fmt.Errorf("invalid object id type %q: %w", value, err)
	}

	id, err := strconv.ParseUint(parts[2], 10, 64)
	if err != nil {
		return out, fmt.Errorf("invalid object id instance %q: %w", value, err)
	}

	out.Space = space
	out.Type = typ
	out.ID = id
	return out, nil
}

func MustParseObjectID(value string) ObjectID {
	out, err := ParseObjectID(value)
	if err != nil {
		panic(err)
	}

	return out
}

// Uint64 packs the object id into the BitShares wire-format integer.
func (o ObjectID) Uint64() uint64 {
	return (o.Space << 56) | (o.Type << 48) | (o.ID & 0x0000ffffffffffff)
}

// MarshalBinary encodes the object id in the BitShares wire format.
func (o ObjectID) MarshalBinary() ([]byte, error) {
	var buf [8]byte
	binary.LittleEndian.PutUint64(buf[:], o.Uint64())
	return buf[:], nil
}

// MarshalBinaryInto writes the object id to a binary writer.
func (o ObjectID) MarshalBinaryInto(w *binaryWriter) error {
	w.writeUint64(o.Uint64())
	return nil
}

// UnmarshalBinary decodes the BitShares wire-format object id.
func (o *ObjectID) UnmarshalBinary(data []byte) error {
	return o.UnmarshalBinaryFrom(newBinaryReader(data))
}

// UnmarshalBinaryFrom reads an object id from a binary reader.
func (o *ObjectID) UnmarshalBinaryFrom(r *binaryReader) error {
	v, err := r.readUint64()
	if err != nil {
		return err
	}
	o.Space = v >> 56
	o.Type = (v >> 48) & 0xff
	o.ID = v & 0x0000ffffffffffff
	return nil
}

func readObjectID(r *binaryReader) (ObjectID, error) {
	var out ObjectID
	return out, out.UnmarshalBinaryFrom(r)
}

func unquote(value string) (string, error) {
	if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
		return strconv.Unquote(value)
	}
	return value, nil
}
