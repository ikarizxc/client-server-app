package main

import (
	"dcsa-lab/internal/entities"
	"dcsa-lab/internal/repository"
	"dcsa-lab/internal/tokens"
	"encoding/json"
	"errors"
	"net/http"

	"strconv"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"
)

var (
	repo = repository.NewRepository()
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", home)

	mux.HandleFunc("/auth/signup", signUp)
	mux.HandleFunc("/auth/login", login)

	mux.Handle("/users", authMiddleware(http.HandlerFunc(getAllUsers)))
	mux.Handle("GET /users/{id}", authMiddleware(http.HandlerFunc(getUserById)))
	mux.Handle("DELETE /users/{id}", authMiddleware(adminRightRequired(http.HandlerFunc(deleteUserById))))
	mux.Handle("PATCH /users/{id}", authMiddleware(adminRightRequired(http.HandlerFunc(updateUserById))))

	s := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	s.ListenAndServe()
}

func home(res http.ResponseWriter, req *http.Request) {
	data := []byte("home page")
	res.WriteHeader(200)
	res.Write(data)
}

func signUp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Write([]byte("invalid method"))
		return
	}
	headerContentTtype := r.Header.Get("Content-Type")
	if headerContentTtype != "application/json" {
		errorResponse(w, "Content Type is not application/json", http.StatusUnsupportedMediaType)
		return
	}

	var user struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var unmarshalErr *json.UnmarshalTypeError

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&user)
	if err != nil {
		if errors.As(err, &unmarshalErr) {
			errorResponse(w, "Bad Request. Wrong Type provided for field "+unmarshalErr.Field, http.StatusBadRequest)
		} else {
			errorResponse(w, "Bad Request "+err.Error(), http.StatusBadRequest)
		}
		return
	}

	if user.Email == "" || user.Password == "" || user.Username == "" {
		errorResponse(w, "fill in all the fields", http.StatusConflict)
		return
	}

	if _, _, ok := repo.GetByEmail(user.Email); ok == nil {
		errorResponse(w, "this email address already taken", http.StatusConflict)
		return
	}
	if _, ok := repo.GetByUsername(user.Username); ok == nil {
		errorResponse(w, "this username already taken", http.StatusConflict)
		return
	}

	userToRepo := &entities.User{
		Username: user.Username,
		Email:    user.Email,
		Password: user.Password,
		IsAdmin:  false,
	}

	repo.Add(userToRepo)

	errorResponse(w, "Success", http.StatusOK)
}

func login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Write([]byte("invalid method"))
		return
	}

	headerContentTtype := r.Header.Get("Content-Type")
	if headerContentTtype != "application/json" {
		errorResponse(w, "Content Type is not application/json", http.StatusUnsupportedMediaType)
		return
	}

	var userLogin struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var unmarshalErr *json.UnmarshalTypeError

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&userLogin)
	if err != nil {
		if errors.As(err, &unmarshalErr) {
			errorResponse(w, "Bad Request. Wrong Type provided for field "+unmarshalErr.Field, http.StatusBadRequest)
		} else {
			errorResponse(w, "Bad Request "+err.Error(), http.StatusBadRequest)
		}
		return
	}

	id, user, err := repo.GetByEmail(userLogin.Email)
	if err != nil {
		errorResponse(w, "the user with this email is not registered", http.StatusUnauthorized)
		return
	}
	if user.Password != userLogin.Password {
		errorResponse(w, "wrong password", http.StatusUnauthorized)
		return
	}

	token, err := tokens.GenerateAccessToken(id)
	if err != nil {
		errorResponse(w, "intenal server error", http.StatusInternalServerError)
		return
	}

	jsonContent, err := jsoniter.Marshal(map[string]string{
		"token": token,
	})
	if err != nil {
		errorResponse(w, "internal server error", 500)
		return
	}

	w.Header()["Content-Type"] = []string{"application/json"}
	w.Write([]byte(jsonContent))
	w.WriteHeader(200)
}

func adminRightRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header["Authorization"]
		headersParts := strings.Split(header[0], " ")
		claims, _ := tokens.ParseJWT(headersParts[1])
		idstr, _ := claims.GetSubject()
		id, _ := strconv.Atoi(idstr)

		currentUser, _ := repo.GetById(id)

		if !currentUser.IsAdmin {
			errorResponse(w, "forbidden", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header["Authorization"]
		if len(header) == 0 {
			errorResponse(w, "missing authorization header", http.StatusUnauthorized)
			return
		}

		headersParts := strings.Split(header[0], " ")
		if len(headersParts) != 2 || headersParts[0] != "Bearer" {
			errorResponse(w, "invalid authorization header", http.StatusUnauthorized)
			return
		}

		claims, err := tokens.ParseJWT(headersParts[1])
		if err != nil {
			errorResponse(w, "invalid authorization token", http.StatusUnauthorized)
			return
		}

		expTime, err := claims.GetExpirationTime()
		if err != nil {
			errorResponse(w, "invalid authorization token", http.StatusUnauthorized)
			return
		}

		if expTime.Compare(time.Now()) == -1 {
			errorResponse(w, "token is expired", http.StatusUnauthorized)
			return
		}

		idstr, err := claims.GetSubject()
		if err != nil {
			errorResponse(w, "invalid authorization token", http.StatusUnauthorized)
			return
		}

		id, err := strconv.Atoi(idstr)
		if err != nil {
			errorResponse(w, "invalid authorization token", http.StatusUnauthorized)
			return
		}

		if user, _ := repo.GetById(id); user == nil {
			errorResponse(w, "invalid authorization token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func getAllUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Write([]byte("invalid method"))
		return
	}

	jsonContent, err := jsoniter.Marshal(repo.GetAll())
	if err != nil {
		errorResponse(w, "internal server error", 500)
		return
	}

	w.Header()["Content-Type"] = []string{"application/json"}
	w.Write(jsonContent)
	w.WriteHeader(200)
}

func getUserById(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Write([]byte("invalid method"))
		return
	}

	idstr := r.PathValue("id")
	if idstr == "" {
		errorResponse(w, "missing parameter 'id'", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idstr)
	if err != nil {
		errorResponse(w, "id must be a number", http.StatusBadRequest)
		return
	}

	user, err := repo.GetById(id)
	if err != nil {
		errorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	jsonContent, err := jsoniter.Marshal(user)
	if err != nil {
		errorResponse(w, "internal server error", 500)
		return
	}

	w.Header()["Content-Type"] = []string{"application/json"}
	w.Write(jsonContent)
	w.WriteHeader(200)
}

func deleteUserById(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		w.Write([]byte("invalid method"))
		return
	}

	idstr := r.PathValue("id")
	if idstr == "" {
		errorResponse(w, "missing parameter 'id'", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idstr)
	if err != nil {
		errorResponse(w, "id must be a number", http.StatusBadRequest)
		return
	}

	repo.Delete(id)
	errorResponse(w, "Success", http.StatusOK)
}

func updateUserById(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		w.Write([]byte("invalid method"))
		return
	}

	idstr := r.PathValue("id")
	if idstr == "" {
		errorResponse(w, "missing parameter 'id'", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idstr)
	if err != nil {
		errorResponse(w, "id must be a number", http.StatusBadRequest)
		return
	}

	headerContentTtype := r.Header.Get("Content-Type")
	if headerContentTtype != "application/json" {
		errorResponse(w, "Content Type is not application/json", http.StatusUnsupportedMediaType)
		return
	}

	var user struct {
		Username string `json:"username,omitempty"`
		Email    string `json:"email,omitempty"`
		Password string `json:"password,omitempty"`
	}

	var unmarshalErr *json.UnmarshalTypeError

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&user)
	if err != nil {
		if errors.As(err, &unmarshalErr) {
			errorResponse(w, "Bad Request. Wrong Type provided for field "+unmarshalErr.Field, http.StatusBadRequest)
		} else {
			errorResponse(w, "Bad Request "+err.Error(), http.StatusBadRequest)
		}
		return
	}

	userToRepo := &entities.User{
		Username: "",
		Email:    "",
		Password: "",
		IsAdmin:  false,
	}

	if user.Username != "" {
		userToRepo.Username = user.Username
	}
	if user.Email != "" {
		userToRepo.Email = user.Email
	}
	if user.Password != "" {
		userToRepo.Password = user.Password
	}

	err = repo.Update(id, userToRepo)
	if err != nil {
		errorResponse(w, err.Error(), http.StatusNotFound)
		return
	}

	errorResponse(w, "Success", http.StatusOK)
}

func errorResponse(w http.ResponseWriter, message string, httpStatusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatusCode)
	resp := make(map[string]string)
	resp["message"] = message
	jsonResp, _ := json.Marshal(resp)
	w.Write(jsonResp)
}
