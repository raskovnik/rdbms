package database

type Index struct {
	ColumnName string
	Data       map[interface{}][]int
}

func NewIndex(columnName string) *Index {
	return &Index{
		ColumnName: columnName,
		Data:       make(map[interface{}][]int),
	}
}

// add a value to the index for a given row index
func (idx *Index) Add(value interface{}, rowIndex int) {
	idx.Data[value] = append(idx.Data[value], rowIndex)
}

// remove a value from the index
func (idx *Index) Remove(value interface{}, rowIndex int) {
	indices := idx.Data[value]

	for i, ind := range indices {
		if ind == rowIndex {
			idx.Data[value] = append(indices[:i], indices[i+1:]...)
			break
		}
	}

	if len(idx.Data[value]) == 0 {
		delete(idx.Data, value)
	}
}

// return row indices for a given value
func (idx *Index) Lookup(value interface{}) []int {
	return idx.Data[value]
}

// check if a value exists in the index
func (idx *Index) Exists(value interface{}) bool {
	_, exists := idx.Data[value]
	return exists
}
