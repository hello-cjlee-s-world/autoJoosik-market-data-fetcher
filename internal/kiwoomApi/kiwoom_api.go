package kiwoomApi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type KiwoomConfig struct {
	AppKey    string `mapstructure:"appKey"`
	SecretKey string `mapstructure:"secretKey"`
	Token     string `mapstructure:"token"`
}

var kwConfig KiwoomConfig

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

func KiwoomInit(kiwoomConfig KiwoomConfig) {
	kwConfig = kiwoomConfig
	token, err := GetToken()
	if err != nil {
		fmt.Println("토큰 요청 실패:", err)
		return
	} else {
		fmt.Println("Token Value : " + token)
		tokenObj := map[string]string{}
		err = json.Unmarshal([]byte(token), &tokenObj)
		if err != nil {
			kwConfig.Token = tokenObj["token"]
			rst, err := GetStockInfo("005930")
			if err != nil {
				fmt.Println("결과값 test 실패 ::", rst)
				return
			} else {
				fmt.Println("결과값 test :: " + rst)
			}
		}
	}
}

func GetToken() (string, error) {
	url := "https://api.kiwoom.com/oauth2/token"
	payload := map[string]string{
		"grant_type": "client_credentials",
		"appkey":     kwConfig.AppKey,
		"secretkey":  kwConfig.SecretKey,
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

func GetStockInfo(stkCd string) (string, error) {
	url := "https://api.kiwoom.com/api/dostk/stkinfo"
	payload := map[string]string{
		"stk_cd": stkCd,
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
	req.Header.Set("api-id", "ka10001")
	req.Header.Set("authorization", fmt.Sprintf("Bearer %s", kwConfig.Token))

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

func GetTradeInfoLog(stkCd string) (string, error) {
	url := "https://api.kiwoom.com/api/dostk/stkinfo"
	payload := map[string]string{
		"stk_cd": stkCd,
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
	req.Header.Set("api-id", "ka10003")
	req.Header.Set("authorization", fmt.Sprintf("Bearer %s", kwConfig.Token))

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

func GetOrderBookLog() (string, error) {
	url := "https://api.kiwoom.com/api/dostk/mrkcond"
	payload := map[string]string{
		"stk_cd": "005930",
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
	req.Header.Set("api-id", "ka10004")
	req.Header.Set("authorization", fmt.Sprintf("Bearer %s", kwConfig.Token))

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

func GetStockDailyLog() (string, error) {
	url := "https://api.kiwoom.com/api/dostk/mrkcond"
	payload := map[string]string{
		"stk_cd": "005930",
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
	req.Header.Set("api-id", "ka10005")
	req.Header.Set("authorization", fmt.Sprintf("Bearer %s", kwConfig.Token))

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

func GetAccountProfitLog() (string, error) {
	url := "https://api.kiwoom.com/api/dostk/acnt"
	payload := map[string]string{
		"stex_tp": "1",
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
	req.Header.Set("api-id", "ka10085")
	req.Header.Set("authorization", fmt.Sprintf("Bearer %s", kwConfig.Token))

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

func GetAtnStlInfr() (string, error) {
	url := "https://api.kiwoom.com/api/dostk/stkinfo"
	payload := map[string]string{
		"stk_cd": "039490", //KRX
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
	req.Header.Set("api-id", "ka10095")
	req.Header.Set("authorization", fmt.Sprintf("Bearer %s", kwConfig.Token))

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
