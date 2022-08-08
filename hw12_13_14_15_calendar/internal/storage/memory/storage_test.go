package memorystorage

import (
	"sync"
	"testing"
	"time"

	storage2 "github.com/LightAir/otus_home_work/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func getEvents(firstID uuid.UUID, secondID uuid.UUID, userID uuid.UUID) map[uuid.UUID]storage2.Event {
	return map[uuid.UUID]storage2.Event{
		firstID: {
			ID:            firstID,
			Title:         "Test title",
			DatetimeStart: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
			DatetimeEnd:   time.Date(2009, time.November, 10, 23, 15, 0, 0, time.UTC),
			Description:   "just description",
			UserID:        userID,
			WhenToNotify:  time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
		},
		secondID: {
			ID:            secondID,
			Title:         "Test title",
			DatetimeStart: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
			DatetimeEnd:   time.Date(2009, time.November, 11, 23, 15, 0, 0, time.UTC),
			Description:   "just description",
			UserID:        userID,
			WhenToNotify:  time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
		},
	}
}

func getDateRange(startDate time.Time) (*storage2.DateRange, error) {
	start, err := time.Parse("2006-01-02", startDate.Format("2006-01-02"))
	if err != nil {
		return nil, err
	}

	end := start.AddDate(0, 0, 1)

	return &storage2.DateRange{
		Start: start,
		End:   end,
	}, nil
}

func TestStorage(t *testing.T) {
	firstID := uuid.New()
	secondID := uuid.New()
	userID := uuid.New()

	t.Run("base tests", func(t *testing.T) {
		storage := New()

		events := getEvents(firstID, secondID, userID)
		err := storage.AddEvent(events[firstID])
		require.Nil(t, err)

		err = storage.AddEvent(events[secondID])
		require.Nil(t, err)

		dateRange, err := getDateRange(events[secondID].DatetimeStart)
		require.Nil(t, err)

		result, err := storage.ListEventsByRange(*dateRange)
		require.Nil(t, err)
		require.Equal(t, events, result)

		err = storage.RemoveEvent(firstID)
		require.Nil(t, err)

		result, err = storage.ListEventsByRange(*dateRange)
		require.Nil(t, err)
		delete(events, firstID)
		require.Equal(t, events, result)

		err = storage.RemoveEvent(secondID)
		require.Nil(t, err)

		result, err = storage.ListEventsByRange(*dateRange)
		require.Nil(t, err)
		delete(events, secondID)
		require.Equal(t, events, result)

		err = storage.Close()
		require.Nil(t, err)
	})

	t.Run("already exist tests", func(t *testing.T) {
		storage := New()

		events := getEvents(firstID, secondID, userID)
		err := storage.AddEvent(events[firstID])
		require.Nil(t, err)

		err = storage.AddEvent(events[firstID])
		require.ErrorIs(t, err, errEventAlreadyExist)

		err = storage.Close()
		require.Nil(t, err)
	})

	t.Run("notfound tests", func(t *testing.T) {
		storage := New()

		err := storage.RemoveEvent(firstID)
		require.ErrorIs(t, err, errEventNotFound)

		events := getEvents(firstID, secondID, userID)
		err = storage.ChangeEvent(firstID, events[firstID])
		require.ErrorIs(t, err, errEventNotFound)

		err = storage.Close()
		require.Nil(t, err)
	})

	t.Run("change event", func(t *testing.T) {
		storage := New()

		events := getEvents(firstID, secondID, userID)
		err := storage.AddEvent(events[firstID])
		require.Nil(t, err)

		event := events[firstID]
		event.Title = "new new text"

		err = storage.ChangeEvent(firstID, event)
		require.Nil(t, err)

		dateRange, err := getDateRange(events[secondID].DatetimeStart)
		require.Nil(t, err)

		result, err := storage.ListEventsByRange(*dateRange)
		require.Nil(t, err)
		require.Equal(t, map[uuid.UUID]storage2.Event{firstID: event}, result)

		err = storage.Close()
		require.Nil(t, err)
	})
}

func TestCacheMultithreading(t *testing.T) {
	storage := New()

	wg := &sync.WaitGroup{}
	wg.Add(2)

	userID := uuid.New()
	timeStart := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	timeEnd := time.Date(2009, time.November, 10, 23, 15, 0, 0, time.UTC)
	timeWhen := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			event := storage2.Event{
				ID:            uuid.New(),
				Title:         "Test title",
				DatetimeStart: timeStart,
				DatetimeEnd:   timeEnd,
				Description:   "just description",
				UserID:        userID,
				WhenToNotify:  timeWhen,
			}
			err := storage.AddEvent(event)
			require.Nil(t, err)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			event := storage2.Event{
				ID:            uuid.New(),
				Title:         "Test title",
				DatetimeStart: timeStart,
				DatetimeEnd:   timeEnd,
				Description:   "just description",
				UserID:        userID,
				WhenToNotify:  timeWhen,
			}
			err := storage.AddEvent(event)
			require.Nil(t, err)
		}
	}()

	wg.Wait()
}
