package entity

import (
	"time"
)

type Yard struct {
	ID        int       `db:"id" json:"id"`
	YardCode  string    `db:"yard_code" json:"yard_code"`
	Name      string    `db:"name" json:"name"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type Block struct {
	ID        int       `db:"id" json:"id"`
	YardID    int       `db:"yard_id" json:"yard_id"`
	BlockCode string    `db:"block_code" json:"block_code"`
	TotalSlot int       `db:"total_slot" json:"total_slot"`
	TotalRow  int       `db:"total_row" json:"total_row"`
	TotalTier int       `db:"total_tier" json:"total_tier"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type YardPlan struct {
	ID              int       `db:"id" json:"id"`
	BlockID         int       `db:"block_id" json:"block_id"`
	SlotStart       int       `db:"slot_start" json:"slot_start"`
	SlotEnd         int       `db:"slot_end" json:"slot_end"`
	RowStart        int       `db:"row_start" json:"row_start"`
	RowEnd          int       `db:"row_end" json:"row_end"`
	ContainerSize   int       `db:"container_size" json:"container_size"`
	ContainerHeight float64   `db:"container_height" json:"container_height"`
	ContainerType   string    `db:"container_type" json:"container_type"`
	CreatedAt       time.Time `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time `db:"updated_at" json:"updated_at"`
}

type ContainerPosition struct {
	ID              int        `db:"id" json:"id"`
	ContainerNumber string     `db:"container_number" json:"container_number"`
	BlockID         int        `db:"block_id" json:"block_id"`
	Slot            int        `db:"slot" json:"slot"`
	Slot2           int        `json:"slot2,omitempty"` // for 40ft containers
	Row             int        `db:"row" json:"row"`
	Tier            int        `db:"tier" json:"tier"`
	PlacedAt        time.Time  `db:"placed_at" json:"placed_at"`
	RemovedAt       *time.Time `db:"removed_at" json:"removed_at"`
}

type SuggestionRequest struct {
	Yard            string  `json:"yard" validate:"required"`
	ContainerNumber string  `json:"container_number" validate:"required"`
	ContainerSize   int     `json:"container_size" validate:"required,oneof=20 40"`
	ContainerHeight float64 `json:"container_height" validate:"required"`
	ContainerType   string  `json:"container_type" validate:"required"`
}

type SuggestedPosition struct {
	Block string `json:"block"`
	Slot  int    `json:"slot"`
	Slot2 int    `json:"slot2,omitempty"` // for 40ft containers
	Row   int    `json:"row"`
	Tier  int    `json:"tier"`
}

type PlacementDTO struct {
	Yard            string `json:"yard" validate:"required"`
	ContainerNumber string `json:"container_number" validate:"required"`
	Block           string `json:"block" validate:"required"`
	Slot            int    `json:"slot" validate:"required"`
	Row             int    `json:"row" validate:"required"`
	Tier            int    `json:"tier" validate:"required"`
}

type PickupDTO struct {
	Yard            string `json:"yard"`
	ContainerNumber string `json:"container_number" validate:"required"`
}

type SuggestRequestDTO struct {
	Yard            string  `json:"yard" validate:"required"`
	ContainerNumber string  `json:"container_number" validate:"required"`
	ContainerSize   int     `json:"container_size" validate:"required,oneof=20 40"`
	ContainerHeight float64 `json:"container_height" validate:"required"`
	ContainerType   string  `json:"container_type" validate:"required"`
}

type SuggestResponseDTO struct {
	SuggestedPosition SuggestedPosition `json:"suggested_position"`
}
