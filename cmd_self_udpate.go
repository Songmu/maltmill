package maltmill

import (
	"context"

	"github.com/Songmu/ghselfupdate"
)

type cmdSelfUpdate struct {
}

var _ runner = (*cmdSelfUpdate)(nil)

func (upd *cmdSelfUpdate) run(ctx context.Context) error {
	return ghselfupdate.Do(version)
}
