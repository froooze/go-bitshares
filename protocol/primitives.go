package protocol

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/froooze/go-bitshares/ecc"
)

// PublicKey stores a BitShares public key string.
type PublicKey string

// ParsePublicKey validates and normalizes a BitShares public key string.
func ParsePublicKey(value string) (PublicKey, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", fmt.Errorf("empty public key")
	}
	if _, err := ecc.PublicKeyFromString(value); err != nil {
		return "", err
	}
	return PublicKey(value), nil
}

// MustPublicKey panics when the string is not a valid public key.
func MustPublicKey(value string) PublicKey {
	out, err := ParsePublicKey(value)
	if err != nil {
		panic(err)
	}
	return out
}

// String returns the canonical key string.
func (k PublicKey) String() string { return string(k) }

func (k PublicKey) MarshalText() ([]byte, error) { return []byte(k.String()), nil }

func (k *PublicKey) UnmarshalText(text []byte) error {
	out, err := ParsePublicKey(string(text))
	if err != nil {
		return err
	}
	*k = out
	return nil
}

func (k PublicKey) MarshalJSON() ([]byte, error) { return json.Marshal(k.String()) }

func (k *PublicKey) UnmarshalJSON(data []byte) error {
	value, err := unquote(string(data))
	if err != nil {
		return err
	}
	out, err := ParsePublicKey(value)
	if err != nil {
		return err
	}
	*k = out
	return nil
}

func (k PublicKey) MarshalBinaryInto(w *binaryWriter) error {
	pub, err := ecc.PublicKeyFromString(k.String())
	if err != nil {
		return err
	}
	_, err = w.Write(pub.Bytes())
	return err
}

func (k *PublicKey) UnmarshalBinaryFrom(r *binaryReader) error {
	buf := make([]byte, 33)
	if _, err := r.Read(buf); err != nil {
		return err
	}
	pub, err := ecc.PublicKeyFromBytes(buf)
	if err != nil {
		return err
	}
	*k = PublicKey(pub.String())
	return nil
}

// Address stores a BitShares address string.
type Address string

// ParseAddress validates and normalizes a BitShares address string.
func ParseAddress(value string) (Address, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", fmt.Errorf("empty address")
	}
	if _, err := ecc.AddressFromString(value); err != nil {
		return "", err
	}
	return Address(value), nil
}

// MustAddress panics when the string is not a valid BitShares address.
func MustAddress(value string) Address {
	out, err := ParseAddress(value)
	if err != nil {
		panic(err)
	}
	return out
}

func (a Address) String() string { return string(a) }

func (a Address) MarshalText() ([]byte, error) { return []byte(a.String()), nil }

func (a *Address) UnmarshalText(text []byte) error {
	out, err := ParseAddress(string(text))
	if err != nil {
		return err
	}
	*a = out
	return nil
}

func (a Address) MarshalJSON() ([]byte, error) { return json.Marshal(a.String()) }

func (a *Address) UnmarshalJSON(data []byte) error {
	value, err := unquote(string(data))
	if err != nil {
		return err
	}
	out, err := ParseAddress(value)
	if err != nil {
		return err
	}
	*a = out
	return nil
}

func (a Address) MarshalBinaryInto(w *binaryWriter) error {
	addr, err := ecc.AddressFromString(a.String())
	if err != nil {
		return err
	}
	_, err = w.Write(addr.Bytes())
	return err
}

func (a *Address) UnmarshalBinaryFrom(r *binaryReader) error {
	buf := make([]byte, 20)
	if _, err := r.Read(buf); err != nil {
		return err
	}
	var raw ecc.Address
	copy(raw[:], buf)
	*a = Address(raw.String())
	return nil
}

// VoteID identifies a vote target.
type VoteID struct {
	Type uint8  `json:"type"`
	ID   uint32 `json:"id"`
}

func ParseVoteID(value string) (VoteID, error) {
	var out VoteID
	parts := strings.Split(strings.TrimSpace(value), ":")
	if len(parts) != 2 {
		return out, fmt.Errorf("invalid vote id %q", value)
	}
	var err error
	if out.Type, err = parseUint8(parts[0]); err != nil {
		return out, err
	}
	if out.ID, err = parseUint32(parts[1]); err != nil {
		return out, err
	}
	return out, nil
}

func MustVoteID(value string) VoteID {
	out, err := ParseVoteID(value)
	if err != nil {
		panic(err)
	}
	return out
}

func (v VoteID) String() string {
	return fmt.Sprintf("%d:%d", v.Type, v.ID)
}

func (v VoteID) MarshalText() ([]byte, error) { return []byte(v.String()), nil }

