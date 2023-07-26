package store

import (
	"context"
	"testing"
)



func Test_contextController(t *testing.T) {

	ctx := context.Background()

	CTX := new_ctx(ctx)

	select {
	case <- CTX.Done():
		t.Errorf("context did not created successfully")
	default:
		break
	}

	CTX.cancel()

	select {
	case <- CTX.Done():
		break
	default:
		t.Errorf("context did not canceled successfully")
	}

}