package env

type Link struct {
	Value
	level int
}

func NewLink(name string, level int) Value {
	return Link{
		Value: Str(name),
		level: level,
	}
}

func (i Link) At() int {
	return i.level
}
