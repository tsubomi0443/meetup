package handler

import (
	"context"
	"errors"
	"fmt"
	infrastructure "meetup/_mac_infrastructure"
	"time"

	"gorm.io/gorm"
)

// Questionを定期監視し、期限が近づいたものをNoticeへと登録する
func (hm *HandlerManager) PollingStart(ctx context.Context) error {
	return hm.checkQuesiton(ctx)
}

func (hm *HandlerManager) checkQuesiton(ctx context.Context) error {
	var ticker = time.NewTicker(30 * time.Minute)
	var checkedQuestions map[int64]any = make(map[int64]any)

	for {
		select {
		case <-ticker.C:
			questions, err := infrastructure.GetQuestions(context.Background(), hm.db)
			if err != nil {
				return err
			}

			for _, question := range questions {
				if _, ok := checkedQuestions[question.ID]; ok {
					continue
				}
				if question.Due == nil {
					checkedQuestions[question.ID] = struct{}{}
					continue
				}
				if question.Due.Add(-24 * time.Hour).Before(time.Now()) {
					// question_idでの検索でRecordNotFoundが出る想定のため、ログ出力を抑えた関数を呼びだしています
					if _, err := infrastructure.GetNoticeByQuestionSilent(context.Background(), hm.db, question); err != nil {
						if errors.Is(err, gorm.ErrRecordNotFound) {
							infrastructure.RegisterNoticeByQuestionID(ctx, hm.db, question.ID)
							checkedQuestions[question.ID] = struct{}{}
							continue
						}
						fmt.Println(err.Error())
					}
				}
				checkedQuestions[question.ID] = struct{}{}
			}
		case <-ctx.Done():
			return nil
		}
	}
}
