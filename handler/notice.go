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

	const noticeLineHour = 72 * time.Hour

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
				if isNearDueForNotice(question, noticeLineHour) {
					overDueProc(ctx, hm.db, question)
				}
				checkedQuestions[question.ID] = struct{}{}
			}
		case question := <-hm.ne.ReceiveUpdateQuestion():
			if question.Support != nil && question.Support.SupportStatusID == 3 {
				removeNoticeProcByQuestionID(ctx, hm.db, hm.hub, question.ID)
				delete(checkedQuestions, question.ID)
			} else {
				if question.Due != nil && isNearDueForNotice(question, noticeLineHour) {
					overDueProc(ctx, hm.db, question)
					checkedQuestions[question.ID] = struct{}{}
				} else {
					removeNoticeProcByQuestionID(ctx, hm.db, hm.hub, question.ID)
					delete(checkedQuestions, question.ID)
				}
			}
		case id := <-hm.ne.ReceiveDeleteQuestion():
			removeNoticeProcByQuestionID(ctx, hm.db, hm.hub, id)
			delete(checkedQuestions, id)
		case <-ctx.Done():
			return nil
		}
	}
}

// isNearDueForNotice は「期限の lead 時間前から」を通知登録の対象とする（従来の Due.Add(-72h).Before(now) と同義）。
func isNearDueForNotice(question infrastructure.Question, lead time.Duration) bool {
	if question.Due == nil {
		return false
	}
	return question.Due.Add(-lead).Before(time.Now())
}

func overDueProc(ctx context.Context, db *gorm.DB, question infrastructure.Question) {
	if question.Support == nil {
		return
	}
	// question_id での検索で RecordNotFound が出る想定のため、ログ出力を抑えた関数を呼びだしています
	if _, err := infrastructure.GetNoticeByQuestionSilent(context.Background(), db, question); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) && question.Support.SupportStatusID != 1 {
			infrastructure.RegisterNoticeByQuestionID(ctx, db, question.ID)
			return
		}
		fmt.Println(err.Error())
	}
}

func removeNoticeProcByQuestionID(ctx context.Context, db *gorm.DB, hub *Hub, id int64) {
	if nid, err := infrastructure.DeleteNoticeByQuestionID(ctx, db, id); err != nil {
		fmt.Println(err.Error())
	} else if nid >= 0 {
		hub.sendDeleteEvent("notice", strconv.FormatInt(nid, 10))
	}
}
