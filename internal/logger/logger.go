package logger

type customLogger struct {
}

type CustomLogger interface {
}

func New(logLevel string) (CustomLogger, error) {
	// Configure logger
	return &customLogger{}, nil
}
