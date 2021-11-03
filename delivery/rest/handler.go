package rest

import (
	"fmt"
	"io"
	"net/http"
	"os"
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
	imageLoc = "upload"
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

// GetProduct godoc
// @Summary Get Product.
// @Description get product.
// @Tags Product
// @Accept */*
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} map[string]interface{}
// @Router /product [get]

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

	dataReq.UrlImage, _ = c.FormFile("fileImage")

	uploadedFile, err := dataReq.UrlImage.Open()
	if err != nil {
		return err
	}
	defer uploadedFile.Close()

	tempFile, err := os.CreateTemp(imageLoc, fmt.Sprintf("%v", dataReq.ID))
	if err != nil {
		fmt.Println(err)
	}
	defer tempFile.Close()

	dataReq.Path = tempFile.Name()

	if _, err = io.Copy(tempFile, uploadedFile); err != nil {
		return err
	}

	if err := c.Bind(&dataReq); err != nil {
		c.JSON(http.StatusBadRequest, responseError{
			Message: "invalid data request",
		})
		return echo.ErrBadRequest
	}

	_, err = h.ProductUsecase.SendProduct(ctx, dataReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responseError{
			Message: "internal error",
		})

		return echo.ErrInternalServerError
	}
	return c.JSON(http.StatusCreated, responseError{
		Message: "success create",
	})
}

func (h *Handler) UpdateProduct(c echo.Context) error {
	ctx := c.Request().Context()
	dataReq := model.Product{}
	productIDParam := c.QueryParam("id")

	userInfo := c.Get("user").(*model.Token)

	if userInfo.Role != isAdmin {
		return echo.ErrUnauthorized
	}

	dataReq.UrlImage, _ = c.FormFile("fileImage")

	uploadedFile, err := dataReq.UrlImage.Open()
	if err != nil {
		return err
	}
	defer uploadedFile.Close()

	tempFile, err := os.CreateTemp(imageLoc, fmt.Sprintf("%v", dataReq.ID))
	if err != nil {
		fmt.Println(err)
	}
	defer tempFile.Close()

	dataReq.Path = tempFile.Name()

	if _, err = io.Copy(tempFile, uploadedFile); err != nil {
		return err
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

	_, err = h.ProductUsecase.UpdateProduct(ctx, dataReq, productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responseError{
			Message: "internal error",
		})
		fmt.Println(err)
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, responseError{
		Message: "update has been successful",
	})
}

func (h *Handler) DeleteProduct(c echo.Context) error {
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