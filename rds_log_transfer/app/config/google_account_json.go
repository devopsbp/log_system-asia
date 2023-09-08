package config

import (
	"encoding/json"
	"os"
)

// GoogleアカウントJSONをパースする

// GoogleAccountJson json構造体
type GoogleAccountJson struct {
	Type                    string `json:"type"`
	ProjectId               string `json:"project_id"`
	PrivateKeyId            string `json:"private_key_id"`
	PrivateKey              string `json:"private_key"`
	ClientEmail             string `json:"client_email"`
	ClientId                string `json:"client_id"`
	AuthUri                 string `json:"auth_uri"`
	TokenUri                string `json:"token_uri"`
	AuthProviderX509CertUrl string `json:"auth_provider_x509_cert_url"`
	ClientX509CertUrl       string `json:"client_x509_cert_url"`
}

const GoogleAccountJsonPath = ".secret/sb_bq_service_account.json"

var googleAccountJsonValues *GoogleAccountJson = nil

func GetParsedGoogleAccountJson() (*GoogleAccountJson, error) {
	if googleAccountJsonValues != nil {
		return googleAccountJsonValues, nil
	}

	// jsonを読み込みポインタセット
	jsonBytes, err := os.ReadFile(GoogleAccountJsonPath)
	if err != nil {
		return googleAccountJsonValues, err
	}

	var jsonValues GoogleAccountJson
	err = json.Unmarshal(jsonBytes, &jsonValues)
	if err != nil {
		return googleAccountJsonValues, err
	}
	googleAccountJsonValues = &jsonValues

	return googleAccountJsonValues, nil
}