func (v *VoteID) UnmarshalText(text []byte) error {
	out, err := ParseVoteID(string(text))
	if err != nil {
		return err
	}
	*v = out
	return nil
}

func (v VoteID) MarshalJSON() ([]byte, error) { return json.Marshal(v.String()) }

func (v *VoteID) UnmarshalJSON(data []byte) error {
	value, err := unquote(string(data))
	if err != nil {
		return err
	}
	out, err := ParseVoteID(value)
	if err != nil {
		return err
	}
	*v = out
	return nil
}

func (v VoteID) MarshalBinaryInto(w *binaryWriter) error {
	w.writeUint32(uint32(v.ID<<8) | uint32(v.Type))
	return nil
}

func (v *VoteID) UnmarshalBinaryFrom(r *binaryReader) error {
	raw, err := r.readUint32()
	if err != nil {
		return err
	}
	v.Type = uint8(raw & 0xff)
	v.ID = raw >> 8
	return nil
}

func parseUint8(value string) (uint8, error) {
	var n uint64
	var err error
	if n, err = parseUint64(value); err != nil {
		return 0, err
	}
	if n > 0xff {
		return 0, fmt.Errorf("value out of range: %s", value)
	}
	return uint8(n), nil
}

func parseUint32(value string) (uint32, error) {
	var n uint64
	var err error
	if n, err = parseUint64(value); err != nil {
		return 0, err
	}
	if n > 0xffffffff {
		return 0, fmt.Errorf("value out of range: %s", value)
	}
	return uint32(n), nil
}

func parseUint64(value string) (uint64, error) {
	var out uint64
	for _, c := range strings.TrimSpace(value) {
		if c < '0' || c > '9' {
			return 0, fmt.Errorf("invalid numeric value %q", value)
		}
		out = out*10 + uint64(c-'0')
	}
	return out, nil
}

func compareStringSlices(a, b []string) int {
	switch {
	case len(a) == 0 && len(b) == 0:
		return 0
	case len(a) == 0:
		return -1
	case len(b) == 0:
		return 1
	}
	limit := len(a)
	if len(b) < limit {
		limit = len(b)
	}
	for i := 0; i < limit; i++ {
		if a[i] == b[i] {
			continue
		}
		if a[i] < b[i] {
			return -1
		}
		return 1
	}
	switch {
	case len(a) < len(b):
		return -1
	case len(a) > len(b):
		return 1
	default:
		return 0
	}
}

func boolToUint8(v bool) uint8 {
	if v {
		return 1
	}
	return 0
}

func uint8ToBool(v uint8) bool {
	return v != 0
}

func cloneRawMessages(values []json.RawMessage) []json.RawMessage {
	if len(values) == 0 {
		return nil
	}
	out := make([]json.RawMessage, len(values))
	for i, v := range values {
		out[i] = append(json.RawMessage(nil), v...)
	}
	return out
}

func normalizeRawMessage(value json.RawMessage) json.RawMessage {
	if len(value) == 0 {
		return nil
	}
	return append(json.RawMessage(nil), value...)
}

func writeVarBytes(w *binaryWriter, value []byte) {
	w.writeBytes(value)
}

func readVarBytes(r *binaryReader) ([]byte, error) {
	return r.readBytes()
}

func writeFixedBytes(w *binaryWriter, value []byte, size int) error {
	if len(value) != size {
		return fmt.Errorf("expected %d bytes, got %d", size, len(value))
	}
	_, err := w.Write(value)
	return err
}

func readFixedBytes(r *binaryReader, size int) ([]byte, error) {
	out := make([]byte, size)
	if _, err := r.Read(out); err != nil {
		return nil, err
	}
	return out, nil
}

func writeBool(w *binaryWriter, value bool) {
	w.writeUint8(boolToUint8(value))
}

func readBool(r *binaryReader) (bool, error) {
	v, err := r.readUint8()
	if err != nil {
		return false, err
	}
	return uint8ToBool(v), nil
}

func writeOptionalJSON(w *binaryWriter, raw json.RawMessage) error {
	if len(raw) == 0 || strings.EqualFold(string(raw), "null") {
		w.writeUint8(0)
		return nil
	}
	w.writeUint8(1)
	w.writeBytes(raw)
	return nil
}

func readOptionalJSON(r *binaryReader) (json.RawMessage, error) {
	present, err := r.readUint8()
	if err != nil {
		return nil, err
	}
	if present == 0 {
		return nil, nil
	}
	raw, err := r.readBytes()
	if err != nil {
		return nil, err
	}
	return append(json.RawMessage(nil), raw...), nil
}

