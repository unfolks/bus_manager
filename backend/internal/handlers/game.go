package handlers

import (
	"net/http"

	"bus-manager/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type GameHandler struct {
	db  *gorm.DB
	rdb *redis.Client
}

func NewGameHandler(db *gorm.DB, rdb *redis.Client) *GameHandler {
	return &GameHandler{
		db:  db,
		rdb: rdb,
	}
}

type CreateCompanyRequest struct {
	Name string `json:"name" binding:"required,min=3,max=100"`
}

type CreateDepotRequest struct {
	Name      string  `json:"name" binding:"required,min=3,max=100"`
	Latitude  float64 `json:"latitude" binding:"required"`
	Longitude float64 `json:"longitude" binding:"required"`
}

type CreateBusRequest struct {
	Name          string  `json:"name" binding:"required,min=3,max=100"`
	Type          string  `json:"type" binding:"required"`
	Capacity      int     `json:"capacity" binding:"required,min=1"`
	ServiceType   string  `json:"service_type" binding:"required"`
	PurchasePrice float64 `json:"purchase_price" binding:"required,min=0"`
}

type CreateTripRequest struct {
	BusID    uint `json:"bus_id" binding:"required"`
	RouteID  uint `json:"route_id" binding:"required"`
	DriverID uint `json:"driver_id"`
}

func (h *GameHandler) GetCompany(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	var company models.Company
	if err := h.db.Where("user_id = ?", userID).Preload("Depots").Preload("Buses").First(&company).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Company not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch company"})
		}
		return
	}

	c.JSON(http.StatusOK, company)
}

func (h *GameHandler) CreateCompany(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	// Check if company already exists
	var existingCompany models.Company
	if err := h.db.Where("user_id = ?", userID).First(&existingCompany).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Company already exists"})
		return
	}

	var req CreateCompanyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	company := models.Company{
		UserID:     userID.(uint),
		Name:       req.Name,
		Money:      1000000, // Starting money 1 million IDR
		Reputation: 0,
		Level:      1,
		Experience: 0,
	}

	if err := h.db.Create(&company).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create company"})
		return
	}

	// Create initial transaction
	transaction := models.Transaction{
		CompanyID:   company.ID,
		Type:        "income",
		Description: "Starting capital",
		Amount:      1000000,
		Balance:     1000000,
	}
	h.db.Create(&transaction)

	c.JSON(http.StatusCreated, company)
}

func (h *GameHandler) GetDepots(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	// Get user's company
	var company models.Company
	if err := h.db.Where("user_id = ?", userID).First(&company).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Company not found"})
		return
	}

	var depots []models.Depot
	if err := h.db.Where("company_id = ?", company.ID).Find(&depots).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch depots"})
		return
	}

	c.JSON(http.StatusOK, depots)
}

func (h *GameHandler) CreateDepot(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	// Get user's company
	var company models.Company
	if err := h.db.Where("user_id = ?", userID).First(&company).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Company not found"})
		return
	}

	var req CreateDepotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	depot := models.Depot{
		CompanyID:    company.ID,
		Name:         req.Name,
		Latitude:     req.Latitude,
		Longitude:    req.Longitude,
		Capacity:     10,
		CurrentBuses: 0,
		Level:        1,
	}

	if err := h.db.Create(&depot).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create depot"})
		return
	}

	c.JSON(http.StatusCreated, depot)
}

func (h *GameHandler) GetBuses(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	// Get user's company
	var company models.Company
	if err := h.db.Where("user_id = ?", userID).First(&company).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Company not found"})
		return
	}

	var buses []models.Bus
	if err := h.db.Where("company_id = ?", company.ID).Preload("Depot").Find(&buses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch buses"})
		return
	}

	c.JSON(http.StatusOK, buses)
}

