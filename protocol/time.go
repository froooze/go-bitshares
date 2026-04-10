package protocol

import (
	"encoding/binary"
	"encoding/json"
	"time"
)

// Time wraps a UTC timestamp with second precision.
type Time struct {
	time.Time
}

func NewTime(t time.Time) Time {
	return Time{Time: t.UTC().Truncate(time.Second)}
}

func (t Time) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.UTC().Format(time.RFC3339))
}

func (t *Time) UnmarshalJSON(data []byte) error {
	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return err
	}

	t.Time = parsed.UTC()
	return nil
}

// NewTimeFromUnix constructs a Time from Unix seconds.
func NewTimeFromUnix(unix uint32) Time {
	return Time{Time: time.Unix(int64(unix), 0).UTC()}
}

// UnixSeconds returns the UTC Unix timestamp as seconds.
func (t Time) UnixSeconds() uint32 {
	if t.Time.IsZero() {
		return 0
	}
	return uint32(t.UTC().Unix())
}

// MarshalBinary encodes the time as a 32-bit Unix timestamp.
func (t Time) MarshalBinary() ([]byte, error) {
	var buf [4]byte
	binary.LittleEndian.PutUint32(buf[:], t.UnixSeconds())
	return buf[:], nil
}

// MarshalBinaryInto writes the time as a 32-bit Unix timestamp.
func (t Time) MarshalBinaryInto(w *binaryWriter) error {
	w.writeUint32(t.UnixSeconds())
	return nil
}

// UnmarshalBinary decodes a 32-bit Unix timestamp.
func (t *Time) UnmarshalBinary(data []byte) error {
	return t.UnmarshalBinaryFrom(newBinaryReader(data))
}

// UnmarshalBinaryFrom reads a 32-bit Unix timestamp from a binary reader.
func (t *Time) UnmarshalBinaryFrom(r *binaryReader) error {
	v, err := r.readUint32()
	if err != nil {
		return err
	}
	t.Time = time.Unix(int64(v), 0).UTC()
	return nil
}

func readTime(r *binaryReader) (Time, error) {
	var out Time
	return out, out.UnmarshalBinaryFrom(r)
}
