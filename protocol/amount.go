package protocol

import "encoding/json"

// AssetAmount represents an amount tied to a specific asset.
type AssetAmount struct {
	Amount  int64    `json:"amount"`
	AssetID ObjectID `json:"asset_id"`
}

func (a AssetAmount) MarshalJSON() ([]byte, error) {
	type alias AssetAmount
	return json.Marshal(alias(a))
}

func (a *AssetAmount) UnmarshalJSON(data []byte) error {
	type alias AssetAmount
	var out alias
	if err := json.Unmarshal(data, &out); err != nil {
		return err
	}

	*a = AssetAmount(out)
	return nil
}

// Price expresses a base/quote ratio.
type Price struct {
	Base  AssetAmount `json:"base"`
	Quote AssetAmount `json:"quote"`
}

func (p Price) MarshalJSON() ([]byte, error) {
	type alias Price
	return json.Marshal(alias(p))
}

func (p *Price) UnmarshalJSON(data []byte) error {
	type alias Price
	var out alias
	if err := json.Unmarshal(data, &out); err != nil {
		return err
	}

	*p = Price(out)
	return nil
}

// MarshalBinary encodes the asset amount in the BitShares wire format.
func (a AssetAmount) MarshalBinary() ([]byte, error) {
	w := newBinaryWriter()
	if err := a.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}

// MarshalBinaryInto writes the asset amount to a binary writer.
func (a AssetAmount) MarshalBinaryInto(w *binaryWriter) error {
	w.writeInt64(a.Amount)
	return a.AssetID.MarshalBinaryInto(w)
}

// UnmarshalBinary decodes the asset amount from the BitShares wire format.
func (a *AssetAmount) UnmarshalBinary(data []byte) error {
	return a.UnmarshalBinaryFrom(newBinaryReader(data))
}

// UnmarshalBinaryFrom reads the asset amount from a binary reader.
func (a *AssetAmount) UnmarshalBinaryFrom(r *binaryReader) error {
	amount, err := r.readInt64()
	if err != nil {
		return err
	}
	asset, err := readObjectID(r)
	if err != nil {
		return err
	}
	a.Amount = amount
	a.AssetID = asset
	return nil
}

func readAssetAmount(r *binaryReader) (AssetAmount, error) {
	var out AssetAmount
	return out, out.UnmarshalBinaryFrom(r)
}

// MarshalBinary encodes the price in the BitShares wire format.
func (p Price) MarshalBinary() ([]byte, error) {
	w := newBinaryWriter()
	if err := p.MarshalBinaryInto(w); err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}

// MarshalBinaryInto writes the price to a binary writer.
func (p Price) MarshalBinaryInto(w *binaryWriter) error {
	if err := p.Base.MarshalBinaryInto(w); err != nil {
		return err
	}
	return p.Quote.MarshalBinaryInto(w)
}

// UnmarshalBinary decodes the price from the BitShares wire format.
func (p *Price) UnmarshalBinary(data []byte) error {
	return p.UnmarshalBinaryFrom(newBinaryReader(data))
}

// UnmarshalBinaryFrom reads the price from a binary reader.
func (p *Price) UnmarshalBinaryFrom(r *binaryReader) error {
	base, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	quote, err := readAssetAmount(r)
	if err != nil {
		return err
	}
	p.Base = base
	p.Quote = quote
	return nil
}

func readPrice(r *binaryReader) (Price, error) {
	var out Price
	return out, out.UnmarshalBinaryFrom(r)
}
