package scoring

import (
	"math"
	"time"
)

//go:generate mockgen -source=$GOFILE -destination=../mock/mock_scoring/mock_$GOFILE -package=mock_scoring

type RepositoryMetrics interface {
	GetStars() int
	GetForks() int
	GetOpenIssues() int
	GetOpenPRs() int
	GetLastCommitDate() time.Time
	GetIsArchived() bool
	GetHasLicense() bool
	GetHasCICD() bool
	GetHasContributing() bool
	GetReleaseCount() int
	GetLastReleaseDate() time.Time
	GetHasReadme() bool
	GetHasCodeOfConduct() bool
	GetHasSecurity() bool
	GetWatchers() int
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

func (s *Scorer) Score(metrics RepositoryMetrics) float64 {
	if metrics.GetIsArchived() {
		return 0.0
	}

	weights := s.config.Weights
	score := 0.0

	starsScore := math.Min(math.Log10(float64(metrics.GetStars()+1))/5.0, 1.0)
	score += starsScore * weights.Stars

	forksScore := math.Min(math.Log10(float64(metrics.GetForks()+1))/4.5, 1.0)
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
	if metrics.GetHasReadme() {
		score += weights.HasReadme
	}
	if metrics.GetHasCodeOfConduct() {
		score += weights.HasCodeOfConduct
	}
	if metrics.GetHasSecurity() {
		score += weights.HasSecurity
	}

	watchersScore := s.calculateWatchersScore(metrics.GetWatchers())
	score += watchersScore * weights.Watchers

	releaseScore := s.calculateReleaseFrequencyScore(metrics.GetReleaseCount(), metrics.GetLastReleaseDate())
	score += releaseScore * weights.ReleaseFrequency

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
	if openIssues == 0 {
		return 1.0
	}

	issueRatio := float64(openIssues)
	return math.Max(0, 1.0-math.Log10(issueRatio+1)/5.0)
}

func (s *Scorer) calculatePRsScore(openPRs int) float64 {
	if openPRs == 0 {
		return 1.0
	}

	prRatio := float64(openPRs)
	return math.Max(0, 1.0-math.Log10(prRatio+1)/4.0)
}

func (s *Scorer) calculateReleaseFrequencyScore(releaseCount int, lastReleaseDate time.Time) float64 {
	if releaseCount == 0 {
		return 0.0
	}

	if lastReleaseDate.IsZero() {
		return 0.0
	}

	daysSinceRelease := time.Since(lastReleaseDate).Hours() / 24

	var recencyScore float64
	switch {
	case daysSinceRelease <= 30:
		recencyScore = 1.0
	case daysSinceRelease <= 90:
		recencyScore = 0.8
	case daysSinceRelease <= 180:
		recencyScore = 0.6
	case daysSinceRelease <= 365:
		recencyScore = 0.4
	default:
		recencyScore = 0.2
	}

	frequencyScore := math.Min(float64(releaseCount)/10.0, 1.0)

	return (recencyScore + frequencyScore) / 2.0
}

func (s *Scorer) calculateWatchersScore(watchers int) float64 {
	return math.Min(math.Log10(float64(watchers+1))/4.0, 1.0)
}
