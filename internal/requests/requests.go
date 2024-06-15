package requests

import (
	"bytes"
	"dcsa-lab/internal/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	jsoniter "github.com/json-iterator/go"
)

var (
	BASE_URL   = "http://localhost:8080/"
	AUTH_TOKEN = ""
)

func SignUp(data models.UserData) (int, string, error) {
	client := &http.Client{}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return -1, "", err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%sauth/signup", BASE_URL), bytes.NewReader(jsonData))
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		return -1, "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return -1, "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return -1, "", err
	}

	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, body, "", "\t")
	if err != nil {
		return -1, "", err
	}

	return resp.StatusCode, string(prettyJSON.Bytes()), nil
}

func Login(data models.LoginData) (int, string, error) {
	client := &http.Client{}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return -1, "", err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%sauth/login", BASE_URL), bytes.NewReader(jsonData))
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		return -1, "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return -1, "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return -1, "", err
	}

	var response struct {
		Token string `json:"token"`
	}

	err = jsoniter.Unmarshal(body, &response)
	if err != nil {
		return -1, "", err
	}

	AUTH_TOKEN = response.Token

	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, body, "", "\t")
	if err != nil {
		return -1, "", err
	}

	return resp.StatusCode, string(prettyJSON.Bytes()), nil
}

func GetAllUsers() (int, string, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", fmt.Sprintf("%susers", BASE_URL), bytes.NewReader([]byte{}))
	req.Header.Add("Authorization", "Bearer "+AUTH_TOKEN)
	if err != nil {
		return -1, "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return -1, "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return -1, "", err
	}

	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, body, "", "\t")
	if err != nil {
		return -1, "", err
	}

	return resp.StatusCode, string(prettyJSON.Bytes()), nil
}

func GetUserById(id int) (int, string, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", fmt.Sprintf("%susers/%d", BASE_URL, id), bytes.NewReader([]byte{}))
	req.Header.Add("Authorization", "Bearer "+AUTH_TOKEN)
	if err != nil {
		return -1, "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return -1, "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return -1, "", err
	}

	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, body, "", "\t")
	if err != nil {
		return -1, "", err
	}

	return resp.StatusCode, string(prettyJSON.Bytes()), nil
}

func DeleteUser(id int) (int, string, error) {
	client := &http.Client{}

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%susers/%d", BASE_URL, id), bytes.NewReader([]byte{}))
	req.Header.Add("Authorization", "Bearer "+AUTH_TOKEN)
	if err != nil {
		return -1, "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return -1, "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return -1, "", err
	}

	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, body, "", "\t")
	if err != nil {
		return -1, "", err
	}

	return resp.StatusCode, string(prettyJSON.Bytes()), nil
}

func UpdateUser(id int, data models.UserData) (int, string, error) {
	client := &http.Client{}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return -1, "", err
	}

	req, err := http.NewRequest("PATCH", fmt.Sprintf("%susers/%d", BASE_URL, id), bytes.NewReader(jsonData))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+AUTH_TOKEN)
	if err != nil {
		return -1, "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return -1, "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return -1, "", err
	}

	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, body, "", "\t")
	if err != nil {
		return -1, "", err
	}

	return resp.StatusCode, string(prettyJSON.Bytes()), nil
}
