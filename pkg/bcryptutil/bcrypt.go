package bcryptutil

import (
	"golang.org/x/crypto/bcrypt"
)

type BcryptUtils interface {
	GenerateHash(s string) (string, error)
	CompareHash(s string, hash string) bool
}

type BcryptUtilsImpl struct{}

func (b *BcryptUtilsImpl) GenerateHash(s string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (b *BcryptUtilsImpl) CompareHash(s string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(s))
	return err == nil
}
