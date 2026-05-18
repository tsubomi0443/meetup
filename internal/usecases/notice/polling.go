package notice

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"meetup/internal/domains/entity"
	domainnotice "meetup/internal/domains/notice"
	domainquestion "meetup/internal/domains/question"

	"gorm.io/gorm"
)

// SSEDeleteNotifier is invoked when a notice row is removed (for SSE fan-out).
type SSEDeleteNotifier func(api, id string)

type Poller struct {
	questions domainquestion.Repository
	notices   domainnotice.Repository
	events    *Event
	notify    SSEDeleteNotifier
}

func NewPoller(
	questions domainquestion.Repository,
	notices domainnotice.Repository,
	events *Event,
	notify SSEDeleteNotifier,
) *Poller {
	return &Poller{
		questions: questions,
		notices:   notices,
		events:    events,
		notify:    notify,
	}
}

func (p *Poller) Run(ctx context.Context) error {
	var ticker = time.NewTicker(30 * time.Minute)
	var checkedQuestions = make(map[int64]any)

	const noticeLineHour = 72 * time.Hour

	for {
		select {
		case <-ticker.C:
			questions, err := p.questions.GetAll(context.Background())
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
					p.overDueProc(ctx, question)
				}
				checkedQuestions[question.ID] = struct{}{}
			}
		case question := <-p.events.ReceiveUpdateQuestion():
			if question.Support != nil && question.Support.SupportStatusID == 3 {
				p.removeNoticeProcByQuestionID(ctx, question.ID)
				delete(checkedQuestions, question.ID)
			} else {
				if question.Due != nil && isNearDueForNotice(question, noticeLineHour) {
					p.overDueProc(ctx, question)
					checkedQuestions[question.ID] = struct{}{}
				} else {
					p.removeNoticeProcByQuestionID(ctx, question.ID)
					delete(checkedQuestions, question.ID)
				}
			}
		case id := <-p.events.ReceiveDeleteQuestion():
			p.removeNoticeProcByQuestionID(ctx, id)
			delete(checkedQuestions, id)
		case <-ctx.Done():
			return nil
		}
	}
}

func isNearDueForNotice(question entity.Question, lead time.Duration) bool {
	if question.Due == nil {
		return false
	}
	return question.Due.Add(-lead).Before(time.Now())
}

func (p *Poller) overDueProc(ctx context.Context, question entity.Question) {
	if question.Support == nil {
		return
	}
	if _, err := p.notices.GetByQuestionSilent(context.Background(), question); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) && question.Support.SupportStatusID != 1 {
			_ = p.notices.RegisterByQuestionID(ctx, question.ID)
			return
		}
		fmt.Println(err.Error())
	}
}

func (p *Poller) removeNoticeProcByQuestionID(ctx context.Context, id int64) {
	if nid, err := p.notices.DeleteByQuestionID(ctx, id); err != nil {
		fmt.Println(err.Error())
	} else if nid >= 0 && p.notify != nil {
		p.notify("notice", strconv.FormatInt(nid, 10))
	}
}
