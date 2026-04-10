package ecc

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"sync"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
)

const (
	defaultAddressPrefix = "BTS"
	privKeyVersion       = 0x80
)

var (
	prefixMu       sync.RWMutex
	addressPrefix  = defaultAddressPrefix
	errNoKey       = errors.New("key is not configured")
	errBadChecksum = errors.New("checksum did not match")
)

// SetAddressPrefix sets the default BitShares address prefix used for public keys and addresses.
func SetAddressPrefix(prefix string) {
	prefixMu.Lock()
	if strings.TrimSpace(prefix) == "" {
		addressPrefix = defaultAddressPrefix
	} else {
		addressPrefix = strings.TrimSpace(prefix)
	}
	prefixMu.Unlock()
}

// GetAddressPrefix returns the current default BitShares address prefix.
func GetAddressPrefix() string {
	prefixMu.RLock()
	defer prefixMu.RUnlock()
	return addressPrefix
}

// NormalizeBrainKey collapses whitespace to a single space and trims the result.
func NormalizeBrainKey(brainKey string) string {
	brainKey = strings.TrimSpace(brainKey)
	if brainKey == "" {
		return ""
	}
	return strings.Join(strings.Fields(brainKey), " ")
}

func normalizeBrainKeyBytes(brainKey []byte) []byte {
	brainKey = bytes.TrimSpace(brainKey)
	if len(brainKey) == 0 {
		return nil
	}
	fields := bytes.Fields(brainKey)
	if len(fields) == 0 {
		return nil
	}
	size := len(fields) - 1
	for _, field := range fields {
		size += len(field)
	}
	out := make([]byte, 0, size)
	for i, field := range fields {
		if i > 0 {
			out = append(out, ' ')
		}
		out = append(out, field...)
	}
	return out
}

// Address is the BitShares blockchain address hash derived from a public key.
type Address [20]byte

// FromPublicAddress returns the BitShares blockchain address hash for a public key.
func FromPublicAddress(pub *PublicKey) Address {
	var out Address
	if pub == nil || pub.key == nil {
		return out
	}

	sum := sha512.Sum512(pub.key.SerializeCompressed())
	h := ripemd160.New()
	_, _ = h.Write(sum[:])
	copy(out[:], h.Sum(nil))
	return out
}

// String returns the BitShares address string representation.
func (a Address) String() string {
	sum := ripemd160.New()
	_, _ = sum.Write(a[:])
	checksum := sum.Sum(nil)
	payload := append(a[:], checksum[:4]...)
	return GetAddressPrefix() + base58.Encode(payload)
}

// AddressFromString parses a BitShares address string.
func AddressFromString(value string) (Address, error) {
	var out Address
	value = strings.TrimSpace(value)
	if value == "" {
		return out, fmt.Errorf("empty address")
	}
	prefix := GetAddressPrefix()
	if !strings.HasPrefix(value, prefix) {
		return out, fmt.Errorf("expected prefix %q, got %q", prefix, value[:min(len(value), len(prefix))])
	}
	raw, err := decodeBitSharesBase58(value[len(prefix):])
	if err != nil {
		return out, err
	}
	if len(raw) != len(out) {
		return out, fmt.Errorf("unexpected address length %d", len(raw))
	}
	copy(out[:], raw)
	return out, nil
}

// MustAddressFromString parses a BitShares address string or panics.
func MustAddressFromString(value string) Address {
	out, err := AddressFromString(value)
	if err != nil {
		panic(err)
	}
	return out
}

// Bytes returns the 20-byte address hash.
func (a Address) Bytes() []byte {
	return append([]byte(nil), a[:]...)
}

// Compare compares two addresses lexicographically.
func (a Address) Compare(other Address) int {
	return bytes.Compare(a[:], other[:])
}

// PublicKey wraps a secp256k1 public key.
type PublicKey struct {
	key *btcec.PublicKey
}

// PublicKeyFromBytes parses a compressed or uncompressed public key.
func PublicKeyFromBytes(buf []byte) (*PublicKey, error) {
	if len(buf) == 0 {
		return nil, fmt.Errorf("empty public key")
	}
	pub, err := btcec.ParsePubKey(buf, btcec.S256())
	if err != nil {
		return nil, err
	}
	return &PublicKey{key: pub}, nil
}

