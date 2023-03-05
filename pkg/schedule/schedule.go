package schedule

import (
	"errors"
	"github.com/miksir/unkatan/pkg/lconfig"
	zlog "github.com/miksir/unkatan/pkg/log"
	"time"
)

type SpecificSchedule struct {
	at time.Time
}

type Schedule struct {
	cfg        lconfig.Reader
	log        zlog.Logger
	daily      []int          // часы когда выкатки закрыты
	shortDaily []int          // часы когда выкатки закрыты перед выходными
	weekly     []time.Weekday // Выходные
	yearly     []string       // праздники
	specific   []SpecificSchedule
}

type TriggerType int

const (
	TriggerNot TriggerType = iota
	TriggerDaily
	TriggerWeekday
	TriggerRegularHolidays
	TriggerSpecialHolidays
	TriggerDailyShortDay
	TriggerExclude
	TriggerNewYear
)

func NewSchedule(cfg lconfig.Reader, log zlog.Logger) *Schedule {
	schd := &Schedule{
		cfg: cfg,
		log: log,
	}

	schd.daily = []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 19, 20, 21, 22, 23}
	schd.shortDaily = []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 17, 18, 19, 20, 21, 22, 23}
	schd.weekly = []time.Weekday{time.Saturday, time.Sunday}
	schd.yearly = []string{
		"01.01", "02.01", "03.01", "04.01", "05.01", "06.01", "07.01", "08.01",
		"23.02",
		"08.03",
		"01.05",
		"09.05",
		"12.06",
		"04.11",
	}

	return schd
}

func (schd *Schedule) FindNextEvent(currentStatus bool, day time.Time) (time.Time, bool, TriggerType, error) {
	day = day.Truncate(1 * time.Hour)
	for i := 0; i < 480; i++ {
		day = day.Add(1 * time.Hour)
		checkStatus, checkReason := schd.CheckSchedule(day)
		if checkStatus != currentStatus {
			return day, checkStatus, checkReason, nil
		}
	}
	return day, false, TriggerNot, errors.New("not found")
}

func (schd *Schedule) CheckSchedule(day time.Time) (bool, TriggerType) {
	if day.Month() == time.December && day.Day() >= 22 {
		return false, TriggerNewYear
	}

	status, reason := schd.checkHolidays(day)
	if status == false {
		return status, reason
	}

	// проверяем можно ли катить в этот час
	tomorrowStatus, _ := schd.checkHolidays(day.Add(24 * time.Hour))
	var dailyCheck []int
	var trigger TriggerType
	if tomorrowStatus == false {
		// если завтра - праздник, то у нас короткий день
		dailyCheck = schd.shortDaily
		trigger = TriggerDailyShortDay
	} else {
		dailyCheck = schd.daily
		trigger = TriggerDaily
	}
	hour := day.Hour()
	for _, h := range dailyCheck {
		if hour == h {
			return false, trigger
		}
	}

	return true, TriggerNot
}

func (schd *Schedule) checkHolidays(day time.Time) (bool, TriggerType) {
	date := day.Format("02.01.2006")
	week := day.Weekday()

	// переносы выходных
	for _, s := range schd.cfg.GetStringSlice("exclude") {
		if date == s {
			return true, TriggerExclude
		}
	}

	// ежегодные праздники
	for _, y := range schd.yearly {
		if y == day.Format("02.01") {
			return false, TriggerRegularHolidays
		}
	}

	// специальные нерабочие дни или просто дни когда выкатки закрыты
	for _, s := range schd.cfg.GetStringSlice("special") {
		if date == s {
			return false, TriggerSpecialHolidays
		}
	}

	// проверяем день недели
	for _, w := range schd.weekly {
		if week == w {
			return false, TriggerWeekday
		}
	}

	return true, TriggerNot
}

func RussianReasonName(reason TriggerType) string {
	var name string
	switch reason {
	case TriggerDaily, TriggerNot, TriggerDailyShortDay:
		name = "расписание рабочего времени"
	case TriggerWeekday:
		name = "выходной"
	case TriggerRegularHolidays, TriggerSpecialHolidays:
		name = "праздники"
	case TriggerExclude:
		name = "перенос выходного"
	case TriggerNewYear:
		name = "с наступающим Новым Годом!"
	default:
		name = "расписание"
	}
	return name
}
