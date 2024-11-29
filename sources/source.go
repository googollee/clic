package sources

import (
	"context"
	"flag"

	"github.com/googollee/clic/structtags"
)

type Source interface {
	Prepare(flagSet *flag.FlagSet, fields []structtags.Field) error
	Parse(ctx context.Context) error
	Error() error
}