func writeCountAndBytes(w *binaryWriter, raw json.RawMessage) error {
	return writeOptionalJSON(w, raw)
}

func readCountAndBytes(r *binaryReader) (json.RawMessage, error) {
	return readOptionalJSON(r)
}

func binaryWriteString(w *binaryWriter, value string) {
	w.writeString(value)
}

func binaryReadString(r *binaryReader) (string, error) {
	return r.readString()
}

func int64Ptr(v int64) *int64 { return &v }

func uint16Ptr(v uint16) *uint16 { return &v }

func uint32Ptr(v uint32) *uint32 { return &v }

func uint8Ptr(v uint8) *uint8 { return &v }

func sortVoteIDs(values []VoteID) []VoteID {
	if len(values) < 2 {
		return values
	}
	out := append([]VoteID(nil), values...)
	for i := 1; i < len(out); i++ {
		j := i
		for j > 0 && compareVoteID(out[j-1], out[j]) > 0 {
			out[j-1], out[j] = out[j], out[j-1]
			j--
		}
	}
	return out
}

func compareVoteID(a, b VoteID) int {
	if a.Type != b.Type {
		if a.Type < b.Type {
			return -1
		}
		return 1
	}
	if a.ID < b.ID {
		return -1
	}
	if a.ID > b.ID {
		return 1
	}
	return 0
}

func sortObjectIDs(values []ObjectID) []ObjectID {
	if len(values) < 2 {
		return values
	}
	out := append([]ObjectID(nil), values...)
	for i := 1; i < len(out); i++ {
		j := i
		for j > 0 && compareObjectID(out[j-1], out[j]) > 0 {
			out[j-1], out[j] = out[j], out[j-1]
			j--
		}
	}
	return out
}

func uniqueSortedObjectIDs(values []ObjectID) []ObjectID {
	sorted := sortObjectIDs(values)
	if len(sorted) < 2 {
		return sorted
	}
	out := sorted[:1]
	for _, value := range sorted[1:] {
		if compareObjectID(out[len(out)-1], value) != 0 {
			out = append(out, value)
		}
	}
	return out
}

func compareObjectID(a, b ObjectID) int {
	if a.Space != b.Space {
		if a.Space < b.Space {
			return -1
		}
		return 1
	}
	if a.Type != b.Type {
		if a.Type < b.Type {
			return -1
		}
		return 1
	}
	if a.ID < b.ID {
		return -1
	}
	if a.ID > b.ID {
		return 1
	}
	return 0
}

func compareStrings(a, b string) int {
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	return 0
}

func sortStrings(values []string) []string {
	if len(values) < 2 {
		return values
	}
	out := append([]string(nil), values...)
	for i := 1; i < len(out); i++ {
		j := i
		for j > 0 && compareStrings(out[j-1], out[j]) > 0 {
			out[j-1], out[j] = out[j], out[j-1]
			j--
		}
	}
	return out
}

func sortPublicKeys(values []PublicKey) []PublicKey {
	if len(values) < 2 {
		return values
	}
	out := append([]PublicKey(nil), values...)
	for i := 1; i < len(out); i++ {
		j := i
		for j > 0 && compareStrings(out[j-1].String(), out[j].String()) > 0 {
			out[j-1], out[j] = out[j], out[j-1]
			j--
		}
	}
	return out
}

func uniqueSortedPublicKeys(values []PublicKey) []PublicKey {
	sorted := sortPublicKeys(values)
	if len(sorted) < 2 {
		return sorted
	}
	out := sorted[:1]
	for _, value := range sorted[1:] {
		if compareStrings(out[len(out)-1].String(), value.String()) != 0 {
			out = append(out, value)
		}
	}
	return out
}

func sortAddresses(values []Address) []Address {
	if len(values) < 2 {
		return values
	}
	out := append([]Address(nil), values...)
	for i := 1; i < len(out); i++ {
		j := i
		for j > 0 && compareStrings(out[j-1].String(), out[j].String()) > 0 {
			out[j-1], out[j] = out[j], out[j-1]
			j--
		}
	}
	return out
}

func uniqueSortedAddresses(values []Address) []Address {
	sorted := sortAddresses(values)
	if len(sorted) < 2 {
		return sorted
	}
	out := sorted[:1]
	for _, value := range sorted[1:] {
		if compareStrings(out[len(out)-1].String(), value.String()) != 0 {
			out = append(out, value)
		}
	}
	return out
}

func writeObjectIDSet(w *binaryWriter, values []ObjectID) error {
	values = uniqueSortedObjectIDs(values)
	w.writeVarUint64(uint64(len(values)))
	for _, value := range values {
		if err := value.MarshalBinaryInto(w); err != nil {
			return err
		}
	}
	return nil
}

