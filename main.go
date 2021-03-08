package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type License struct {
	Type string `json:"type"`
	Id   string `json:"id"`
}

type ValidationParams struct {
	Key string `json:"key"`
}

type ValidationResult struct {
	Valid bool   `json:"valid"`
	Code  string `json:"constant"`
}

type ValidationRequest struct {
	Meta ValidationParams `json:"meta"`
}

type ValidationResponse struct {
	Result  ValidationResult `json:"meta"`
	License *License         `json:"data,omitempty"`
}

func validateLicenseKey(key string) (*ValidationResponse, error) {
	req, err := json.Marshal(ValidationRequest{ValidationParams{key}})
	if err != nil {
		return nil, err
	}

	body := bytes.NewBuffer(req)
	res, err := http.Post(
		fmt.Sprintf("https://api.keygen.sh/v1/accounts/%s/licenses/actions/validate-key", os.Getenv("KEYGEN_ACCOUNT_ID")),
		"application/vnd.api+json",
		body,
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(res.Body)

		return nil, fmt.Errorf("An API error occurred: %s", body)
	}

	var v *ValidationResponse
	json.NewDecoder(res.Body).Decode(&v)

	return v, nil
}

func promptForLicenseKey() string {
	fmt.Print("Enter license key: ")

	stdin := bufio.NewReader(os.Stdin)
	input, _, _ := stdin.ReadLine()

	return fmt.Sprintf("%s", input)
}

func main() {
	licenseKey := promptForLicenseKey()
	validation, err := validateLicenseKey(licenseKey)
	if err != nil {
		fmt.Println(err)

		os.Exit(1)
	}

	if validation.Result.Valid {
		fmt.Printf("License key is valid: code=%s id=%s\n", validation.Result.Code, validation.License.Id)
	} else {
		fmt.Printf("License key is invalid: code=%s\n", validation.Result.Code)
	}
}
