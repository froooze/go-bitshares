package protocol

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type binaryWriter struct {
	bytes.Buffer
}

func newBinaryWriter() *binaryWriter {
	return &binaryWriter{}
}

func (w *binaryWriter) writeUint8(v uint8) {
	_ = w.WriteByte(v)
}

func (w *binaryWriter) writeUint16(v uint16) {
	var buf [2]byte
	binary.LittleEndian.PutUint16(buf[:], v)
	_, _ = w.Write(buf[:])
}

func (w *binaryWriter) writeUint32(v uint32) {
	var buf [4]byte
	binary.LittleEndian.PutUint32(buf[:], v)
	_, _ = w.Write(buf[:])
}

func (w *binaryWriter) writeUint64(v uint64) {
	var buf [8]byte
	binary.LittleEndian.PutUint64(buf[:], v)
	_, _ = w.Write(buf[:])
}

func (w *binaryWriter) writeInt64(v int64) {
	w.writeUint64(uint64(v))
}

func (w *binaryWriter) writeVarUint64(v uint64) {
	var buf [10]byte
	n := binary.PutUvarint(buf[:], v)
	_, _ = w.Write(buf[:n])
}

func (w *binaryWriter) writeVarInt32(v int32) {
	var buf [10]byte
	n := binary.PutVarint(buf[:], int64(v))
	_, _ = w.Write(buf[:n])
}

func (w *binaryWriter) writeBytes(data []byte) {
	w.writeVarUint64(uint64(len(data)))
	_, _ = w.Write(data)
}

func (w *binaryWriter) writeString(value string) {
	w.writeBytes([]byte(value))
}

type binaryReader struct {
	*bytes.Reader
}

func newBinaryReader(data []byte) *binaryReader {
	return &binaryReader{Reader: bytes.NewReader(data)}
}

func (r *binaryReader) readUint8() (uint8, error) {
	return r.ReadByte()
}

func (r *binaryReader) readUint16() (uint16, error) {
	var buf [2]byte
	if _, err := r.Read(buf[:]); err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint16(buf[:]), nil
}

func (r *binaryReader) readUint32() (uint32, error) {
	var buf [4]byte
	if _, err := r.Read(buf[:]); err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint32(buf[:]), nil
}

func (r *binaryReader) readUint64() (uint64, error) {
	var buf [8]byte
	if _, err := r.Read(buf[:]); err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint64(buf[:]), nil
}

func (r *binaryReader) readInt64() (int64, error) {
	v, err := r.readUint64()
	return int64(v), err
}

func (r *binaryReader) readVarUint64() (uint64, error) {
	return binary.ReadUvarint(r.Reader)
}

func (r *binaryReader) readVarInt32() (int32, error) {
	v, err := binary.ReadVarint(r.Reader)
	return int32(v), err
}

func (r *binaryReader) readBytes() ([]byte, error) {
	n, err := r.readVarUint64()
	if err != nil {
		return nil, err
	}
	if n == 0 {
		return []byte{}, nil
	}
	if n > uint64(r.Len()) {
		return nil, fmt.Errorf("invalid length %d", n)
	}
	out := make([]byte, n)
	if _, err := r.Read(out); err != nil {
		return nil, err
	}
	return out, nil
}

func (r *binaryReader) readString() (string, error) {
	data, err := r.readBytes()
	if err != nil {
		return "", err
	}
	return string(data), nil
}
