package scoring

import (
	"math"
	"time"
)

type RepoMetrics interface {
	GetStars() int
	GetForks() int
	GetOpenIssues() int
	GetOpenPRs() int
	GetLastCommitDate() time.Time
	GetIsArchived() bool
	GetHasLicense() bool
	GetHasCICD() bool
	GetHasContributing() bool
}

type Scorer struct {
	config *Config
}

func NewScorer(config *Config) *Scorer {
	if config == nil {
		config = DefaultConfig()
	}
	return &Scorer{config: config}
}

func (s *Scorer) Score(metrics RepoMetrics) float64 {
	if metrics.GetIsArchived() {
		return 0.0
	}

	weights := s.config.Weights
	score := 0.0

	starsScore := math.Min(math.Log10(float64(metrics.GetStars()+1))/4.0, 1.0)
	score += starsScore * weights.Stars

	forksScore := math.Min(math.Log10(float64(metrics.GetForks()+1))/3.0, 1.0)
	score += forksScore * weights.Forks

	activityScore := s.calculateActivityScore(metrics.GetLastCommitDate())
	score += activityScore * weights.RecentActivity

	issuesScore := s.calculateIssuesScore(metrics.GetOpenIssues())
	score += issuesScore * weights.OpenIssues

	prsScore := s.calculatePRsScore(metrics.GetOpenPRs())
	score += prsScore * weights.OpenPRs

	if metrics.GetHasLicense() {
		score += weights.HasLicense
	}
	if metrics.GetHasCICD() {
		score += weights.HasCICD
	}
	if metrics.GetHasContributing() {
		score += weights.HasContributing
	}

	// Normalize to 0-100 scale
	return math.Min(score*100, 100)
}

func (s *Scorer) calculateActivityScore(lastCommitDate time.Time) float64 {
	if lastCommitDate.IsZero() {
		return 0.0
	}

	daysSinceCommit := time.Since(lastCommitDate).Hours() / 24

	switch {
	case daysSinceCommit <= 7:
		return 1.0
	case daysSinceCommit <= 30:
		return 0.8
	case daysSinceCommit <= 90:
		return 0.6
	case daysSinceCommit <= 180:
		return 0.4
	case daysSinceCommit <= 365:
		return 0.2
	default:
		return 0.0
	}
}

func (s *Scorer) calculateIssuesScore(openIssues int) float64 {
	switch {
	case openIssues == 0:
		return 1.0
	case openIssues <= 5:
		return 0.9
	case openIssues <= 10:
		return 0.7
	case openIssues <= 25:
		return 0.5
	case openIssues <= 50:
		return 0.3
	case openIssues <= 100:
		return 0.1
	default:
		return 0.0
	}
}

func (s *Scorer) calculatePRsScore(openPRs int) float64 {
	switch {
	case openPRs == 0:
		return 1.0
	case openPRs <= 3:
		return 0.8
	case openPRs <= 5:
		return 0.6
	case openPRs <= 10:
		return 0.4
	case openPRs <= 20:
		return 0.2
	default:
		return 0.0
	}
}
