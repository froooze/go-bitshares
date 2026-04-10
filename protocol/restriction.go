package protocol

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// VariantAssertArgument mirrors restriction::variant_assert_argument_type in core.
type VariantAssertArgument struct {
	Tag          int64         `json:"tag"`
	Restrictions []Restriction `json:"restrictions"`
}

// RestrictionArgument mirrors graphene::protocol::restriction::argument_type.
type RestrictionArgument struct {
	Kind            uint16
	Bool            *bool
	Int64           *int64
	String          *string
	Time            *Time
	PublicKey       *PublicKey
	SHA256          *string
	ObjectID        *ObjectID
	BoolSet         []bool
	Int64Set        []int64
	StringSet       []string
	TimeSet         []Time
	PublicKeySet    []PublicKey
	SHA256Set       []string
	ObjectIDSet     []ObjectID
	Restrictions    []Restriction
	RestrictionSets [][]Restriction
	VariantAssert   *VariantAssertArgument
}

func (a RestrictionArgument) MarshalJSON() ([]byte, error) {
	switch a.Kind {
	case 0:
		return json.Marshal([]any{uint16(0), struct{}{}})
	case 1:
		if a.Bool == nil {
			return nil, fmt.Errorf("missing restriction bool argument")
		}
		return json.Marshal([]any{uint16(1), *a.Bool})
	case 2:
		if a.Int64 == nil {
			return nil, fmt.Errorf("missing restriction int64 argument")
		}
		return json.Marshal([]any{uint16(2), *a.Int64})
	case 3:
		if a.String == nil {
			return nil, fmt.Errorf("missing restriction string argument")
		}
		return json.Marshal([]any{uint16(3), *a.String})
	case 4:
		if a.Time == nil {
			return nil, fmt.Errorf("missing restriction time argument")
		}
		return json.Marshal([]any{uint16(4), a.Time})
	case 5:
		if a.PublicKey == nil {
			return nil, fmt.Errorf("missing restriction public key argument")
		}
		return json.Marshal([]any{uint16(5), a.PublicKey})
	case 6:
		if a.SHA256 == nil {
			return nil, fmt.Errorf("missing restriction sha256 argument")
		}
		return json.Marshal([]any{uint16(6), *a.SHA256})
	case 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19:
		if a.ObjectID == nil {
			return nil, fmt.Errorf("missing restriction object id argument")
		}
		return json.Marshal([]any{a.Kind, a.ObjectID})
	case 20:
		return json.Marshal([]any{uint16(20), a.BoolSet})
	case 21:
		return json.Marshal([]any{uint16(21), a.Int64Set})
	case 22:
		return json.Marshal([]any{uint16(22), a.StringSet})
	case 23:
		return json.Marshal([]any{uint16(23), a.TimeSet})
	case 24:
		return json.Marshal([]any{uint16(24), a.PublicKeySet})
	case 25:
		return json.Marshal([]any{uint16(25), a.SHA256Set})
	case 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38:
		return json.Marshal([]any{a.Kind, a.ObjectIDSet})
	case 39:
		return json.Marshal([]any{uint16(39), a.Restrictions})
	case 40:
		return json.Marshal([]any{uint16(40), a.RestrictionSets})
	case 41:
		if a.VariantAssert == nil {
			return nil, fmt.Errorf("missing restriction variant_assert argument")
		}
		return json.Marshal([]any{uint16(41), a.VariantAssert})
	default:
		return nil, fmt.Errorf("unsupported restriction argument type %d", a.Kind)
	}
}

