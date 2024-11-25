package sources

import (
	"context"
	"reflect"
)

type Source interface {
	ConfigWith(t reflect.Type) error
	ParseTo(ctx context.Context, value any) error
}
