package jwt

import (
	"errors"
	"gin-scaffold/internal/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
	"time"
)

const (
	DefaultExpire = time.Second * 86400 * 7
)

var (
	DefaultAlg = jwt.SigningMethodHS256
)

type Token struct {
	Key    string
	Alg    jwt.SigningMethod
	Body   map[string]interface{}
	Claims jwt.StandardClaims
}

type OptionFunc func(*Token)

type Claims struct {
	Body map[string]interface{} `json:"body"`
	jwt.StandardClaims
}

func WithKey(key string) OptionFunc {
	return func(j *Token) {
		j.Key = key
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

	if viper.IsSet("jwt.key") && viper.GetString("jwt.key") != "" {
		j.Key = viper.GetString("jwt.key")
	}

	if viper.IsSet("jwt.expire") && viper.GetInt64("jwt.expire") > 0 {
		j.Claims.ExpiresAt = time.Now().Unix() + viper.GetInt64("jwt.expire")
	}

	for _, f := range options {
		f(j)
	}

	if j.Key == "" {
		return nil, errors.New("缺少密钥")
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
			err = errors.New("token 格式错误")
		case jwt.ValidationErrorUnverifiable:
			err = errors.New("签名无效，无法验证令牌")
		case jwt.ValidationErrorSignatureInvalid:
			err = errors.New("签名验证失败")
		case jwt.ValidationErrorAudience:
			err = errors.New("token 身份验证失败")
		case jwt.ValidationErrorExpired:
			err = errors.New("token 已过期")
		case jwt.ValidationErrorIssuedAt:
			err = errors.New("token 签发时间验证失败")
		case jwt.ValidationErrorIssuer:
			err = errors.New("token 签发身份验证失败")
		case jwt.ValidationErrorNotValidYet:
			err = errors.New("token 暂不可用")
		case jwt.ValidationErrorId:
			err = errors.New("token 标识验证失败")
		case jwt.ValidationErrorClaimsInvalid:
			err = errors.New("token 结构体验证失败")
		default:
			err = validateError
		}
	} else {
		err = validateError
	}

	return
}
