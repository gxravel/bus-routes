package dataprovider

const (
	MaxPerPage     = 100
	DefaultPerPage = 20
)

type Paginator struct {
	offset uint64
	limit  uint64
}

func NewPaginator(offset, limit uint64) *Paginator {
	return &Paginator{
		offset: offset,
		limit:  normalizeLimit(limit),
	}
}

func normalizeLimit(limit uint64) uint64 {
	switch {
	case limit == 0:
		return DefaultPerPage
	case limit > MaxPerPage:
		return MaxPerPage
	default:
		return limit
	}
}

func (p *Paginator) Offset() uint64 {
	return p.offset
}

func (p *Paginator) Limit() uint64 {
	return p.limit
}
