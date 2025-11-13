package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Email     string         `json:"email" gorm:"uniqueIndex;not null"`
	Username  string         `json:"username" gorm:"uniqueIndex;not null"`
	Password  string         `json:"-" gorm:"not null"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relations
	Company Company `json:"company" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

type Company struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	UserID     uint      `json:"user_id" gorm:"not null"`
	Name       string    `json:"name" gorm:"not null"`
	Money      float64   `json:"money" gorm:"default:1000000"` // Starting money in IDR
	Reputation int       `json:"reputation" gorm:"default:0"`
	Level      int       `json:"level" gorm:"default:1"`
	Experience int       `json:"experience" gorm:"default:0"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`

	// Relations
	Depots []Depot `json:"depots"`
	Buses  []Bus   `json:"buses"`
}

type Depot struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	CompanyID    uint      `json:"company_id" gorm:"not null"`
	Name         string    `json:"name" gorm:"not null"`
	Latitude     float64   `json:"latitude" gorm:"not null"`
	Longitude    float64   `json:"longitude" gorm:"not null"`
	Capacity     int       `json:"capacity" gorm:"default:10"` // Maximum buses
	CurrentBuses int       `json:"current_buses" gorm:"default:0"`
	Level        int       `json:"level" gorm:"default:1"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	// Relations
	Company Company `json:"company" gorm:"foreignKey:CompanyID"`
	Buses   []Bus   `json:"buses"`
}

type Bus struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	CompanyID     uint      `json:"company_id" gorm:"not null"`
	DepotID       uint      `json:"depot_id" gorm:"not null"`
	Name          string    `json:"name" gorm:"not null"`
	Type          string    `json:"type" gorm:"default:normal"` // normal, high_decker, super_high_decker, etc.
	Capacity      int       `json:"capacity" gorm:"default:40"`
	FuelCapacity  float64   `json:"fuel_capacity" gorm:"default:100"` // liters
	CurrentFuel   float64   `json:"current_fuel" gorm:"default:100"`
	Range         float64   `json:"range" gorm:"default:500"`            // km
	ServiceType   string    `json:"service_type" gorm:"default:economy"` // economy, business, executive, night
	Status        string    `json:"status" gorm:"default:available"`     // available, on_trip, maintenance
	Condition     float64   `json:"condition" gorm:"default:100"`        // percentage
	PurchasePrice float64   `json:"purchase_price" gorm:"default:0"`
	OperatingCost float64   `json:"operating_cost" gorm:"default:0"` // per km
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`

	// Relations
	Company  Company      `json:"company" gorm:"foreignKey:CompanyID"`
	Depot    Depot        `json:"depot" gorm:"foreignKey:DepotID"`
	Trips    []Trip       `json:"trips"`
	Upgrades []BusUpgrade `json:"upgrades"`
}

type Route struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"not null"`
	Origin      string    `json:"origin" gorm:"not null"`
	Destination string    `json:"destination" gorm:"not null"`
	OriginLat   float64   `json:"origin_lat" gorm:"not null"`
	OriginLng   float64   `json:"origin_lng" gorm:"not null"`
	DestLat     float64   `json:"dest_lat" gorm:"not null"`
	DestLng     float64   `json:"dest_lng" gorm:"not null"`
	Distance    float64   `json:"distance" gorm:"not null"`      // km
	Duration    int       `json:"duration" gorm:"not null"`      // minutes
	Popularity  int       `json:"popularity" gorm:"default:50"`  // 1-100
	Type        string    `json:"type" gorm:"default:intercity"` // intercity, interprovince
	MinBusType  string    `json:"min_bus_type" gorm:"default:normal"`
	BaseFare    float64   `json:"base_fare" gorm:"not null"` // IDR
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relations
	Trips []Trip `json:"trips"`
}

type Trip struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	BusID       uint      `json:"bus_id" gorm:"not null"`
	RouteID     uint      `json:"route_id" gorm:"not null"`
	DriverID    uint      `json:"driver_id"`
	Status      string    `json:"status" gorm:"default:planned"` // planned, active, completed, cancelled
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	ActualStart time.Time `json:"actual_start"`
	ActualEnd   time.Time `json:"actual_end"`
	Passengers  int       `json:"passengers" gorm:"default:0"`
	Revenue     float64   `json:"revenue" gorm:"default:0"`
	Cost        float64   `json:"cost" gorm:"default:0"`
	Profit      float64   `json:"profit" gorm:"default:0"`
	CurrentLat  float64   `json:"current_lat"`
	CurrentLng  float64   `json:"current_lng"`
	Progress    float64   `json:"progress" gorm:"default:0"` // percentage 0-100
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relations
	Bus    Bus    `json:"bus" gorm:"foreignKey:BusID"`
	Route  Route  `json:"route" gorm:"foreignKey:RouteID"`
	Driver Driver `json:"driver" gorm:"foreignKey:DriverID"`
}

type Driver struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	CompanyID   uint      `json:"company_id" gorm:"not null"`
	Name        string    `json:"name" gorm:"not null"`
	Age         int       `json:"age" gorm:"not null"`
	Experience  int       `json:"experience" gorm:"default:0"`     // years
	Skill       int       `json:"skill" gorm:"default:50"`         // 1-100
	Salary      float64   `json:"salary" gorm:"default:2000000"`   // IDR per month
	Status      string    `json:"status" gorm:"default:available"` // available, driving, rest
	Energy      float64   `json:"energy" gorm:"default:100"`       // percentage
	LicenseType string    `json:"license_type" gorm:"default:B"`   // B, B1, B2
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relations
	Company Company `json:"company" gorm:"foreignKey:CompanyID"`
	Trips   []Trip  `json:"trips"`
}

type BusUpgrade struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	BusID       uint      `json:"bus_id" gorm:"not null"`
	Type        string    `json:"type" gorm:"not null"` // engine, bathroom, multimedia, etc.
	Name        string    `json:"name" gorm:"not null"`
	Description string    `json:"description"`
	Cost        float64   `json:"cost" gorm:"not null"`
	Benefit     string    `json:"benefit"` // JSON string for various benefits
	CreatedAt   time.Time `json:"created_at"`

	// Relations
	Bus Bus `json:"bus" gorm:"foreignKey:BusID"`
}

type Transaction struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	CompanyID   uint      `json:"company_id" gorm:"not null"`
	Type        string    `json:"type" gorm:"not null"` // income, expense, purchase, sale
	Description string    `json:"description" gorm:"not null"`
	Amount      float64   `json:"amount" gorm:"not null"`
	Balance     float64   `json:"balance" gorm:"not null"` // Company balance after transaction
	CreatedAt   time.Time `json:"created_at"`

	// Relations
	Company Company `json:"company" gorm:"foreignKey:CompanyID"`
}
