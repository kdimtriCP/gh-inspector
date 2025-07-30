package metrics

import "time"

type Repository struct {
	Owner           string
	Name            string
	Stars           int
	Forks           int
	OpenIssues      int
	OpenPRs         int
	LastCommitDate  time.Time
	Description     string
	PrimaryLanguage string
	IsArchived      bool
	HasCICD         bool
	HasLicense      bool
	HasContributing bool
	ReleaseCount    int
	LastReleaseDate time.Time
	Score           float64
}

func (m *Repository) GetStars() int                 { return m.Stars }
func (m *Repository) GetForks() int                 { return m.Forks }
func (m *Repository) GetOpenIssues() int            { return m.OpenIssues }
func (m *Repository) GetOpenPRs() int               { return m.OpenPRs }
func (m *Repository) GetLastCommitDate() time.Time  { return m.LastCommitDate }
func (m *Repository) GetIsArchived() bool           { return m.IsArchived }
func (m *Repository) GetHasLicense() bool           { return m.HasLicense }
func (m *Repository) GetHasCICD() bool              { return m.HasCICD }
func (m *Repository) GetHasContributing() bool      { return m.HasContributing }
func (m *Repository) GetReleaseCount() int          { return m.ReleaseCount }
func (m *Repository) GetLastReleaseDate() time.Time { return m.LastReleaseDate }

func (m *Repository) DaysSinceLastCommit() int {
	if m.LastCommitDate.IsZero() {
		return -1
	}
	return int(time.Since(m.LastCommitDate).Hours() / 24)
}

func (m *Repository) FullName() string {
	return m.Owner + "/" + m.Name
}
