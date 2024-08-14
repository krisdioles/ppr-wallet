package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/krisdioles/ppr-wallet/app/domain"
	"github.com/krisdioles/ppr-wallet/app/domain/errors"
)

type UserBalanceController struct {
	UserBalanceUsecase domain.UserBalanceUsecase
}

func NewUserBalanceController(userBalanceUsecase domain.UserBalanceUsecase) *UserBalanceController {
	return &UserBalanceController{
		UserBalanceUsecase: userBalanceUsecase,
	}
}

func (c *UserBalanceController) GetUserBalanceByID(gc *gin.Context) {
	ctx := gc.Request.Context()

	idParam, err := strconv.Atoi(gc.Param("id"))
	if err != nil {
		gc.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": errors.ErrInvalidParameter.Error(),
		})
		return
	}

	userBalance, err := c.UserBalanceUsecase.GetUserBalanceByID(ctx, int64(idParam))
	if err != nil {
		switch err {
		case errors.ErrUserNotFound:
			gc.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
		default:
			gc.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": errors.ErrInternalServerError.Error(),
			})
		}

		return
	}

	gc.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "success",
		"data":    userBalance,
	})
}

func (c *UserBalanceController) DisburseBalance(gc *gin.Context) {
	ctx := gc.Request.Context()

	idParam, err := strconv.Atoi(gc.Param("id"))
	if err != nil {
		gc.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": errors.ErrInvalidParameter.Error(),
		})
		return
	}

	err = c.UserBalanceUsecase.DisburseBalance(ctx, int64(idParam))
	if err != nil {
		switch err {
		case errors.ErrInsufficientBalance:
			gc.JSON(http.StatusUnprocessableEntity, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
		case errors.ErrUserNotFound:
			gc.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
		default:
			gc.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": errors.ErrInternalServerError.Error(),
			})
		}
		return
	}

	gc.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "success",
	})
}
