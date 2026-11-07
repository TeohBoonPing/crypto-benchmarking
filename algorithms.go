package main

import (
	"bytes"
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"

	"golang.org/x/crypto/chacha20poly1305"
)

// ============================================================
// SHA-256
// ============================================================

var sha256Result [32]byte

type SHA256Target struct {
	inputSize int
	input     []byte
}

func NewSHA256Target(inputSize int) *SHA256Target {
	return &SHA256Target{
		inputSize: inputSize,
	}
}

func (t *SHA256Target) Name() string {
	return "SHA-256"
}

func (t *SHA256Target) Operation() string {
	return "hash"
}

func (t *SHA256Target) InputSize() int {
	return t.inputSize
}

func (t *SHA256Target) Setup() error {
	t.input = bytes.Repeat([]byte{1}, t.inputSize)

	return nil
}

func (t *SHA256Target) Run() error {
	sha256Result = sha256.Sum256(t.input)

	return nil
}

// ============================================================
// AES-GCM
// ============================================================

var aesGCMResult []byte

type AESGCMTarget struct {
	inputSize int
	input     []byte
	nonce     []byte
	aead      cipher.AEAD
	output    []byte
}

func NewAESGCMTarget(inputSize int) *AESGCMTarget {
	return &AESGCMTarget{
		inputSize: inputSize,
	}
}

func (t *AESGCMTarget) Name() string {
	return "AES-GCM"
}

func (t *AESGCMTarget) Operation() string {
	return "encrypt"
}

func (t *AESGCMTarget) InputSize() int {
	return t.inputSize
}

func (t *AESGCMTarget) Setup() error {
	t.input = bytes.Repeat([]byte{1}, t.inputSize)

	key := bytes.Repeat([]byte{2}, 32)

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	t.aead, err = cipher.NewGCM(block)
	if err != nil {
		return err
	}

	t.nonce = bytes.Repeat(
		[]byte{3},
		t.aead.NonceSize(),
	)

	t.output = make(
		[]byte,
		0,
		t.inputSize+t.aead.Overhead(),
	)

	return nil
}

func (t *AESGCMTarget) Run() error {
	aesGCMResult = t.aead.Seal(
		t.output[:0],
		t.nonce,
		t.input,
		nil,
	)

	return nil
}

// ============================================================
// ChaCha20-Poly1305
// ============================================================

var chaChaResult []byte

type ChaCha20Poly1305Target struct {
	inputSize int
	input     []byte
	nonce     []byte
	aead      cipher.AEAD
	output    []byte
}

func NewChaCha20Poly1305Target(
	inputSize int,
) *ChaCha20Poly1305Target {
	return &ChaCha20Poly1305Target{
		inputSize: inputSize,
	}
}

func (t *ChaCha20Poly1305Target) Name() string {
	return "ChaCha20-Poly1305"
}

func (t *ChaCha20Poly1305Target) Operation() string {
	return "encrypt"
}

func (t *ChaCha20Poly1305Target) InputSize() int {
	return t.inputSize
}

func (t *ChaCha20Poly1305Target) Setup() error {
	t.input = bytes.Repeat([]byte{1}, t.inputSize)

	key := bytes.Repeat(
		[]byte{4},
		chacha20poly1305.KeySize,
	)

	var err error

	t.aead, err = chacha20poly1305.New(key)
	if err != nil {
		return err
	}

	t.nonce = bytes.Repeat(
		[]byte{5},
		t.aead.NonceSize(),
	)

	t.output = make(
		[]byte,
		0,
		t.inputSize+t.aead.Overhead(),
	)

	return nil
}

func (t *ChaCha20Poly1305Target) Run() error {
	chaChaResult = t.aead.Seal(
		t.output[:0],
		t.nonce,
		t.input,
		nil,
	)

	return nil
}

// ============================================================
// RSA shared data
// ============================================================

type RSAData struct {
	PrivateKey *rsa.PrivateKey
	Digest     [32]byte
	Signature  []byte
}

func NewRSAData() (*RSAData, error) {
	privateKey, err := rsa.GenerateKey(
		rand.Reader,
		2048,
	)
	if err != nil {
		return nil, err
	}

	message := []byte("crypto benchmark")
	digest := sha256.Sum256(message)

	signature, err := rsa.SignPSS(
		rand.Reader,
		privateKey,
		crypto.SHA256,
		digest[:],
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &RSAData{
		PrivateKey: privateKey,
		Digest:     digest,
		Signature:  signature,
	}, nil
}

// ============================================================
// RSA signing
// ============================================================

var rsaSignResult []byte

type RSASignTarget struct {
	data *RSAData
}

func NewRSASignTarget(data *RSAData) *RSASignTarget {
	return &RSASignTarget{
		data: data,
	}
}

func (t *RSASignTarget) Name() string {
	return "RSA"
}

func (t *RSASignTarget) Operation() string {
	return "sign"
}

func (t *RSASignTarget) InputSize() int {
	return len(t.data.Digest)
}

func (t *RSASignTarget) Setup() error {
	return nil
}

func (t *RSASignTarget) Run() error {
	signature, err := rsa.SignPSS(
		rand.Reader,
		t.data.PrivateKey,
		crypto.SHA256,
		t.data.Digest[:],
		nil,
	)
	if err != nil {
		return err
	}

	rsaSignResult = signature

	return nil
}

// ============================================================
// RSA verification
// ============================================================

type RSAVerifyTarget struct {
	data *RSAData
}

func NewRSAVerifyTarget(data *RSAData) *RSAVerifyTarget {
	return &RSAVerifyTarget{
		data: data,
	}
}

func (t *RSAVerifyTarget) Name() string {
	return "RSA"
}

func (t *RSAVerifyTarget) Operation() string {
	return "verify"
}

func (t *RSAVerifyTarget) InputSize() int {
	return len(t.data.Digest)
}

func (t *RSAVerifyTarget) Setup() error {
	return nil
}

func (t *RSAVerifyTarget) Run() error {
	return rsa.VerifyPSS(
		&t.data.PrivateKey.PublicKey,
		crypto.SHA256,
		t.data.Digest[:],
		t.data.Signature,
		nil,
	)
}