// PublicKeyFromString parses a BitShares public key string.
func PublicKeyFromString(value string) (*PublicKey, error) {
	if value == "" {
		return nil, fmt.Errorf("empty public key")
	}
	prefix := GetAddressPrefix()
	if !strings.HasPrefix(value, prefix) {
		return nil, fmt.Errorf("expected prefix %q, got %q", prefix, value[:min(len(value), len(prefix))])
	}

	raw, err := decodeBitSharesBase58(value[len(prefix):])
	if err != nil {
		return nil, err
	}
	return PublicKeyFromBytes(raw)
}

// MustPublicKeyFromString parses a public key string or panics.
func MustPublicKeyFromString(value string) *PublicKey {
	pub, err := PublicKeyFromString(value)
	if err != nil {
		panic(err)
	}
	return pub
}

// Bytes returns the compressed public key encoding.
func (p *PublicKey) Bytes() []byte {
	if p == nil || p.key == nil {
		return nil
	}
	return p.key.SerializeCompressed()
}

// String returns the BitShares public key string.
func (p *PublicKey) String() string {
	if p == nil || p.key == nil {
		return ""
	}
	return GetAddressPrefix() + encodeBitSharesBase58(p.Bytes())
}

// ToAddressString returns the BitShares blockchain address string for this public key.
func (p *PublicKey) ToAddressString() string {
	return FromPublicAddress(p).String()
}

// BlockchainAddress returns the address hash for comparison purposes.
func (p *PublicKey) BlockchainAddress() Address {
	return FromPublicAddress(p)
}

// Compare compares the public key blockchain addresses.
func (p *PublicKey) Compare(other *PublicKey) int {
	if p == nil || other == nil {
		switch {
		case p == nil && other == nil:
			return 0
		case p == nil:
			return -1
		default:
			return 1
		}
	}
	return p.BlockchainAddress().Compare(other.BlockchainAddress())
}

// ToECDSA returns the underlying key as an ECDSA public key.
func (p *PublicKey) ToECDSA() *btcec.PublicKey {
	if p == nil {
		return nil
	}
	return p.key
}

// PrivateKey wraps a secp256k1 private key.
type PrivateKey struct {
	key *btcec.PrivateKey
}

// Wipe clears the private-key wrapper and drops the underlying handle.
// This is best-effort only; the btcec internals are not memory-zeroed here.
func (p *PrivateKey) Wipe() {
	if p == nil {
		return
	}
	if p.key != nil {
		secret := p.key.Serialize()
		zeroBytes(secret)
	}
	p.key = nil
}

// RandomPrivateKey returns a random private key.
func RandomPrivateKey() (*PrivateKey, error) {
	k, err := btcec.NewPrivateKey(btcec.S256())
	if err != nil {
		return nil, err
	}
	return &PrivateKey{key: k}, nil
}

// PrivateKeyFromBytes creates a private key from a 32-byte seed.
func PrivateKeyFromBytes(seed []byte) *PrivateKey {
	priv, _ := btcec.PrivKeyFromBytes(btcec.S256(), seed)
	return &PrivateKey{key: priv}
}

// PrivateKeyFromSeed creates a private key by hashing seed bytes with SHA-256.
// The caller owns the seed buffer.
func PrivateKeyFromSeed(seed []byte) *PrivateKey {
	sum := sha256.Sum256(seed)
	return PrivateKeyFromBytes(sum[:])
}

// PrivateKeyFromWIF parses a WIF-encoded private key from bytes.
// The caller owns the input buffer.
func PrivateKeyFromWIF(wif []byte) (*PrivateKey, error) {
	payload, version, err := base58.CheckDecode(string(bytes.TrimSpace(wif)))
	if err != nil {
		return nil, err
	}
	if version != privKeyVersion {
		return nil, fmt.Errorf("unexpected WIF version %d", version)
	}
	if len(payload) != 32 && len(payload) != 33 {
		return nil, fmt.Errorf("unexpected WIF payload length %d", len(payload))
	}
	if len(payload) == 33 {
		payload = payload[:32]
	}
	defer zeroBytes(payload)
	return &PrivateKey{key: PrivateKeyFromBytes(payload).key}, nil
}

// MustPrivateKeyFromWIF parses a WIF or panics.
func MustPrivateKeyFromWIF(wif []byte) *PrivateKey {
	priv, err := PrivateKeyFromWIF(wif)
	if err != nil {
		panic(err)
	}
	return priv
}

