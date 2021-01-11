package gcpss

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/BESTSELLER/go-vault/models"
	"github.com/hashicorp/vault/api"
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

	var s models.Auth

	err = json.NewDecoder(resp.Body).Decode(&s)
	if err != nil {
		return "", err
	}

	if len(s.Errors) > 0 {
		return "", fmt.Errorf(s.Errors[0])
	}

	return s.ClientToken, nil
}

// FetchVaultSecret returns secret from Hashicorp Vault.
func FetchVaultSecret(vaultAddr string, vaultSecret string, vaultRole string) (secret string, err error) {

	jwt, err := fetchJWT(vaultRole)
	if err != nil {
		return "", err
	}
	fmt.Println("jwt:", jwt)

	token, err := fetchVaultToken(vaultAddr, jwt, vaultRole)
	if err != nil {
		return "", err
	}
	fmt.Println("token:", token)

	client, err := api.NewClient(&api.Config{
		Address: vaultAddr,
	})
	if err != nil {
		panic(err)
	}
	client.SetToken(token)

	sec, err := client.Logical().Read(vaultSecret)
	if err != nil {
		return "", err
	}

	data, err := json.Marshal(sec.Data["data"])
	if err != nil {
		return "", err
	}

	return string(data), nil

}
