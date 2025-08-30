package kiwoomApi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type KiwoomApiConfig struct {
	AppKey    string `mapstructure:"appKey"`
	SecretKey string `mapstructure:"secretKey"`
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

func (fc *KiwoomApiConfig) Initialize() {
	token, err := fc.GetToken()
	if err != nil {
		fmt.Println("토큰 요청 실패:", err)
		return
	} else {
		fmt.Println("Token Value : " + token)
	}
}

func (fc *KiwoomApiConfig) GetToken() (string, error) {
	url := "https://api.kiwoom.com/oauth2/token"
	payload := map[string]string{
		"grant_type": "client_credentials",
		"appkey":     fc.AppKey,
		"secretkey":  fc.SecretKey,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	//var token TokenResponse
	//if err := json.Unmarshal(body, &token); err != nil {
	//	return nil, err
	//}

	return string(body), nil
}
