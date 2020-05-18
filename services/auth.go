package services

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
)

const CookieName = "millwheatck"

// Auth is the service tasked with creating and validating jwt tokens
type Auth struct {
	secretToken []byte
}

func NewAuth(secretToken []byte) *Auth {
	return &Auth{
		secretToken: secretToken,
	}
}

// UserClaims is a custom struct used as jwt payload
type UserClaims struct {
	UserID string
	*jwt.StandardClaims
}

func (u UserClaims) Valid() error {
	if u.UserID == "" {
		return errors.New("user id missing from claims")
	}
	return nil
}

func (a *Auth) tokenFunc(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}

	return a.secretToken, nil
}

func (a *Auth) Create(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaims{
		UserID: userID,
		StandardClaims: &jwt.StandardClaims{
			NotBefore: time.Now().UTC().Unix(),
		},
	})

	return token.SignedString(a.secretToken)
}

func (a *Auth) ReadFromRequest(r *http.Request) (*UserClaims, error) {
	ext := request.MultiExtractor([]request.Extractor{
		request.OAuth2Extractor,
		CookieExtractor{},
	})
	token, err := request.ParseFromRequest(r, ext, a.tokenFunc, request.WithClaims(&UserClaims{}))
	if err != nil {
		return nil, fmt.Errorf("failed to parse authentication request: %v", err)
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid user claims")
	}

	if !claims.VerifyNotBefore(time.Now().UTC().Unix(), false) {
		return nil, fmt.Errorf("token not valid before: %v", time.Unix(claims.NotBefore, 0))
	}

	return claims, nil
}

type CookieExtractor struct{}

func (c CookieExtractor) ExtractToken(r *http.Request) (string, error) {
	cookie, err := r.Cookie(CookieName)
	if err != nil {
		return "", err
	}

	return cookie.Value, nil
}
