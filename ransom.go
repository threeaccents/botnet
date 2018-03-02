package botnet

//Ransomer is
type Ransomer interface {
	Encrypt(dir string) error
	Key() []byte
	Decrypt(dir string) error
}
