package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	ntlmssp "github.com/Azure/go-ntlmssp"
	"github.com/koltyakov/gosip"
	"github.com/koltyakov/gosip/auth/basic"
)

func main() {
	configPath := "./config/private.basic.json"
	auth := &basic.AuthCnfg{}

	err := auth.ReadConfig(configPath)
	if err != nil {
		fmt.Printf("Unable to get config: %v\n", err)
		return
	}

	client := &gosip.SPClient{
		AuthCnfg: auth,
	}

	// NTML + Negotiation
	client.Transport = ntlmssp.Negotiator{
		RoundTripper: &http.Transport{},
	}

	apiEndpoint := auth.GetSiteURL() + "/_api/web?$select=Title"
	req, err := http.NewRequest("GET", apiEndpoint, nil)
	if err != nil {
		fmt.Printf("Unable to create a request: %v", err)
		return
	}

	req.Header.Set("Accept", "application/json;odata=verbose")

	fmt.Printf("Requesting api endpoint: %s\n", apiEndpoint)
	resp, err := client.Execute(req)
	if err != nil {
		fmt.Printf("Unable to request api: %v\n", err)
		return
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Unable to read a response: %v\n", err)
		return
	}

	fmt.Printf("Raw data: %s", string(data))
}
