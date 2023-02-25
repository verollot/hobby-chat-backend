package business

import (
	"context"

	"github.com/frisk038/livechat/business/models"
	"github.com/google/uuid"
)

type repo interface {
	InsertUser(ctx context.Context, userID, firstName, lastName string) error
	InsertHobbies(ctx context.Context, hobbies []string) error
	InsertUserHobbies(ctx context.Context, userID string, hobbies []string) error
	GetUserHobbies(ctx context.Context, userID string) ([]models.Hobby, error)
	DelUserHobbies(ctx context.Context, userID string, hobbyID uuid.UUID) error
}

type BusinessProfile struct {
	repo repo
}

func NewBusinessProfile(repo repo) BusinessProfile {
	return BusinessProfile{repo: repo}
}

func (bp *BusinessProfile) CreateUser(ctx context.Context, user models.User) error {
	err := bp.repo.InsertUser(context.Background(), user.ID, user.FirsttName, user.LastName)
	if err != nil {
		return err
	}
	return nil
}

func (bp *BusinessProfile) SetHobbies(ctx context.Context, userID, hobby string) error {
	err := bp.repo.InsertHobbies(ctx, []string{hobby})
	if err != nil {
		return err
	}

	return bp.repo.InsertUserHobbies(ctx, userID, []string{hobby})
}

func (bp *BusinessProfile) GetHobbies(ctx context.Context, userID string) ([]models.Hobby, error) {
	return bp.repo.GetUserHobbies(ctx, userID)
}

func (bp *BusinessProfile) DelHobbies(ctx context.Context, userID string, hobbyID uuid.UUID) error {
	return bp.repo.DelUserHobbies(ctx, userID, hobbyID)
}
