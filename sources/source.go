package sources

import (
	"context"

	"github.com/googollee/clic/structtags"
)

type Source interface {
	Prepare(flagSet FlagSet, fields []structtags.Field) error
	Parse(ctx context.Context) error
	Error() error
}
