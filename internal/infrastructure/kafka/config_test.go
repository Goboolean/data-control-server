package kafka

import "testing"

func Test_packTopic(t *testing.T) {
	type args struct {
		topic string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "topic name asdf test",
			args: args{topic: "asdf"},
			want: "broker.asdf",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := packTopic(tt.args.topic); got != tt.want {
				t.Errorf("packTopic() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_unpackTopic(t *testing.T) {
	type args struct {
		topic string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "topic name asdf test",
			args: args{topic: "broker.asdf"},
			want: "asdf",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := unpackTopic(tt.args.topic); got != tt.want {
				t.Errorf("unpackTopic() = %v, want %v", got, tt.want)
			}
		})
	}
}


func Test_All(t *testing.T) {
	type args struct {
		topic string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "topic name asdf test",
			args: args{topic: "asdf"},
			want: "asdf",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			packed := packTopic(tt.args.topic)
			if got := unpackTopic(packed); got != tt.want {
				t.Errorf("unpackTopic() = %v, want %v", got, tt.want)
			}
		})
	}
}