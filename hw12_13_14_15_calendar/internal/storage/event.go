package storage

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID            uuid.UUID      // Уникальный идентификатор события (можно воспользоваться UUID)
	Title         string         // Короткий текст
	DatetimeStart time.Time      // Дата и время события
	DatetimeEnd   time.Time      // Дата и время окончания события
	Description   sql.NullString // Описание события - длинный текст, опционально
	UserID        uuid.UUID      // ID пользователя, владельца события
	WhenToNotify  sql.NullString // За сколько времени высылать уведомление, опционально.
}
