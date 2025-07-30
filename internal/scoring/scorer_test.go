package scoring

import (
	"math"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/kdimtriCP/gh-inspector/internal/mock/mock_scoring"
)

func TestNewScorer(t *testing.T) {
	t.Run("with nil config", func(t *testing.T) {
		scorer := NewScorer(nil)
		require.NotNil(t, scorer)
		require.NotNil(t, scorer.config, "NewScorer() should set default config")
	})

	t.Run("with custom config", func(t *testing.T) {
		config := &Config{
			Weights: Weights{
				Stars: 0.5,
			},
		}
		scorer := NewScorer(config)
		require.Equal(t, config, scorer.config)
	})
}

func TestScore(t *testing.T) {
	scorer := NewScorer(DefaultConfig())

	tests := []struct {
		name      string
		setupMock func(*mock_scoring.MockRepositoryMetrics)
		wantMin   float64
		wantMax   float64
	}{
		{
			name: "archived repository",
			setupMock: func(m *mock_scoring.MockRepositoryMetrics) {
				m.EXPECT().GetIsArchived().Return(true)
			},
			wantMin: 0.0,
			wantMax: 0.0,
		},
		{
			name: "perfect repository",
			setupMock: func(m *mock_scoring.MockRepositoryMetrics) {
				m.EXPECT().GetIsArchived().Return(false)
				m.EXPECT().GetStars().Return(10000)
				m.EXPECT().GetForks().Return(1000)
				m.EXPECT().GetOpenIssues().Return(0)
				m.EXPECT().GetOpenPRs().Return(0)
				m.EXPECT().GetLastCommitDate().Return(time.Now())
				m.EXPECT().GetHasLicense().Return(true)
				m.EXPECT().GetHasCICD().Return(true)
				m.EXPECT().GetHasContributing().Return(true)
			},
			wantMin: 85.0,
			wantMax: 100.0,
		},
		{
			name: "inactive repository",
			setupMock: func(m *mock_scoring.MockRepositoryMetrics) {
				m.EXPECT().GetIsArchived().Return(false)
				m.EXPECT().GetStars().Return(100)
				m.EXPECT().GetForks().Return(10)
				m.EXPECT().GetOpenIssues().Return(50)
				m.EXPECT().GetOpenPRs().Return(20)
				m.EXPECT().GetLastCommitDate().Return(time.Now().AddDate(-2, 0, 0))
				m.EXPECT().GetHasLicense().Return(false)
				m.EXPECT().GetHasCICD().Return(false)
				m.EXPECT().GetHasContributing().Return(false)
			},
			wantMin: 0.0,
			wantMax: 30.0,
		},
		{
			name: "medium activity repository",
			setupMock: func(m *mock_scoring.MockRepositoryMetrics) {
				m.EXPECT().GetIsArchived().Return(false)
				m.EXPECT().GetStars().Return(1000)
				m.EXPECT().GetForks().Return(100)
				m.EXPECT().GetOpenIssues().Return(10)
				m.EXPECT().GetOpenPRs().Return(5)
				m.EXPECT().GetLastCommitDate().Return(time.Now().AddDate(0, -1, 0))
				m.EXPECT().GetHasLicense().Return(true)
				m.EXPECT().GetHasCICD().Return(true)
				m.EXPECT().GetHasContributing().Return(false)
			},
			wantMin: 40.0,
			wantMax: 70.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockMetrics := mock_scoring.NewMockRepositoryMetrics(ctrl)
			tt.setupMock(mockMetrics)

			got := scorer.Score(mockMetrics)
			require.GreaterOrEqual(t, got, tt.wantMin, "Score() should be >= %v", tt.wantMin)
			require.LessOrEqual(t, got, tt.wantMax, "Score() should be <= %v", tt.wantMax)
			require.GreaterOrEqual(t, got, 0.0, "Score() should be >= 0")
			require.LessOrEqual(t, got, 100.0, "Score() should be <= 100")
		})
	}
}

func TestCalculateActivityScore(t *testing.T) {
	scorer := NewScorer(DefaultConfig())

	tests := []struct {
		name           string
		lastCommitDate time.Time
		want           float64
	}{
		{
			name:           "committed today",
			lastCommitDate: time.Now(),
			want:           1.0,
		},
		{
			name:           "committed 5 days ago",
			lastCommitDate: time.Now().AddDate(0, 0, -5),
			want:           1.0,
		},
		{
			name:           "committed 20 days ago",
			lastCommitDate: time.Now().AddDate(0, 0, -20),
			want:           0.8,
		},
		{
			name:           "committed 60 days ago",
			lastCommitDate: time.Now().AddDate(0, 0, -60),
			want:           0.6,
		},
		{
			name:           "committed 120 days ago",
			lastCommitDate: time.Now().AddDate(0, 0, -120),
			want:           0.4,
		},
		{
			name:           "committed 200 days ago",
			lastCommitDate: time.Now().AddDate(0, 0, -200),
			want:           0.2,
		},
		{
			name:           "committed 400 days ago",
			lastCommitDate: time.Now().AddDate(0, 0, -400),
			want:           0.0,
		},
		{
			name:           "zero time",
			lastCommitDate: time.Time{},
			want:           0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := scorer.calculateActivityScore(tt.lastCommitDate)
			require.InDelta(t, tt.want, got, 0.01, "calculateActivityScore() = %v, want %v", got, tt.want)
		})
	}
}

func TestCalculateIssuesScore(t *testing.T) {
	scorer := NewScorer(DefaultConfig())

	tests := []struct {
		name       string
		openIssues int
		want       float64
	}{
		{name: "0 issues", openIssues: 0, want: 1.0},
		{name: "3 issues", openIssues: 3, want: 1.0 - math.Log10(4)/5.0},       // ~0.8796
		{name: "8 issues", openIssues: 8, want: 1.0 - math.Log10(9)/5.0},       // ~0.8092
		{name: "15 issues", openIssues: 15, want: 1.0 - math.Log10(16)/5.0},    // ~0.7592
		{name: "35 issues", openIssues: 35, want: 1.0 - math.Log10(36)/5.0},    // ~0.6887
		{name: "75 issues", openIssues: 75, want: 1.0 - math.Log10(76)/5.0},    // ~0.6238
		{name: "200 issues", openIssues: 200, want: 1.0 - math.Log10(201)/5.0}, // ~0.5394
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := scorer.calculateIssuesScore(tt.openIssues)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestCalculatePRsScore(t *testing.T) {
	scorer := NewScorer(DefaultConfig())

	tests := []struct {
		name    string
		openPRs int
		want    float64
	}{
		{name: "0 PRs", openPRs: 0, want: 1.0},
		{name: "2 PRs", openPRs: 2, want: 1.0 - math.Log10(3)/4.0},    // ~0.8807
		{name: "4 PRs", openPRs: 4, want: 1.0 - math.Log10(5)/4.0},    // ~0.8253
		{name: "7 PRs", openPRs: 7, want: 1.0 - math.Log10(8)/4.0},    // ~0.7742
		{name: "15 PRs", openPRs: 15, want: 1.0 - math.Log10(16)/4.0}, // ~0.6990
		{name: "30 PRs", openPRs: 30, want: 1.0 - math.Log10(31)/4.0}, // ~0.6272
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := scorer.calculatePRsScore(tt.openPRs)
			require.Equal(t, tt.want, got)
		})
	}
}