// FromBrainKey derives a private key from normalized brain key bytes and a sequence number.
// The caller owns the input buffer.
func FromBrainKey(brainKey []byte, sequence int) (*PrivateKey, error) {
	if sequence < 0 {
		return nil, fmt.Errorf("invalid sequence")
	}
	normalized := normalizeBrainKeyBytes(brainKey)
	if len(normalized) == 0 {
		return nil, fmt.Errorf("empty brain key")
	}
	seq := strconv.AppendInt(make([]byte, 0, 20), int64(sequence), 10)
	material := make([]byte, 0, len(normalized)+1+len(seq))
	material = append(material, normalized...)
	material = append(material, ' ')
	material = append(material, seq...)
	brain := sha512.Sum512(material)
	zeroBytes(material)
	zeroBytes(seq)
	zeroBytes(normalized)
	keySeed := sha256.Sum256(brain[:])
	return PrivateKeyFromBytes(keySeed[:]), nil
}

// WIF returns the private key in Wallet Import Format.
// The returned string is a fresh copy and should be treated as secret material.
func (p *PrivateKey) WIF() string {
	if p == nil || p.key == nil {
		return ""
	}
	return base58.CheckEncode(p.key.Serialize(), privKeyVersion)
}

// PublicKey returns the corresponding public key.
func (p *PrivateKey) PublicKey() *PublicKey {
	if p == nil || p.key == nil {
		return nil
	}
	return &PublicKey{key: p.key.PubKey()}
}

// Bytes returns the 32-byte private key.
// The returned slice is a fresh copy and should be treated as secret material.
func (p *PrivateKey) Bytes() []byte {
	if p == nil || p.key == nil {
		return nil
	}
	return p.key.Serialize()
}

// SharedSecret derives the BitShares ECIES shared secret.
func (p *PrivateKey) SharedSecret(pub *PublicKey, legacy bool) ([]byte, error) {
	if p == nil || p.key == nil || pub == nil || pub.key == nil {
		return nil, errNoKey
	}

	x, _ := pub.key.Curve.ScalarMult(pub.key.X, pub.key.Y, p.key.Serialize())
	secret := x.Bytes()
	defer zeroBytes(secret)
	if !legacy && len(secret) < 32 {
		pad := make([]byte, 32-len(secret))
		secret = append(pad, secret...)
	}

	sum := sha512.Sum512(secret)
	out := make([]byte, len(sum))
	copy(out, sum[:])
	return out, nil
}

// Child derives a child key from the current key and 32 bytes of offset entropy.
func (p *PrivateKey) Child(offset []byte) (*PrivateKey, error) {
	if len(offset) != 32 {
		return nil, fmt.Errorf("offset length")
	}
	if p == nil || p.key == nil {
		return nil, errNoKey
	}
	seed := append(append([]byte(nil), p.PublicKey().Bytes()...), offset...)
	defer zeroBytes(seed)
	sum := sha256.Sum256(seed)
	n := new(big.Int).SetBytes(sum[:])
	if n.Cmp(btcec.S256().N) >= 0 {
		return nil, fmt.Errorf("child offset went out of bounds")
	}
	derived := new(big.Int).Add(p.key.D, n)
	derived.Mod(derived, btcec.S256().N)
	if derived.Sign() == 0 {
		return nil, fmt.Errorf("child offset derived to an invalid key")
	}
	priv, _ := btcec.PrivKeyFromBytes(btcec.S256(), derived.Bytes())
	return &PrivateKey{key: priv}, nil
}

// SignCompact signs a 32-byte hash and returns a compact recoverable signature.
func (p *PrivateKey) SignCompact(hash []byte) (*Signature, error) {
	if p == nil || p.key == nil {
		return nil, errNoKey
	}
	sig, err := btcec.SignCompact(btcec.S256(), p.key, hash, true)
	if err != nil {
		return nil, err
	}
	return &Signature{data: sig}, nil
}

// Signature is a recoverable compact secp256k1 signature.
type Signature struct {
	data []byte
}

// SignatureFromHex parses a compact signature hex string.
func SignatureFromHex(value string) (*Signature, error) {
	raw, err := hex.DecodeString(strings.TrimSpace(value))
	if err != nil {
		return nil, err
	}
	if len(raw) != 65 {
		return nil, fmt.Errorf("invalid signature length %d", len(raw))
	}
	return &Signature{data: raw}, nil
}

