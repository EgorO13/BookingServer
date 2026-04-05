package Service

import (
	"fmt"
	"time"

	"task/internal/Models"
	"task/internal/Repository"

	"github.com/google/uuid"
)

type IScheduleService interface {
	CreateSchedule(schedule *Models.Schedule) error
}

type ScheduleService struct {
	scheduleRepo Repository.IScheduleRepository
	slotRepo     Repository.ISlotRepository
	roomRepo     Repository.IRoomRepository
}

func NewScheduleService(scheduleRepo Repository.IScheduleRepository, slotRepo Repository.ISlotRepository, roomRepo Repository.IRoomRepository) IScheduleService {
	return &ScheduleService{
		scheduleRepo: scheduleRepo,
		slotRepo:     slotRepo,
		roomRepo:     roomRepo,
	}
}

func (s *ScheduleService) CreateSchedule(schedule *Models.Schedule) error {
	exists, err := s.roomRepo.Exists(schedule.RoomID)
	if err != nil {
		return fmt.Errorf("failed to check room: %w", err)
	}
	if !exists {
		return ErrRoomNotFound
	}
	exists, err = s.scheduleRepo.ExistsForRoom(schedule.RoomID)
	if err != nil {
		return fmt.Errorf("failed to check schedule: %w", err)
	}
	if exists {
		return ErrScheduleAlreadyExists
	}
	startTimeParsed, err := time.Parse("15:04", schedule.StartTime)
	if err != nil {
		return fmt.Errorf("invalid start time format: %w", err)
	}
	endTimeParsed, err := time.Parse("15:04", schedule.EndTime)
	if err != nil {
		return fmt.Errorf("invalid end time format: %w", err)
	}
	if !startTimeParsed.Before(endTimeParsed) {
		return fmt.Errorf("start time must be before end time")
	}
	if err := s.scheduleRepo.Create(schedule); err != nil {
		return fmt.Errorf("failed to create schedule: %w", err)
	}
	if err := s.generateSlots(schedule.RoomID, schedule.DaysOfWeek, schedule.StartTime, schedule.EndTime); err != nil {
		return fmt.Errorf("schedule created but failed to generate slots: %w", err)
	}
	return nil
}

func (s *ScheduleService) generateSlots(roomID uuid.UUID, daysOfWeek []int, startTime, endTime string) error {
	if err := s.slotRepo.DeleteByRoomID(roomID); err != nil {
		return fmt.Errorf("failed to delete existing slots: %w", err)
	}
	startTimeParsed, err := time.Parse("15:04", startTime)
	if err != nil {
		return fmt.Errorf("invalid start time format: %w", err)
	}
	endTimeParsed, err := time.Parse("15:04", endTime)
	if err != nil {
		return fmt.Errorf("invalid end time format: %w", err)
	}
	now := time.Now().UTC()
	startDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	endDate := startDate.Add(30 * 24 * time.Hour)
	var slots []Models.Slot
	for d := startDate; d.Before(endDate); d = d.Add(24 * time.Hour) {
		weekday := int(d.Weekday())
		if weekday == 0 {
			weekday = 7
		}
		if !contains(daysOfWeek, weekday) {
			continue
		}
		slotStart := time.Date(d.Year(), d.Month(), d.Day(), startTimeParsed.Hour(), startTimeParsed.Minute(), 0, 0, time.UTC)
		slotEnd := time.Date(d.Year(), d.Month(), d.Day(), endTimeParsed.Hour(), endTimeParsed.Minute(), 0, 0, time.UTC)
		for t := slotStart; t.Before(slotEnd); t = t.Add(30 * time.Minute) {
			slot := Models.Slot{
				ID:     uuid.New(),
				RoomID: roomID,
				Start:  t,
				End:    t.Add(30 * time.Minute),
			}
			slots = append(slots, slot)
		}
	}
	if len(slots) > 0 {
		if err := s.slotRepo.CreateBatch(slots); err != nil {
			return fmt.Errorf("failed to create slots: %w", err)
		}
	}
	return nil
}

func contains(slice []int, val int) bool {
	for _, v := range slice {
		if v == val {
			return true
		}
	}
	return false
}
