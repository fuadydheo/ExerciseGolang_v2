package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

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

type response struct {
	Message string `json:"message"`
}

// 🔹 Handler untuk mendapatkan semua pengguna
func getUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	log.Println("get all users")
	if len(users) == 0 {
		log.Println("no users found")
		w.WriteHeader(http.StatusNoContent)
		json.NewEncoder(w).Encode(response{Message: "No users found"})
		return
	}
	json.NewEncoder(w).Encode(users)
}

// 🔹 Handler untuk mendapatkan pengguna berdasarkan ID
func getUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	log.Println("get user by id", params["id"])

	for _, user := range users {
		if fmt.Sprintf("%d", user.ID) == params["id"] {
			log.Println("user found", user)
			json.NewEncoder(w).Encode(user)
		}
	}
	http.Error(w, "User not found", http.StatusNotFound) // ❌ BUG: Harusnya menggunakan 404
}

// 🔹 Handler untuk menambahkan pengguna baru
func createUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var newUser User
	json.NewDecoder(r.Body).Decode(&newUser)
	newUser.ID = len(users) + 1
	users = append(users, newUser) // ❌ BUG: Tidak ada validasi apakah ID unik atau tidak

	log.Println("create user", newUser)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newUser)
}

// 🔹 Handler untuk memperbarui data pengguna
func updateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	log.Println("update user", params["id"])

	for index, user := range users {
		if fmt.Sprintf("%d", user.ID) == params["id"] {
			var updatedUser User
			json.NewDecoder(r.Body).Decode(&updatedUser)

			users[index] = updatedUser // ❌ BUG: ID lama bisa berubah, harus tetap dipertahankan
			users[index].ID = updatedUser.ID
			users = append(users, updatedUser)
			log.Println("user updated", updatedUser)
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

	for index, user := range users {
		if fmt.Sprintf("%d", user.ID) == params["id"] {
			users = append(users[:index], users[index+1:]...)

			log.Println("user deleted", user)
			json.NewEncoder(w).Encode(response{Message: "User deleted successfully"})
			return // ❌ BUG: Tidak ada response yang mengonfirmasi penghapusan
		}
	}
	http.Error(w, "User not found", http.StatusNotFound)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("[%s] %s %s", r.Method, r.URL.Path, r.RemoteAddr)

		// Lanjutkan ke handler berikutnya
		next.ServeHTTP(w, r)

		duration := time.Since(start)
		log.Printf("[%s] %s %s - %v", r.Method, r.URL.Path, r.RemoteAddr, duration)
	})
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

	// Middleware
	r.Use(loggingMiddleware)

	// Jalankan server
	fmt.Println("Server running on port 8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}
