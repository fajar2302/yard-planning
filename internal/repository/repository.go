package repository

import (
	"context"
	"database/sql"
	"errors"
	"yard-planning/internal/entity"
)

type Repo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) GetPlansBySpec(ctx context.Context, yardCode string, size int, height float64, ctype string) ([]entity.YardPlan, error) {
	query := `
        SELECT yp.id, yp.block_id, yp.slot_start, yp.slot_end,
               yp.row_start, yp.row_end, yp.container_size,
               yp.container_height, yp.container_type,
               yp.created_at, yp.updated_at
        FROM yard_plans yp
        JOIN blocks b ON yp.block_id = b.id
        JOIN yards y ON b.yard_id = y.id
        WHERE y.yard_code = $1
          AND yp.container_size = $2
          AND yp.container_height = $3
          AND yp.container_type = $4
        ORDER BY yp.block_id, yp.slot_start;
    `

	rows, err := r.db.QueryContext(ctx, query, yardCode, size, height, ctype)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var plans []entity.YardPlan

	for rows.Next() {
		var p entity.YardPlan
		if err := rows.Scan(
			&p.ID, &p.BlockID, &p.SlotStart, &p.SlotEnd,
			&p.RowStart, &p.RowEnd, &p.ContainerSize,
			&p.ContainerHeight, &p.ContainerType,
			&p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, err
		}
		plans = append(plans, p)
	}

	return plans, nil
}

func (r *Repo) GetBlockByID(ctx context.Context, id int) (*entity.Block, error) {
	query := `
        SELECT id, yard_id, block_code, total_slot,
               total_row, total_tier, created_at, updated_at
        FROM blocks WHERE id = $1
    `
	var b entity.Block
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&b.ID, &b.YardID, &b.BlockCode, &b.TotalSlot,
		&b.TotalRow, &b.TotalTier, &b.CreatedAt, &b.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *Repo) GetBlockByCode(ctx context.Context, code string) (*entity.Block, error) {
	query := `
        SELECT id, yard_id, block_code, total_slot,
               total_row, total_tier, created_at, updated_at
        FROM blocks WHERE block_code = $1 LIMIT 1
    `
	var b entity.Block
	err := r.db.QueryRowContext(ctx, query, code).Scan(
		&b.ID, &b.YardID, &b.BlockCode, &b.TotalSlot,
		&b.TotalRow, &b.TotalTier, &b.CreatedAt, &b.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *Repo) IsPositionOccupied(ctx context.Context, blockID, slot, row, tier int) (bool, error) {
	var count int
	q := `
        SELECT COUNT(*) FROM container_positions
        WHERE block_id=$1 AND slot=$2 AND row=$3 AND tier=$4 AND removed_at IS NULL
    `
	err := r.db.QueryRowContext(ctx, q, blockID, slot, row, tier).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *Repo) ArePositionsOccupiedFor40ft(ctx context.Context, blockID, slot, row, tier int) (bool, error) {
	var count int
	q := `
        SELECT COUNT(*) FROM container_positions
        WHERE block_id=$1 AND row=$2 AND tier=$3
          AND (slot=$4 OR slot=$4+1)
          AND removed_at IS NULL
    `
	err := r.db.QueryRowContext(ctx, q, blockID, row, tier, slot).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *Repo) PlaceContainer(ctx context.Context, cp *entity.ContainerPosition) error {
	q := `
        INSERT INTO container_positions
        (container_number, block_id, slot, row, tier, placed_at)
        VALUES ($1, $2, $3, $4, $5, NOW())
    `
	_, err := r.db.ExecContext(ctx, q, cp.ContainerNumber, cp.BlockID, cp.Slot, cp.Row, cp.Tier)
	return err
}

func (r *Repo) FindContainer(ctx context.Context, number string) (*entity.ContainerPosition, error) {
	q := `
        SELECT id, container_number, block_id, slot, row, tier, placed_at, removed_at
        FROM container_positions WHERE container_number=$1 AND removed_at IS NULL
    `
	var cp entity.ContainerPosition
	err := r.db.QueryRowContext(ctx, q, number).Scan(
		&cp.ID, &cp.ContainerNumber, &cp.BlockID,
		&cp.Slot, &cp.Row, &cp.Tier,
		&cp.PlacedAt, &cp.RemovedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &cp, nil
}

func (r *Repo) PickupContainer(ctx context.Context, number string) error {
	q := `UPDATE container_positions SET removed_at=NOW() WHERE container_number=$1 AND removed_at IS NULL`
	res, err := r.db.ExecContext(ctx, q, number)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New("container not found")
	}
	return nil
}