func (a *RestrictionArgument) UnmarshalJSON(data []byte) error {
	var body []json.RawMessage
	if err := json.Unmarshal(data, &body); err != nil {
		return err
	}
	if len(body) != 2 {
		return fmt.Errorf("invalid restriction argument")
	}
	var kind uint16
	if err := json.Unmarshal(body[0], &kind); err != nil {
		return err
	}
	a.Kind = kind
	switch kind {
	case 0:
	case 1:
		var value bool
		if err := json.Unmarshal(body[1], &value); err != nil {
			return err
		}
		a.Bool = &value
	case 2:
		var value int64
		if err := json.Unmarshal(body[1], &value); err != nil {
			return err
		}
		a.Int64 = &value
	case 3:
		var value string
		if err := json.Unmarshal(body[1], &value); err != nil {
			return err
		}
		a.String = &value
	case 4:
		var value Time
		if err := json.Unmarshal(body[1], &value); err != nil {
			return err
		}
		a.Time = &value
	case 5:
		var value PublicKey
		if err := json.Unmarshal(body[1], &value); err != nil {
			return err
		}
		a.PublicKey = &value
	case 6:
		var value string
		if err := json.Unmarshal(body[1], &value); err != nil {
			return err
		}
		a.SHA256 = &value
	case 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19:
		var value ObjectID
		if err := json.Unmarshal(body[1], &value); err != nil {
			return err
		}
		a.ObjectID = &value
	case 20:
		return json.Unmarshal(body[1], &a.BoolSet)
	case 21:
		return json.Unmarshal(body[1], &a.Int64Set)
	case 22:
		return json.Unmarshal(body[1], &a.StringSet)
	case 23:
		return json.Unmarshal(body[1], &a.TimeSet)
	case 24:
		return json.Unmarshal(body[1], &a.PublicKeySet)
	case 25:
		return json.Unmarshal(body[1], &a.SHA256Set)
	case 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38:
		return json.Unmarshal(body[1], &a.ObjectIDSet)
	case 39:
		return json.Unmarshal(body[1], &a.Restrictions)
	case 40:
		return json.Unmarshal(body[1], &a.RestrictionSets)
	case 41:
		var value VariantAssertArgument
		if err := json.Unmarshal(body[1], &value); err != nil {
			return err
		}
		a.VariantAssert = &value
	default:
		return fmt.Errorf("unsupported restriction argument type %d", kind)
	}
	return nil
}

func (a RestrictionArgument) MarshalBinary() ([]byte, error) {
	w := newBinaryWriter()
	w.writeVarUint64(uint64(a.Kind))
	switch a.Kind {
	case 0:
	case 1:
		if a.Bool == nil {
			return nil, fmt.Errorf("missing restriction bool argument")
		}
		writeBool(w, *a.Bool)
	case 2:
		if a.Int64 == nil {
			return nil, fmt.Errorf("missing restriction int64 argument")
		}
		w.writeInt64(*a.Int64)
	case 3:
		if a.String == nil {
			return nil, fmt.Errorf("missing restriction string argument")
		}
		w.writeString(*a.String)
	case 4:
		if a.Time == nil {
			return nil, fmt.Errorf("missing restriction time argument")
		}
		if err := a.Time.MarshalBinaryInto(w); err != nil {
			return nil, err
		}
	case 5:
		if a.PublicKey == nil {
			return nil, fmt.Errorf("missing restriction public key argument")
		}
		if err := a.PublicKey.MarshalBinaryInto(w); err != nil {
			return nil, err
		}
	case 6:
		if a.SHA256 == nil {
			return nil, fmt.Errorf("missing restriction sha256 argument")
		}
		if err := writeSHA256Hex(w, *a.SHA256); err != nil {
			return nil, err
		}
	case 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19:
		if a.ObjectID == nil {
			return nil, fmt.Errorf("missing restriction object id argument")
		}
		if err := a.ObjectID.MarshalBinaryInto(w); err != nil {
			return nil, err
		}
	case 20:
		writeBoolSet(w, a.BoolSet)
	case 21:
		writeInt64Set(w, a.Int64Set)
	case 22:
		writeStringSet(w, a.StringSet)
	case 23:
		if err := writeTimeSet(w, a.TimeSet); err != nil {
			return nil, err
		}
	case 24:
		if err := writePublicKeySet(w, a.PublicKeySet); err != nil {
			return nil, err
		}
	case 25:
		if err := writeSHA256HexSet(w, a.SHA256Set); err != nil {
			return nil, err
		}
	case 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38:
		if err := writeObjectIDSet(w, a.ObjectIDSet); err != nil {
			return nil, err
		}
	case 39:
		if err := writeRestrictionArray(w, a.Restrictions); err != nil {
			return nil, err
		}
	case 40:
		if err := writeRestrictionListArray(w, a.RestrictionSets); err != nil {
			return nil, err
		}
	case 41:
		if a.VariantAssert == nil {
			return nil, fmt.Errorf("missing restriction variant_assert argument")
		}
		w.writeInt64(a.VariantAssert.Tag)
		if err := writeRestrictionArray(w, a.VariantAssert.Restrictions); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unsupported restriction argument type %d", a.Kind)
	}
	return w.Bytes(), nil
}

