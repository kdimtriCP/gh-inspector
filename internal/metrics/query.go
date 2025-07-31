package metrics

import "github.com/shurcooL/githubv4"

type Owner struct {
	Login githubv4.String
}

type Language struct {
	Name githubv4.String
}

type IssuesConnection struct {
	TotalCount githubv4.Int
}

type PullRequestsConnection struct {
	TotalCount githubv4.Int
}

type CommitNode struct {
	CommittedDate githubv4.DateTime
}

type CommitEdge struct {
	Node CommitNode
}

type CommitHistory struct {
	Edges []CommitEdge
}

type Commit struct {
	History CommitHistory `graphql:"history(first: 1)"`
}

type GitObject struct {
	Commit Commit `graphql:"... on Commit"`
}

type Ref struct {
	Target GitObject
}

type License struct {
	Key githubv4.String
}

type TreeEntry struct {
	Name githubv4.String
}

type Tree struct {
	Entries []TreeEntry
}

type TreeObject struct {
	Tree Tree `graphql:"... on Tree"`
}

type ReleaseNode struct {
	PublishedAt githubv4.DateTime
}

type ReleaseEdge struct {
	Node ReleaseNode
}

type ReleasesConnection struct {
	TotalCount githubv4.Int
	Edges      []ReleaseEdge
}

type RepositoryGraphQL struct {
	Owner            Owner
	Name             githubv4.String
	Description      githubv4.String
	StargazerCount   githubv4.Int
	ForkCount        githubv4.Int
	IsArchived       githubv4.Boolean
	PrimaryLanguage  *Language
	Issues           IssuesConnection       `graphql:"issues(states: OPEN)"`
	PullRequests     PullRequestsConnection `graphql:"pullRequests(states: OPEN)"`
	DefaultBranchRef *Ref
	LicenseInfo      *License
	Object           TreeObject         `graphql:"object(expression: \"HEAD:\")"`
	Releases         ReleasesConnection `graphql:"releases(first: 10, orderBy: {field: CREATED_AT, direction: DESC})"`
	Watchers         struct {
		TotalCount githubv4.Int
	}
}

type RepositoryQuery struct {
	Repository RepositoryGraphQL `graphql:"repository(owner: $owner, name: $name)"`
}
