package metrics

import (
	"crypto/rsa"
	"io/ioutil"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type IamContext interface {
	SignedToken() (signed string, err error)
}

type IamContextImpl struct {
	keyId            string
	serviceAccountID string
	keyFileName      string
}

func CreateIamContext(keyId string, serviceAccountId string, keyFileName string) IamContext {
	return &IamContextImpl{
		keyId:            keyId,
		serviceAccountID: serviceAccountId,
		keyFileName:      keyFileName,
	}
}

// Формирование JWT.
func (this *IamContextImpl) SignedToken() (signed string, err error) {
	claims := jwt.RegisteredClaims{
		Issuer:    this.serviceAccountID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
		Audience:  []string{"https://iam.api.cloud.yandex.net/iam/v1/tokens"},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodPS256, claims)
	token.Header["kid"] = this.keyId

	privateKey, err := this.loadPrivateKey()
	if err != nil {
		return
	}
	signed, err = token.SignedString(privateKey)
	return
}

func (this *IamContextImpl) loadPrivateKey() (rsaPrivateKey *rsa.PrivateKey, err error) {
	data, err := ioutil.ReadFile(this.keyFileName)
	if err != nil {
		return
	}
	rsaPrivateKey, err = jwt.ParseRSAPrivateKeyFromPEM(data)
	return
}
