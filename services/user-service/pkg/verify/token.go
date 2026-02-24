package verify

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/peterouob/seckill_service/services/user-service/pkg/configs"
	"github.com/peterouob/seckill_service/utils/logs"

	"sync/atomic"
	"time"
)

var (
	err        error
	TokenKey   atomic.Value
	RefreshKey atomic.Value
)

type TokenInterface interface {
	CreateToken()
	CreateRefreshToken()
}

type Token struct {
	UserId       int64         `json:"user_id"`
	AccessId     string        `json:"access_id"`
	AccessToken  string        `json:"access_token"`
	RefreshId    string        `json:"refresh_id"`
	RefreshToken string        `json:"refresh_token"`
	Token        configs.Token `json:"token"`
}

var _ TokenInterface = (*Token)(nil)

func NewToken(id int64) *Token {
	TokenKey.Store(configs.Config.GetString("token.token_key"))
	RefreshKey.Store(configs.Config.GetString("token.refresh_key"))
	token := &configs.Token{}
	token.AccessUuid = uuid.NewString()
	token.RefreshUuid = uuid.NewString()
	token.AtExpires = time.Now().Add(time.Hour * 2).Unix()
	token.RefreshAtExpires = time.Now().Add(time.Hour * 24 * 7 * 2).Unix()
	return &Token{
		UserId: id,
		Token:  *token,
	}
}

// CreateToken  不存Redis單純驗證
func (t *Token) CreateToken() {
	claims := jwt.MapClaims{
		"access_id": t.Token.AccessUuid,
		"exp":       t.Token.AtExpires,
		"type":      "access",
		"userId":    t.UserId,
		"jti":       t.UserId,
		"iat":       time.Now().Unix(),
	}

	tk := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t.AccessToken, err = tk.SignedString([]byte(TokenKey.Load().(string)))
	t.AccessId = claims["access_id"].(string)
}

func (t *Token) CreateRefreshToken() {
	claims := jwt.MapClaims{
		"refresh_id": t.Token.RefreshUuid,
		"exp":        t.Token.RefreshAtExpires,
		"type":       "refresh",
		"userId":     t.UserId,
		"jti":        t.UserId,
		"iat":        time.Now().Unix(),
	}
	tk := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t.RefreshToken = string(fmt.Appendf(tk.Signature, fmt.Sprintf("%s%d", RefreshKey.Load().(string), t.UserId)))

	logs.HandelError("create refresh token error", err)
	t.RefreshId = claims["refresh_id"].(string)
}

func TokenVerify(tokenString string) *jwt.Token {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			logs.HandelError("parse token error type", err)
		}
		return []byte(TokenKey.Load().(string)), nil
	})
	logs.HandelError("parse token error", err)
	switch {
	case token.Valid:
		logs.Log("valid success token")
	case errors.Is(err, jwt.ErrTokenMalformed):
		logs.Log("error in Malformed token type")
	case errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet):
		logs.Log("error in token expired")
	default:
		logs.HandelError("couldn't handle this token", err)
	}
	return token
}
