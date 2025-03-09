package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

var users = []User{
	{ID: 1, Name: "Alice", Email: "alice@example.com", Age: 25},
	{ID: 2, Name: "Bob", Email: "bob@example.com", Age: 30},
}

// 🔹 Handler untuk mendapatkan semua pengguna
func getUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// 🔹 Handler untuk mendapatkan pengguna berdasarkan ID
func getUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for _, user := range users {
		if fmt.Sprintf("%d", user.ID) == params["id"] {
			json.NewEncoder(w).Encode(user)
			return
		}
	}
	http.Error(w, "User not found", http.StatusNotFound)
}

// 🔹 Handler untuk menambahkan pengguna baru
func createUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var newUser User
	json.NewDecoder(r.Body).Decode(&newUser)

	for _, user := range users {
		if user.ID == newUser.ID {
			http.Error(w, "User ID already exists", http.StatusConflict)
			return
		}
	}

	users = append(users, newUser)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newUser)
}

// 🔹 Handler untuk memperbarui data pengguna
func updateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for index, user := range users {
		if fmt.Sprintf("%d", user.ID) == params["id"] {
			var updatedUser User
			json.NewDecoder(r.Body).Decode(&updatedUser)
			updatedUser.ID = user.ID
			users[index] = updatedUser
			json.NewEncoder(w).Encode(updatedUser)
			return
		}
	}
	http.Error(w, "User not found", 404)
}

// 🔹 Handler untuk menghapus pengguna berdasarkan ID
func deleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for index, user := range users {
		if fmt.Sprintf("%d", user.ID) == params["id"] {
			users = append(users[:index], users[index+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	http.Error(w, "User not found", 404)
}

// 🔹 Fungsi utama untuk menjalankan server
func main() {
	r := mux.NewRouter()

	// Rute API
	r.HandleFunc("/users", getUsers).Methods("GET")
	r.HandleFunc("/users/{id}", getUser).Methods("GET")
	r.HandleFunc("/users", createUser).Methods("POST")
	r.HandleFunc("/users/{id}", updateUser).Methods("PUT")
	r.HandleFunc("/users/{id}", deleteUser).Methods("DELETE")

	// Jalankan server
	fmt.Println("Server running on port 9000")
	log.Fatal(http.ListenAndServe(":9000", r))
}
