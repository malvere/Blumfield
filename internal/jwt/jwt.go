package jwts

import (
	"fmt"
	"log"
	"time"

	"github.com/lestrrat-go/jwx/v2/jwt"
)

func ParseAndCheckToken(tokenString string) error {
	// Parse the token
	token, err := jwt.ParseString(tokenString, jwt.WithVerify(false))
	if err != nil {
		return fmt.Errorf("error parsing token: %w", err)
	}
	log.Print("Tokens is valid for: ", time.Until(token.Expiration()))
	return nil
}

func GetTokenEXP(tokenString string) (time.Duration, error) {
	token, err := jwt.ParseString(tokenString, jwt.WithVerify(false))
	if err != nil {
		return 0, fmt.Errorf("error parsing token: %w", err)
	}
	return time.Until(token.Expiration()), nil
}
