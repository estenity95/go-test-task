package subscription

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	List(ctx context.Context, limit, offset int) (Subscriptions, error)
	Create(ctx context.Context, s *Subscription) (*Subscription, error)
	Get(ctx context.Context, id int64) (*Subscription, error)
	Update(ctx context.Context, s *Subscription) (int64, error)
	Delete(ctx context.Context, id int64) error
	Summary(ctx context.Context, q SummaryQuery) (int64, error)
}

type Repo struct {
	db *gorm.DB
}

type SummaryQuery struct {
	From        time.Time
	To          time.Time
	UserID      *uuid.UUID
	ServiceName *string
}

func NewRepo(db *gorm.DB) *Repo {
	return &Repo{
		db: db,
	}
}

func (r *Repo) List(ctx context.Context, limit, offset int) (Subscriptions, error) {
	subscriptions := make([]*Subscription, 0, limit)

	if err := r.db.
		WithContext(ctx).
		Order("id ASC").
		Limit(limit).
		Offset(offset).
		Find(&subscriptions).Error; err != nil {
		return nil, err
	}

	return subscriptions, nil
}

func (r *Repo) Create(ctx context.Context, s *Subscription) (*Subscription, error) {
	if err := r.db.WithContext(ctx).Create(s).Error; err != nil {
		return nil, err
	}

	return s, nil
}

func (r *Repo) Get(ctx context.Context, id int64) (*Subscription, error) {
	subscription := &Subscription{}
	if err := r.db.WithContext(ctx).First(subscription, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return subscription, nil
}

func (r *Repo) Update(ctx context.Context, s *Subscription) (int64, error) {
	res := r.db.WithContext(ctx).Model(&Subscription{}).
		Where("id = ?", s.ID).
		Updates(map[string]any{
			"service_name": s.ServiceName,
			"price":        s.Price,
			"user_id":      s.UserID,
			"start_month":  s.StartMonth,
			"end_month":    s.EndMonth,
			"updated_at":   time.Now().UTC(),
		})

	if res.Error != nil {
		return 0, res.Error
	}

	if res.RowsAffected == 0 {
		return 0, ErrNotFound
	}

	return res.RowsAffected, nil
}

func (r *Repo) Delete(ctx context.Context, id int64) error {
	res := r.db.WithContext(ctx).Delete(&Subscription{}, id)

	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *Repo) Summary(ctx context.Context, q SummaryQuery) (int64, error) {
	db := r.db.WithContext(ctx).Model(&Subscription{}).
		Where("start_month <= ?", q.To).
		Where("(end_month IS NULL OR end_month >= ?)", q.From)

	if q.UserID != nil {
		db = db.Where("user_id = ?", *q.UserID)
	}

	if q.ServiceName != nil && *q.ServiceName != "" {
		db = db.Where("service_name = ?", *q.ServiceName)
	}

	var total int64
	if err := db.Select("COALESCE(SUM(price), 0)::bigint AS total").Scan(&total).Error; err != nil {
		return 0, err
	}

	return total, nil
}
