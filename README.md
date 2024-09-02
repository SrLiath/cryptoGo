# cryptoGo

### Encrypted Notes Application

This Go application provides a simple graphical user interface \\(GUI\\) for encrypting and decrypting notes. It uses AES encryption with GCM mode for secure storage and retrieval of notes.

#### Features
- **Encrypt Notes**: Encrypts and stores notes in a specified file.
- **Decrypt Notes**: Decrypts and displays notes from a specified file.
- **Password Protection**: Requires a password for encryption and decryption.
- **File Management**: Lists encrypted files and manages encryption/decryption attempts.

#### Dependencies
- `fyne.io/fyne/v2`: For building the GUI.
- `golang.org/x/crypto/scrypt`: For key derivation.
#### Installation
1. Ensure you have Go installed on your system.
2. Install the required Go packages using:
   ```bash
   go get fyne.io/fyne/v2
   go get golang.org/x/crypto/scrypt
   ```
3. Run with 
   ```bash
    go run main.go
   ```

-- Or just Download the release option