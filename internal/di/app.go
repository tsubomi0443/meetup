package di

import (
	authuc "meetup/internal/usecases/auth"
	masteruc "meetup/internal/usecases/master"
	noticeuc "meetup/internal/usecases/notice"
	questionuc "meetup/internal/usecases/question"
	taguc "meetup/internal/usecases/tag"
	useruc "meetup/internal/usecases/user"

	"meetup/internal/infrastructures/crypto"
	"meetup/internal/infrastructures/persistence/postgres"
	"meetup/internal/interfaces"
	"meetup/internal/interfaces/sse"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

// App wires layers and exposes the HTTP router.
type App struct {
	Router *interfaces.Router
}

func NewApp(db *gorm.DB, e *echo.Echo) *App {
	hasher := crypto.NewHasherAdapter()

	userRepo := postgres.NewUserRepository(db)
	questionRepo := postgres.NewQuestionRepository(db)
	tagRepo := postgres.NewTagRepository(db)
	noticeRepo := postgres.NewNoticeRepository(db)
	masterRepo := postgres.NewMasterRepository(db)

	authUC := authuc.NewUseCase(userRepo, hasher)
	userUC := useruc.NewUseCase(userRepo, hasher)
	questionUC := questionuc.NewUseCase(questionRepo)
	tagUC := taguc.NewUseCase(tagRepo)
	noticeUC := noticeuc.NewUseCase(noticeRepo)
	masterUC := masteruc.NewUseCase(masterRepo)

	hub := sse.NewHub()
	noticeEvents := noticeuc.NewEvent()
	noticePoller := noticeuc.NewPoller(questionRepo, noticeRepo, noticeEvents, hub.SendDeleteEvent)

	router := interfaces.NewRouter(e, interfaces.Deps{
		Hub:          hub,
		NoticeEvents: noticeEvents,
		NoticePoller: noticePoller,
		Auth:         authUC,
		User:         userUC,
		Question:     questionUC,
		Tag:          tagUC,
		Notice:       noticeUC,
		Master:       masterUC,
	})

	return &App{Router: router}
}
