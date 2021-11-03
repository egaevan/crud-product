package rest

import (
	"fmt"
	"net/http"
	"strconv"

	"crud-product/model"
	"crud-product/usecase"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	ProductUsecase usecase.ProductUsecase
	UserUsecase    usecase.UserUsecase
}

type responseError struct {
	Message string `json:"message"`
}

const (
	isAdmin int = 1
)

func NewHandler(e *echo.Echo, productUsecase usecase.ProductUsecase, userUsecase usecase.UserUsecase) {
	handler := &Handler{
		ProductUsecase: productUsecase,
		UserUsecase:    userUsecase,
	}

	// Routing Product
	e.GET("/product", handler.GetProduct, JwtVerify)
	e.GET("/product/brand", handler.GetProductAll, JwtVerify)
	e.POST("/product", handler.SendProduct, JwtVerify)
	e.PATCH("/product", handler.UpdateProduct, JwtVerify)
	e.DELETE("/product", handler.DeleteProduct, JwtVerify)

	// Routing User
	e.POST("/login", handler.Login)
	e.POST("/register", handler.Register)
}

func (h *Handler) GetProduct(c echo.Context) error {
	ctx := c.Request().Context()
	productIDParam := c.QueryParam("id")

	userInfo := c.Get("user").(*model.Token)

	if userInfo.Role != isAdmin {
		return echo.ErrUnauthorized
	}

	if productIDParam == "" {
		c.JSON(http.StatusBadRequest, responseError{
			Message: "invalid parameter",
		})

		return echo.ErrBadRequest
	}

	productID, err := strconv.Atoi(productIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, responseError{
			Message: "invalid parameter",
		})

		return echo.ErrBadRequest
	}

	res, err := h.ProductUsecase.GetProduct(ctx, productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responseError{
			Message: fmt.Sprint(err),
		})

		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, res)
}

func (h *Handler) GetProductAll(c echo.Context) error {
	ctx := c.Request().Context()
	brandIDParam := c.QueryParam("id")

	userInfo := c.Get("user").(*model.Token)

	if userInfo.Role != isAdmin {
		return echo.ErrUnauthorized
	}

	if brandIDParam == "" {
		c.JSON(http.StatusBadRequest, responseError{
			Message: "invalid parameter",
		})

		return echo.ErrBadRequest
	}

	brandID, err := strconv.Atoi(brandIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, responseError{
			Message: "invalid parameter",
		})
		return echo.ErrBadRequest
	}

	res, err := h.ProductUsecase.GetProductAll(ctx, brandID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responseError{
			Message: "internal error",
		})
		return echo.ErrInternalServerError
	}
	return c.JSON(http.StatusOK, res)
}

func (h *Handler) SendProduct(c echo.Context) error {
	ctx := c.Request().Context()
	dataReq := model.Product{}

	userInfo := c.Get("user").(*model.Token)

	if userInfo.Role != isAdmin {
		return echo.ErrUnauthorized
	}

	if err := c.Bind(&dataReq); err != nil {
		c.JSON(http.StatusBadRequest, responseError{
			Message: "invalid data request",
		})
		return echo.ErrBadRequest
	}

	res, err := h.ProductUsecase.SendProduct(ctx, dataReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responseError{
			Message: "internal error",
		})

		return echo.ErrInternalServerError
	}
	return c.JSON(http.StatusCreated, res)
}

func (h *Handler) UpdateProduct(c echo.Context) error {
	ctx := c.Request().Context()
	dataReq := model.Product{}
	productIDParam := c.QueryParam("id")

	userInfo := c.Get("user").(*model.Token)

	if userInfo.Role != isAdmin {
		return echo.ErrUnauthorized
	}

	if productIDParam == "" {
		c.JSON(http.StatusBadRequest, responseError{
			Message: "invalid parameter",
		})

		return echo.ErrBadRequest
	}

	productID, err := strconv.Atoi(productIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, responseError{
			Message: "invalid parameter",
		})

		return echo.ErrBadRequest
	}

	if err := c.Bind(&dataReq); err != nil {
		c.JSON(http.StatusBadRequest, responseError{
			Message: "invalid data request",
		})
		return echo.ErrBadRequest
	}

	res, err := h.ProductUsecase.UpdateProduct(ctx, dataReq, productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responseError{
			Message: "internal error",
		})

		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, res)
}

func (h *Handler) DeleteProduct(c echo.Context) error {
	ctx := c.Request().Context()
	productIDParam := c.Param("productID")

	userInfo := c.Get("user").(*model.Token)

	if userInfo.Role != isAdmin {
		return echo.ErrUnauthorized
	}

	if productIDParam == "" {
		c.JSON(http.StatusBadRequest, responseError{
			Message: "invalid parameter",
		})

		return echo.ErrBadRequest
	}

	productID, err := strconv.Atoi(productIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, responseError{
			Message: "invalid parameter",
		})

		return echo.ErrBadRequest
	}

	err = h.ProductUsecase.DeleteProduct(ctx, productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responseError{
			Message: "internal error",
		})

		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, responseError{
		Message: "Product has been deleted",
	})
}

func (h *Handler) Login(c echo.Context) error {
	dataReq := model.User{}
	if err := c.Bind(&dataReq); err != nil {
		c.JSON(http.StatusBadRequest, responseError{
			Message: "invalid data request",
		})
		return echo.ErrBadRequest
	}

	user, err := h.UserUsecase.Login(c.Request().Context(), dataReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, responseError{
			Message: err.Error(),
		})
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "logged in",
		"token":   user.Token,
	})
}

func (h *Handler) Register(c echo.Context) error {
	dataReq := model.User{}
	if err := c.Bind(&dataReq); err != nil {
		c.JSON(http.StatusBadRequest, responseError{
			Message: "invalid data request",
		})

		return echo.ErrBadRequest
	}

	err := h.UserUsecase.CreateUser(c.Request().Context(), dataReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, responseError{
			Message: err.Error(),
		})
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusCreated, "success")
}