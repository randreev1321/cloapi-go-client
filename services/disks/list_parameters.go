package disks

type PaginatorOptions struct {
	Limit  int
	Offset int
}

type FilteringField struct {
	FieldName string
	Condition string
	Value     string
}
