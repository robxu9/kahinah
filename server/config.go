package main

type Config struct {
	Version   int
	SecretKey string

	GlobalNotice string
	Admin        admin

	HTTP      string
	DevMode   bool
	DebugMode bool

	Database   database
	Connectors connectors
	Karma      karma
	Advisory   advisory
}

type admin struct {
	Permanent []string
}

type database struct {
	Dialect string
	Params  string
}

type connectors struct {
	ABF abfstruct
}

type abfstruct struct {
	Enabled       bool
	PlatformIds   []string
	User          string
	APIKey        string
	CheckEveryMin int64
}

type karma struct {
	PassLimit int
	FailLimit int

	AddPassKarma int
	AddFailKarma int

	AddOverrideKarma int
	AddBlockKarma    int

	OverrideHours int64
	BlockHours    int64
}

type advisory struct {
	Families []string
}

func DefaultConfig() *Config {
	return &Config{
		Version:   VERSION,
		SecretKey: "ChangeMeToSomethingRandomForSecurity!",
		HTTP:      ":3000",
		DevMode:   true,
		DebugMode: false,
		Database: database{
			Dialect: "sqlite3",
			Params:  ":memory:",
		},
		Karma: karma{
			PassLimit:        3,
			FailLimit:        -3,
			AddPassKarma:     1,
			AddFailKarma:     -1,
			AddOverrideKarma: 6,
			AddBlockKarma:    -6,
			OverrideHours:    168,
			BlockHours:       12,
		},
	}
}
