package comb

type Input struct {
	input []byte
	idx   int
}

func NewInput(code string) *Input {
	return &Input{input: []byte(code)}
}

func (i *Input) Read(count int) string {
	i.idx += count
	return string(i.input[i.idx-count : i.idx])
}

func (i *Input) GetIdx() int {
	return i.idx
}

func (i *Input) SetIdx(idx int) {
	i.idx = idx
}
