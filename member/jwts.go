package member

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"neckless.adviser.com/key"
)

type PublicMemberClaim struct {
	JsonPublicMember
	jwt.StandardClaims
}

func MakePublicMemberJWT(signer *key.PrivateKey, pm *PublicMember) (string, error) {
	claims := &PublicMemberClaim{
		JsonPublicMember: *pm.AsJSON(),
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: pm.ValidUntil.Unix(),
			Id:        pm.Id,
			IssuedAt:  time.Now().Unix(),
			Subject:   "PublicMemberClaim",
		},
	}
	// js, err := json.Marshal(claims)
	// return string(js), err
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(signer.Key.Raw[:])
	return tokenString, err
}

func VerifyJWT(pk *key.PrivateKey, tknStr string) (*PublicMemberClaim, *jwt.Token, error) {
	claims := PublicMemberClaim{}
	token, err := jwt.ParseWithClaims(tknStr, &claims, func(token *jwt.Token) (interface{}, error) {
		return pk.Key.Raw[:], nil
	})
	return &claims, token, err
}
