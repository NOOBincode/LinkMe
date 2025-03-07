package repository

import (
	"context"

	"github.com/GoSimplicity/LinkMe/internal/domain"
	"github.com/GoSimplicity/LinkMe/internal/repository/dao"
)

type ActivityRepository interface {
	GetRecentActivity(ctx context.Context) ([]domain.RecentActivity, error)
	SetRecentActivity(ctx context.Context, dr domain.RecentActivity) error
}

type activityRepository struct {
	dao dao.ActivityDAO
}

func NewActivityRepository(dao dao.ActivityDAO) ActivityRepository {
	return &activityRepository{
		dao: dao,
	}
}

func (a *activityRepository) GetRecentActivity(ctx context.Context) ([]domain.RecentActivity, error) {
	activity, err := a.dao.GetRecentActivity(ctx)
	if err != nil {
		return nil, err
	}
	return toDomainActivities(activity), nil
}

func (a *activityRepository) SetRecentActivity(ctx context.Context, dr domain.RecentActivity) error {
	err := a.dao.SetRecentActivity(ctx, fromDomainActivity(dr))
	if err != nil {
		return err
	}
	return nil
}

// 将领域层对象转为dao层对象
func fromDomainActivity(dr domain.RecentActivity) dao.RecentActivity {
	return dao.RecentActivity{
		ID:          dr.ID,
		Description: dr.Description,
		Time:        dr.Time,
		UserID:      dr.UserID,
	}
}

// 将dao层对象转为领域层对象
func toDomainActivities(mrList []dao.RecentActivity) []domain.RecentActivity {
	domainList := make([]domain.RecentActivity, len(mrList))
	for i, mr := range mrList {
		domainList[i] = domain.RecentActivity{
			ID:          mr.ID,
			Description: mr.Description,
			Time:        mr.Time,
			UserID:      mr.UserID,
		}
	}
	return domainList
}