func readObjectIDSet(r *binaryReader) ([]ObjectID, error) {
	count, err := r.readVarUint64()
	if err != nil {
		return nil, err
	}
	out := make([]ObjectID, 0, count)
	for i := uint64(0); i < count; i++ {
		value, err := readObjectID(r)
		if err != nil {
			return nil, err
		}
		out = append(out, value)
	}
	return out, nil
}

func writePublicKeySet(w *binaryWriter, values []PublicKey) error {
	values = uniqueSortedPublicKeys(values)
	w.writeVarUint64(uint64(len(values)))
	for _, value := range values {
		if err := value.MarshalBinaryInto(w); err != nil {
			return err
		}
	}
	return nil
}

func readPublicKeySet(r *binaryReader) ([]PublicKey, error) {
	count, err := r.readVarUint64()
	if err != nil {
		return nil, err
	}
	out := make([]PublicKey, 0, count)
	for i := uint64(0); i < count; i++ {
		value, err := readPublicKey(r)
		if err != nil {
			return nil, err
		}
		out = append(out, value)
	}
	return out, nil
}

func writeAddressSet(w *binaryWriter, values []Address) error {
	values = uniqueSortedAddresses(values)
	w.writeVarUint64(uint64(len(values)))
	for _, value := range values {
		if err := value.MarshalBinaryInto(w); err != nil {
			return err
		}
	}
	return nil
}

func readAddressSet(r *binaryReader) ([]Address, error) {
	count, err := r.readVarUint64()
	if err != nil {
		return nil, err
	}
	out := make([]Address, 0, count)
	for i := uint64(0); i < count; i++ {
		value, err := readAddress(r)
		if err != nil {
			return nil, err
		}
		out = append(out, value)
	}
	return out, nil
}

func writeVoteIDSet(w *binaryWriter, values []VoteID) error {
	values = sortVoteIDs(values)
	w.writeVarUint64(uint64(len(values)))
	for _, value := range values {
		if err := value.MarshalBinaryInto(w); err != nil {
			return err
		}
	}
	return nil
}

func readVoteIDSet(r *binaryReader) ([]VoteID, error) {
	count, err := r.readVarUint64()
	if err != nil {
		return nil, err
	}
	out := make([]VoteID, 0, count)
	for i := uint64(0); i < count; i++ {
		value, err := readVoteID(r)
		if err != nil {
			return nil, err
		}
		out = append(out, value)
	}
	return out, nil
}

func writeOptionalObjectID(w *binaryWriter, value *ObjectID) error {
	if value == nil {
		w.writeUint8(0)
		return nil
	}
	w.writeUint8(1)
	return value.MarshalBinaryInto(w)
}

func readOptionalObjectID(r *binaryReader) (*ObjectID, error) {
	present, err := r.readUint8()
	if err != nil {
		return nil, err
	}
	if present == 0 {
		return nil, nil
	}
	value, err := readObjectID(r)
	if err != nil {
		return nil, err
	}
	return &value, nil
}

func writeOptionalPublicKey(w *binaryWriter, value *PublicKey) error {
	if value == nil {
		w.writeUint8(0)
		return nil
	}
	w.writeUint8(1)
	return value.MarshalBinaryInto(w)
}

func readOptionalPublicKey(r *binaryReader) (*PublicKey, error) {
	present, err := r.readUint8()
	if err != nil {
		return nil, err
	}
	if present == 0 {
		return nil, nil
	}
	value, err := readPublicKey(r)
	if err != nil {
		return nil, err
	}
	return &value, nil
}

func writeOptionalAddress(w *binaryWriter, value *Address) error {
	if value == nil {
		w.writeUint8(0)
		return nil
	}
	w.writeUint8(1)
	return value.MarshalBinaryInto(w)
}

func readOptionalAddress(r *binaryReader) (*Address, error) {
	present, err := r.readUint8()
	if err != nil {
		return nil, err
	}
	if present == 0 {
		return nil, nil
	}
	value, err := readAddress(r)
	if err != nil {
		return nil, err
	}
	return &value, nil
}

func readPublicKey(r *binaryReader) (PublicKey, error) {
	var out PublicKey
	return out, out.UnmarshalBinaryFrom(r)
}

func readAddress(r *binaryReader) (Address, error) {
	var out Address
	return out, out.UnmarshalBinaryFrom(r)
}

func readVoteID(r *binaryReader) (VoteID, error) {
	var out VoteID
	return out, out.UnmarshalBinaryFrom(r)
}

// assertUint helpers keep JSON parsing error messages small and direct.
