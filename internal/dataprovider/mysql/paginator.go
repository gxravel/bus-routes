package mysql

import (
	"github.com/gxravel/bus-routes/internal/dataprovider"

	"github.com/Masterminds/squirrel"
)

func withPaginator(b squirrel.SelectBuilder, p *dataprovider.Paginator) squirrel.SelectBuilder {
	return b.Limit(p.Limit()).Offset(p.Offset())
}
