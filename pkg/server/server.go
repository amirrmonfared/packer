package server

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/amirrmonfared/packer/pkg/fulfillment"
	"github.com/amirrmonfared/packer/pkg/store"
)

type PackHandler struct {
	Store *store.MemoryStore
}

func NewPackHandler(s *store.MemoryStore) *PackHandler {
	return &PackHandler{Store: s}
}

func RegisterRoutes(r *gin.Engine, memoryStore *store.MemoryStore) {
	handler := NewPackHandler(memoryStore)

	api := r.Group("/api/v1")
	{
		api.POST("/packs", handler.UpdatePackSizes)
		api.GET("/packs", handler.GetPackSizes)
		api.POST("/calculate", handler.CalculatePacks)
	}
}

// UpdatePackSizesRequest is the expected input for updating available pack sizes.
type UpdatePackSizesRequest struct {
	Packs []int `json:"packs"`
}

// UpdatePackSizesResponse is the output after updating pack sizes.
type UpdatePackSizesResponse struct {
	Packs []int `json:"packs"`
}

// UpdatePackSizes updates available pack sizes.
//
// @Summary     Update pack sizes
// @Description Updates the list of available pack sizes (e.g., 250, 500, 1000, etc.).
// @Tags        Packs
// @Accept      json
// @Produce     json
// @Param       packs body UpdatePackSizesRequest true "New pack sizes"
// @Success     200 {object} UpdatePackSizesResponse
// @Failure     400 {object} map[string]string "Invalid JSON body or no packs provided"
// @Failure     500 {object} map[string]string "Failed to update pack sizes"
// @Router      /packs [post]
func (h *PackHandler) UpdatePackSizes(c *gin.Context) {
	var req UpdatePackSizesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logrus.WithError(err).Error("Failed to bind JSON for UpdatePackSizes")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON body"})
		return
	}

	if len(req.Packs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No packs provided"})
		return
	}

	if err := h.Store.UpdatePackSizes(context.Background(), req.Packs); err != nil {
		logrus.WithError(err).Error("Failed to update pack sizes in memory")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update pack sizes"})
		return
	}

	c.JSON(http.StatusOK, UpdatePackSizesResponse(req))
}

// GetPackSizesResponse is what we return from listing the current available pack sizes.
type GetPackSizesResponse struct {
	Packs []int `json:"packs"`
}

// GetPackSizes returns the currently stored pack sizes.
//
// @Summary     Get pack sizes
// @Description Retrieves the list of currently available pack sizes.
// @Tags        Packs
// @Produce     json
// @Success     200 {object} GetPackSizesResponse
// @Failure     500 {object} map[string]string "Failed to retrieve pack sizes"
// @Router      /packs [get]
func (h *PackHandler) GetPackSizes(c *gin.Context) {
	packs, err := h.Store.GetPackSizes(c)
	if err != nil {
		logrus.WithError(err).Error("Failed to retrieve pack sizes from memory")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve pack sizes"})
		return
	}

	c.JSON(http.StatusOK, GetPackSizesResponse{Packs: packs})
}

// CalculatePacksRequest is the expected input for calculating how many packs to ship.
type CalculatePacksRequest struct {
	Items int `json:"items" binding:"required"`
}

// CalculatePacksResponse describes the result of packing an order.
type CalculatePacksResponse struct {
	Order             int         `json:"order"`
	Leftover          int         `json:"leftover"`
	TotalPacks        int         `json:"total_packs"`
	Distribution      map[int]int `json:"distribution"`
	TotalItemsShipped int         `json:"total_items_shipped"`
}

// CalculatePacks calculates how many packs to ship for the requested number of items.
//
// @Summary     Calculate how many packs are needed
// @Description Given a number of items, this calculates how many packs and which sizes to use.
// @Tags        Packs
// @Accept      json
// @Produce     json
// @Param       items body CalculatePacksRequest true "Number of items to pack"
// @Success     200 {object} CalculatePacksResponse
// @Failure     400 {object} map[string]string "Invalid JSON body or items < 1"
// @Failure     500 {object} map[string]string "Failed to retrieve pack sizes"
// @Router      /calculate [post]
func (h *PackHandler) CalculatePacks(c *gin.Context) {
	var req CalculatePacksRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON body"})
		return
	}

	if req.Items < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Items must be > 0"})
		return
	}

	packs, err := h.Store.GetPackSizes(context.Background())
	if err != nil {
		logrus.WithError(err).Error("Could not retrieve pack sizes for calculation")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve pack sizes"})
		return
	}

	plan := fulfillment.CalculateShipmentPlan(req.Items, packs)

	shipped := 0
	for size, count := range plan.BoxCounts {
		shipped += size * count
	}

	resp := CalculatePacksResponse{
		Order:             req.Items,
		Leftover:          plan.ExtraUnits,
		TotalPacks:        plan.TotalBoxes,
		Distribution:      plan.BoxCounts,
		TotalItemsShipped: shipped,
	}
	c.JSON(http.StatusOK, resp)
}
