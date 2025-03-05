package main

import (
	"fmt"

	"aidanwoods.dev/go-paseto"
)

func main() {
	k := paseto.NewV4AsymmetricSecretKey()
	sk := k.ExportHex()
	pk := k.Public().ExportHex()
	encrypt := paseto.NewV4SymmetricKey().ExportHex()
	fmt.Printf("PUBLIC_KEY=%s\nPRIVATE_KEY=%s\nSYMMETRIC_KEY=%s\n", pk, sk, encrypt)
}
