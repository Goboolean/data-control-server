package withintime

import (
	"reflect"
	"testing"
	"time"
)


// Integration test for WithinDurationChecker struct is consist of following Tests:
// Test_verifyCondition
// Test_verifyOption
// Test_checkCondition
// Test_checkConditionList
// Test_Constructor



func Test_verifyCondition(t *testing.T) {

	type args struct {
		dateSpecification DateSpecification
		condition         *Condition
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid dateSpecification type: repeated weekday",
			args: args{
				dateSpecification: time.Monday,
				condition: &Condition{Allday: true},
			},
			wantErr: false,
		},
		{
			name: "valid dateSpecification type: one day",
			args: args{
				// ?????
				condition: &Condition{Allday: true},
			},
			wantErr: false,
		},
		{
			name: "valid dateSpecification type: repeated day of month",
			args: args{
				dateSpecification: -1,
				condition: &Condition{Allday: true},
			},
			wantErr: false,
		},
		{
			name: "invalid dateSpecification type: nil",
			args: args{
				dateSpecification: nil,
				condition: &Condition{Allday: true},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO: Need to implement		
		})
	}
}



func Test_verifyOption(t *testing.T) {

	// Listed testcase only cover the case related to location.
	// Case related to duration is covered by Test_verifyCondition.

	tests := []struct {
		name    string
		m       *WithinDurationChecker
		wantErr bool
	}{
		{
			name: "no location",
			m: &WithinDurationChecker{
					option: &Option{},
				},
			wantErr: true,
		},
		{
			name: "invalid location",
			m: &WithinDurationChecker{
					option: &Option{
						Location: "invalid location",
					},
				},
			wantErr: true,
		},
		{
			name: "valid location: korea",
			m: &WithinDurationChecker{
					option: &Option{
						Location: "Asia/Seoul",
					},
				},
			wantErr: false,
		},
		{
			name: "valid location: usa",
			m: &WithinDurationChecker{
					option: &Option{
						Location: "America/New_York",
					},
				},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.m.verifyOption(); (err != nil) != tt.wantErr {
				t.Errorf("WithinDurationChecker.verifyOption() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}


func Test_chechCondition(t *testing.T) {

}



func Test_checkConditionList(t *testing.T) {

}



func Test_Constructor(t *testing.T) {
	type args struct {
		o *Option
		t TimeProvider
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
			got, err := New(tt.args.o, tt.args.t)
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