package dto

import "D/Go/messenger/internal/chat/domain"

func ToDomainPrivateChat(r CreatePrivateChatRequest, creatorId int64) domain.PrivateChat {
	return domain.PrivateChat{
		FirstUId:  creatorId,
		SecondUId: r.ParticipantId,
	}
}

func ToDomainGroupChat(r CreateGroupChatRequest, ownerId int64) domain.GroupChat {
	return domain.GroupChat{
		OwnerId: ownerId,
		Title:   r.Title,
		Users:   r.Users,
	}
}

func ToDomainParticipant(r ParticipantRequest) (*domain.Participant, error) {
	var role domain.ParticipantRole
	switch r.Role {
	case "member":
		role = domain.ParticipantRoleMember
	case "admin":
		role = domain.ParticipantRoleAdmin
	default:
		return nil, domain.ErrInvalidRole
	}

	return &domain.Participant{
		ChatId: r.ChatId,
		UserId: r.UserId,
		Role:   role,
	}, nil
}

func ToDomainRemoveParticipant(r RemovePartRequest) domain.Participant {
	return domain.Participant{
		ChatId: r.ChatId,
		UserId: r.UserId,
	}
}

func ToChatResponse(c domain.Chat) ChatResponse {
	var chatType string
	if c.Type == domain.ChatTypePrivate {
		chatType = "private"
	} else {
		chatType = "group"
	}
	return ChatResponse{
		Id:          c.Id,
		Type:        chatType,
		Title:       c.Title,
		OwnerId:     c.OwnerId,
		LastMsgText: c.LastMsgText,
		LastMsgTime: c.LastMsgTime,
		UnreadCount: c.UnreadCount,
	}
}

func ToChatListResponse(chats []domain.Chat, cursor *domain.Cursor) ChatListResponse {
	resp := ChatListResponse{
		Chats: make([]ChatResponse, 0, len(chats)),
	}

	for _, c := range chats {
		resp.Chats = append(resp.Chats, ToChatResponse(c))
	}

	if cursor != nil {
		encoded := EncodeCursor(cursor)
		resp.NextCursor = &encoded
	}

	return resp
}

func ToParticipantsResponse(p []domain.Participant) []ParticipantResponse {
	res := make([]ParticipantResponse, 0, len(p))

	for _, v := range p {
		var role string
		if v.Role == domain.ParticipantRoleAdmin {
			role = "admin"
		} else {
			role = "member"
		}
		res = append(res, ParticipantResponse{
			ChatId: v.ChatId,
			UserId: v.UserId,
			Role:   role,
		})
	}

	return res
}
