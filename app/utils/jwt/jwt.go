package jwt

import (
	"errors"
	"gin-scaffold/core/utils"
	"github.com/dgrijalva/jwt-go"
	"time"
)

const (
	DefaultExpire = time.Second * 86400 * 7
)

var (
	DefaultAlg = jwt.SigningMethodHS256

	ErrMissingKey       = errors.New("token missing key")
	ErrMalformed        = errors.New("token 格式错误")
	ErrUnverifiable     = errors.New("签名无效，无法验证令牌")
	ErrSignatureInvalid = errors.New("签名验证失败")
	ErrAudience         = errors.New("token 身份验证失败")
	ErrExpired          = errors.New("token 已过期")
	ErrIssuedAt         = errors.New("token 签发时间验证失败")
	ErrIssuer           = errors.New("token 签发身份验证失败")
	ErrNotValidYet      = errors.New("token 暂不可用")
	ErrId               = errors.New("token 标识验证失败")
	ErrClaimsInvalid    = errors.New("token 结构体验证失败")
)

type Token struct {
	Key    string
	Expire time.Duration
	Alg    jwt.SigningMethod
	Body   map[string]interface{}
	Claims jwt.StandardClaims
}

type OptionFunc func(*Token)

type Claims struct {
	Body map[string]interface{} `json:"body"`
	jwt.StandardClaims
}

func WithKey(k string) OptionFunc {
	return func(j *Token) {
		j.Key = k
	}
}

func WithExpire(e time.Duration) OptionFunc {
	return func(j *Token) {
		j.Claims.ExpiresAt = time.Now().Add(time.Second * e).Unix()
	}
}

func WithAlg(alg jwt.SigningMethod) OptionFunc {
	return func(j *Token) {
		j.Alg = alg
	}
}

func WithBody(body map[string]interface{}) OptionFunc {
	return func(j *Token) {
		j.Body = body
	}
}

func WithClaims(claims jwt.StandardClaims) OptionFunc {
	return func(j *Token) {
		j.Claims = claims
	}
}

func NewToken(options ...OptionFunc) (*Token, error) {
	j := &Token{
		Key:  "",
		Alg:  DefaultAlg,
		Body: map[string]interface{}{},
		Claims: jwt.StandardClaims{
			Audience:  "server",
			ExpiresAt: time.Now().Add(DefaultExpire).Unix(),
			Id:        utils.RandomString(64),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "app",
			NotBefore: time.Now().Unix(),
			Subject:   "authentication token",
		},
	}

	for _, f := range options {
		f(j)
	}

	if j.Key == "" {
		return nil, ErrMissingKey
	}

	return j, nil
}

func (t *Token) Create() (string, error) {
	token := jwt.NewWithClaims(t.Alg, Claims{
		Body: t.Body,
		StandardClaims: jwt.StandardClaims{
			Audience:  t.Claims.Audience,
			ExpiresAt: t.Claims.ExpiresAt,
			Id:        t.Claims.Id,
			IssuedAt:  t.Claims.IssuedAt,
			Issuer:    t.Claims.Issuer,
			NotBefore: t.Claims.NotBefore,
			Subject:   t.Claims.Subject,
		},
	})

	signedToken, err := token.SignedString([]byte(t.Key))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (t *Token) Parse(token string) (*jwt.Token, *Claims, error) {
	claims := &Claims{}

	parsedToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(t.Key), nil
	})

	if err != nil {
		err = handleParseError(err)
	}

	return parsedToken, claims, err
}

func handleParseError(parseError error) (err error) {
	validateError, ok := parseError.(*jwt.ValidationError)

	if ok {
		switch validateError.Errors {
		case jwt.ValidationErrorMalformed:
			err = ErrMalformed
		case jwt.ValidationErrorUnverifiable:
			err = ErrUnverifiable
		case jwt.ValidationErrorSignatureInvalid:
			err = ErrSignatureInvalid
		case jwt.ValidationErrorAudience:
			err = ErrAudience
		case jwt.ValidationErrorExpired:
			err = ErrExpired
		case jwt.ValidationErrorIssuedAt:
			err = ErrIssuedAt
		case jwt.ValidationErrorIssuer:
			err = ErrIssuer
		case jwt.ValidationErrorNotValidYet:
			err = ErrNotValidYet
		case jwt.ValidationErrorId:
			err = ErrId
		case jwt.ValidationErrorClaimsInvalid:
			err = ErrClaimsInvalid
		default:
			err = validateError
		}
	} else {
		err = validateError
	}

	return
}
