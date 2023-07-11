package withintime

import (
	"time"
)



// It is not recommended to use this package on production.
// Except for the fact that configuring the condition is a bit complicated,
// checking the market is open or not should be dependant to its platform, not to time, such in case of platform server shutting down, 
// Use this package only for testing purpose.

// It is a checker that checks whether the current time is within the specified duration.
// Since holyday can be varied year by year, it is recommeded to renew the condition every year manually.
type WithinDurationChecker struct {
	option *Option
	timePrivider TimeProvider
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
// If TimeProvider is nil, then it will be initialized as CurrentTime, which uses time.Now()
// If you want to specify time for testing purpose,
// then you can use FixedTime with initializing time field value as you want
func New(o *Option, t TimeProvider) (*WithinDurationChecker, error) {

	if t == nil {
		t = &CurrentTime{}
	}

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
//   time.Monday:                            &Condition{},
//   time.Parse("2006-01-02", "0000-00-00"): &Condition{},
//   1 :                                     &Condition{},
// }
type ConditionList map[DateSpecification]*Condition

type Condition struct {
	Allday    bool   // equals to 00:00:00 to 23:59:59 if given value is true
	StartTime string // string format: 09:30:00
	EndTime   string // string format: 15:30:00

	startTime time.Duration // will be initialized by constructor referring to StartTime
	endTime   time.Duration // will be initialized by constructor referring to EndTime
}

func (m *WithinDurationChecker) verifyCondition(dateSpecification DateSpecification, condition *Condition) error {

	// check if interface{} defined type: dateSpecification is whether time.Weekday, time.Duration, int
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

	return nil
}



func (m *WithinDurationChecker) verifyOption() error {

	for dateSpecification, condition := range m.option.Inclusion {
		if err := m.verifyCondition(dateSpecification, condition); err != nil {
			return err
		}
	}

	for dateSpecification, condition := range m.option.Exclusion {
		if err := m.verifyCondition(dateSpecification, condition); err != nil {
			return err
		}
	}

	return nil
}

func (m *WithinDurationChecker) checkCondition(s *Condition, timeNow time.Duration) bool {

	if s.Allday {
		return true
	} else if s.startTime <= timeNow && timeNow <= s.endTime {
		return true
	}

	return false
}

func (m *WithinDurationChecker) checkConditionList(conditionList ConditionList) bool {

	timeNow := m.timePrivider.Time().In(&m.option.location) // current time
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
	// 1. check whether condition matches with Exclusion list to be excluded
	// 2. check whether condition matches with Inclusion list to be included
	// 3. default is false

	if matches := m.checkConditionList(m.option.Exclusion); matches {
		return false
	}

	if matches := m.checkConditionList(m.option.Inclusion); matches {
		return true
	}

	return false
}