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

![image](https://github.com/user-attachments/assets/9931b815-39d5-49bb-baa4-567cc1bc7c4e)
![image](https://github.com/user-attachments/assets/8d01a535-bff5-4980-81f8-1da55f3fa602)
![image](https://github.com/user-attachments/assets/a8763fce-5dbc-4c87-8dc3-357b55be24e8)
