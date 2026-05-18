package notice

import "meetup/internal/domains/entity"

type Event struct {
	Question struct {
		updateNotice chan entity.Question
		deleteNotice chan int64
	}
}

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

func (ne *Event) UpdateQuestion(model entity.Question) {
	ne.Question.updateNotice <- model
}

func (ne *Event) DeleteQuestion(qid int64) {
	ne.Question.deleteNotice <- qid
}

func (ne *Event) ReceiveUpdateQuestion() chan entity.Question {
	return ne.Question.updateNotice
}

func (ne *Event) ReceiveDeleteQuestion() chan int64 {
	return ne.Question.deleteNotice
}
