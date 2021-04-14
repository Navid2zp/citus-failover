package logging

import "go.uber.org/zap"

func (l *Logger) APIStartFailed(err error) {
	l.l.Error("failed to start api",
		zap.String("service", l.service), zap.Error(err))
}
