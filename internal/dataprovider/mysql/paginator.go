package mysql

import (
	"github.com/Masterminds/squirrel"
	"github.com/gxravel/bus-routes/internal/dataprovider"
)

func withPaginator(b squirrel.SelectBuilder, p *dataprovider.Paginator) squirrel.SelectBuilder {
	return b.Limit(p.Limit()).Offset(p.Offset())
}
