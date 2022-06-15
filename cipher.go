package mainframe

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"io"
	"os"
	"path"
)

// Makes it easier for user to supply an encoding
// format we use for keys or encrypted date
type Encoding string

const (
	PEM    Encoding = "pem"
	Hex    Encoding = "hex"
	Base32 Encoding = "base32"
	Base64 Encoding = "base64"
)

type Key struct {
	Bytes []byte
}

func GenerateKey(keyLen int) (*Key, error) {
	key := make([]byte, keyLen)
	_, err := rand.Read(key)
	if err != nil {
		return nil, fmt.Errorf("could not generate encryption key: %w", err)
	}

	return &Key{key}, nil
}

func LoadKey(reader io.Reader) (*Key, error) {
	return LoadEncodedKey(reader, Base32)
}

func LoadEncodedKey(reader io.Reader, encoding Encoding) (*Key, error) {
	if encoding == PEM {
		pemData, err := io.ReadAll(reader)
		if err != nil {
			return nil, err
		}
		block, _ := pem.Decode(pemData)
		if block == nil {
			return nil, fmt.Errorf("could not decode PEM data")
		}
		return &Key{block.Bytes}, nil
	}

	var decoder io.Reader
	switch encoding {
	case Hex:
		decoder = hex.NewDecoder(reader)
	case Base32:
		decoder = base32.NewDecoder(base32.StdEncoding, reader)
	case Base64:
		decoder = base64.NewDecoder(base64.StdEncoding, reader)
	}

	bytes, err := io.ReadAll(decoder)
	if err != nil {
		return nil, fmt.Errorf("unable to read key data "+
			"with assumed encoding %s because %w", encoding, err)
	}

	return &Key{bytes}, nil
}

func LoadKeyFromFile(filename string) (*Key, error) {
	return LoadEncodedKeyFromFile(filename, Base32)
}

func LoadEncodedKeyFromFile(filename string, encoding Encoding) (*Key, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("could not load key from file: %w", err)
	}
	defer file.Close()
	return LoadEncodedKey(file, encoding)
}

// Save file with base32 encoding to writer
func (key *Key) Save(writer io.Writer) {
	var encoder io.WriteCloser = base32.NewEncoder(base32.StdEncoding, writer)
	encoder.Write(key.Bytes)
	encoder.Close()
}

// SaveToFile stores key at file name filename
func (key *Key) SaveToFile(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("could not save cipher key: %w", err)
	}
	defer file.Close()

	key.Save(file)

	return nil
}

type Cipher struct {
	block cipher.Block
}

// New cipher for symetrical encryption
func NewCipher(key *Key) (*Cipher, error) {
	block, err := aes.NewCipher(key.Bytes)
	if err != nil {
		return nil, err
	}
	return &Cipher{block}, nil
}

// Encrypt plain text using AES encryption. Will pad if plaintext message not
// the size of an AES block
func (cip *Cipher) Encrypt(plaintext []byte) ([]byte, error) {
	block := cip.block

	var msg []byte
	if len(plaintext)%aes.BlockSize != 0 {
		nblocks := 1 + len(plaintext)/aes.BlockSize
		msg = make([]byte, nblocks*aes.BlockSize)
	} else {
		msg = make([]byte, len(plaintext))
	}
	copy(msg, plaintext)

	ciphertext := make([]byte, aes.BlockSize+len(msg))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	// To allow decryption of multiple blocks, not just one
	var mode cipher.BlockMode = cipher.NewCBCEncrypter(block, iv)

	// CryptBlocks can work in-place if the two arguments are the same.
	mode.CryptBlocks(ciphertext[aes.BlockSize:], msg)

	return ciphertext, nil
}

// Decrypt cipher text using AES. We assume data stored is UTF-8 text
// and hence a 0, would mark the end of the text
func (cip *Cipher) Decrypt(ciphertext []byte) ([]byte, error) {

	if len(ciphertext) < aes.BlockSize {
		return nil, fmt.Errorf("ciphertext too short. Needs to be larger than a block")
	}

	// The initialization vector is the first block. Also called a nonce. This
	// works a little bit like a salt. It is not secret but adds randomness so the same
	// data does not get encrypted the same way repeatedly.
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	// In cipher-block chaining (CBC) mode we always work with whole blocks,
	// not partial blocks
	if len(ciphertext)%aes.BlockSize != 0 {
		return nil, fmt.Errorf("ciphertext not a multiple of the AES block size")
	}

	// To allow decryption of multiple blocks, not just one
	var mode cipher.BlockMode = cipher.NewCBCDecrypter(cip.block, iv)

	// CryptBlocks can work in-place if the two arguments are the same.
	mode.CryptBlocks(ciphertext, ciphertext)

	// remove padding
	n := bytes.IndexByte(ciphertext, byte(0))
	if n != -1 {
		return ciphertext[:n], nil
	}

	return ciphertext, nil
}

// Decrypt contents of file with filename using AES algorithm operating in block chaining mode
// We assume the encrypted data is stored on file in base32 encoding
func (cip *Cipher) DecryptFile(filename string) (message []byte, err error) {
	filepath := path.Join("data", filename)
	file, err := storage.Open(filepath)
	if err != nil {
		err = fmt.Errorf("could not open encrypted file %s: %w", filename, err)
		return
	}
	defer file.Close()

	var decoder io.Reader = base32.NewDecoder(base32.StdEncoding, file)

	ciphertext, err := io.ReadAll(decoder)
	if err != nil {
		err = fmt.Errorf("unable to decode encrypted file %s: %w", filename, err)
	}

	message, err = cip.Decrypt(ciphertext)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to decrypt %s: %v\n", filename, err)
		os.Exit(1)
	}
	return
}