func (h *GameHandler) CreateBus(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	// Get user's company
	var company models.Company
	if err := h.db.Where("user_id = ?", userID).First(&company).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Company not found"})
		return
	}

	var req CreateBusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if company has enough money
	if company.Money < req.PurchasePrice {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient funds"})
		return
	}

	// Get first depot (for now)
	var depot models.Depot
	if err := h.db.Where("company_id = ?", company.ID).First(&depot).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No depot found. Create a depot first."})
		return
	}

	// Check depot capacity
	if depot.CurrentBuses >= depot.Capacity {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Depot is at full capacity"})
		return
	}

	bus := models.Bus{
		CompanyID:     company.ID,
		DepotID:       depot.ID,
		Name:          req.Name,
		Type:          req.Type,
		Capacity:      req.Capacity,
		FuelCapacity:  100,
		CurrentFuel:   100,
		Range:         500,
		ServiceType:   req.ServiceType,
		Status:        "available",
		Condition:     100,
		PurchasePrice: req.PurchasePrice,
		OperatingCost: 1000, // Default operating cost per km
	}

	// Start transaction
	tx := h.db.Begin()

	// Create bus
	if err := tx.Create(&bus).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create bus"})
		return
	}

	// Update company money
	company.Money -= req.PurchasePrice
	if err := tx.Save(&company).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update company funds"})
		return
	}

	// Update depot current buses
	depot.CurrentBuses++
	if err := tx.Save(&depot).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update depot"})
		return
	}

	// Create transaction record
	transaction := models.Transaction{
		CompanyID:   company.ID,
		Type:        "expense",
		Description: "Purchased bus: " + req.Name,
		Amount:      -req.PurchasePrice,
		Balance:     company.Money,
	}
	tx.Create(&transaction)

	// Commit transaction
	tx.Commit()

	c.JSON(http.StatusCreated, bus)
}

func (h *GameHandler) GetRoutes(c *gin.Context) {
	var routes []models.Route
	if err := h.db.Find(&routes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch routes"})
		return
	}

	c.JSON(http.StatusOK, routes)
}

func (h *GameHandler) CreateTrip(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	// Get user's company
	var company models.Company
	if err := h.db.Where("user_id = ?", userID).First(&company).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Company not found"})
		return
	}

	var req CreateTripRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate bus ownership
	var bus models.Bus
	if err := h.db.Where("id = ? AND company_id = ?", req.BusID, company.ID).First(&bus).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bus not found or not owned by company"})
		return
	}

	// Check if bus is available
	if bus.Status != "available" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bus is not available"})
		return
	}

	// Get route
	var route models.Route
	if err := h.db.First(&route, req.RouteID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Route not found"})
		return
	}

	// Check if bus has enough fuel
	if bus.CurrentFuel < route.Distance/10 { // Assuming 1 liter per 10 km
		c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient fuel"})
		return
	}

	// Calculate revenue based on passengers and fare
	passengers := int(float64(bus.Capacity) * float64(route.Popularity) / 100.0)
	revenue := float64(passengers) * route.BaseFare
	cost := route.Distance * bus.OperatingCost
	profit := revenue - cost

	trip := models.Trip{
		BusID:      req.BusID,
		RouteID:    req.RouteID,
		DriverID:   req.DriverID,
		Status:     "planned",
		Passengers: passengers,
		Revenue:    revenue,
		Cost:       cost,
		Profit:     profit,
		Progress:   0,
	}

	if err := h.db.Create(&trip).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create trip"})
		return
	}

	// Update bus status
	bus.Status = "on_trip"
	h.db.Save(&bus)

	c.JSON(http.StatusCreated, trip)
}

func (h *GameHandler) GetActiveTrips(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	// Get user's company
	var company models.Company
	if err := h.db.Where("user_id = ?", userID).First(&company).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Company not found"})
		return
	}

	var trips []models.Trip
	if err := h.db.Where("status IN ? AND bus_id IN (SELECT id FROM buses WHERE company_id = ?)",
		[]string{"planned", "active"}, company.ID).
		Preload("Bus").
		Preload("Route").
		Preload("Driver").
		Find(&trips).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch active trips"})
		return
	}

	c.JSON(http.StatusOK, trips)
}
