package app

type Config struct {
	Logging  bool
	Trace    bool
	Bind     string
	Database databaseConfig
}

func (c *Config) Initialize() error {
	if c.Database.File == "" {
		c.Database.File = "./events.db"
	}

	return nil
}
