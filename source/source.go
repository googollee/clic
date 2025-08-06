package source

import (
	"context"

	"github.com/googollee/clic/structtags"
)

type Source interface {
	Register(fset FlagSet, fields []structtags.Field) error
	Parse(ctx context.Context, args []string) error
	Error() error
}
