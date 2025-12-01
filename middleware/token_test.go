package middleware

import (
	"org_chart/initializers"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	initializers.InitViper("../config.yaml")

}

func TestValidateToken(t *testing.T) {
	//Private Key Check
	prvKey, err := getPrvKey()
	require.NoError(t, err, "private.pem should be readable")
	require.NotNil(t, prvKey, "getPrvKey should not return nil value for private key")
	//Public Key Check
	pubKey, err := getPubKey()
	assert.NoError(t, err, "public.pem should be readable")
	assert.NotNil(t, pubKey, "getPubKey should not return nil value for Public key")

	/////////////////////////////////////////////////////////////////////////////
	// GenerateSignedToken with NIL claims - Error
	tokenStr, err := generateSignedToken(
		nil,
		false,
		prvKey,
		time.Now().Add(5*time.Minute),
	)
	assert.Error(t, err, "Nil claims should return error")
	assert.Empty(t, tokenStr, "There should be no token generation")
	// GenerateSignedToken without claims - Error
	tokenStr, err = generateSignedToken(
		&Claims{},
		false,
		prvKey,
		time.Now().Add(1*time.Second),
	)
	assert.Error(t, err, "Empty claims should return error")
	assert.Empty(t, tokenStr, "There should be no token generation")
	//GenerateSignedToken with claims - No error
	claims := &Claims{
		UserID: 123,
	}
	tokenStr, err = generateSignedToken(
		claims,
		false,
		prvKey,
		time.Now().Add(1*time.Second),
	)
	require.NoError(t, err, "No errors should be thrown")
	require.NotEmpty(t, tokenStr, "A JWT token should be returned")

	//////////////////////////////////////////////////////////////////////////////////////////
	// ValidateToke with empty token - Error
	validated, err := ValidateToken("")
	assert.Error(t, err, "Empty token should return error")
	assert.Nil(t, validated, "Nil claims should be retured for empty token")
	// ValidateToken with invaid token - Error
	validated, err = ValidateToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMTIzIiwiZW1haWwiOiJjaHJpc0BtYWlsLmNvbSIsImlhdCI6MTUxNjIzOTAyMn0.cKYye5weZzMhLfVS1ApBYjjGzIjfAJRdsfx1quHvm0Q")
	assert.Error(t, err, "The token doesn't use private key for signing. So the error is valid")
	assert.Nil(t, validated, "Nil claims should be retured for invalid token")
	// ValidateToken with valid token - No Error
	validated, err = ValidateToken(tokenStr)
	require.NoError(t, err, "ValidateToken returns error")
	require.NotNil(t, validated, "ValidateToken returns Nil claims")
	require.Equal(t, claims.UserID, validated.UserID)
	require.Equal(t, false, validated.IsRefresh)
	require.True(t, validated.ExpiresAt.After(time.Now()))
}
