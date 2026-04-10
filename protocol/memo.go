package protocol

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/froooze/go-bitshares/ecc"
)

// MemoData is the JSON-backed memo payload used by the wallet helpers.
type MemoData struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Nonce   string `json:"nonce"`
	Message string `json:"message"`
}

func (m MemoData) MarshalBinary() ([]byte, error) {
	w := newBinaryWriter()
	if err := m.writeBinary(w); err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}

func (m *MemoData) UnmarshalBinary(data []byte) error {
	return m.readBinary(newBinaryReader(data))
}

func (m MemoData) writeBinary(w *binaryWriter) error {
	from, err := ecc.PublicKeyFromString(strings.TrimSpace(m.From))
	if err != nil {
		return fmt.Errorf("memo from key: %w", err)
	}
	to, err := ecc.PublicKeyFromString(strings.TrimSpace(m.To))
	if err != nil {
		return fmt.Errorf("memo to key: %w", err)
	}
	nonce, err := strconv.ParseUint(strings.TrimSpace(m.Nonce), 10, 64)
	if err != nil {
		return fmt.Errorf("memo nonce: %w", err)
	}
	message, err := hex.DecodeString(strings.TrimSpace(m.Message))
	if err != nil {
		return fmt.Errorf("memo message: %w", err)
	}

	if _, err := w.Write(from.Bytes()); err != nil {
		return err
	}
	if _, err := w.Write(to.Bytes()); err != nil {
		return err
	}
	w.writeUint64(nonce)
	w.writeBytes(message)
	return nil
}

func (m *MemoData) readBinary(r *binaryReader) error {
	fromBytes := make([]byte, 33)
	if _, err := r.Read(fromBytes); err != nil {
		return err
	}
	toBytes := make([]byte, 33)
	if _, err := r.Read(toBytes); err != nil {
		return err
	}
	nonce, err := r.readUint64()
	if err != nil {
		return err
	}
	message, err := r.readBytes()
	if err != nil {
		return err
	}
	from, err := ecc.PublicKeyFromBytes(fromBytes)
	if err != nil {
		return err
	}
	to, err := ecc.PublicKeyFromBytes(toBytes)
	if err != nil {
		return err
	}
	m.From = from.String()
	m.To = to.String()
	m.Nonce = strconv.FormatUint(nonce, 10)
	m.Message = hex.EncodeToString(message)
	return nil
}

func memoDataFromRaw(data json.RawMessage) (*MemoData, error) {
	trimmed := bytes.TrimSpace(data)
	if len(trimmed) == 0 || bytes.Equal(trimmed, []byte("null")) {
		return nil, nil
	}
	var memo MemoData
	if err := json.Unmarshal(trimmed, &memo); err != nil {
		return nil, err
	}
	return &memo, nil
}
