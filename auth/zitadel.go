package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateClientAssertion(clientId, keyId, privateKeyPEM string, audience string) (string, error) {
	// 1. 解析 PEM 格式的私钥
	block, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKeyPEM))
	if err != nil {
		return "", fmt.Errorf("解析私钥失败: %w", err)
	}

	// 2. 构建 Claims (断言负载)
	// 根据 OAuth2 规范，client_assertion 必须包含以下字段
	now := time.Now()
	claims := jwt.MapClaims{
		"iss": clientId,                          // 签发者是你的 ClientID
		"sub": clientId,                          // 面向的对象也是你的 ClientID
		"aud": audience,                          // 接收者是 ZITADEL 的域名 (e.g., https://auth.pdnode.com)
		"iat": now.Unix(),                        // 签发时间
		"exp": now.Add(time.Minute).Unix(),       // 过期时间 (建议设短一点，1分钟足够)
		"jti": fmt.Sprintf("%d", now.UnixNano()), // 唯一标识符
	}

	// 3. 创建 Token 并设置 Header 中的 kid (ZITADEL 靠它找对应的公钥)
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = keyId

	// 4. 使用 RSA 私钥签名
	signedToken, err := token.SignedString(block)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

type ZitadelConfig struct {
	Type     string `json:"type"`
	KeyID    string `json:"keyId"`
	Key      string `json:"key"` // 这里的 key 就是 RSA 私钥字符串
	AppID    string `json:"appId"`
	ClientID string `json:"clientId"`
}