func (a *RestrictionArgument) UnmarshalBinaryFrom(r *binaryReader) error {
	kind, err := r.readVarUint64()
	if err != nil {
		return err
	}
	a.Kind = uint16(kind)
	switch a.Kind {
	case 0:
	case 1:
		value, err := readBool(r)
		if err != nil {
			return err
		}
		a.Bool = &value
	case 2:
		value, err := r.readInt64()
		if err != nil {
			return err
		}
		a.Int64 = &value
	case 3:
		value, err := r.readString()
		if err != nil {
			return err
		}
		a.String = &value
	case 4:
		value, err := readTime(r)
		if err != nil {
			return err
		}
		a.Time = &value
	case 5:
		value, err := readPublicKey(r)
		if err != nil {
			return err
		}
		a.PublicKey = &value
	case 6:
		value, err := readSHA256Hex(r)
		if err != nil {
			return err
		}
		a.SHA256 = &value
	case 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19:
		value, err := readObjectID(r)
		if err != nil {
			return err
		}
		a.ObjectID = &value
	case 20:
		values, err := readBoolSet(r)
		if err != nil {
			return err
		}
		a.BoolSet = values
	case 21:
		values, err := readInt64Set(r)
		if err != nil {
			return err
		}
		a.Int64Set = values
	case 22:
		values, err := readStringSet(r)
		if err != nil {
			return err
		}
		a.StringSet = values
	case 23:
		values, err := readTimeSet(r)
		if err != nil {
			return err
		}
		a.TimeSet = values
	case 24:
		values, err := readPublicKeySet(r)
		if err != nil {
			return err
		}
		a.PublicKeySet = values
	case 25:
		values, err := readSHA256HexSet(r)
		if err != nil {
			return err
		}
		a.SHA256Set = values
	case 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38:
		values, err := readObjectIDSet(r)
		if err != nil {
			return err
		}
		a.ObjectIDSet = values
	case 39:
		values, err := readRestrictionArray(r)
		if err != nil {
			return err
		}
		a.Restrictions = values
	case 40:
		values, err := readRestrictionListArray(r)
		if err != nil {
			return err
		}
		a.RestrictionSets = values
	case 41:
		tag, err := r.readInt64()
		if err != nil {
			return err
		}
		restrictions, err := readRestrictionArray(r)
		if err != nil {
			return err
		}
		a.VariantAssert = &VariantAssertArgument{Tag: tag, Restrictions: restrictions}
	default:
		return fmt.Errorf("unsupported restriction argument type %d", a.Kind)
	}
	return nil
}

// Restriction mirrors graphene::protocol::restriction.
type Restriction struct {
	MemberIndex     uint64              `json:"member_index"`
	RestrictionType uint64              `json:"restriction_type"`
	Argument        RestrictionArgument `json:"argument"`
	Extensions      []json.RawMessage   `json:"extensions"`
}

func (r Restriction) MarshalJSON() ([]byte, error) {
	type alias Restriction
	if r.Extensions == nil {
		r.Extensions = []json.RawMessage{}
	}
	return json.Marshal(alias(r))
}

