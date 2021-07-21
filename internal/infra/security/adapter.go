package security

import (
	"backend/internal/domain"
	"bytes"
	crand "crypto/rand"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/jwtauth"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/argon2"
	"io/ioutil"
	"os"
	"strconv"
	"time"
)

type adapter struct {
	logger *logrus.Logger
	config *Config

	// jwt
	jwtAuth *jwtauth.JWTAuth
}

func NewAdapter(logger *logrus.Logger, config *Config) (domain.Security, error) {
	a := &adapter{
		logger: logger,
		config: config,
	}

	// Read JWT signing key
	fileWithAccessToken, err := os.Open(a.config.JWTPrivateKey)
	if err != nil {
		a.logger.WithError(err).Error("Error while opening file!")
		return nil, err
	}
	//noinspection ALL
	defer fileWithAccessToken.Close()

	bts, err := ioutil.ReadAll(fileWithAccessToken)
	if err != nil {
		a.logger.WithError(err).Error("Error while reading file!")
		return nil, err
	}

	jwtAuth := jwtauth.New(jwt.SigningMethodHS256.Name, bts, nil)
	a.jwtAuth = jwtAuth

	return a, nil
}

func (a *adapter) GenerateNewJWT(id int, duration time.Duration) (string, error) {
	claims := jwt.StandardClaims{
		ExpiresAt: time.Now().Add(duration).Unix(),
		IssuedAt:  time.Now().Unix(),
		Issuer:    "sharito",
		Subject:   strconv.Itoa(id),
	}

	_, tokenString, err := a.jwtAuth.Encode(claims)
	if err != nil {
		a.logger.WithError(err).Error("Error while encoding JWT token!")
		return "", domain.ErrJWT
	}

	return tokenString, nil
}

func (a *adapter) HashPassword(password string) ([]byte, []byte, error) {
	salt, err := getRandomBytes(64)
	if err != nil {
		a.logger.WithError(err).Error("cannot get random bytes for salt")
		return nil, nil, domain.ErrInternalSecurity
	}

	return getHashPassword(salt, password), salt, nil
}

func getRandomBytes(length int) ([]byte, error) {
	salt := make([]byte, length)

	_, err := crand.Read(salt)
	if err != nil {
		return nil, err
	}

	return salt, nil
}

func (a *adapter) VerifyPassword(salt []byte, passwordHash []byte, password string) bool {
	return bytes.Equal(getHashPassword(salt, password), passwordHash)
}

func getHashPassword(salt []byte, password string) []byte {
	return argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 64)
}
