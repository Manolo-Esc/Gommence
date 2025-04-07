package jwt

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const testUserName = "testUser"

func TestHappyPath(t *testing.T) {
	token, err := CreateToken(testUserName)
	assert.Nil(t, err)
	claims, err := ValidateToken(token)
	assert.Nil(t, err)
	assert.Equal(t, testUserName, claims["user"])
}

func TestTokenExpired(t *testing.T) {
	expirated := time.Now().Add(-time.Second * 5).Unix()
	token, err := createToken(testUserName, expirated)
	assert.Nil(t, err)
	_, err = ValidateToken(token)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "token is expired")
}

func TestBadEncryptionKey(t *testing.T) {
	token, err := CreateToken(testUserName)
	assert.Nil(t, err)
	secretKey = []byte("bad_key")
	_, err = ValidateToken(token)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "signature is invalid")
}

// Token created with another key in https://jwt.io/#debugger-io with alg HS256
func TestBadToken1(t *testing.T) {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
	_, err := ValidateToken(token)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "signature is invalid")
}

// Token created with another alg (RS256) in https://jwt.io/#debugger-io
func TestBadToken2(t *testing.T) {
	token := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.NHVaYe26MbtOYhSKkoKYdFVomg4i8ZJd8_-RU8VNbftc4TSMb4bXP3l3YlNWACwyXPGffz5aXHc6lty1Y2t4SWRqGteragsVdZufDn5BlnJl9pdR_kdVFUsra2rWKEofkZeIC4yWytE58sMIihvo9H1ScmmVwBcQP6XETqYd0aSHp1gOa9RdUPDvoXQ5oqygTqVtxaDr6wUFKrKItgBMzWIdNZ6y7O9E0DhEPTbE9rfBo6KTFsHAZnMg4k68CDp2woYIaXbmYTWcvbzIuHO7_37GT79XdIwkm95QJ7hYC9RiwrV7mesbY4PAahERJawntho0my942XheVLmGwLMBkQ"
	_, err := ValidateToken(token)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Unexpected signature method")
}
