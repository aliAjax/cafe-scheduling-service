package roster

import (
	"context"
	"encoding/json"
	"io"
)

func EncodeWeek(ctx context.Context, writer io.Writer, view WeekView) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return json.NewEncoder(writer).Encode(cloneWeek(view))
}
