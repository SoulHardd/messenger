package participant

import (
	"D/Go/messenger/internal/message/domain"
	"D/Go/messenger/internal/message/repository"
	"D/Go/messenger/internal/message/service"
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	q repository.Queryer
}

func New(p *pgxpool.Pool) *Repository {
	return &Repository{q: p}
}

func (r *Repository) WithTx(tx pgx.Tx) service.ParticipantRepository {
	return &Repository{
		q: tx,
	}
}

func (r *Repository) GetParticipants(ctx context.Context, chatId int64) ([]int64, error) {
	query, args, err := squirrel.
		Select("user_id").
		From("chat_participants").
		Where(squirrel.Eq{"chat_id": chatId}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return nil, domain.ErrDatabase
	}

	rows, err := r.q.Query(ctx, query, args...)
	if err != nil {
		return nil, domain.ErrDatabase
	}
	defer rows.Close()

	participants := make([]int64, 0)

	for rows.Next() {
		var id int64
		err := rows.Scan(&id)
		if err != nil {
			return nil, domain.ErrDatabase
		}
		participants = append(participants, id)
	}

	if err = rows.Err(); err != nil {
		return nil, domain.ErrDatabase
	}

	return participants, nil
}

func (r *Repository) IsParticipant(ctx context.Context, userId int64, chatId int64) (bool, error) {
	query, args, err := squirrel.
		Select("COUNT(*)").
		From("chat_participants").
		Where(squirrel.And{
			squirrel.Eq{"user_id": userId},
			squirrel.Eq{"chat_id": chatId},
		}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return false, domain.ErrDatabase
	}

	var count int
	err = r.q.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return false, domain.ErrDatabase
	}

	return count > 0, nil
}
