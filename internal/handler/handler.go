package handler

import (
	"net/http"

	"yard-planning/internal/entity"
	"yard-planning/internal/services"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	svc *services.Service
	val *validator.Validate
}

func NewHandler(svc *services.Service) *Handler {
	return &Handler{svc: svc, val: validator.New()}
}

// POST /suggestion
func (h *Handler) Suggest(c echo.Context) error {
	var req entity.SuggestRequestDTO
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid body"})
	}
	if err := h.val.Struct(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	sr := entity.SuggestionRequest{
		Yard:            req.Yard,
		ContainerNumber: req.ContainerNumber,
		ContainerSize:   req.ContainerSize,
		ContainerHeight: req.ContainerHeight,
		ContainerType:   req.ContainerType,
	}
	pos, err := h.svc.SuggestPosition(c.Request().Context(), sr)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"suggested_position": pos,
	})
}

// POST /placement
func (h *Handler) Placement(c echo.Context) error {
	var req entity.PlacementDTO
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid body"})
	}
	if err := h.val.Struct(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	block, err := h.svc.RepoGetBlockByCode(c.Request().Context(), req.Block)
	if err != nil || block == nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "block not found"})
	}

	cp := &entity.ContainerPosition{
		ContainerNumber: req.ContainerNumber,
		BlockID:         block.ID,
		Slot:            req.Slot,
		Row:             req.Row,
		Tier:            req.Tier,
	}

	if err := h.svc.PlaceContainer(c.Request().Context(), cp); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Success"})
}

// POST /pickup
func (h *Handler) Pickup(c echo.Context) error {
	var req entity.PickupDTO
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid body"})
	}
	if err := h.val.Struct(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	if err := h.svc.PickupContainer(c.Request().Context(), req.ContainerNumber); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Success"})
}
