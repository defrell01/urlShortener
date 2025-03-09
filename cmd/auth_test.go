package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"urlshortener/internal/auth"
	"urlshortener/internal/user"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func initDb() *gorm.DB {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	db, err := gorm.Open(postgres.Open(os.Getenv("DSN")), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	return db
}

func initData(db *gorm.DB) {
	db.Create(&user.User{
		Email:    "a2@2a.ru",
		Password: "$2a$10$vyF4VBhpqBV5cHST19eb/eV31xVLh9by2xfeWS6GEmnZmWxG6UyQG",
		Name:     "Vasya",
	})
}

func removeData(db *gorm.DB) {
	db.Unscoped().
		Where("email = ?", "a2@2a.ru").
		Delete(&user.User{})
}

func TestLoginSuccess(t *testing.T) {
	/// prepare
	db := initDb()
	initData(db)

	ts := httptest.NewServer(App())
	defer ts.Close()

	data, _ := json.Marshal(&auth.LoginRequest{
		Email:    "a2@2a.ru",
		Password: "testing123",
	})

	response, err := http.Post(ts.URL+"/auth/login", "application/json", bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}

	if response.StatusCode != http.StatusOK {
		t.Fatalf("Expected %d got %d", http.StatusOK, response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}

	var resData auth.LoginResponse
	err = json.Unmarshal(body, &resData)

	if err != nil {
		t.Fatal(err)
	}

	if resData.Token == "" {
		t.Fatal("Token empty")
	}

	removeData(db)
}

func TestLoginFailed(t *testing.T) {
	db := initDb()
	initData(db)

	ts := httptest.NewServer(App())
	defer ts.Close()

	data, _ := json.Marshal(&auth.LoginRequest{
		Email:    "a2@2a.ru",
		Password: "testing1234",
	})

	response, err := http.Post(ts.URL+"/auth/login", "application/json", bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}

	if response.StatusCode != http.StatusUnauthorized {
		t.Fatalf("Expected %d got %d", http.StatusUnauthorized, response.StatusCode)
	}

	removeData(db)
}
