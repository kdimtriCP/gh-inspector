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
	HasReadme        float64 `yaml:"has_readme"`
	HasCodeOfConduct float64 `yaml:"has_code_of_conduct"`
	HasSecurity      float64 `yaml:"has_security"`
	Watchers         float64 `yaml:"watchers"`
}

func DefaultConfig() *Config {
	return &Config{
		Weights: Weights{
			Stars:            0.20,
			Forks:            0.08,
			RecentActivity:   0.18,
			OpenIssues:       0.08,
			OpenPRs:          0.04,
			HasLicense:       0.04,
			HasCICD:          0.04,
			HasContributing:  0.04,
			ReleaseFrequency: 0.12,
			HasReadme:        0.03,
			HasCodeOfConduct: 0.03,
			HasSecurity:      0.03,
			Watchers:         0.09,
		},
	}
}
