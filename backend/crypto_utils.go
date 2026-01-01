package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

// KeyManager 密钥管理器
type KeyManager struct {
	masterKey []byte
}

// NewKeyManager 创建新的密钥管理器
func NewKeyManager(masterKeyPath string) (*KeyManager, error) {
	// 从环境变量或文件读取主密钥
	masterKey := os.Getenv("MASTER_KEY")
	if masterKey == "" {
		// 如果没有设置，使用默认的（生产环境应该使用环境变量）
		masterKey = "inarbit-master-key-32-bytes-long!"
	}

	return &KeyManager{
		masterKey: []byte(masterKey)[:32], // 确保32字节
	}, nil
}

// EncryptAPIKey 加密API密钥
func (km *KeyManager) EncryptAPIKey(apiKey string) (string, error) {
	block, err := aes.NewCipher(km.masterKey)
	if err != nil {
		return "", fmt.Errorf("创建cipher失败: %v", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("创建GCM失败: %v", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("生成nonce失败: %v", err)
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(apiKey), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptAPIKey 解密API密钥
func (km *KeyManager) DecryptAPIKey(encryptedKey string) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(encryptedKey)
	if err != nil {
		return "", fmt.Errorf("解码失败: %v", err)
	}

	block, err := aes.NewCipher(km.masterKey)
	if err != nil {
		return "", fmt.Errorf("创建cipher失败: %v", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("创建GCM失败: %v", err)
	}

	nonceSize := gcm.NonceSize()
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("解密失败: %v", err)
	}

	return string(plaintext), nil
}

// HashPassword 密码哈希（使用bcrypt）
func HashPassword(password string) (string, error) {
	// 这里应该使用bcrypt库
	// 示例实现
	return hex.EncodeToString([]byte(password)), nil
}

// VerifyPassword 验证密码
func VerifyPassword(hashedPassword, password string) bool {
	// 这里应该使用bcrypt库验证
	return hex.EncodeToString([]byte(password)) == hashedPassword
}

// GenerateSecureToken 生成安全令牌
func GenerateSecureToken(length int) (string, error) {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
