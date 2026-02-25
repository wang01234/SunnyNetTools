package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/crypto/bcrypt"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

type Config struct {
	DBHost     string `json:"db_host"`
	DBPort     string `json:"db_port"`
	DBUser     string `json:"db_user"`
	DBPassword string `json:"db_password"`
	DBDatabase string `json:"db_database"`
	APIPort    string `json:"api_port"`
}

var config Config

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Success  bool   `json:"success"`
	Message  string `json:"message,omitempty"`
	Username string `json:"username,omitempty"`
	Role     string `json:"role,omitempty"`
}

type User struct {
	Username string `json:"username"`
	Role     string `json:"role"`
}

type UsersResponse struct {
	Success bool   `json:"success"`
	Users   []User `json:"users,omitempty"`
	Message string `json:"message,omitempty"`
}

type ActionResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

type CreateUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

func loadConfig() error {
	data, err := os.ReadFile("config.json")
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &config)
}

func initDB() error {
	dsn := config.DBUser + ":" + config.DBPassword + "@tcp(" + config.DBHost + ":" + config.DBPort + ")/" + config.DBDatabase + "?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	return db.Ping()
}

func hashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

func checkPassword(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		resp := LoginResponse{Success: false, Message: "请求格式错误"}
		json.NewEncoder(w).Encode(resp)
		return
	}

	var storedPassword, storedRole string
	err := db.QueryRow("SELECT password, COALESCE(role, 'user') FROM users WHERE username = ?", req.Username).Scan(&storedPassword, &storedRole)
	if err != nil {
		if err == sql.ErrNoRows {
			resp := LoginResponse{Success: false, Message: "用户不存在"}
			json.NewEncoder(w).Encode(resp)
			return
		}
		resp := LoginResponse{Success: false, Message: "查询失败"}
		json.NewEncoder(w).Encode(resp)
		return
	}

	if checkPassword(req.Password, storedPassword) {
		resp := LoginResponse{Success: true, Username: req.Username, Role: storedRole}
		json.NewEncoder(w).Encode(resp)
		return
	}

	resp := LoginResponse{Success: false, Message: "密码错误"}
	json.NewEncoder(w).Encode(resp)
}

func getUsersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	currentUser := r.Header.Get("X-Current-User")

	// Get current user's role
	var myRole string
	err := db.QueryRow("SELECT COALESCE(role, 'user') FROM users WHERE username = ?", currentUser).Scan(&myRole)
	if err != nil {
		myRole = "user" // Default to regular user if not found
	}

	rows, err := db.Query("SELECT username, COALESCE(role, 'user') as role FROM users")
	if err != nil {
		resp := UsersResponse{Success: false, Message: "查询失败"}
		json.NewEncoder(w).Encode(resp)
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var username, role string
		if err := rows.Scan(&username, &role); err != nil {
			continue
		}
		// Admin sees all users, regular user sees only themselves
		if myRole != "admin" && username != currentUser {
			continue
		}
		users = append(users, User{Username: username, Role: role})
	}

	resp := UsersResponse{Success: true, Users: users}
	json.NewEncoder(w).Encode(resp)
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	currentUser := r.Header.Get("X-Current-User")

	// Check if current user is admin
	var myRole string
	err := db.QueryRow("SELECT COALESCE(role, 'user') FROM users WHERE username = ?", currentUser).Scan(&myRole)
	if err != nil {
		myRole = "user"
	}

	// Only admin can create users
	if myRole != "admin" {
		resp := ActionResponse{Success: false, Message: "无权限添加账号"}
		json.NewEncoder(w).Encode(resp)
		return
	}

	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		resp := ActionResponse{Success: false, Message: "请求格式错误"}
		json.NewEncoder(w).Encode(resp)
		return
	}

	if req.Username == "" || req.Password == "" {
		resp := ActionResponse{Success: false, Message: "用户名或密码不能为空"}
		json.NewEncoder(w).Encode(resp)
		return
	}

	// Default role is user if not specified
	if req.Role == "" {
		req.Role = "user"
	}

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", req.Username).Scan(&count)
	if err != nil {
		resp := ActionResponse{Success: false, Message: "查询失败"}
		json.NewEncoder(w).Encode(resp)
		return
	}
	if count > 0 {
		resp := ActionResponse{Success: false, Message: "账号已存在"}
		json.NewEncoder(w).Encode(resp)
		return
	}

	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		resp := ActionResponse{Success: false, Message: "密码加密失败"}
		json.NewEncoder(w).Encode(resp)
		return
	}

	_, err = db.Exec("INSERT INTO users (username, password, role) VALUES (?, ?, ?)", req.Username, hashedPassword, req.Role)
	if err != nil {
		resp := ActionResponse{Success: false, Message: "添加账号失败"}
		json.NewEncoder(w).Encode(resp)
		return
	}

	resp := ActionResponse{Success: true}
	json.NewEncoder(w).Encode(resp)
}

func updatePasswordHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/api/users/")
	username := strings.Split(path, "/")[0]

	var req struct {
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		resp := ActionResponse{Success: false, Message: "请求格式错误"}
		json.NewEncoder(w).Encode(resp)
		return
	}

	if username == "" || req.Password == "" {
		resp := ActionResponse{Success: false, Message: "用户名或密码不能为空"}
		json.NewEncoder(w).Encode(resp)
		return
	}

	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", username).Scan(&count)
	if err != nil || count == 0 {
		resp := ActionResponse{Success: false, Message: "账号不存在"}
		json.NewEncoder(w).Encode(resp)
		return
	}

	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		resp := ActionResponse{Success: false, Message: "密码加密失败"}
		json.NewEncoder(w).Encode(resp)
		return
	}

	_, err = db.Exec("UPDATE users SET password = ? WHERE username = ?", hashedPassword, username)
	if err != nil {
		resp := ActionResponse{Success: false, Message: "修改密码失败"}
		json.NewEncoder(w).Encode(resp)
		return
	}

	resp := ActionResponse{Success: true}
	json.NewEncoder(w).Encode(resp)
}

func deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/api/users/")
	username := strings.Split(path, "/")[0]

	currentUser := r.Header.Get("X-Current-User")

	// Prevent user from deleting themselves
	if currentUser != "" && username == currentUser {
		resp := ActionResponse{Success: false, Message: "不能删除自己的账号"}
		json.NewEncoder(w).Encode(resp)
		return
	}

	if username == "admin" {
		resp := ActionResponse{Success: false, Message: "默认账号不能删除"}
		json.NewEncoder(w).Encode(resp)
		return
	}


	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", username).Scan(&count)
	if err != nil || count == 0 {
		resp := ActionResponse{Success: false, Message: "账号不存在"}
		json.NewEncoder(w).Encode(resp)
		return
	}

	_, err = db.Exec("DELETE FROM users WHERE username = ?", username)
	if err != nil {
		resp := ActionResponse{Success: false, Message: "删除账号失败"}
		json.NewEncoder(w).Encode(resp)
		return
	}

	resp := ActionResponse{Success: true}
	json.NewEncoder(w).Encode(resp)
}

func main() {
	if err := loadConfig(); err != nil {
		log.Println("未找到配置文件，使用默认配置")
		config = Config{
			DBHost:     "47.108.25.208",
			DBPort:     "3306",
			DBUser:     "sunny_db",
			DBPassword: "LnczeYcjKssZesdm",
			DBDatabase: "sunny_db",
			APIPort:    "8080",
		}
		data, _ := json.MarshalIndent(config, "", "  ")
		os.WriteFile("config.json", data, 0644)
		log.Println("已创建默认配置文件 config.json，请根据需要修改")
	}

	if err := initDB(); err != nil {
		log.Fatal("数据库连接失败:", err)
	}

	log.Println("数据库连接成功")

	http.HandleFunc("/api/login", loginHandler)
	http.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getUsersHandler(w, r)
		case http.MethodPost:
			createUserHandler(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/api/users/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut {
			updatePasswordHandler(w, r)
		} else if r.Method == http.MethodDelete {
			deleteUserHandler(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	addr := ":" + config.APIPort
	log.Printf("API服务启动于 %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
