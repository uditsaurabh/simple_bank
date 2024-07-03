package api

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/uditsaurabh/simple_bank/db/sqlc"
	t "github.com/uditsaurabh/simple_bank/token"
	util "github.com/uditsaurabh/simple_bank/util"
)

type LoginWithPassword struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (server *Server) LoginWithPassword(ctx *gin.Context) {
	var req LoginWithPassword
	var user db.User
	var err error
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if user, err = server.store.GetUser(ctx, req.Username); err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		} else {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}
	if res := util.DoPasswordMatch(user.HashedPassword, req.Password); res {
		config, err := util.GetConfig(".")
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(errors.New("config details not found")))
			return
		}
		m, err := t.NewPasetoMaker(config.EncryptionKey)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(errors.New("paesto maker object not created")))
			return
		}
		token, _, err := m.CreateToken(req.Username, "Admin", 10)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusCreated, token)
		return
	} else {
		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("authentication error")))
		return
	}

}

type CreateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required"`
}
type UserResponse struct {
	Username          string    `json:"username"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

type LoginUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
type LoginUserResponse struct {
	AccessToken string `json:"access_token"`
	User        string `json:"password" binding:"required"`
}

func (server *Server) CreateUser(ctx *gin.Context) {
	var req CreateUserRequest
	var arg db.CreateUserParams
	var hashed_password string

	var err error
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if hashed_password, err = util.HashPassword(req.Password); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg = db.CreateUserParams{
		Username:       req.Username,
		HashedPassword: hashed_password,
		FullName:       req.FullName,
		Email:          req.Email,
	}

	if _, err = server.store.CreateUser(ctx, arg); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusCreated, "user created")
}
