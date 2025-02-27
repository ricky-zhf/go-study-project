package felo

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
)

// Apple JWK URL
const appleJWKURL = "https://appleid.apple.com/auth/keys"

// Apple 服务器通知结构体
type AppStoreNotification struct {
	NotificationType string `json:"notificationType"`
	Subtype          string `json:"subtype,omitempty"`
	Version          string `json:"version"`
	SignedPayload    string `json:"signedPayload"`
}

// Apple 公钥结构体
type AppleJWK struct {
	Keys []AppleJWKKey `json:"keys"`
}

type AppleJWKKey struct {
	Kid string   `json:"kid"`
	Alg string   `json:"alg"`
	Kty string   `json:"kty"`
	Use string   `json:"use"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	X5c []string `json:"x5c"`
}

// 获取苹果公钥
func getApplePublicKey1(kid string) (*rsa.PublicKey, error) {
	resp, err := http.Get(appleJWKURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch JWKS: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read JWKS response: %v", err)
	}

	var jwks struct {
		Keys []struct {
			Kty string `json:"kty"`
			Kid string `json:"kid"`
			Use string `json:"use"`
			Alg string `json:"alg"`
			N   string `json:"n"`
			E   string `json:"e"`
		} `json:"keys"`
	}

	if err := json.Unmarshal(body, &jwks); err != nil {
		return nil, fmt.Errorf("failed to parse JWKS: %v", err)
	}

	for _, key := range jwks.Keys {
		if key.Kid == kid && key.Use == "sig" && key.Alg == "RS256" {
			nBytes, err := base64.RawURLEncoding.DecodeString(key.N)
			if err != nil {
				return nil, fmt.Errorf("failed to decode modulus: %v", err)
			}

			eBytes, err := base64.RawURLEncoding.DecodeString(key.E)
			if err != nil {
				return nil, fmt.Errorf("failed to decode exponent: %v", err)
			}

			if len(eBytes) < 4 {
				eBytes = append(make([]byte, 4-len(eBytes)), eBytes...)
			}
			e := int(uint32(eBytes[0])<<24 | uint32(eBytes[1])<<16 | uint32(eBytes[2])<<8 | uint32(eBytes[3]))

			return &rsa.PublicKey{
				N: new(big.Int).SetBytes(nBytes),
				E: e,
			}, nil
		}
	}

	return nil, fmt.Errorf("public key not found for kid: %s", kid)
}

// 处理 Apple Webhook
func appleWebhookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// 解析 Apple 通知
	var notification AppStoreNotification
	err := json.NewDecoder(r.Body).Decode(&notification)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// 解析 `SignedPayload`
	signedClaims, err := parseJWT(notification.SignedPayload)
	if err != nil {
		log.Printf("Error parsing SignedPayload: %v", err)
		http.Error(w, "Invalid JWT", http.StatusBadRequest)
		return
	}

	// 解析 `JWSTransaction`
	dataJSON, _ := json.Marshal(signedClaims["data"])
	var appStoreData AppStoreData
	json.Unmarshal(dataJSON, &appStoreData)

	transactionClaims, err := parseJWT(appStoreData.JWSTransaction)
	if err != nil {
		log.Printf("Error parsing JWSTransaction: %v", err)
		http.Error(w, "Invalid Transaction JWT", http.StatusBadRequest)
		return
	}

	// 解析 `JWSTransactionDecodedPayload`
	transactionJSON, _ := json.Marshal(transactionClaims)
	var transaction JWSTransactionDecodedPayload
	json.Unmarshal(transactionJSON, &transaction)

	log.Printf("Received Transaction: %+v", transaction)

	w.WriteHeader(http.StatusOK)
}

func main() {
	http.HandleFunc("/apple-webhook", appleWebhookHandler)
	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
