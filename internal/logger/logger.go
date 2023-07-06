package logger

import "go.uber.org/zap"

var Log *zap.SugaredLogger = zap.NewNop().Sugar()

func Initialize(level string) error {
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return err
	}

	cfg := zap.NewDevelopmentConfig()
	cfg.Level = lvl
	zl := zap.Must(cfg.Build())
	Log = zl.Sugar()
	return nil
}
