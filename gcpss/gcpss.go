package gcpss

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/BESTSELLER/go-vault/models"
)

func fetchJWT(vaultRole string) (jwt string, err error) {
	client := new(http.Client)

	url := "http://metadata/computeMetadata/v1/instance/service-accounts/default/identity?audience=http://vault/" + vaultRole + "&format=full"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Metadata-Flavor", "Google")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func fetchVaultToken(vaultAddr string, jwt string, vaultRole string) (vaultToken string, err error) {
	client := new(http.Client)

	j := `{"role":"` + vaultRole + `", "jwt":"` + jwt + `"}`

	req, err := http.NewRequest(http.MethodPost, vaultAddr+"/v1/auth/gcp/login", bytes.NewBufferString(j))
	if err != nil {
		return "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var s models.Login

	err = json.NewDecoder(resp.Body).Decode(&s)
	if err != nil {
		return "", err
	}

	if len(s.Auth.Errors) > 0 {
		return "", fmt.Errorf(s.Auth.Errors[0])
	}

	return s.Auth.ClientToken, nil
}

func readSecret(vaultAddr string, vaultToken string, vaultSecret string) (secret string, err error) {
	client := new(http.Client)
	req, err := http.NewRequest(http.MethodGet, vaultAddr+"/v1/"+vaultSecret, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("X-Vault-Token", vaultToken)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	s := struct {
		Data struct {
			Data string `json:"data"`
		} `json:"data"`
	}{}

	if err := json.NewDecoder(resp.Body).Decode(&s); err != nil {
		return "", err
	}

	return s.Data.Data, nil

}

// FetchVaultSecret returns secret from Hashicorp Vault.
func FetchVaultSecret(vaultAddr string, vaultSecret string, vaultRole string) (secret string, err error) {

	jwt, err := fetchJWT(vaultRole)
	if err != nil {
		return "", err
	}

	token, err := fetchVaultToken(vaultAddr, jwt, vaultRole)
	if err != nil {
		return "", err
	}

	data, err := readSecret(vaultAddr, token, vaultSecret)
	if err != nil {
		return "", err
	}

	return string(data), nil

}
