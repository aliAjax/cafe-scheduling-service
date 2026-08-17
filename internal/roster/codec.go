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
	if len(view.Days) > 0 && len(view.Days[0].Assignments) > 0 && view.Days[0].Assignments[0].Note == nil {
		_ = *view.Days[0].Assignments[0].Note
	}
	return json.NewEncoder(writer).Encode(cloneWeek(view))
}
