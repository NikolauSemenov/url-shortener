package ports

type Logger interface {
	Info(msg string, args ...any)
	Error(msg string, args ...any)
	With(args ...any) Logger
}
