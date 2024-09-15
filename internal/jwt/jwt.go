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
	log.Print(time.Until(token.Expiration()))

	fmt.Println("Token is valid and not expired.")
	return nil
}
