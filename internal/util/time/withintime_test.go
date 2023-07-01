package time

import (
	"reflect"
	"testing"
	"time"
)

func Test_Constructor(t *testing.T) {
	type args struct {
		o *Option
	}
	tests := []struct {
		name    string
		args    args
		want    *WithinDurationChecker
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.o)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}


func TestWithinDurationChecker_verifyOption(t *testing.T) {
	tests := []struct {
		name    string
		m       *WithinDurationChecker
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.m.verifyOption(); (err != nil) != tt.wantErr {
				t.Errorf("WithinDurationChecker.verifyOption() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWithinDurationChecker_checkCondition(t *testing.T) {
	type args struct {
		s       *StatusOption
		timeNow time.Duration
	}
	tests := []struct {
		name string
		m    *WithinDurationChecker
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.checkCondition(tt.args.s, tt.args.timeNow); got != tt.want {
				t.Errorf("WithinDurationChecker.checkCondition() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithinDurationChecker_checkConditionList(t *testing.T) {
	type args struct {
		conditionList ConditionList
	}
	tests := []struct {
		name string
		m    *WithinDurationChecker
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.checkConditionList(tt.args.conditionList); got != tt.want {
				t.Errorf("WithinDurationChecker.checkConditionList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithinDurationChecker_IsTimeNowWithinDuration(t *testing.T) {
	tests := []struct {
		name string
		m    *WithinDurationChecker
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.IsTimeNowWithinDuration(); got != tt.want {
				t.Errorf("WithinDurationChecker.IsTimeNowWithinDuration() = %v, want %v", got, tt.want)
			}
		})
	}
}
