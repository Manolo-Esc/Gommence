package jwt

import (
	"fmt"
	"time"

	"github.com/Manolo-Esc/gommence/src/pkg/logger"
	"github.com/golang-jwt/jwt/v5"
)

// Secret key for signing the JWT tokens
// "what_makes_a_good_secret": "secret should 64 chars or more with lower case, upper case, numbers and special characters. crypto.randomBytes(64).toString('hex') is a good way to generate a secret",
var secretKey = []byte("mlcfZIl97A930yvsVuoR171aMS3tXBbWqbE1IscEWzvn2w2AzSNF9RnA1Pcnc46DPimhwEGRLC2UaQY6hBow7u")

func createToken(user string, exp int64) (string, error) {
	claims := jwt.MapClaims{
		"user": user,
		"exp":  exp,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func CreateToken(user string) (string, error) {
	expiration := time.Now().Add(time.Hour * 24).Unix() // 24 hours
	return createToken(user, expiration)
}

func ValidateToken(tokenString string) (map[string]string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signature method: %v", token.Header["alg"]) // the text is used in tests!
		}
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Check expiration
		exp, ok := claims["exp"].(float64)
		if !ok || time.Now().Unix() > int64(exp) {
			logger.Assert(false) // jwt.Parse() already checks the expiration (if there is a "exp" claim) so we should not have reached this point
			return nil, fmt.Errorf("expiration not valid")
		}
		claimsMap := make(map[string]string)
		// Convert claims to string map
		for key, value := range claims {
			switch v := value.(type) {
			case string:
				claimsMap[key] = v
			case float64:
				claimsMap[key] = fmt.Sprintf("%.0f", v)
			default:
				claimsMap[key] = fmt.Sprintf("%v", v)
			}
		}

		return claimsMap, nil
	}

	return nil, fmt.Errorf("invalid token")
}