// Bytes returns the raw compact signature bytes.
func (s *Signature) Bytes() []byte {
	if s == nil {
		return nil
	}
	return append([]byte(nil), s.data...)
}

// Hex returns the signature as hex.
func (s *Signature) Hex() string {
	if s == nil {
		return ""
	}
	return hex.EncodeToString(s.data)
}

// String returns the hex representation.
func (s *Signature) String() string {
	return s.Hex()
}

// RecoverPublicKey recovers the public key used to create the signature.
func (s *Signature) RecoverPublicKey(hash []byte) (*PublicKey, bool, error) {
	if s == nil || len(s.data) != 65 {
		return nil, false, fmt.Errorf("invalid signature")
	}
	pub, compressed, err := btcec.RecoverCompact(btcec.S256(), s.data, hash)
	if err != nil {
		return nil, false, err
	}
	return &PublicKey{key: pub}, compressed, nil
}

// Aes provides symmetric AES-256-CBC encryption using the BitShares seed derivation.
type Aes struct {
	key []byte
	iv  []byte
}

// FromSeed builds an AES helper from a seed value.
// The caller owns the seed buffer.
func FromSeed(seed []byte) *Aes {
	sum := sha512.Sum512(seed)
	return &Aes{
		key: append([]byte(nil), sum[:32]...),
		iv:  append([]byte(nil), sum[32:48]...),
	}
}

// Encrypt encrypts raw bytes using PKCS7 padding.
func (a *Aes) Encrypt(plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(a.key)
	if err != nil {
		return nil, err
	}
	padded := pkcs7Pad(plaintext, aes.BlockSize)
	ciphertext := make([]byte, len(padded))
	mode := cipher.NewCBCEncrypter(block, a.iv)
	mode.CryptBlocks(ciphertext, padded)
	return ciphertext, nil
}

// Decrypt decrypts raw bytes using PKCS7 padding.
func (a *Aes) Decrypt(ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(a.key)
	if err != nil {
		return nil, err
	}
	if len(ciphertext)%aes.BlockSize != 0 {
		return nil, fmt.Errorf("invalid ciphertext size")
	}
	plaintext := make([]byte, len(ciphertext))
	mode := cipher.NewCBCDecrypter(block, a.iv)
	mode.CryptBlocks(plaintext, ciphertext)
	return pkcs7Unpad(plaintext, aes.BlockSize)
}

// EncryptHex encrypts a hex string and returns hex ciphertext.
func (a *Aes) EncryptHex(plaintext string) (string, error) {
	raw, err := hex.DecodeString(strings.TrimSpace(plaintext))
	if err != nil {
		return "", err
	}
	out, err := a.Encrypt(raw)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(out), nil
}

// DecryptHex decrypts a hex ciphertext and returns hex plaintext.
func (a *Aes) DecryptHex(ciphertext string) (string, error) {
	raw, err := hex.DecodeString(strings.TrimSpace(ciphertext))
	if err != nil {
		return "", err
	}
	out, err := a.Decrypt(raw)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(out), nil
}

// DecryptHexToBuffer decrypts a hex ciphertext and returns raw bytes.
func (a *Aes) DecryptHexToBuffer(ciphertext string) ([]byte, error) {
	raw, err := hex.DecodeString(strings.TrimSpace(ciphertext))
	if err != nil {
		return nil, err
	}
	return a.Decrypt(raw)
}

// Wipe zeroes the derived AES key material.
func (a *Aes) Wipe() {
	if a == nil {
		return
	}
	zeroBytes(a.key)
	zeroBytes(a.iv)
	a.key = nil
	a.iv = nil
}

// EncryptWithChecksum encrypts a memo payload with a checksum prefix.
func EncryptWithChecksum(privateKey *PrivateKey, publicKey *PublicKey, nonce string, message []byte) ([]byte, error) {
	if privateKey == nil || publicKey == nil {
		return nil, errNoKey
	}
	secret, err := privateKey.SharedSecret(publicKey, false)
	if err != nil {
		return nil, err
	}
	secretHex := make([]byte, hex.EncodedLen(len(secret)))
	hex.Encode(secretHex, secret)
	seed := append([]byte(nonce), secretHex...)
	aes := FromSeed(seed)
	defer aes.Wipe()
	defer zeroBytes(seed)
	defer zeroBytes(secretHex)
	zeroBytes(secret)
	checksum := sha256.Sum256(message)
	payload := append(checksum[:4], message...)
	return aes.Encrypt(payload)
}

