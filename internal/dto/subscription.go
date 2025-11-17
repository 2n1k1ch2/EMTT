package dto

import (
	"EMTT/internal/repository"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func SubscriptionToDTO(s *repository.Subscription) SubscriptionDTO {
	dto := SubscriptionDTO{
		ServiceName: s.ServiceName,
		Price:       s.Price,
	}

	if s.ID.Valid {
		dto.ID = uuid.UUID(s.ID.Bytes).String()
	}

	if s.UserID.Valid {
		dto.UserID = uuid.UUID(s.UserID.Bytes).String()
	}

	if s.StartDate.Valid {
		dto.StartDate = s.StartDate.Time.Format("2006-01-02")
	}

	if s.FinishDate.Valid {
		dto.FinishDate = s.FinishDate.Time.Format("2006-01-02")
	}

	return dto
}

func DTOToSubscription(dto *SubscriptionDTO) (*repository.Subscription, error) {
	sub := &repository.Subscription{
		ServiceName: dto.ServiceName,
		Price:       dto.Price,
	}

	var err error

	if dto.ID != "" {
		idUUID, err := uuid.Parse(dto.ID)
		if err != nil {
			return nil, err
		}
		sub.ID = pgtype.UUID{Bytes: idUUID, Valid: true}
	} else {
		sub.ID = pgtype.UUID{Bytes: uuid.New(), Valid: true}
	}

	userUUID, err := uuid.Parse(dto.UserID)
	if err != nil {
		return nil, err
	}
	sub.UserID = pgtype.UUID{Bytes: userUUID, Valid: true}

	startTime, err := time.Parse("2006-01-02", dto.StartDate)
	if err != nil {
		return nil, err
	}
	sub.StartDate = pgtype.Timestamptz{Time: startTime, Valid: true}

	if dto.FinishDate != "" {
		finishTime, err := time.Parse("2006-01-02", dto.FinishDate)
		if err != nil {
			return nil, err
		}
		sub.FinishDate = pgtype.Timestamptz{Time: finishTime, Valid: true}
	} else {
		sub.FinishDate = pgtype.Timestamptz{Valid: false}
	}

	return sub, nil
}

func SubscriptionsToDTO(list []*repository.Subscription) []SubscriptionDTO {
	result := make([]SubscriptionDTO, 0, len(list))
	for _, s := range list {
		result = append(result, SubscriptionToDTO(s))
	}
	return result
}
func StringToUUID(s string) (pgtype.UUID, error) {
	if s == "" {
		return pgtype.UUID{Valid: false}, nil
	}
	u, err := uuid.Parse(s)
	if err != nil {
		return pgtype.UUID{}, err
	}
	return pgtype.UUID{Bytes: u, Valid: true}, nil
}

type SubscriptionDTO struct {
	ID          string `json:"id,omitempty"`
	ServiceName string `json:"service_name"`
	Price       int    `json:"price"`
	UserID      string `json:"user_id"`
	StartDate   string `json:"start_date"`
	FinishDate  string `json:"finish_date,omitempty"`
}
