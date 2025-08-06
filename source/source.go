package source

import (
	"context"

	"github.com/googollee/clic/structtags"
)

type Source interface {
	Register(fields []structtags.Field) error
	Parse(ctx context.Context, args []string) error
	Error() error
}

// var Default = []Source{
// 	Flag(FlagPrefix(""), FlagSplitter(".")),
// 	File(FilePathFlag("config", ""), FileFormat(JSON{})),
// 	Env(EnvPrefix("CLIC"), EnvSplitter("_")),
// }
