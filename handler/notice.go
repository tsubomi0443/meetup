package handler

import (
	"context"
	"errors"
	"fmt"
	infrastructure "meetup/_mac_infrastructure"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type NoticeEvent struct {
	Question struct {
		updateNotice chan infrastructure.Question
		deleteNotice chan int64
	}
}

func NewNoticeEvent() *NoticeEvent {
	return &NoticeEvent{
		Question: struct {
			updateNotice chan infrastructure.Question
			deleteNotice chan int64
		}{
			updateNotice: make(chan infrastructure.Question),
			deleteNotice: make(chan int64),
		},
	}
}

func (ne *NoticeEvent) UpdateQuestion(model infrastructure.Question) {
	ne.Question.updateNotice <- model
}

func (ne *NoticeEvent) DeleteQuestion(qid int64) {
	ne.Question.deleteNotice <- qid
}

func (ne *NoticeEvent) ReceiveUpdateQuestion() chan infrastructure.Question {
	return ne.Question.updateNotice
}

func (ne *NoticeEvent) ReceiveDeleteQuestion() chan int64 {
	return ne.Question.deleteNotice
}

// Question を定期監視し、期限が近づいたものを Notice へと登録する（バックグラウンド）
func (hm *HandlerManager) PollingStart(ctx context.Context) error {
	return hm.checkQuesiton(ctx)
}

func (hm *HandlerManager) checkQuesiton(ctx context.Context) error {
	var ticker = time.NewTicker(30 * time.Minute)
	var checkedQuestions = make(map[int64]any)

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
				if question.Due.Add(-72 * time.Hour).Before(time.Now()) {
					// question_id での検索で RecordNotFound が出る想定のため、ログ出力を抑えた関数を呼びだしています
					if _, err := infrastructure.GetNoticeByQuestionSilent(context.Background(), hm.db, question); err != nil {
						if errors.Is(err, gorm.ErrRecordNotFound) && question.Support.SupportStatusID != 1 {
							infrastructure.RegisterNoticeByQuestionID(ctx, hm.db, question.ID)
							checkedQuestions[question.ID] = struct{}{}
							continue
						}
						fmt.Println(err.Error())
					}
				}
				checkedQuestions[question.ID] = struct{}{}
			}
		case q := <-hm.ne.ReceiveUpdateQuestion():
			if q.Support.SupportStatusID == 3 {
				if nid, err := infrastructure.DeleteNoticeByQuestion(ctx, hm.db, q); err != nil {
					fmt.Println(err.Error())
				} else if nid >= 0 {
					hm.hub.sendDeleteEvent("notice", strconv.FormatInt(nid, 10))
				}
				delete(checkedQuestions, q.ID)
			}
		case id := <-hm.ne.ReceiveDeleteQuestion():
			if nid, err := infrastructure.DeleteNoticeByQuestionID(ctx, hm.db, id); err != nil {
				fmt.Println(err.Error())
			} else if nid >= 0 {
				hm.hub.sendDeleteEvent("notice", strconv.FormatInt(nid, 10))
			}
			delete(checkedQuestions, id)
		case <-ctx.Done():
			return nil
		}
	}
}
