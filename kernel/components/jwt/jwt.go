package jwt

import (
	"errors"
	"gin-scaffold/kernel/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
	"time"
)

type Jwt struct {
	Key       string
	Method    jwt.SigningMethod
	Audience  string
	ExpiresAt int64
	Id        string
	IssuedAt  int64
	Issuer    string
	NotBefore int64
	Subject   string
	Body      map[string]interface{}
}

type Claims struct {
	Body map[string]interface{} `json:"body"`
	jwt.StandardClaims
}

type OptionFunc func(*Jwt)

func WithKey(key string) OptionFunc {
	return func(j *Jwt) {
		j.Key = key
	}
}

func WithMethod(alg jwt.SigningMethod) OptionFunc {
	return func(j *Jwt) {
		j.Method = alg
	}
}

func WithAudience(aud string) OptionFunc {
	return func(j *Jwt) {
		j.Audience = aud
	}
}

func WithExpiresAt(exp int64) OptionFunc {
	return func(j *Jwt) {
		j.ExpiresAt = exp
	}
}

func WithId(jti string) OptionFunc {
	return func(j *Jwt) {
		j.Id = jti
	}
}

func WithIssuedAt(issuedAt int64) OptionFunc {
	return func(j *Jwt) {
		j.IssuedAt = issuedAt
	}
}

func WithIssuer(iss string) OptionFunc {
	return func(j *Jwt) {
		j.Issuer = iss
	}
}

func WithNotBefore(nbf int64) OptionFunc {
	return func(j *Jwt) {
		j.NotBefore = nbf
	}
}

func WithSubject(sub string) OptionFunc {
	return func(j *Jwt) {
		j.Subject = sub
	}
}

func WithBody(body map[string]interface{}) OptionFunc {
	return func(j *Jwt) {
		j.Body = body
	}
}

func NewJwt(options ...OptionFunc) (*Jwt, error) {
	jwt := &Jwt{
		Key:       "",
		Method:    jwt.SigningMethodHS256,
		Audience:  "server",
		ExpiresAt: time.Now().Add(time.Hour * 24 * 7).Unix(),
		Id:        utils.RandomString(64),
		IssuedAt:  time.Now().Unix(),
		Issuer:    "app",
		NotBefore: time.Now().Unix(),
		Subject:   "authentication token",
		Body:      map[string]interface{}{},
	}

	if viper.IsSet("jwt.key") && len(viper.GetString("jwt.key")) > 0 {
		jwt.Key = viper.GetString("jwt.key")
	}

	if viper.IsSet("jwt.expire") && viper.GetInt64("jwt.expire") > 0 {
		jwt.ExpiresAt = time.Now().Unix() + viper.GetInt64("jwt.expire")
	}

	for _, f := range options {
		f(jwt)
	}

	if len(jwt.Key) == 0 {
		return nil, errors.New("缺少密钥")
	}

	return jwt, nil
}

func (t *Jwt) Make() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		Body: t.Body,
		StandardClaims: jwt.StandardClaims{
			Audience:  t.Audience,
			ExpiresAt: t.ExpiresAt,
			Id:        t.Id,
			IssuedAt:  t.IssuedAt,
			Issuer:    t.Issuer,
			NotBefore: t.NotBefore,
			Subject:   t.Subject,
		},
	})

	signedToken, err := token.SignedString([]byte(t.Key))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (t *Jwt) Parse(token string) (*jwt.Token, *Claims, error) {
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
	validateError := parseError.(*jwt.ValidationError)

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
	}

	return
}
