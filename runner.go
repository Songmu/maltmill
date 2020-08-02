package maltmill

import "context"

type runner interface {
	run(context.Context) error
}
