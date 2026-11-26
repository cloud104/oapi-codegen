package examples

import (
	"time"

	openapi_types "github.com/oapi-codegen/runtime/types"
)

const (
	Charging     ChargingSessionStatus = "charging"
	Completed    ChargingSessionStatus = "completed"
	Failed       ChargingSessionStatus = "failed"
	Initializing ChargingSessionStatus = "initializing"
	Stopping     ChargingSessionStatus = "stopping"
	Suspended    ChargingSessionStatus = "suspended"
)

func (e ChargingSessionStatus) Valid() bool {
	switch e {
	case Charging:
		return true
	case Completed:
		return true
	case Failed:
		return true
	case Initializing:
		return true
	case Stopping:
		return true
	case Suspended:
		return true
	default:
		return false
	}
}

const (
	Corporate DriverMembership = "corporate"
	Premium   DriverMembership = "premium"
	Standard  DriverMembership = "standard"
)

func (e DriverMembership) Valid() bool {
	switch e {
	case Corporate:
		return true
	case Premium:
		return true
	case Standard:
		return true
	default:
		return false
	}
}

const (
	Draft   InvoiceStatus = "draft"
	Issued  InvoiceStatus = "issued"
	Overdue InvoiceStatus = "overdue"
	Paid    InvoiceStatus = "paid"
	Void    InvoiceStatus = "void"
)

func (e InvoiceStatus) Valid() bool {
	switch e {
	case Draft:
		return true
	case Issued:
		return true
	case Overdue:
		return true
	case Paid:
		return true
	case Void:
		return true
	default:
		return false
	}
}

const (
	Active      VehicleStatus = "active"
	Inactive    VehicleStatus = "inactive"
	Maintenance VehicleStatus = "maintenance"
)

func (e VehicleStatus) Valid() bool {
	switch e {
	case Active:
		return true
	case Inactive:
		return true
	case Maintenance:
		return true
	default:
		return false
	}
}

const (
	Operational GetHealth200JSONResponseBodyStatus = "operational"
)

func (e GetHealth200JSONResponseBodyStatus) Valid() bool {
	switch e {
	case Operational:
		return true
	default:
		return false
	}
}

type ChargingSession struct {
	BatteryPercentage  *int                  `json:"batteryPercentage,omitempty"`
	ConnectorID        string                `json:"connectorId"`
	CostMinor          int                   `json:"costMinor"`
	Currency           Currency              `json:"currency"`
	CurrentPowerKw     *float32              `json:"currentPowerKw,omitempty"`
	EndedAt            *time.Time            `json:"endedAt"`
	EnergyDeliveredKwh float32               `json:"energyDeliveredKwh"`
	ID                 string                `json:"id"`
	StartedAt          time.Time             `json:"startedAt"`
	StationID          string                `json:"stationId"`
	Status             ChargingSessionStatus `json:"status"`
	VehicleID          string                `json:"vehicleId"`
}

type ChargingSessionStatus string

type Currency = string

type Driver struct {
	Currency           *Currency           `json:"currency,omitempty"`
	Email              openapi_types.Email `json:"email"`
	FullName           string              `json:"fullName"`
	ID                 string              `json:"id"`
	Membership         DriverMembership    `json:"membership"`
	WalletBalanceMinor *int                `json:"walletBalanceMinor,omitempty"`
}

type DriverMembership string

type Invoice struct {
	BillingPeriod     string              `json:"billingPeriod"`
	ChargingSessions  *int                `json:"chargingSessions,omitempty"`
	Currency          Currency            `json:"currency"`
	DownloadURL       *string             `json:"downloadUrl,omitempty"`
	DueDate           *openapi_types.Date `json:"dueDate,omitempty"`
	EnergyConsumedKwh *float32            `json:"energyConsumedKwh,omitempty"`
	ID                string              `json:"id"`
	Status            InvoiceStatus       `json:"status"`
	TotalMinor        int                 `json:"totalMinor"`
}

type InvoiceStatus string

type Problem struct {
	Code     *string `json:"code,omitempty"`
	Detail   *string `json:"detail,omitempty"`
	Instance *string `json:"instance,omitempty"`
	Status   int     `json:"status"`
	Title    string  `json:"title"`
	TraceID  *string `json:"traceId,omitempty"`
	Type     string  `json:"type"`
}

type StartSessionRequest struct {
	ConnectorID string  `json:"connectorId"`
	CostCenter  *string `json:"costCenter,omitempty"`
	FleetID     *string `json:"fleetId,omitempty"`
	StationID   string  `json:"stationId"`
	VehicleID   string  `json:"vehicleId"`
}

type Station struct {
	AvailableConnectors int      `json:"availableConnectors"`
	Currency            Currency `json:"currency"`
	ID                  string   `json:"id"`
	Latitude            float32  `json:"latitude"`
	Longitude           float32  `json:"longitude"`
	MaxPowerKw          float32  `json:"maxPowerKw"`
	Name                string   `json:"name"`
	PricePerKwhMinor    int      `json:"pricePerKwhMinor"`
	TotalConnectors     int      `json:"totalConnectors"`
}

type Vehicle struct {
	BatteryCapacityKwh *float32       `json:"batteryCapacityKwh,omitempty"`
	ID                 string         `json:"id"`
	LicensePlate       string         `json:"licensePlate"`
	Manufacturer       string         `json:"manufacturer"`
	Model              string         `json:"model"`
	Status             *VehicleStatus `json:"status,omitempty"`
}

type VehicleStatus string

type FleetID = string

type IdempotencyKey = openapi_types.UUID

type SessionID = string

type BadRequest = Problem

type Forbidden = Problem

type NotFound = Problem

type RateLimited = Problem

type Unauthorized = Problem

type StartChargingSessionParams struct {
	IdempotencyKey IdempotencyKey `json:"Idempotency-Key"`
}

type StopChargingSessionParams struct {
	IdempotencyKey IdempotencyKey `json:"Idempotency-Key"`
}

type ListFleetInvoicesParams struct {
	BillingPeriod string `form:"billingPeriod" json:"billingPeriod"`
}

type ListFleetVehiclesParams struct {
	Limit  *int    `form:"limit,omitempty" json:"limit,omitempty"`
	Cursor *string `form:"cursor,omitempty" json:"cursor,omitempty"`
}

type GetHealth200JSONResponseBodyStatus string

type SearchChargingStationsParams struct {
	Latitude      float32  `form:"latitude" json:"latitude"`
	Longitude     float32  `form:"longitude" json:"longitude"`
	RadiusKm      *float32 `form:"radiusKm,omitempty" json:"radiusKm,omitempty"`
	AvailableOnly *bool    `form:"availableOnly,omitempty" json:"availableOnly,omitempty"`
}

type StartChargingSessionJSONRequestBody = StartSessionRequest
