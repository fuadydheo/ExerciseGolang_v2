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
	log.Println("Fetching all users")
	json.NewEncoder(w).Encode(users)
}

// 🔹 Handler untuk mendapatkan pengguna berdasarkan ID
func getUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	log.Println("Fetching user with ID:", params["id"])
	for _, user := range users {
		if fmt.Sprintf("%d", user.ID) == params["id"] {
			json.NewEncoder(w).Encode(user)
		}
	}
	http.Error(w, "User not found", http.StatusNotFound) // ❌ BUG: Harusnya menggunakan 404
}

// 🔹 Handler untuk menambahkan pengguna baru
func createUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var newUser User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Println("Adding new User:", newUser.Name)
	users = append(users, newUser) // ❌ BUG: Tidak ada validasi apakah ID unik atau tidak
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newUser)
}

// 🔹 Handler untuk memperbarui data pengguna
func updateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	log.Println("Updating user with ID:", params["id"])
	for index, user := range users {
		if fmt.Sprintf("%d", user.ID) == params["id"] {
			var updatedUser User
			err := json.NewDecoder(r.Body).Decode(&updatedUser)
			if err != nil {
				log.Println("Failed to decode request body:", err)
				http.Error(w, "Invalid input", http.StatusBadRequest)
				return
			}
			users[index] = updatedUser // ❌ BUG: ID lama bisa berubah, harus tetap dipertahankan
			json.NewEncoder(w).Encode(updatedUser)
			return
		}
	}
	http.Error(w, "User not found", http.StatusNotFound)
}

// 🔹 Handler untuk menghapus pengguna berdasarkan ID
func deleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	log.Println("Deleting user with ID:", params["id"])
	for index, user := range users {
		if fmt.Sprintf("%d", user.ID) == params["id"] {
			users = append(users[:index], users[index+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return // ❌ BUG: Tidak ada response yang mengonfirmasi penghapusan
		}
	}
	http.Error(w, "User not found", http.StatusNotFound)
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
	fmt.Println("Server running on port 8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}
