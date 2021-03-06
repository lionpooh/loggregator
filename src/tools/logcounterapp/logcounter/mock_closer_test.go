// This file was generated by github.com/nelsam/hel.  Do not
// edit this code by hand unless you *really* know what you're
// doing.  Expect any changes made manually to be overwritten
// the next time hel regenerates this file.

package logcounter_test

type mockCloser struct {
	CloseCalled chan bool
	CloseOutput struct {
		Ret0 chan error
	}
}

func newMockCloser() *mockCloser {
	m := &mockCloser{}
	m.CloseCalled = make(chan bool, 100)
	m.CloseOutput.Ret0 = make(chan error, 100)
	return m
}
func (m *mockCloser) Close() error {
	m.CloseCalled <- true
	return <-m.CloseOutput.Ret0
}
