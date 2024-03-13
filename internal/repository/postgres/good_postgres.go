package postgres

import (
	"context"
	"fmt"
	"github.com/blazee5/hezzl-test/internal/domain"
	"github.com/blazee5/hezzl-test/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type GoodRepository struct {
	db *pgxpool.Pool
}

func NewGoodRepository(db *pgxpool.Pool) *GoodRepository {
	return &GoodRepository{db: db}
}

func (repo *GoodRepository) Get(ctx context.Context, limit, offset int) (domain.GoodList, error) {
	rows, err := repo.db.Query(ctx, fmt.Sprintf(`SELECT id, project_id, name, description, priority, removed, created_at FROM goods LIMIT %d OFFSET %d`, limit, offset))

	if err != nil {
		return domain.GoodList{}, err
	}

	Goods, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Good])

	if err != nil {
		return domain.GoodList{}, err
	}

	var removedCount int

	err = repo.db.QueryRow(ctx, `SELECT COUNT(id) FROM goods WHERE removed = true`).Scan(&removedCount)

	if err != nil {
		return domain.GoodList{}, err
	}

	return domain.GoodList{
		Meta: domain.Meta{
			Total:   len(Goods),
			Removed: removedCount,
			Limit:   limit,
			Offset:  offset,
		},
		Goods: Goods,
	}, nil
}

func (repo *GoodRepository) Create(ctx context.Context, projectID int, input domain.CreateGoodRequest) (models.Good, error) {
	rows, err := repo.db.Query(ctx, `INSERT INTO goods (project_id, name) VALUES ($1, $2)
		RETURNING id, project_id, name, description, priority, removed, created_at`, projectID, input.Name)

	if err != nil {
		return models.Good{}, err
	}

	Good, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.Good])

	if err != nil {
		return models.Good{}, err
	}

	return Good, nil
}

func (repo *GoodRepository) Update(ctx context.Context, projectID, id int, input domain.UpdateGoodRequest) (models.Good, error) {
	tx, err := repo.db.Begin(ctx)
	defer tx.Rollback(ctx)

	if err != nil {
		return models.Good{}, err
	}

	rows, err := tx.Query(ctx, `UPDATE goods SET name = $1, description = COALESCE(NULLIF($2, ''), description) WHERE id = $3 AND project_id = $4
		RETURNING id, project_id, name, description, priority, removed, created_at`, input.Name, input.Description, id, projectID)

	if err != nil {
		return models.Good{}, err
	}

	Good, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.Good])

	if err != nil {
		return models.Good{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return models.Good{}, err
	}

	return Good, nil
}

func (repo *GoodRepository) Delete(ctx context.Context, projectID, id int) (models.DeletedGood, error) {
	tx, err := repo.db.Begin(ctx)
	defer tx.Rollback(ctx)

	if err != nil {
		return models.DeletedGood{}, err
	}

	rows, err := tx.Query(ctx, `UPDATE goods SET removed = true WHERE id = $1 AND project_id = $2
		RETURNING id, project_id, removed`, id, projectID)

	if err != nil {
		return models.DeletedGood{}, err
	}

	Good, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.DeletedGood])

	if err != nil {
		return models.DeletedGood{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return models.DeletedGood{}, err
	}

	return Good, nil
}

func (repo *GoodRepository) Reprioritize(ctx context.Context, projectID, id int, input domain.ReprioritizeRequest) (models.GoodPriorities, error) {
	tx, err := repo.db.Begin(ctx)

	if err != nil {
		return models.GoodPriorities{}, err
	}

	defer tx.Rollback(ctx)

	var priority int

	err = tx.QueryRow(ctx, `SELECT priority FROM goods WHERE id = $1 AND project_id = $2`, id, projectID).Scan(&priority)

	if err != nil {
		return models.GoodPriorities{}, err
	}

	rows, err := tx.Query(ctx, `SELECT id, priority FROM goods WHERE priority >= $1 ORDER BY id ASC`, priority)

	if err != nil {
		return models.GoodPriorities{}, err
	}

	goods, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Priorities])

	if err != nil {
		return models.GoodPriorities{}, err
	}

	batch := &pgx.Batch{}
	for _, good := range goods {
		batch.Queue(`UPDATE goods SET priority = $1 WHERE id = $2 RETURNING id, priority`, input.NewPriority, good.ID)
		input.NewPriority++
	}

	br := tx.SendBatch(ctx, batch)
	defer br.Close()

	var updatedGoods models.GoodPriorities
	for range goods {
		var updatedGood models.Priorities
		err = br.QueryRow().Scan(&updatedGood.ID, &updatedGood.Priority)

		if err != nil {
			return models.GoodPriorities{}, err
		}

		updatedGoods.Priorities = append(updatedGoods.Priorities, updatedGood)
	}

	defer tx.Commit(ctx)

	return updatedGoods, nil
}
