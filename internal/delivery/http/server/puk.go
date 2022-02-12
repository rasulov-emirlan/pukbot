package server

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/rasulov-emirlan/pukbot/internal/puk"
	"github.com/rasulov-emirlan/pukbot/pkg/logger"
)

type (
	pukListResponse struct {
		Pagination struct {
			Next string `json:"next"`
			Back string `json:"back"`
		} `json:"pagination"`
		Data struct {
			Puks []*puk.Puk `json:"puks"`
		} `json:"data"`
	}
)

// pukList docs
// @Tags		puks
// @Summary		List puks
// @Package server
// @Description	Returns you the links for the puks we have in our database
// @Accept		json
// @Produce		json
// @Param       page		query	int		true	"Page number, !first page is 0 not 1"
// @Param       limit		query	int		false	"Size of the page you want to get, !if not passed limit will be set to 10"
// @Success		200 		{object} pukListResponse
// @Router		/api/puks	[get]
func pukList(pukService puk.Service, l logger.Logger) echo.HandlerFunc {
	return func(c echo.Context) error {
		page, err := strconv.Atoi(c.QueryParam("page"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, errors.New("page queryparam was not provided").Error())
		}
		limit, err := strconv.Atoi(c.QueryParam("limit"))
		if err != nil {
			limit = 10
		}

		puks, err := pukService.List(c.Request().Context(), page, limit)
		if err != nil {
			l.Infof("error occured in 'server/pukList()' due to: %v", err)
			return c.JSON(http.StatusInternalServerError, err)
		}

		result := &pukListResponse{}
		result.Data.Puks = puks
		result.Pagination.Next = fmt.Sprintf("puks?page=%d&limit=%d", page+1, limit)
		result.Pagination.Back = fmt.Sprintf("puks?page=%d&limit=%d", page-1, limit)
		if page == 0 {
			result.Pagination.Back = ""
		}
		if limit != len(puks) {
			result.Pagination.Next = ""
		}
		return c.JSON(http.StatusOK, result)
	}
}
