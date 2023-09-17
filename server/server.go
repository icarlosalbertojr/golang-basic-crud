package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"my-db-module/db"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type user struct {
	Id    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func getDatabaseConnection() *sql.DB {
	conn, err := db.Connect()

	if err != nil {
		log.Fatal(err)
	}
	return conn
}

func CreateNewUser(w http.ResponseWriter, r *http.Request) {

	conn := getDatabaseConnection()

	defer conn.Close()

	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte("Falha ao ler o body da requisição"))
		return
	}

	var user user

	err = json.Unmarshal(requestBody, &user)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		w.Write([]byte("Falha ao deserializar json"))
		return
	}

	stmt, err := conn.Prepare("INSERT INTO users (name, email) VALUES (?, ?)")
	if err != nil {
		w.Write([]byte("Erro ao criar o statement"))
		fmt.Println(err)
		return
	}
	defer stmt.Close()

	insert, err := stmt.Exec(user.Name, user.Email)
	if err != nil {
		w.Write([]byte("Erro ao inserir dados"))
		fmt.Println(err)
		return
	}

	id, err := insert.LastInsertId()
	user.Id = id

	data, err := json.Marshal(user)

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(201)
	w.Write([]byte(data))
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	conn := getDatabaseConnection()
	defer conn.Close()

	rows, err := conn.Query("SELECT * FROM users")
	defer rows.Close()

	if err != nil {
		w.Write([]byte("Erro ao buscar usuários"))
		fmt.Println(err)
	}

	var users []user
	for rows.Next() {
		var user user
		if err = rows.Scan(&user.Id, &user.Name, &user.Email); err != nil {
			w.Write([]byte("Erro ao escanear resuldato da consulta"))
			return
		}
		users = append(users, user)
		fmt.Println(user)
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(users); err != nil {
		w.Write([]byte("Erro ao converter resposta para json"))
	}
}

func GetUsersById(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userId, err := strconv.ParseUint(params["id"], 10, 32)
	if err != nil {
		fmt.Println(err)
		w.Write([]byte("Erro ao capturar parametro"))
	}

	conn := getDatabaseConnection()
	defer conn.Close()

	rows, err := conn.Query("SELECT * FROM users WHERE id = ?", userId)
	defer rows.Close()

	if rows.Next() {
		var user user
		if err = rows.Scan(&user.Id, &user.Name, &user.Email); err != nil {
			w.Write([]byte("Erro ao ler os dados"))
			return
		}

		w.Header().Add("content-type", "application/json")
		if err = json.NewEncoder(w).Encode(user); err != nil {
			w.Write([]byte("Erro ao serializar dados"))
			return
		}
	}
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte("Erro ao ler os dados"))
		return
	}

	var user user
	if err = json.Unmarshal(body, &user); err != nil {
		w.Write([]byte("Erro ao ler body"))
		return
	}

	userId, err := strconv.ParseUint(params["id"], 10, 32)
	if err != nil {
		w.Write([]byte("Erro ao ler parametro"))
		return
	}

	conn := getDatabaseConnection()
	defer conn.Close()

	stmt, err := conn.Prepare("UPDATE users SET name = ?, email = ? WHERE id = ?")
	if err != nil {
		w.Write([]byte("Erro ao criar update"))
		return
	}
	defer stmt.Close()

	if _, err = stmt.Exec(user.Name, user.Email, userId); err != nil {
		fmt.Println(err)
		w.Write([]byte("Erro ao executar update"))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func DeleteUserById(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	userId, err := strconv.ParseUint(params["id"], 10, 32)
	if err != nil {
		w.Write([]byte("Erro ao ler parametro"))
		return
	}

	conn := getDatabaseConnection()
	defer conn.Close()

	stmt, err := conn.Prepare("DELETE FROM users WHERE id = ?")
	defer stmt.Close()
	if err != nil {
		w.Write([]byte("Erro ao criar statement"))
		return
	}

	if _, err = stmt.Exec(userId); err != nil {
		w.Write([]byte("Erro ao executar exclusão"))
		return
	}

	w.WriteHeader(204)
}
