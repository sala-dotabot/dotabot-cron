package metrics

import (
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type IamContext interface {
	IamToken() (signed string, err error)
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
func (this *IamContextImpl) signedToken() (signed string, err error) {
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

func (this *IamContextImpl) IamToken() (iam string, err error) {
	jwt, err := this.signedToken()
	if err != nil {
		return
	}
	resp, err := http.Post(
		"https://iam.api.cloud.yandex.net/iam/v1/tokens",
		"application/json",
		strings.NewReader(fmt.Sprintf(`{"jwt":"%s"}`, jwt)),
	)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		err = errors.New(fmt.Sprintf("%s: %s", resp.Status, body))
	}
	var data struct {
		IAMToken string `json:"iamToken"`
	}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return
	}
	iam = data.IAMToken
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
