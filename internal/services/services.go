package services

import (
	"context"
	"errors"
	"fmt"

	"yard-planning/internal/entity"
	"yard-planning/internal/repository"
)

type Service struct {
	repo repository.RepoInterface
}

type ServiceInterface interface {
	SuggestPosition(ctx context.Context, req entity.SuggestionRequest) (*entity.SuggestedPosition, error)
	PlaceContainer(ctx context.Context, cp *entity.ContainerPosition) (string, error)
	PickupContainer(ctx context.Context, containerNumber string) (string, error)
	RepoGetBlockByCode(ctx context.Context, code string) (*entity.Block, error)
}

func NewService(r repository.RepoInterface) ServiceInterface {
	return &Service{repo: r}
}

func (s *Service) SuggestPosition(ctx context.Context, req entity.SuggestionRequest) (*entity.SuggestedPosition, error) {
	plans, err := s.repo.GetPlansBySpec(ctx, req.Yard, req.ContainerSize, req.ContainerHeight, req.ContainerType)
	if err != nil {
		return nil, err
	}
	if len(plans) == 0 {
		return nil, errors.New("no plan found for given spec")
	}

	for _, p := range plans {
		block, err := s.repo.GetBlockByID(ctx, p.BlockID)
		if err != nil {
			return nil, err
		}
		// iterate slots, rows, tiers (priority: slot ascending, row ascending, tier ascending)
		for slot := p.SlotStart; slot <= p.SlotEnd; slot++ {
			// If 40ft, ensure slot+1 <= block.TotalSlot and slot+1 <= p.SlotEnd
			if req.ContainerSize == 40 {
				if slot+1 > p.SlotEnd || slot+1 > block.TotalSlot {
					continue
				}
			}
			for row := p.RowStart; row <= p.RowEnd; row++ {
				for tier := 1; tier <= block.TotalTier; tier++ {
					// check occupancy
					if req.ContainerSize == 40 {
						occ, err := s.repo.ArePositionsOccupiedFor40ft(ctx, block.ID, slot, row, tier)
						if err != nil {
							return nil, err
						}
						if occ {
							continue
						}
						// free => suggest slot (slot) (we can return slot as the left slot)
						return &entity.SuggestedPosition{
							Block: block.BlockCode,
							Slot:  slot,
							Slot2: slot + 1,
							Row:   row,
							Tier:  tier,
						}, nil
					}
					// 20ft
					occ, err := s.repo.IsPositionOccupied(ctx, block.ID, slot, row, tier)
					if err != nil {
						return nil, err
					}
					if occ {
						continue
					}
					return &entity.SuggestedPosition{
						Block: block.BlockCode,
						Slot:  slot,
						Row:   row,
						Tier:  tier,
					}, nil
				}
			}
		}
	}
	return nil, errors.New("no free position available in plans")
}

func (s *Service) PlaceContainer(ctx context.Context, cp *entity.ContainerPosition) (string, error) {

	existing, _ := s.repo.FindContainer(ctx, cp.ContainerNumber)
	if existing != nil {
		return "", fmt.Errorf("container %s already placed", cp.ContainerNumber)
	}

	occ, err := s.repo.IsPositionOccupied(ctx, cp.BlockID, cp.Slot, cp.Row, cp.Tier)
	if err != nil {
		return "", err
	}
	if occ {
		return "", errors.New("target position already occupied")
	}

	if cp.Slot2 > 0 {
		// check occupancy slot2 for 40ft
		occ2, err := s.repo.IsPositionOccupied(ctx, cp.BlockID, cp.Slot2, cp.Row, cp.Tier)
		if err != nil {
			return "", err
		}
		if occ2 {
			return "", errors.New("target position already occupied")
		}
	}

	block, err := s.repo.GetBlockByID(ctx, cp.BlockID)
	if err != nil {
		return "", err
	}

	//check limit
	if cp.Tier > block.TotalTier || cp.Row > block.TotalRow || cp.Slot > block.TotalSlot {
		return "", errors.New("target position exceeds block limits")
	}
	if cp.Slot2 > 0 && cp.Slot2 > block.TotalSlot {
		return "", errors.New("target position exceeds block limits")
	}
	// place
	res, err := s.repo.PlaceContainer(ctx, cp)
	if err != nil {
		return "", err
	}
	return res, nil
}

func (s *Service) PickupContainer(ctx context.Context, containerNumber string) (string, error) {
	// check exist
	existing, err := s.repo.FindContainer(ctx, containerNumber)
	if err != nil {
		return "", err
	}
	if existing == nil {
		return "", errors.New("container not found or already removed")
	}
	return s.repo.PickupContainer(ctx, containerNumber)
}

func (s *Service) RepoGetBlockByCode(ctx context.Context, code string) (*entity.Block, error) {
	return s.repo.GetBlockByCode(ctx, code)
}
