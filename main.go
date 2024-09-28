package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/spf13/viper"
)

// User 구조체 정의
type User struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

var db *pg.DB

// Viper를 사용하여 설정 로드 (환경 변수 사용)
func initConfig() {
	viper.AutomaticEnv() // 환경 변수 사용을 자동으로 설정
}

// 데이터베이스 초기화
func initDB() {
	db = pg.Connect(&pg.Options{
		User:     viper.GetString("SENDMIND_POSTGRES_USER"),
		Password: viper.GetString("SENDMIND_POSTGRES_PASSWORD"),
		Database: viper.GetString("SENDMIND_POSTGRES_DB"),
	})

	err := createSchema()
	if err != nil {
		log.Fatalf("Error creating schema: %v\n", err)
	}
}

// 테이블 스키마 생성
func createSchema() error {
	return db.Model((*User)(nil)).CreateTable(&orm.CreateTableOptions{
		IfNotExists: true,
	})
}

// 사용자 생성 API
func createUser(w http.ResponseWriter, r *http.Request) {
	var user User
	json.NewDecoder(r.Body).Decode(&user)

	_, err := db.Model(&user).Insert()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(user)
}

// 모든 사용자 조회 API
func getUsers(w http.ResponseWriter, r *http.Request) {
	var users []User
	err := db.Model(&users).Select()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(users)
}

func main() {
	// 설정 초기화 (환경 변수 로드)
	initConfig()

	// 데이터베이스 초기화
	initDB()

	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			createUser(w, r)
		} else if r.Method == http.MethodGet {
			getUsers(w, r)
		}
	})

	fmt.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
