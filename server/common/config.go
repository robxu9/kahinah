package common

// Config represents server configuration
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

// DefaultConfig represents the default configuration with the specified version
func DefaultConfig(version int) *Config {
	return &Config{
		Version:   version,
		SecretKey: "ChangeMeToSomethingRandomForSecurity!",

		GlobalNotice: "",
		Admin: admin{
			Permanent: []string{
				"permanent@admin.email",
			},
		},

		HTTP:      ":3000",
		DevMode:   true,
		DebugMode: false,

		Database: database{
			Dialect: "sqlite3",
			Params:  ":memory:",
		},

		Connectors: connectors{
			ABF: abfstruct{
				Enabled:       false,
				PlatformIds:   []string{},
				User:          "user",
				APIKey:        "apikey",
				CheckEveryMin: 60,
			},
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

		Advisory: advisory{
			Families: []string{
				"myfamily",
			},
		},
	}
}