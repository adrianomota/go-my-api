package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/adrianomota/fullcycle/my-api/internal/dto"
	"github.com/adrianomota/fullcycle/my-api/internal/entity"
	"github.com/adrianomota/fullcycle/my-api/internal/infra/database"
	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
)

type Error struct {
	Message string `json:"message"`
}

type UserHandler struct {
	UserDB       database.UserInterface
	Jwt          *jwtauth.JWTAuth
	JwtExpiresIn int
}

func NewUserHandler(db database.UserInterface, jwt *jwtauth.JWTAuth, jwtExpiresIn int) *UserHandler {
	return &UserHandler{UserDB: db,
		Jwt:          jwt,
		JwtExpiresIn: jwtExpiresIn,
	}
}

// GetJWT godoc
// @Summary 			Get a user JWT
// @Description  	Get a user JWT
// @Tags					users
// @Accept				json
// @Produce				json
// @Param					request		body				dto.GetJwtInput 	true		"user credentials"
// @Success				200				{object} 		dto.GetJwtOutput
// @Failure				404
// @Failure				500				{object}		Error
// @Router				/users/token	[post]
func (h *UserHandler) GetJWT(w http.ResponseWriter, r *http.Request) {
	var user dto.GetJwtInput
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	u, err := h.UserDB.FindByEmail(user.Email)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if !u.ValidatePassword(user.Password) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	_, tokenString, _ := h.Jwt.Encode(map[string]interface{}{
		"sub": u.ID,
		"exp": time.Now().Add(time.Second * time.Duration(h.JwtExpiresIn)).Unix(),
	})

	accessToken := dto.GetJwtOutput{AccessToken: tokenString}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(accessToken)
}

// Create user godoc
// @Summary 			Create User
// @Description  	Create User
// @Tags					users
// @Accept				json
// @Produce				json
// @Param					request		body			dto.CreateUserInput 	true		"user request"
// @Success				201
// @Failure				500				{object}	Error
// @Router				/users	[post]
func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var user dto.CreateUserInput
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	u, err := entity.NewUser(user.Name, user.Email, user.Password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.UserDB.Create(u)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errorMessage := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(errorMessage)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *UserHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Invalid param")
		return
	}
	user, err := h.UserDB.FindById(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		errorMessage := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(errorMessage)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}
