package service

import (
	"D/Go/messenger/internal/chat/domain"
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ChatService struct {
	ChatRepo ChatRepository
	PartRepo ParticipantRepository
	pool     *pgxpool.Pool
}

func New(cr ChatRepository, pr ParticipantRepository, p *pgxpool.Pool) *ChatService {
	return &ChatService{
		ChatRepo: cr,
		PartRepo: pr,
		pool:     p,
	}
}

func (s *ChatService) CreatePrivateChat(ctx context.Context, chat domain.PrivateChat) (int64, error) {
	if chat.FirstUId == 0 || chat.SecondUId == 0 {
		return 0, domain.ErrInvalidInput
	}
	if chat.FirstUId == chat.SecondUId {
		return 0, domain.ErrInvalidInput
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return 0, domain.ErrDatabase
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	ChatRepoTx := s.ChatRepo.WithTx(tx)
	PartRepoTx := s.PartRepo.WithTx(tx)

	chExists, err := ChatRepoTx.FindPrivateChat(ctx, chat.FirstUId, chat.SecondUId)
	if err != nil {
		return 0, err
	}
	if chExists == true {
		return 0, domain.ErrPrivateChatExists
	}
	ch := &domain.Chat{
		Type:    domain.ChatTypePrivate,
		Title:   "private",
		OwnerId: chat.FirstUId,
	}

	chatId, err := ChatRepoTx.Create(ctx, ch)
	if err != nil {
		return 0, err
	}

	partA := domain.Participant{
		ChatId: chatId,
		UserId: chat.FirstUId,
		Role:   domain.ParticipantRoleAdmin,
	}
	partB := domain.Participant{
		ChatId: chatId,
		UserId: chat.SecondUId,
		Role:   domain.ParticipantRoleAdmin,
	}

	err = PartRepoTx.Add(ctx, partA)
	if err != nil {
		return 0, err
	}
	err = PartRepoTx.Add(ctx, partB)
	if err != nil {
		return 0, err
	}

	err = ChatRepoTx.CreateState(ctx, chatId, chat.FirstUId)
	if err != nil {
		return 0, domain.ErrDatabase
	}
	err = ChatRepoTx.CreateState(ctx, chatId, chat.SecondUId)
	if err != nil {
		return 0, domain.ErrDatabase
	}

	err = tx.Commit(ctx)
	if err != nil {
		return 0, domain.ErrDatabase
	}

	return chatId, nil
}

func (s *ChatService) CreateGroupChat(ctx context.Context, chat domain.GroupChat) (int64, error) {
	if chat.OwnerId == 0 {
		return 0, domain.ErrInvalidInput
	}
	if len(chat.Users) == 0 {
		return 0, domain.ErrInvalidInput
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return 0, domain.ErrDatabase
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	ChatRepoTx := s.ChatRepo.WithTx(tx)
	PartRepoTx := s.PartRepo.WithTx(tx)

	ch := &domain.Chat{
		Type:    domain.ChatTypeGroup,
		Title:   chat.Title,
		OwnerId: chat.OwnerId,
	}

	chatId, err := ChatRepoTx.Create(ctx, ch)
	if err != nil {
		return 0, err
	}
	owner := domain.Participant{
		ChatId: chatId,
		UserId: chat.OwnerId,
		Role:   domain.ParticipantRoleAdmin,
	}
	err = PartRepoTx.Add(ctx, owner)
	if err != nil {
		return 0, err
	}

	err = ChatRepoTx.CreateState(ctx, chatId, owner.UserId)
	if err != nil {
		return 0, domain.ErrDatabase
	}

	seen := map[int64]struct{}{}
	for _, uid := range chat.Users {
		if uid == chat.OwnerId {
			continue
		}
		if _, ok := seen[uid]; ok {
			continue
		}
		seen[uid] = struct{}{}

		p := domain.Participant{
			ChatId: chatId,
			UserId: uid,
			Role:   domain.ParticipantRoleMember,
		}
		err := PartRepoTx.Add(ctx, p)
		if err != nil {
			return 0, err
		}
		err = ChatRepoTx.CreateState(ctx, p.ChatId, p.UserId)
		if err != nil {
			return 0, domain.ErrDatabase
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return 0, domain.ErrDatabase
	}

	return chatId, nil
}

func (s *ChatService) GetUserChats(ctx context.Context, userId int64, limit int, cursor *domain.Cursor) ([]domain.Chat, *domain.Cursor, error) {
	chats, err := s.ChatRepo.GetUserChats(ctx, userId, limit, cursor)
	if err != nil {
		return nil, nil, err
	}

	hasMore := len(chats) > limit
	if hasMore {
		chats = chats[:limit]
	}

	var nextCursor *domain.Cursor

	if hasMore && len(chats) > 0 {
		last := chats[len(chats)-1]

		nextCursor = &domain.Cursor{
			Time: last.LastMsgTime,
			Id:   last.Id,
		}
	}

	return chats, nextCursor, nil
}

func (s *ChatService) GetUserChatById(ctx context.Context, userId int64, chatId int64) (*domain.Chat, error) {
	chat, err := s.ChatRepo.GetUserChat(ctx, userId, chatId)
	if err != nil {
		return nil, err
	}
	return chat, nil
}

func (s *ChatService) GetChatParticipants(ctx context.Context, userId int64, chatId int64) ([]domain.Participant, error) {
	_, err := s.PartRepo.GetParticipant(ctx, userId, chatId)
	if err != nil {
		if errors.Is(err, domain.ErrChatOrParticipantNotFound) {
			return nil, domain.ErrPermissionDenied
		}
		return nil, err
	}
	return s.PartRepo.GetParticipants(ctx, chatId)
}

func (s *ChatService) AddParticipant(ctx context.Context, userId int64, p domain.Participant) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return domain.ErrDatabase
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	PartRepoTx := s.PartRepo.WithTx(tx)
	ChatRepoTx := s.ChatRepo.WithTx(tx)

	changer, err := PartRepoTx.GetParticipant(ctx, p.ChatId, userId)
	if err != nil {
		if errors.Is(err, domain.ErrChatOrParticipantNotFound) {
			return domain.ErrPermissionDenied
		}
		return err
	}
	if changer.Role != domain.ParticipantRoleAdmin {
		return domain.ErrPermissionDenied
	}

	chat, err := ChatRepoTx.GetUserChat(ctx, userId, p.ChatId)
	if err != nil {
		return err
	}
	if chat.Type == domain.ChatTypePrivate {
		return domain.ErrPermissionDenied
	}

	err = PartRepoTx.Add(ctx, p)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return domain.ErrDatabase
	}

	return nil
}

func (s *ChatService) RemoveParticipant(ctx context.Context, userId int64, p domain.Participant) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return domain.ErrDatabase
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	PartRepoTx := s.PartRepo.WithTx(tx)
	ChatRepoTx := s.ChatRepo.WithTx(tx)

	changer, err := PartRepoTx.GetParticipant(ctx, p.ChatId, userId)
	if err != nil {
		if errors.Is(err, domain.ErrParticipantNotFound) {
			return domain.ErrPermissionDenied
		}
		return err
	}
	if changer.Role != domain.ParticipantRoleAdmin {
		return domain.ErrPermissionDenied
	}

	chat, err := ChatRepoTx.GetUserChat(ctx, p.UserId, p.ChatId)
	if err != nil {
		return err
	}
	if chat.OwnerId == p.UserId || chat.Type == domain.ChatTypePrivate {
		return domain.ErrPermissionDenied
	}

	err = PartRepoTx.Remove(ctx, p)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return domain.ErrDatabase
	}

	return nil
}

func (s *ChatService) ChangeParticipantRole(ctx context.Context, userId int64, p domain.Participant) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return domain.ErrDatabase
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	PartRepoTx := s.PartRepo.WithTx(tx)
	ChatRepoTx := s.ChatRepo.WithTx(tx)

	changer, err := PartRepoTx.GetParticipant(ctx, p.ChatId, userId)
	if err != nil {
		if errors.Is(err, domain.ErrParticipantNotFound) {
			return domain.ErrPermissionDenied
		}
		return err
	}
	if changer.Role != domain.ParticipantRoleAdmin {
		return domain.ErrPermissionDenied
	}

	chat, err := ChatRepoTx.GetUserChat(ctx, p.UserId, p.ChatId)
	if err != nil {
		return err
	}
	if chat.OwnerId == p.UserId {
		return domain.ErrPermissionDenied
	}

	err = PartRepoTx.ChangeParticipantRole(ctx, p)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return domain.ErrDatabase
	}

	return nil
}
