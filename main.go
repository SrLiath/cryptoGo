package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"golang.org/x/crypto/scrypt"
)

const (
	saltSize      = 16
	nonceSize     = 12
	keyLen        = 32
	maxAttempts   = 10
	maxFails      = 3
	lockoutPeriod = 10 * time.Second
)

var (
	attempts     int
	failures     int
	lastAttempt  time.Time
	encryptedDir = "cryptos/" // Directory to store encrypted files
	fileList     *widget.List // Make fileList a global variable
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Notas Criptografadas")

	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetPlaceHolder("Digite a senha")

	fileNameEntry := widget.NewEntry()
	fileNameEntry.SetPlaceHolder("Nome do arquivo criptografado")

	encryptButton := widget.NewButton("Criptografar", func() {
		handleEncryption(myWindow, passwordEntry.Text, fileNameEntry.Text)
	})

	decryptButton := widget.NewButton("Descriptografar", func() {
		handleDecryption(myWindow, passwordEntry.Text, fileNameEntry.Text, fileNameEntry)
	})

	fileList = widget.NewList(
		func() int { return len(listFiles()) }, // Number of items
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			item.(*widget.Label).SetText(listFiles()[id])
		},
	)

	content := container.NewVBox(
		widget.NewLabel("Senha:"),
		passwordEntry,
		widget.NewLabel("Nome do Arquivo:"),
		fileNameEntry,
		container.NewHBox(encryptButton, decryptButton),
		widget.NewLabel("Arquivos Criptografados:"),
		fileList,
	)

	myWindow.SetContent(content)
	myWindow.Resize(fyne.NewSize(600, 250))
	myWindow.ShowAndRun()
}

func handleEncryption(window fyne.Window, password, fileName string) {
	if attemptsExceeded() {
		dialog.ShowInformation("Erro", "Muitas tentativas erradas. Tente novamente mais tarde.", window)
		return
	}

	if fileName == "" {
		dialog.ShowInformation("Erro", "Nome do arquivo não pode estar vazio.", window)
		return
	}

	err := encryptContentToFile(fileName, "Conteúdo da nota", password)
	if err != nil {
		dialog.ShowInformation("Erro", fmt.Sprintf("Erro ao criptografar: %v", err), window)
	} else {
		dialog.ShowInformation("Sucesso", "Notas criptografadas com sucesso.", window)
		updateFileList()
		resetAttempts()
	}
}

func handleDecryption(window fyne.Window, password, fileName string, fileNameEntry *widget.Entry) {
	if attemptsExceeded() {
		dialog.ShowInformation("Erro", "Muitas tentativas erradas. Tente novamente mais tarde.", window)
		return
	}

	if fileName == "" {
		dialog.ShowInformation("Erro", "Nome do arquivo não pode estar vazio.", window)
		return
	}

	notes, err := decryptFile(fileName, password)
	if err != nil {
		attempts++
		failures++
		if failures >= maxFails {
			dialog.ShowInformation("Erro", fmt.Sprintf("Tentativas excedidas. Arquivo será apagado: %v", err), window)
			os.Remove(encryptedDir + "/" + fileName) // Verifique o caminho correto
			resetAttempts()
		} else {
			if attempts >= maxAttempts {
				dialog.ShowInformation("Erro", "Muitas tentativas erradas. Tente novamente mais tarde.", window)
				time.Sleep(lockoutPeriod)
				resetAttempts()
			} else {
				dialog.ShowInformation("Erro", fmt.Sprintf("Senha incorreta."), window)
			}
		}
		return
	}

	notesWindow(window, password, notes, fileNameEntry)
	resetAttempts()
}

func notesWindow(parent fyne.Window, password, notes string, fileNameEntry *widget.Entry) {
	window := fyne.CurrentApp().NewWindow("Notas")
	notesEntry := widget.NewMultiLineEntry()
	notesEntry.SetText(notes)
	notesEntry.SetPlaceHolder("Digite suas notas aqui...")
	notesEntry.Wrapping = fyne.TextWrapWord
	// Removed SetTextStyle

	saveButton := widget.NewButton("Salvar", func() {
		notes := notesEntry.Text
		fileName := fileNameEntry.Text
		fmt.Printf("Senha usada para salvar: %s\n", password) // Debug: senha usada para salvar

		err := encryptContentToFile(fileName, notes, password)
		if err != nil {
			dialog.ShowInformation("Erro", fmt.Sprintf("Erro ao salvar as notas: %v", err), window)
		} else {
			dialog.ShowInformation("Sucesso", "Notas salvas com sucesso.", window)
		}
		updateFileList()
	})

	content := container.NewVBox(
		widget.NewLabel("Notas:"),
		notesEntry,
		saveButton,
	)

	window.SetContent(content)
	window.Resize(fyne.NewSize(600, 400))
	window.Show()
}

func encryptContentToFile(outputPath, content, password string) error {
	fmt.Println("Iniciando criptografia...")

	// Create directory if it does not exist
	if err := os.MkdirAll(encryptedDir, os.ModePerm); err != nil {
		return err
	}

	salt := deriveSalt()
	fmt.Printf("Salt derivado: %s\n", hex.EncodeToString(salt))

	key := deriveKey([]byte(password), salt)
	fmt.Printf("Chave derivada: %s\n", hex.EncodeToString(key))

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	nonce := make([]byte, nonceSize)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return err
	}
	fmt.Printf("Nonce gerado: %s\n", hex.EncodeToString(nonce))

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(content), nil)
	fmt.Printf("Texto criptografado: %s\n", hex.EncodeToString(ciphertext))

	dataToSave := append(salt, ciphertext...)
	return os.WriteFile(encryptedDir+"/"+outputPath, dataToSave, 0644)
}

func decryptFile(inputPath, password string) (string, error) {
	fmt.Println("Iniciando descriptografia...")

	data, err := os.ReadFile(encryptedDir + "/" + inputPath)
	if err != nil {
		return "", err
	}
	fmt.Printf("Texto criptografado lido: %s\n", hex.EncodeToString(data))

	if len(data) < saltSize+nonceSize {
		return "", fmt.Errorf("dados corrompidos")
	}

	salt, ciphertext := data[:saltSize], data[saltSize:]
	fmt.Printf("Salt extraído: %s\n", hex.EncodeToString(salt))

	key := deriveKey([]byte(password), salt)
	fmt.Printf("Chave derivada: %s\n", hex.EncodeToString(key))

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	if len(ciphertext) < nonceSize {
		return "", fmt.Errorf("dados corrompidos")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	fmt.Printf("Nonce extraído: %s\n", hex.EncodeToString(nonce))
	fmt.Printf("Texto criptografado restante: %s\n", hex.EncodeToString(ciphertext))

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

func deriveKey(password, salt []byte) []byte {
	key, err := scrypt.Key(password, salt, 16384, 8, 1, keyLen)
	if err != nil {
		panic(err)
	}
	return key
}

func deriveSalt() []byte {
	salt := make([]byte, saltSize)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		panic(err)
	}
	return salt
}

func attemptsExceeded() bool {
	if time.Since(lastAttempt) > lockoutPeriod {
		resetAttempts()
	}
	return attempts >= maxAttempts
}

func resetAttempts() {
	attempts = 0
	failures = 0
}

func listFiles() []string {
	files, err := os.ReadDir(encryptedDir)
	if err != nil {
		return nil
	}
	var fileNames []string
	for _, file := range files {
		if !file.IsDir() {
			fileNames = append(fileNames, file.Name())
		}
	}
	return fileNames
}

func updateFileList() {
	fileList.Refresh() // Refresh the list view to update the file list
}
