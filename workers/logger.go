package workers

type nopLogger struct {
}

func (nopLogger) Log(keyvals ...interface{}) error {
	return nil
}
