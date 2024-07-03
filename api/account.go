package api

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	db "github.com/uditsaurabh/simple_bank/db/sqlc"
)

type CreateAccountRequest struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required"`
}

func (server *Server) createAccount(ctx *gin.Context) {

	var req CreateAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	arg := db.CreateAccountParams{
		Owner:    req.Owner,
		Currency: req.Currency,
		Balance:  0,
	}
	store := server.store
	account, err := store.CreateAccount(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, account)
}

type GetAccountRequest struct {
	ID string `uri:"id" binding:"required,min=2" `
}

func (server *Server) getAccount(ctx *gin.Context) {
	var err error
	var req GetAccountRequest
	var num int
	var account db.Account

	if err = ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if num, err = strconv.Atoi(req.ID); err != nil {
		// Handle error if conversion fails
		log.Println("Error: Input cannot be converted to num.", err)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	store := server.store
	if account, err = store.GetAccount(ctx, int64(num)); err != nil {

		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, account)
}

type ListAccountRequest struct {
	PageNumber int64 `form:"page_id"`
	PageSize   int64 `form:"page_size"`
}

func (server *Server) listAccount(ctx *gin.Context) {
	var err error
	var req ListAccountRequest

	var accounts []db.Account

	if err = ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	page_number := int32(req.PageNumber)
	page_size := int32(req.PageSize)
	arg := db.ListAccountParams{
		Limit:  page_size,
		Offset: (page_number - 1) * page_size,
	}
	store := server.store
	if accounts, err = store.ListAccount(ctx, arg); err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, accounts)
}
