package participant

import (
	"D/Go/messenger/internal/chat/domain"
	"D/Go/messenger/internal/chat/repository"
	"D/Go/messenger/internal/chat/service"
	"context"
	"errors"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const duplicateKeyCode = "23505"
const fkViolationKeyCode = "23503"

type Repository struct {
	q repository.Queryer
}

func New(q repository.Queryer) *Repository {
	return &Repository{q: q}
}

func (r *Repository) WithTx(tx pgx.Tx) service.ParticipantRepository {
	return &Repository{
		q: tx,
	}
}

func (r *Repository) Add(ctx context.Context, p domain.Participant) error {
	participant := ParticipantDB{
		UserId: p.UserId,
		ChatId: p.ChatId,
		Role:   int(p.Role),
	}
	query, args, err := squirrel.
		Insert("chat_participants").
		Columns("chat_id", "user_id", "role").
		Values(participant.ChatId, participant.UserId, participant.Role).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return domain.ErrDatabase
	}

	_, err = r.q.Exec(ctx, query, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == duplicateKeyCode {
			return domain.ErrParticipantAlreadyExists
		}
		if errors.As(err, &pgErr) && pgErr.Code == fkViolationKeyCode {
			return domain.ErrUserNotFound
		}
		return domain.ErrDatabase
	}

	return nil
}

func (r *Repository) Remove(ctx context.Context, p domain.Participant) error {
	participant := ParticipantDB{
		UserId: p.UserId,
		ChatId: p.ChatId,
	}
	query, args, err := squirrel.
		Delete("chat_participants").
		Where(squirrel.And{
			squirrel.Eq{"chat_id": participant.ChatId},
			squirrel.Eq{"user_id": participant.UserId},
		}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return domain.ErrDatabase
	}

	result, err := r.q.Exec(ctx, query, args...)
	if err != nil {
		return domain.ErrDatabase
	}
	if result.RowsAffected() == 0 {
		return domain.ErrChatOrParticipantNotFound
	}

	return nil
}

func (r *Repository) ChangeParticipantRole(ctx context.Context, p domain.Participant) error {
	participant := ParticipantDB{
		UserId: p.UserId,
		ChatId: p.ChatId,
		Role:   int(p.Role),
	}

	query, args, err := squirrel.
		Update("chat_participants").
		Set("role", participant.Role).
		Where(squirrel.And{
			squirrel.Eq{"chat_id": participant.ChatId},
			squirrel.Eq{"user_id": participant.UserId},
		}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return domain.ErrDatabase
	}

	result, err := r.q.Exec(ctx, query, args...)
	if err != nil {
		return domain.ErrDatabase
	}
	if result.RowsAffected() == 0 {
		return domain.ErrChatOrParticipantNotFound
	}

	return nil
}

func (r *Repository) GetParticipants(ctx context.Context, chatId int64) ([]domain.Participant, error) {

	query, args, err := squirrel.
		Select("chat_id", "user_id", "role").
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

	participants := make([]domain.Participant, 0)

	for rows.Next() {
		var participant domain.Participant
		err := rows.Scan(
			&participant.ChatId,
			&participant.UserId,
			&participant.Role,
		)
		if err != nil {
			return nil, domain.ErrDatabase
		}

		participants = append(participants, participant)
	}

	return participants, nil
}

func (r *Repository) GetParticipant(ctx context.Context, chatId int64, userId int64) (*domain.Participant, error) {
	var participant ParticipantDB

	query, args, err := squirrel.
		Select("chat_id", "user_id", "role").
		From("chat_participants").
		Where(squirrel.And{
			squirrel.Eq{"chat_id": chatId},
			squirrel.Eq{"user_id": userId},
		}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return nil, domain.ErrDatabase
	}

	err = r.q.QueryRow(ctx, query, args...).Scan(
		&participant.ChatId,
		&participant.UserId,
		&participant.Role,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrChatOrParticipantNotFound
		}
		return nil, domain.ErrDatabase
	}

	return &domain.Participant{
		ChatId: participant.ChatId,
		UserId: participant.UserId,
		Role:   domain.ParticipantRole(participant.Role),
	}, nil
}