func (r *Restriction) UnmarshalJSON(data []byte) error {
	type alias Restriction
	var payload alias
	if err := json.Unmarshal(data, &payload); err != nil {
		return err
	}
	*r = Restriction(payload)
	return nil
}

func (r Restriction) MarshalBinary() ([]byte, error) {
	w := newBinaryWriter()
	if err := r.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}

func (r Restriction) MarshalBinaryInto(w *binaryWriter) error {
	w.writeVarUint64(r.MemberIndex)
	w.writeVarUint64(r.RestrictionType)
	raw, err := r.Argument.MarshalBinary()
	if err != nil {
		return err
	}
	if _, err := w.Write(raw); err != nil {
		return err
	}
	if !extensionsEmpty(r.Extensions) {
		return fmt.Errorf("restriction extensions are not supported in binary serialization")
	}
	w.writeVarUint64(0)
	return nil
}

func (r *Restriction) UnmarshalBinaryFrom(rd *binaryReader) error {
	memberIndex, err := rd.readVarUint64()
	if err != nil {
		return err
	}
	restrictionType, err := rd.readVarUint64()
	if err != nil {
		return err
	}
	var argument RestrictionArgument
	if err := argument.UnmarshalBinaryFrom(rd); err != nil {
		return err
	}
	extCount, err := rd.readVarUint64()
	if err != nil {
		return err
	}
	if extCount != 0 {
		return fmt.Errorf("restriction extensions are not supported in binary serialization")
	}
	r.MemberIndex = memberIndex
	r.RestrictionType = restrictionType
	r.Argument = argument
	r.Extensions = nil
	return nil
}

func writeRestrictionArray(w *binaryWriter, values []Restriction) error {
	w.writeVarUint64(uint64(len(values)))
	for _, value := range values {
		if err := value.MarshalBinaryInto(w); err != nil {
			return err
		}
	}
	return nil
}

func readRestrictionArray(r *binaryReader) ([]Restriction, error) {
	count, err := r.readVarUint64()
	if err != nil {
		return nil, err
	}
	out := make([]Restriction, 0, count)
	for i := uint64(0); i < count; i++ {
		var value Restriction
		if err := value.UnmarshalBinaryFrom(r); err != nil {
			return nil, err
		}
		out = append(out, value)
	}
	return out, nil
}

func writeRestrictionListArray(w *binaryWriter, values [][]Restriction) error {
	w.writeVarUint64(uint64(len(values)))
	for _, value := range values {
		if err := writeRestrictionArray(w, value); err != nil {
			return err
		}
	}
	return nil
}

func readRestrictionListArray(r *binaryReader) ([][]Restriction, error) {
	count, err := r.readVarUint64()
	if err != nil {
		return nil, err
	}
	out := make([][]Restriction, 0, count)
	for i := uint64(0); i < count; i++ {
		value, err := readRestrictionArray(r)
		if err != nil {
			return nil, err
		}
		out = append(out, value)
	}
	return out, nil
}

func writeSHA256Hex(w *binaryWriter, value string) error {
	raw, err := hex.DecodeString(strings.TrimSpace(value))
	if err != nil {
		return err
	}
	return writeFixedBytes(w, raw, 32)
}

func readSHA256Hex(r *binaryReader) (string, error) {
	raw, err := readFixedBytes(r, 32)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(raw), nil
}

func writeSHA256HexSet(w *binaryWriter, values []string) error {
	if len(values) < 2 {
		w.writeVarUint64(uint64(len(values)))
		for _, value := range values {
			if err := writeSHA256Hex(w, value); err != nil {
				return err
			}
		}
		return nil
	}
	sorted := append([]string(nil), values...)
	sort.Slice(sorted, func(i, j int) bool { return compareStrings(sorted[i], sorted[j]) < 0 })
	w.writeVarUint64(uint64(len(sorted)))
	for _, value := range sorted {
		if err := writeSHA256Hex(w, value); err != nil {
			return err
		}
	}
	return nil
}

