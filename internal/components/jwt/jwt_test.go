package jwt

import (
	"github.com/dgrijalva/jwt-go"
	"math"
	"testing"
	"time"
)

var TestKey = "123456"

func TestNewToken(t *testing.T) {
	excepts := map[string]map[string]string{
		"withoutKey": {"key": ""},
		"withKey":    {"key": TestKey},
	}

	for k, v := range excepts {
		if k == "withoutKey" {
			t.Run(k, func(t *testing.T) {
				token, err := NewToken()
				if token != nil || err != ErrMissingKey {
					t.Errorf("创建 Token 对象出错，对象：%v，错误信息：%v", token, err)
				}
			})
		} else if k == "withKey" {
			t.Run(k, func(t *testing.T) {
				token, err := NewToken(WithKey(v["key"]))
				if token == nil || err != nil {
					t.Errorf("创建 Token 对象出错，对象：%v，错误信息：%v", token, err)
				}
			})
		}
	}
}

func TestTokenCreate(t *testing.T) {
	now, err := time.Parse("2006-01-02 15:04:05", "1995-01-01 00:00:00")
	if err != nil {
		t.Error(err)
	}
	exceptToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJib2R5Ijp7Im5hbWUiOiJUb20iLCJzZXgiOiJtYWxlIn0sImF1ZCI6InNlcnZlciIsImV4cCI6NTA4Mzg4NTY5NSwianRpIjoiYWJjZGUiLCJpYXQiOjc4ODkxODQwMCwiaXNzIjoiYXBwIiwibmJmIjo3ODg5MTg0MDAsInN1YiI6ImF1dGhlbnRpY2F0aW9uIHRva2VuIn0.KFq7Nea4aJEPuBr_sEGSHZKbvJheqb8uUp4ilPwdP7o"

	token, err := NewToken(
		WithKey(TestKey),
		WithAlg(DefaultAlg),
		WithBody(map[string]interface{}{
			"name": "Tom",
			"sex":  "male",
		}),
		WithClaims(jwt.StandardClaims{
			Audience:  "server",
			ExpiresAt: now.Add(time.Second * time.Duration(math.MaxUint32)).Unix(),
			Id:        "abcde",
			IssuedAt:  now.Unix(),
			Issuer:    "app",
			NotBefore: now.Unix(),
			Subject:   "authentication token",
		}),
	)
	if err != nil {
		t.Errorf("创建 Token 对象出错，错误信息：%v", err)
	}

	t.Run("createToken", func(t *testing.T) {
		signedToken, err := token.Create()
		if err != nil {
			t.Errorf("创建 JWT 字符串出错，错误信息：%v", err)
		}

		if signedToken != exceptToken {
			t.Errorf("创建 JWT 字符串出错，实际的字符串：%s", signedToken)
		}
	})
}

func TestTokenParse(t *testing.T) {
	signedToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJib2R5Ijp7Im5hbWUiOiJUb20iLCJzZXgiOiJtYWxlIn0sImF1ZCI6InNlcnZlciIsImV4cCI6NTA4Mzg4NTY5NSwianRpIjoiYWJjZGUiLCJpYXQiOjc4ODkxODQwMCwiaXNzIjoiYXBwIiwibmJmIjo3ODg5MTg0MDAsInN1YiI6ImF1dGhlbnRpY2F0aW9uIHRva2VuIn0.KFq7Nea4aJEPuBr_sEGSHZKbvJheqb8uUp4ilPwdP7o"

	token, err := NewToken(
		WithKey(TestKey),
	)
	if err != nil {
		t.Errorf("创建 Token 对象出错，错误信息：%v", err)
	}

	t.Run("parseToken", func(t *testing.T) {
		parsedToken, claims, err := token.Parse(signedToken)
		if err != nil {
			t.Errorf("解析 Token 字符串出错，错误信息：%v", err)
		}

		if parsedToken == nil {
			t.Error("验证 Token 字符串失败，期待返回解析后的 Token 对象")
		} else {
			if !parsedToken.Valid {
				t.Error("验证 Token 字符串失败，期待验证通过")
			}
		}

		if claims.Body["name"] != "Tom" || claims.Body["sex"] != "male" {
			t.Errorf("验证解析后的 Claims 出错，Claims：%v", claims)
		}
	})
}

func TestTokenExpire(t *testing.T) {
	now, err := time.Parse("2006-01-02 15:04:05", "1995-01-01 00:00:00")
	if err != nil {
		t.Error(err)
	}

	token, err := NewToken(
		WithKey(TestKey),
		WithAlg(DefaultAlg),
		WithBody(map[string]interface{}{
			"name": "Tom",
			"sex":  "male",
		}),
		WithClaims(jwt.StandardClaims{
			Audience:  "server",
			ExpiresAt: now.Unix(),
			Id:        "abcde",
			IssuedAt:  now.Unix(),
			Issuer:    "app",
			NotBefore: now.Unix(),
			Subject:   "authentication token",
		}),
	)
	if err != nil {
		t.Errorf("创建 Token 对象出错，错误信息：%v", err)
	}

	signedToken, err := token.Create()
	if err != nil {
		t.Errorf("创建 JWT 字符串出错，错误信息：%v", err)
	}

	t.Run("tokenExpire", func(t *testing.T) {
		token, err := NewToken(
			WithKey(TestKey),
		)
		if err != nil {
			t.Errorf("创建 Token 对象出错，错误信息：%v", err)
		}

		parsedToken, _, err := token.Parse(signedToken)

		if err == nil {
			t.Error("验证 Token 字符串过期失败，期待返回错误")
		}

		if err != ErrorExpired {
			t.Error("验证 Token 字符串过期失败，期待返回过期错误")
		}

		if parsedToken == nil {
			t.Error("验证 Token 字符串过期失败，期待返回解析后的 Token 对象")
		} else {
			if parsedToken.Valid {
				t.Error("验证 Token 字符串过期失败，期待验证不通过")
			}
		}
	})
}

func TestTokenNotValidYet(t *testing.T) {
	token, err := NewToken(
		WithKey(TestKey),
		WithAlg(DefaultAlg),
		WithBody(map[string]interface{}{
			"name": "Tom",
			"sex":  "male",
		}),
		WithClaims(jwt.StandardClaims{
			Audience:  "server",
			ExpiresAt: time.Now().Add(time.Second * time.Duration(math.MaxUint32)).Unix(),
			Id:        "abcde",
			IssuedAt:  time.Now().Unix(),
			Issuer:    "app",
			NotBefore: time.Now().Add(time.Second * 3600).Unix(),
			Subject:   "authentication token",
		}),
	)
	if err != nil {
		t.Errorf("创建 Token 对象出错，错误信息：%v", err)
	}

	signedToken, err := token.Create()
	if err != nil {
		t.Errorf("创建 JWT 字符串出错，错误信息：%v", err)
	}

	t.Run("tokenNotValidYet", func(t *testing.T) {
		token, err := NewToken(
			WithKey(TestKey),
		)
		if err != nil {
			t.Errorf("创建 Token 对象出错，错误信息：%v", err)
		}

		parsedToken, _, err := token.Parse(signedToken)

		if err == nil {
			t.Error("验证 Token 字符串暂不可用失败，期待返回错误")
		}

		if err != ErrorNotValidYet {
			t.Error("验证 Token 字符串暂不可用失败，期待返回暂不可用错误")
		}

		if parsedToken == nil {
			t.Error("验证 Token 字符串暂不可用失败，期待返回解析后的 Token 对象")
		} else {
			if parsedToken.Valid {
				t.Error("验证 Token 字符串暂不可用失败，期待验证不通过")
			}
		}
	})
}
