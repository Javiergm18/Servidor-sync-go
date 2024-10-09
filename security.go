package main

import (
	"os/exec"
)

// Cifrado de archivos usando GPG
func EncryptFile(filePath, recipient string) error {
	cmd := exec.Command("gpg", "--encrypt", "--recipient", recipient, filePath)
	return cmd.Run()
}

// Descifrar archivos
func DecryptFile(filePath, passphrase string) error {
	cmd := exec.Command("gpg", "--decrypt", "--passphrase", passphrase, filePath)
	return cmd.Run()
}