func readSHA256HexSet(r *binaryReader) ([]string, error) {
	count, err := r.readVarUint64()
	if err != nil {
		return nil, err
	}
	out := make([]string, 0, count)
	for i := uint64(0); i < count; i++ {
		value, err := readSHA256Hex(r)
		if err != nil {
			return nil, err
		}
		out = append(out, value)
	}
	return out, nil
}

func writeBoolSet(w *binaryWriter, values []bool) {
	if len(values) < 2 {
		w.writeVarUint64(uint64(len(values)))
		for _, value := range values {
			writeBool(w, value)
		}
		return
	}
	sorted := append([]bool(nil), values...)
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i] == sorted[j] {
			return false
		}
		return !sorted[i] && sorted[j]
	})
	w.writeVarUint64(uint64(len(sorted)))
	for _, value := range sorted {
		writeBool(w, value)
	}
}

func readBoolSet(r *binaryReader) ([]bool, error) {
	count, err := r.readVarUint64()
	if err != nil {
		return nil, err
	}
	out := make([]bool, 0, count)
	for i := uint64(0); i < count; i++ {
		value, err := readBool(r)
		if err != nil {
			return nil, err
		}
		out = append(out, value)
	}
	return out, nil
}

func writeInt64Set(w *binaryWriter, values []int64) {
	if len(values) < 2 {
		w.writeVarUint64(uint64(len(values)))
		for _, value := range values {
			w.writeInt64(value)
		}
		return
	}
	sorted := append([]int64(nil), values...)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i] < sorted[j] })
	w.writeVarUint64(uint64(len(sorted)))
	for _, value := range sorted {
		w.writeInt64(value)
	}
}

func readInt64Set(r *binaryReader) ([]int64, error) {
	count, err := r.readVarUint64()
	if err != nil {
		return nil, err
	}
	out := make([]int64, 0, count)
	for i := uint64(0); i < count; i++ {
		value, err := r.readInt64()
		if err != nil {
			return nil, err
		}
		out = append(out, value)
	}
	return out, nil
}

func writeStringSet(w *binaryWriter, values []string) {
	if len(values) < 2 {
		w.writeVarUint64(uint64(len(values)))
		for _, value := range values {
			w.writeString(value)
		}
		return
	}
	sorted := append([]string(nil), values...)
	sort.Slice(sorted, func(i, j int) bool { return compareStrings(sorted[i], sorted[j]) < 0 })
	w.writeVarUint64(uint64(len(sorted)))
	for _, value := range sorted {
		w.writeString(value)
	}
}

func readStringSet(r *binaryReader) ([]string, error) {
	count, err := r.readVarUint64()
	if err != nil {
		return nil, err
	}
	out := make([]string, 0, count)
	for i := uint64(0); i < count; i++ {
		value, err := r.readString()
		if err != nil {
			return nil, err
		}
		out = append(out, value)
	}
	return out, nil
}

func writeTimeSet(w *binaryWriter, values []Time) error {
	if len(values) < 2 {
		w.writeVarUint64(uint64(len(values)))
		for _, value := range values {
			if err := value.MarshalBinaryInto(w); err != nil {
				return err
			}
		}
		return nil
	}
	sorted := append([]Time(nil), values...)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].UnixSeconds() < sorted[j].UnixSeconds() })
	w.writeVarUint64(uint64(len(sorted)))
	for _, value := range sorted {
		if err := value.MarshalBinaryInto(w); err != nil {
			return err
		}
	}
	return nil
}

func readTimeSet(r *binaryReader) ([]Time, error) {
	count, err := r.readVarUint64()
	if err != nil {
		return nil, err
	}
	out := make([]Time, 0, count)
	for i := uint64(0); i < count; i++ {
		value, err := readTime(r)
		if err != nil {
			return nil, err
		}
		out = append(out, value)
	}
	return out, nil
}
