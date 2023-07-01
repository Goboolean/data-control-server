package time

import (
	"time"
)

/* mulmuri's comment:
 I'm not sure using this package is a good idea.
 Except for the fact that configuring the condition is a bit complicated,
 checking the market is open or not may be dependant to its platform, not to time.
 Regarding this issue, test code is not written yet.
*/

// It is a checker that checks whether the current time is within the specified duration.
// Since holyday can be varied year by year, it is recommeded to renew the condition every year.
type WithinDurationChecker struct {
	option *Option
}

// Option is an argument struct of Constructor of WithinDurationChecker
// a condition that is not included in both Inclusion and Exclusion is regarded as Exclusion
type Option struct {
	Location  string        // Location of time zone
	Exclusion ConditionList // List of conditions to exclude
	Inclusion ConditionList // List of conditions to include

	location  time.Location // It will be initialized by constructor referring to Location
}

// WithinDurationChecker checks validity of option when it is generated
func New(o *Option) (*WithinDurationChecker, error) {

	instance := &WithinDurationChecker{option: o}
	return instance, instance.verifyOption()
}

// Here are accepted type list
// time.Weekday  : to specify weekday repeated on weekday
// time.Time : to specify Date
// int           : to specify day repteted on month
// or else it throws panic
type DateSpecification interface{}

// Example of usage is as follows:
// ConditionList{
//   time.Monday:                            &StatusOption{},
//   time.Parse("2006-01-02", "0000-00-00"): &StatusOption{},
//   1 :                                     &StatusOption{},
// }
type ConditionList map[DateSpecification]*StatusOption

type StatusOption struct {
	Allday    bool   // equals to 00:00:00 to 23:59:59 if given value is true
	StartTime string // string format: 09:30:00
	EndTime   string // string format: 15:30:00

	startTime time.Duration // will be initialized by constructor referring to StartTime
	endTime   time.Duration // will be initialized by constructor referring to EndTime
}




func forDomestic() {

	New(&Option{
		location: time.Location{},
		Inclusion: ConditionList{
			time.Monday:    &StatusOption{StartTime: "09:30:00", EndTime: "15:30:00"},
			time.Tuesday:   &StatusOption{StartTime: "09:30:00", EndTime: "15:30:00"},
			time.Wednesday: &StatusOption{StartTime: "09:30:00", EndTime: "15:30:00"},
			time.Thursday:  &StatusOption{StartTime: "09:30:00", EndTime: "15:30:00"},
			time.Friday:    &StatusOption{StartTime: "09:30:00", EndTime: "15:30:00"},
		},
	})
}

func (m *WithinDurationChecker) verifyOption() error {

	// check if interface{} defined type: dateSpecification is whether time.Weekday, time.Duration, int
	for dateSpecification, condition := range m.option.Inclusion {
		switch dateSpecification.(type) {
		case time.Weekday:  break
		case time.Duration: break
		case int:           break
		default: return errInvalidDateSpecificationType
		}
		
		startTime, err := time.ParseDuration(condition.StartTime)
		if err != nil {
			return err
		}

		endTime, err := time.ParseDuration(condition.EndTime)
		if err != nil {
			return err
		}

		condition.startTime = startTime
		condition.endTime = endTime
	}

	// check if string format of StartTime and EndTime inside StatusOption is valid

	return nil
}

func (m *WithinDurationChecker) checkCondition(s *StatusOption, timeNow time.Duration) bool {

	if s.Allday {
		return true
	} else if s.startTime <= timeNow && timeNow <= s.endTime {
		return true
	}

	return false
}

func (m *WithinDurationChecker) checkConditionList(conditionList ConditionList) bool {

	timeNow := time.Now().In(&m.option.location) // current time
	timeToday := timeNow.Round(24 * time.Hour) // time of 00:00:00 today
	durationNow := timeNow.Sub(timeToday) // duration from 00:00:00 today to current time	

	for dateSpecification, condition := range conditionList {

		switch dateSpecification.(type) {

		// case for repeated weekday
		case time.Weekday:
			if timeNow.Weekday() == dateSpecification.(time.Weekday) {
				if m.checkCondition(condition, durationNow) {
					return true
				}
			}

		// case for the specific date
		case time.Time:
			if timeToday == dateSpecification.(time.Time) {
				if m.checkCondition(condition, durationNow) {
					return true
				}
			}

		// case for repeated day on month
		case int:
			if timeNow.Day() == dateSpecification.(int) {
				if m.checkCondition(condition, durationNow) {
					return false
				}
			}
		}
	}

	return false
}



func (m *WithinDurationChecker) IsTimeNowWithinDuration() bool {

	// checks whether condition matches in the excluded list
	// if true, then TimeNowWithinDuration is false
	if contains := m.checkConditionList(m.option.Exclusion); contains {
		return false
	}

	// checks whether condition matches in the included list
	// if true, then TimeNowWithinDuration is true
	if contains := m.checkConditionList(m.option.Inclusion); contains {
		return true
	}

	// default return value is false
	return false
}