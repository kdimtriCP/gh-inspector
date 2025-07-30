package scoring

type Config struct {
	Weights Weights `yaml:"weights"`
}

type Weights struct {
	Stars            float64 `yaml:"stars"`
	Forks            float64 `yaml:"forks"`
	RecentActivity   float64 `yaml:"recent_activity"`
	OpenIssues       float64 `yaml:"open_issues"`
	OpenPRs          float64 `yaml:"open_prs"`
	HasLicense       float64 `yaml:"has_license"`
	HasCICD          float64 `yaml:"has_cicd"`
	HasContributing  float64 `yaml:"has_contributing"`
	ReleaseFrequency float64 `yaml:"release_frequency"`
}

func DefaultConfig() *Config {
	return &Config{
		Weights: Weights{
			Stars:            0.25,
			Forks:            0.10,
			RecentActivity:   0.20,
			OpenIssues:       0.10,
			OpenPRs:          0.05,
			HasLicense:       0.05,
			HasCICD:          0.05,
			HasContributing:  0.05,
			ReleaseFrequency: 0.15,
		},
	}
}
