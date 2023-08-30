package kafka

import "strings"

// config is deprecated

// every topic that regards to stock data should be prefixed with this
var prefix string = "broker."

// make topic string prefixed
// every topic access should be done through the packTopic function
func packTopic(topic string) string {
	return prefix + topic
}

// remove prefix from topic string
// every topic queried from the broker should be unpacked with the unpackTopic function
func unpackTopic(topic string) string {
	return strings.TrimPrefix(topic, prefix)
}
