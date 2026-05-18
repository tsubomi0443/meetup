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

// SSEDeleteNotifier は通知行削除時に SSE へ削除イベントを配信するためのコールバック型。
//
// args:
//   - api string: API 種別（例: "notice"）
//   - id string: 削除対象リソース ID
type SSEDeleteNotifier func(api, id string)

// Poller は期限接近・質問更新に応じて通知の作成・削除を行うバックグラウンドワーカー。
type Poller struct {
	questions domainquestion.Repository
	notices   domainnotice.Repository
	events    *Event
	notify    SSEDeleteNotifier
}

// NewPoller は通知ポーラーを生成する。
//
// args:
//   - questions domainquestion.Repository: 質問リポジトリ
//   - notices domainnotice.Repository: 通知リポジトリ
//   - events *Event: 質問更新・削除イベントハブ
//   - notify SSEDeleteNotifier: 通知削除時の SSE コールバック（nil 可）
//
// return:
//   - *Poller: 生成したポーラー
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

// Run は定期スキャンとイベント受信のメインループを実行する。ctx がキャンセルされると終了する。
//
// args:
//   - ctx context.Context: キャンセル用コンテキスト
//
// return:
//   - error: 定期スキャン中の取得エラー（キャンセル時は nil）
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

// isNearDueForNotice は期限の lead 時間前に入ったか（通知対象に近いか）を判定する。
//
// args:
//   - question entity.Question: 判定対象の質問
//   - lead time.Duration: 期限前の余裕時間（例: 72 時間）
//
// return:
//   - bool: 期限が近く通知対象とみなす場合 true
func isNearDueForNotice(question entity.Question, lead time.Duration) bool {
	if question.Due == nil {
		return false
	}
	return question.Due.Add(-lead).Before(time.Now())
}

// overDueProc は期限接近かつ未対応以外の質問に対し、通知がなければ登録する。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//   - question entity.Question: 対象質問
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

// removeNoticeProcByQuestionID は質問 ID に紐づく通知を削除し、必要なら SSE 削除イベントを送る。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//   - id int64: 質問 ID
func (p *Poller) removeNoticeProcByQuestionID(ctx context.Context, id int64) {
	if nid, err := p.notices.DeleteByQuestionID(ctx, id); err != nil {
		fmt.Println(err.Error())
	} else if nid >= 0 && p.notify != nil {
		p.notify("notice", strconv.FormatInt(nid, 10))
	}
}
