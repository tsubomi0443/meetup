package notice

import "meetup/internal/domains/entity"

// Event は質問の更新・削除をポーラーへ通知するためのチャネルハブを表す。
type Event struct {
	Question struct {
		updateNotice chan entity.Question
		deleteNotice chan int64
	}
}

// NewEvent は質問更新・削除用チャネルを初期化したイベントハブを生成する。
//
// return:
//   - *Event: 生成したイベントハブ
func NewEvent() *Event {
	return &Event{
		Question: struct {
			updateNotice chan entity.Question
			deleteNotice chan int64
		}{
			updateNotice: make(chan entity.Question),
			deleteNotice: make(chan int64),
		},
	}
}

// UpdateQuestion は質問更新イベントをチャネルへ送信する。
//
// args:
//   - model entity.Question: 更新された質問
func (ne *Event) UpdateQuestion(model entity.Question) {
	ne.Question.updateNotice <- model
}

// DeleteQuestion は質問削除イベントをチャネルへ送信する。
//
// args:
//   - qid int64: 削除された質問 ID
func (ne *Event) DeleteQuestion(qid int64) {
	ne.Question.deleteNotice <- qid
}

// ReceiveUpdateQuestion は質問更新通知を受け取るチャネルを返す。
//
// return:
//   - chan entity.Question: 更新された質問を受け取るチャネル
func (ne *Event) ReceiveUpdateQuestion() chan entity.Question {
	return ne.Question.updateNotice
}

// ReceiveDeleteQuestion は質問削除通知を受け取るチャネルを返す。
//
// return:
//   - chan int64: 削除された質問 ID を受け取るチャネル
func (ne *Event) ReceiveDeleteQuestion() chan int64 {
	return ne.Question.deleteNotice
}
