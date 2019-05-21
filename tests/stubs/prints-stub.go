package stubs

func NewPrintsStub() *PrintsStub {
	return &PrintsStub{0}
}

type PrintsStub struct {
	called int
}

func (p *PrintsStub) Printf(format string, v ...interface{}) {
	p.called++
}
