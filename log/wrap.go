package log

import "fmt"

type Logger2 interface {
	Printf(format string, args ...interface{})
}

func WrapLogger2(l Logger2) Logger {
	if l == nil {
		return nil
	}
	return &logger2Wrapper{
		l: l,
	}
}

type logger2Wrapper struct {
	l Logger2
}

func (w *logger2Wrapper) Errorf(format string, args ...interface{}) {
	w.l.Printf(format, args...)
}

type Logger3 interface {
	Error(args ...interface{})
}

func WrapLogger3(l Logger3) Logger {
	if l == nil {
		return nil
	}
	return &logger3Wrapper{
		l: l,
	}
}

type logger3Wrapper struct {
	l Logger3
}

func (w *logger3Wrapper) Errorf(format string, args ...interface{}) {
	w.l.Error(fmt.Sprintf(format, args...))
}

type Logger4 interface {
	Error(msg string, fields ...interface{})
}

func WrapLogger4(l Logger4) Logger {
	if l == nil {
		return nil
	}
	return &logger4Wrapper{
		l: l,
	}
}

type logger4Wrapper struct {
	l Logger4
}

func (w *logger4Wrapper) Errorf(format string, args ...interface{}) {
	w.l.Error(fmt.Sprintf(format, args...))
}
