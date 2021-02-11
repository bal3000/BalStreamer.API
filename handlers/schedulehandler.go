package handlers

import (
	"github.com/bal3000/BalStreamer.API/infrastructure"
	"net/http"

	"github.com/labstack/echo/v4"
)

// ScheduleHandler is the handler struct for schedule endpoints
type ScheduleHandler struct {
	RabbitMQ infrastructure.RabbitMQ
}

// NewScheduleHandler creates a new pointer to schedule
func NewScheduleHandler(rabbit infrastructure.RabbitMQ) *ScheduleHandler {
	return &ScheduleHandler{RabbitMQ: rabbit}
}

// AddEventToSchedule sends the event to the schedule app and logs a copy
func (handler *ScheduleHandler) AddEventToSchedule(c echo.Context) error {
	// Get info from post object and create a rabbit message

	// Send message to rabbit and also save to db if needed

	// Return success
	return c.String(http.StatusOK, "Successfully added event to schedule")
}