// DecryptWithChecksum decrypts a memo payload and verifies its checksum.
func DecryptWithChecksum(privateKey *PrivateKey, publicKey *PublicKey, nonce string, message []byte, legacy bool) ([]byte, error) {
	if privateKey == nil || publicKey == nil {
		return nil, errNoKey
	}
	secret, err := privateKey.SharedSecret(publicKey, legacy)
	if err != nil {
		return nil, err
	}
	secretHex := make([]byte, hex.EncodedLen(len(secret)))
	hex.Encode(secretHex, secret)
	seed := append([]byte(nonce), secretHex...)
	aes := FromSeed(seed)
	defer aes.Wipe()
	defer zeroBytes(seed)
	defer zeroBytes(secretHex)
	zeroBytes(secret)
	plain, err := aes.Decrypt(message)
	if err != nil {
		return nil, err
	}
	if len(plain) < 4 {
		return nil, fmt.Errorf("invalid key, could not decrypt message")
	}

	checksum := plain[:4]
	payload := plain[4:]
	verify := sha256.Sum256(payload)
	if !bytes.Equal(checksum, verify[:4]) {
		return nil, errBadChecksum
	}
	return payload, nil
}

// GenerateKeys produces BitShares private/public key pairs for the requested roles.
// The caller owns the password buffer.
func GenerateKeys(accountName string, password []byte, roles []string, prefix string) (map[string]*PrivateKey, map[string]string, error) {
	if strings.TrimSpace(accountName) == "" || len(password) == 0 {
		return nil, nil, fmt.Errorf("account name or password required")
	}
	if len(password) < 12 {
		return nil, nil, fmt.Errorf("password must have at least 12 characters")
	}
	if strings.TrimSpace(prefix) != "" {
		SetAddressPrefix(prefix)
	}

	uniq := map[string]struct{}{}
	if len(roles) == 0 {
		roles = []string{"active", "owner", "memo"}
	}
	keys := make(map[string]*PrivateKey, len(roles))
	pubs := make(map[string]string, len(roles))
	for _, role := range roles {
		if _, ok := uniq[role]; ok {
			continue
		}
		uniq[role] = struct{}{}
		seed := make([]byte, 0, len(accountName)+len(role)+len(password))
		seed = append(seed, accountName...)
		seed = append(seed, role...)
		seed = append(seed, password...)
		normalized := normalizeBrainKeyBytes(seed)
		zeroBytes(seed)
		key := PrivateKeyFromSeed(normalized)
		zeroBytes(normalized)
		keys[role] = key
		pubs[role] = key.PublicKey().String()
	}
	return keys, pubs, nil
}

func pkcs7Pad(src []byte, blockSize int) []byte {
	padLen := blockSize - (len(src) % blockSize)
	if padLen == 0 {
		padLen = blockSize
	}
	return append(src, bytes.Repeat([]byte{byte(padLen)}, padLen)...)
}

func pkcs7Unpad(src []byte, blockSize int) ([]byte, error) {
	if len(src) == 0 || len(src)%blockSize != 0 {
		return nil, fmt.Errorf("invalid padding")
	}
	padLen := int(src[len(src)-1])
	if padLen == 0 || padLen > blockSize || padLen > len(src) {
		return nil, fmt.Errorf("invalid padding")
	}
	for _, b := range src[len(src)-padLen:] {
		if int(b) != padLen {
			return nil, fmt.Errorf("invalid padding")
		}
	}
	return append([]byte(nil), src[:len(src)-padLen]...), nil
}

func encodeBitSharesBase58(payload []byte) string {
	sum := ripemd160.New()
	_, _ = sum.Write(payload)
	checksum := sum.Sum(nil)
	return base58.Encode(append(payload, checksum[:4]...))
}

func decodeBitSharesBase58(value string) ([]byte, error) {
	raw := base58.Decode(value)
	if len(raw) < 5 {
		return nil, fmt.Errorf("invalid base58 payload")
	}
	payload := raw[:len(raw)-4]
	sum := ripemd160.New()
	_, _ = sum.Write(payload)
	checksum := sum.Sum(nil)
	if !bytes.Equal(raw[len(raw)-4:], checksum[:4]) {
		return nil, errBadChecksum
	}
	return append([]byte(nil), payload...), nil
}

func zeroBytes(buf []byte) {
	for i := range buf {
		buf[i] = 0
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
