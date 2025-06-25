package main

import (
	"encoding/json"
	"net/http"
)

func main() {
	http.HandleFunc("/api/v1/students/4", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id":         4,
			"name":       "Alumno Demo",
			"email":      "alumno.demo@example.com",
			"class_name": "Grade 1",
		})
	})

	http.HandleFunc("/api/v1/auth/login", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Set-Cookie", "accessToken=fake-token")
		w.Header().Add("Set-Cookie", "refreshToken=fake-refresh")
		w.Header().Add("Set-Cookie", "csrfToken=fake-csrf")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success": true}`))
	})

	println("Mock backend en http://localhost:5007")
	http.ListenAndServe(":5007", nil)
}
