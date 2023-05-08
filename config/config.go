package config

func Init() {
	logger := NewLogger("config")
	if err := initializeEnvironment(); err != nil {
		logger.Warnf("error initializing environment: %v", err)
	}
}
