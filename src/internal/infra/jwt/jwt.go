package jwt

import (
	"fmt"
	"time"

	"github.com/Manolo-Esc/gommence/src/pkg/logger"
	"github.com/golang-jwt/jwt/v5"
)

// Clave secreta para firmar el token
// "what_makes_a_good_secret": "secret should 64 chars or more with lower case, upper case, numbers and special characters. crypto.randomBytes(64).toString('hex') is a good way to generate a secret",
var secretKey = []byte("mlcfZIl97A930yvsVuoR171aMS3tXBbWqbE1IscEWzvn2w2AzSNF9RnA1Pcnc46DPimhwEGRLC2UaQY6hBow7u") // xxxx

func generarToken(usuario string, exp int64) (string, error) {
	claims := jwt.MapClaims{
		"user": usuario,
		"exp":  exp,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func GenerarToken(usuario string) (string, error) {
	expiration := time.Now().Add(time.Hour * 24).Unix() // 24 hours
	return generarToken(usuario, expiration)
	/*
		claims := jwt.MapClaims{
			"user": usuario,
			"exp":  time.Now().Add(time.Hour * 24).Unix(), // Expira en 24 horas
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(secretKey)
		if err != nil {
			return "", err
		}
		return tokenString, nil
	*/
}

func ValidarToken(tokenString string) (map[string]string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de firma inesperado: %v", token.Header["alg"]) // the text is used in tests!
		}
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Verificar la expiración
		exp, ok := claims["exp"].(float64)
		if !ok || time.Now().Unix() > int64(exp) {
			logger.Assert(false) // jwt.Parse() already checks the expiration (if there is a "exp" claim) so we should not reach this point
			return nil, fmt.Errorf("expiration not valid")
		}
		claimsMap := make(map[string]string)
		// Convertir los claims a un mapa de strings
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

/*
func verificarToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Verificar que el método de firma es el esperado
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			// xxx log it (invalid signature method) -> seguramente nos estan atacando
			return nil, fmt.Errorf("Invalid token 2: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}

func extraerClaims(tokenString string) (jwt.MapClaims, error) {
	token, err := verificarToken(tokenString)
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token 1")
}
*/
