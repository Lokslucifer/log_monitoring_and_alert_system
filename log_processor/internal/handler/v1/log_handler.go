package v1

import (
	"net/http"
	"strings"
	"log_processor/internal/dto"
	"log_processor/internal/utils"
	"github.com/gin-gonic/gin"
)

func (h *Handler) FilterLogsHandler(c *gin.Context) {
	levelsParam := c.Query("levels")
	search := strings.ToLower(c.Query("search"))
	fromStr := c.Query("from")
	toStr := c.Query("to")
	limit := utils.ParseInt(c.Query("limit"), 50)
	offset := utils.ParseInt(c.Query("offset"), 0)

	// Parse time strings
	fromTime, err := utils.ParseTimeString(fromStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid 'from' time"})
		return
	}

	toTime, err := utils.ParseTimeString(toStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid 'to' time"})
		return
	}

	// Validate time range
	if toTime.Before(fromTime) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "'to' time cannot be before 'from' time"})
		return
	}

	// Parse log levels
	levels := utils.ParseLevels(levelsParam)

	// Build filter DTO
	filterDTO := dto.FilterDTO{
		Levels: levels,
		Search: search,
		From:   fromTime,
		To:     toTime,
		Limit:  limit,
		Offset: offset,
	}

	// Call service layer
	data, err := h.ser.FilterLogs(filterDTO)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Send response
	c.JSON(http.StatusOK, gin.H{
		"message": "Success",
		"data":    data,
	})
}
