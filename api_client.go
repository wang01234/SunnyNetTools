package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// APIConfig holds the API server configuration
type APIConfig struct {
	URL string
}

// GetAPIConfig returns the API configuration from GlobalConfig
func GetAPIConfig() APIConfig {
	return APIConfig{
		URL: GlobalConfig.APIURL,
	}
}

// LoginRequest represents a login API request
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse represents a login API response
type LoginResponse struct {
	Success  bool   `json:"success"`
	Message  string `json:"message,omitempty"`
	Username string `json:"username,omitempty"`
	Role     string `json:"role,omitempty"`
}

// User represents a user with role
type User struct {
	Username string `json:"username"`
	Role     string `json:"role"`
}

// UsersResponse represents a users list API response
type UsersResponse struct {
	Success bool   `json:"success"`
	Users   []User `json:"users,omitempty"`
	Message string `json:"message,omitempty"`
}

// ActionResponse represents a generic action API response
type ActionResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

// CreateUserRequest represents a create user API request
type CreateUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

// APICall performs an HTTP request to the API server
func APICall(method, endpoint string, body interface{}) ([]byte, error) {
	config := GetAPIConfig()

	var reqBody []byte
	if body != nil {
		reqBody, _ = json.Marshal(body)
	}

	url := config.URL + endpoint
	req, err := http.NewRequest(method, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	result := make([]byte, 1024)
	n, _ := resp.Body.Read(result)
	return result[:n], nil
}

// Login calls the login API
func Login(username, password string) (bool, string, string, error) {
	config := GetAPIConfig()

	req := LoginRequest{
		Username: username,
		Password: password,
	}

	jsonData, _ := json.Marshal(req)
	resp, err := http.Post(config.URL+"/api/login", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return false, "", "", err
	}
	defer resp.Body.Close()

	var result LoginResponse
	json.NewDecoder(resp.Body).Decode(&result)

	if result.Success {
		return true, result.Username, result.Role, nil
	}
	return false, result.Message, "", nil
}

// GetUsers calls the get users API
func GetUsers(currentUser, currentRole string) ([]User, error) {
	config := GetAPIConfig()

	req, err := http.NewRequest("GET", config.URL+"/api/users", nil)
	if err != nil {
		return nil, err
	}
	if currentUser != "" {
		req.Header.Set("X-Current-User", currentUser)
	}
	if currentRole != "" {
		req.Header.Set("X-Current-Role", currentRole)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result UsersResponse
	json.NewDecoder(resp.Body).Decode(&result)

	return result.Users, nil
}

// CreateUser calls the create user API
func CreateUser(username, password, role string, currentUser string) (bool, string, error) {
	config := GetAPIConfig()

	req := CreateUserRequest{
		Username: username,
		Password: password,
		Role:     role,
	}

	jsonData, _ := json.Marshal(req)

	httpReq, err := http.NewRequest("POST", config.URL+"/api/users", bytes.NewBuffer(jsonData))
	if err != nil {
		return false, "", err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	if currentUser != "" {
		httpReq.Header.Set("X-Current-User", currentUser)
	}

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return false, "", err
	}
	defer resp.Body.Close()

	var result ActionResponse
	json.NewDecoder(resp.Body).Decode(&result)

	return result.Success, result.Message, nil
}

// UpdatePassword calls the update password API
func UpdatePassword(username, password string) (bool, string, error) {
	config := GetAPIConfig()

	req := map[string]string{"password": password}
	jsonData, _ := json.Marshal(req)

	url := fmt.Sprintf("%s/api/users/%s", config.URL, username)
	httpClient := &http.Client{}
	reqHTTP, _ := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	reqHTTP.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(reqHTTP)
	if err != nil {
		return false, "", err
	}
	defer resp.Body.Close()

	var result ActionResponse
	json.NewDecoder(resp.Body).Decode(&result)

	return result.Success, result.Message, nil
}

func DeleteUser(username, currentUser string) (bool, string, error) {
	config := GetAPIConfig()

	url := fmt.Sprintf("%s/api/users/%s", config.URL, username)
	httpClient := &http.Client{}
	reqHTTP, _ := http.NewRequest("DELETE", url, nil)
	if currentUser != "" {
		reqHTTP.Header.Set("X-Current-User", currentUser)
	}

	resp, err := httpClient.Do(reqHTTP)
func DeleteUser(username string) (bool, string, error) {
	config := GetAPIConfig()

	url := fmt.Sprintf("%s/api/users/%s", config.URL, username)
	httpClient := &http.Client{}
	reqHTTP, _ := http.NewRequest("DELETE", url, nil)

	resp, err := httpClient.Do(reqHTTP)
	if err != nil {
		return false, "", err
	}
	defer resp.Body.Close()

	var result ActionResponse
	json.NewDecoder(resp.Body).Decode(&result)

	return result.Success, result.Message, nil
}
