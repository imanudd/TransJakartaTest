package handler

import (
	"app/internal/delivery/http/helper"
	"app/internal/domain"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// GetVehicleHistory handler
// @Summary get vehicle latest location by vehicle_id
// @Description get vehicle latest location by vehicle_id
// @Tags vehicle
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param vehicle_id path string true "vehicle_id"
// @Query
// @Success 200 {object} helper.JSONResponse
// @Failure 400 {object} helper.JSONResponse
// @Failure 500 {object} helper.JSONResponse
// @Router /vehicles/{vehicle_id}/history [GET]
func (h *Handler) GetVehicleHistory(c *gin.Context) {
	var req domain.GetVehicleHistoryRequest

	if err := c.ShouldBind(&req); err != nil {
		helper.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	vehicleId := c.Param("vehicle_id")
	if vehicleId == "" {
		helper.Error(c, http.StatusBadRequest, "vehicle id is required")
		return
	}

	if req.Start == 0 {
		req.Start = time.Now().Unix()
	}

	if req.End == 0 {
		req.End = time.Now().Unix()
	}

	data, respCode := h.usecase.GetVehicleUseCase().GetHistory(c, vehicleId, &req)
	if respCode != nil {
		helper.Error(c, respCode.Code, respCode.Message)
		return
	}

	helper.Success(c, http.StatusOK, data)
}

// GetVehicleLocationLatest handler
// @Summary get vehicle latest location by vehicle_id
// @Description get vehicle latest location by vehicle_id
// @Tags vehicle
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param vehicle_id path string true "vehicle_id"
// @Success 200 {object} helper.JSONResponse
// @Failure 400 {object} helper.JSONResponse
// @Failure 500 {object} helper.JSONResponse
// @Router /vehicles/{vehicle_id}/location [GET]
func (h *Handler) GetVehicleLocationLatest(c *gin.Context) {
	vehicleId := c.Param("vehicle_id")
	if vehicleId == "" {
		helper.Error(c, http.StatusBadRequest, "vehicle id is required")
		return
	}

	data, respCode := h.usecase.GetVehicleUseCase().GetLatestLocation(c, vehicleId)
	if respCode != nil {
		helper.Error(c, respCode.Code, respCode.Message)
		return
	}

	helper.Success(c, http.StatusOK, data)
}
