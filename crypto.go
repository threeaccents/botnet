package botnet

// CryptoService is the best for me
type CryptoService interface {
	Encrypt(plaintext []byte) ([]byte, error)
	EncryptToString(plaintext []byte) (string, error)
	Decrypt(cipherText []byte) ([]byte, error)
	DecryptString(cipherText string) ([]byte, error)
}