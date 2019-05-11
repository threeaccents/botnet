package main

import (
	"log"

	"github.com/threeaccents/botnet/aes"
	"github.com/threeaccents/botnet/attacks/ransom"
)

func main() {
	cryptoService := &aes.CryptoService{
		Key: []byte("O36febOjxrZPZEodlTEQQzEzuurhtMm8"),
	}

	r := ransom.RansomService{
		CryptoService: cryptoService,
		InitialDir:    "./cmd/ransomtest/data",
	}

	if err := r.Reverse(); err != nil {
		log.Println(err)
	}
}
