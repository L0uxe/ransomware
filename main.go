package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
)

// Geneate a crypted key of 32 octets
func generateKey() ([]byte, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}
	hexKey := make([]byte, hex.EncodedLen(len(key)))
	hex.Encode(hexKey, key)
	err = os.WriteFile("key.txt", key, 0644)
	if err != nil {
		fmt.Println("Erreur lors de l'enregistrement de la clé :", err)
	}
	return hexKey, nil

}

func main() {
	// Generate key
	key, err := generateKey()
	if err != nil {
		fmt.Println("Erreur lors de la génération de la clé :", err)
		return
	}
	err = os.WriteFile("key.txt", key, 0644)
	if err != nil {
		fmt.Println("Erreur lors de l'enregistrement de la clé :", err)
		return
	}

	// Compile encrypt.go
	encryptCmd := exec.Command("go", "build", "-o", "encrypt", "./encryptF/encrypt.go")
	encryptCmd.Stdout = os.Stdout
	encryptCmd.Stderr = os.Stderr
	encryptCmd.Env = append(os.Environ(), "GOOS=windows", "GOARCH=amd64", fmt.Sprintf("GO_ARGS=%s", key))
	err = encryptCmd.Run()
	if err != nil {
		fmt.Println("Erreur lors de la compilation de encrypt.go :", err)
		return
	}

	// Compile decrypt.go
	decryptCmd := exec.Command("go", "build", "-o", "decrypt", "./decryptF/decrypt.go")
	decryptCmd.Stdout = os.Stdout
	decryptCmd.Stderr = os.Stderr
	decryptCmd.Env = append(os.Environ(), "GOOS=windows", "GOARCH=amd64", fmt.Sprintf("GO_ARGS=%s", string(key)))
	err = decryptCmd.Run()
	if err != nil {
		fmt.Println("Erreur lors de la compilation de decrypt.go :", err)
		return
	}

	fmt.Println("Exécutables encrypt et decrypt créés avec succès.")
}
