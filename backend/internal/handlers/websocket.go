package handlers

import (
	"log"
	"net/http"
	"time"

	"bus-manager/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for development
	},
}

type WSMessage struct {
	Type   string      `json:"type"`
	Data   interface{} `json:"data"`
	TripID uint        `json:"trip_id,omitempty"`
	BusID  uint        `json:"bus_id,omitempty"`
}

type WSClient struct {
	hub    *WSHub
	conn   *websocket.Conn
	send   chan WSMessage
	userID uint
}

type WSHub struct {
	clients    map[*WSClient]bool
	broadcast  chan WSMessage
	register   chan *WSClient
	unregister chan *WSClient
	db         *gorm.DB
	rdb        *redis.Client
}

func NewWSHub(db *gorm.DB, rdb *redis.Client) *WSHub {
	return &WSHub{
		clients:    make(map[*WSClient]bool),
		broadcast:  make(chan WSMessage),
		register:   make(chan *WSClient),
		unregister: make(chan *WSClient),
		db:         db,
		rdb:        rdb,
	}
}

func (h *WSHub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			log.Printf("Client connected. Total clients: %d", len(h.clients))

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				log.Printf("Client disconnected. Total clients: %d", len(h.clients))
			}

		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

func (c *WSClient) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	for {
		var message WSMessage
		err := c.conn.ReadJSON(&message)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		// Handle incoming messages
		switch message.Type {
		case "subscribe_trip":
			// Handle trip subscription
			if tripID, ok := message.Data.(float64); ok {
				c.subscribeToTrip(uint(tripID))
			}
		}
	}
}

func (c *WSClient) writePump() {
	defer func() {
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.conn.WriteJSON(message); err != nil {
				log.Printf("WebSocket write error: %v", err)
				return
			}
		}
	}
}

func (c *WSClient) subscribeToTrip(tripID uint) {
	// Send initial trip data
	var trip models.Trip
	if err := c.hub.db.Preload("Bus").Preload("Route").Preload("Driver").First(&trip, tripID).Error; err == nil {
		message := WSMessage{
			Type: "trip_update",
			Data: trip,
		}
		c.send <- message
	}
}

func HandleWebSocket(c *gin.Context) {
	// This is a placeholder - in a real implementation, you would:
	// 1. Authenticate the WebSocket connection
	// 2. Get the user ID from the JWT token
	// 3. Create a new client and register it with the hub

	// For now, we'll create a simple echo WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	client := &WSClient{
		conn: conn,
		send: make(chan WSMessage, 256),
	}

	// Start goroutines
	go client.writePump()
	go client.readPump()

	// Send welcome message
	welcome := WSMessage{
		Type: "welcome",
		Data: "Connected to Bus Manager WebSocket",
	}
	client.send <- welcome
}

// Trip simulation functions
func (h *WSHub) StartTripSimulation(tripID uint) {
	go h.simulateTrip(tripID)
}

func (h *WSHub) simulateTrip(tripID uint) {
	var trip models.Trip
	if err := h.db.Preload("Bus").Preload("Route").First(&trip, tripID).Error; err != nil {
		return
	}

	// Update trip status to active
	trip.Status = "active"
	trip.ActualStart = time.Now()
	h.db.Save(&trip)

	// Broadcast trip start
	message := WSMessage{
		Type:   "trip_started",
		Data:   trip,
		TripID: tripID,
	}
	h.broadcast <- message

	// Simulate trip progress
	steps := 10
	for i := 0; i <= steps; i++ {
		progress := float64(i) / float64(steps) * 100
		trip.Progress = progress

		// Calculate current position (linear interpolation)
		if i > 0 {
			ratio := float64(i) / float64(steps)
			trip.CurrentLat = trip.Route.OriginLat + (trip.Route.DestLat-trip.Route.OriginLat)*ratio
			trip.CurrentLng = trip.Route.OriginLng + (trip.Route.DestLng-trip.Route.OriginLng)*ratio
		}

		// Update database
		h.db.Save(&trip)

		// Broadcast update
		message := WSMessage{
			Type:   "trip_progress",
			Data:   trip,
			TripID: tripID,
		}
		h.broadcast <- message

		// Sleep for simulation (trip duration / steps)
		time.Sleep(time.Duration(trip.Route.Duration/steps) * time.Second)
	}

	// Complete trip
	trip.Status = "completed"
	trip.ActualEnd = time.Now()
	trip.Progress = 100
	trip.CurrentLat = trip.Route.DestLat
	trip.CurrentLng = trip.Route.DestLng
	h.db.Save(&trip)

	// Update bus status
	var bus models.Bus
	if err := h.db.First(&bus, trip.BusID).Error; err == nil {
		bus.Status = "available"
		bus.CurrentFuel -= trip.Route.Distance / 10 // Consume fuel
		h.db.Save(&bus)
	}

	// Broadcast completion
	completionMessage := WSMessage{
		Type:   "trip_completed",
		Data:   trip,
		TripID: tripID,
	}
	h.broadcast <- completionMessage
}
