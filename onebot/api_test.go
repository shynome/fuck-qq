package onebot

import (
	"context"
	"testing"
)

func TestActive(t *testing.T) {
	ctx := context.Background()
	err := active(ctx, "194519409")
	t.Log(err)
}
