// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    getIssueCommentsReturns, err := UnmarshalGetIssueCommentsReturns(bytes)
//    bytes, err = getIssueCommentsReturns.Marshal()

package getissuecomments

import "encoding/json"

func UnmarshalGetIssueCommentsReturns(data []byte) (GetIssueCommentsReturns, error) {
	var r GetIssueCommentsReturns
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *GetIssueCommentsReturns) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type GetIssueCommentsReturns struct {
	Comments []IssueComment `json:"comments,omitempty"`
}

type IssueComment struct {
	AuthorAssociation *string    `json:"author_association,omitempty"`
	Body              *string    `json:"body,omitempty"`
	CreatedAt         *string    `json:"created_at,omitempty"`
	HTMLURL           *string    `json:"html_url,omitempty"`
	ID                *int64     `json:"id,omitempty"`
	IssueURL          *string    `json:"issue_url,omitempty"`
	NodeID            *string    `json:"node_id,omitempty"`
	Reactions         *Reactions `json:"reactions,omitempty"`
	UpdatedAt         *string    `json:"updated_at,omitempty"`
	URL               *string    `json:"url,omitempty"`
	User              *User      `json:"user,omitempty"`
}

type Reactions struct {
	The1       *int64  `json:"+1,omitempty"`
	Reactions1 *int64  `json:"-1,omitempty"`
	Confused   *int64  `json:"confused,omitempty"`
	Eyes       *int64  `json:"eyes,omitempty"`
	Heart      *int64  `json:"heart,omitempty"`
	Hooray     *int64  `json:"hooray,omitempty"`
	Laugh      *int64  `json:"laugh,omitempty"`
	Rocket     *int64  `json:"rocket,omitempty"`
	TotalCount *int64  `json:"total_count,omitempty"`
	URL        *string `json:"url,omitempty"`
}

type User struct {
	AvatarURL               *string         `json:"avatar_url,omitempty"`
	Bio                     *string         `json:"bio,omitempty"`
	Blog                    *string         `json:"blog,omitempty"`
	Collaborators           *int64          `json:"collaborators,omitempty"`
	Company                 *string         `json:"company,omitempty"`
	CreatedAt               *Timestamp      `json:"created_at,omitempty"`
	DiskUsage               *int64          `json:"disk_usage,omitempty"`
	Email                   *string         `json:"email,omitempty"`
	EventsURL               *string         `json:"events_url,omitempty"`
	Followers               *int64          `json:"followers,omitempty"`
	FollowersURL            *string         `json:"followers_url,omitempty"`
	Following               *int64          `json:"following,omitempty"`
	FollowingURL            *string         `json:"following_url,omitempty"`
	GistsURL                *string         `json:"gists_url,omitempty"`
	GravatarID              *string         `json:"gravatar_id,omitempty"`
	Hireable                *bool           `json:"hireable,omitempty"`
	HTMLURL                 *string         `json:"html_url,omitempty"`
	ID                      *int64          `json:"id,omitempty"`
	LDAPDN                  *string         `json:"ldap_dn,omitempty"`
	Location                *string         `json:"location,omitempty"`
	Login                   *string         `json:"login,omitempty"`
	Name                    *string         `json:"name,omitempty"`
	NodeID                  *string         `json:"node_id,omitempty"`
	OrganizationsURL        *string         `json:"organizations_url,omitempty"`
	OwnedPrivateRepos       *int64          `json:"owned_private_repos,omitempty"`
	Permissions             map[string]bool `json:"permissions,omitempty"`
	Plan                    *Plan           `json:"plan,omitempty"`
	PrivateGists            *int64          `json:"private_gists,omitempty"`
	PublicGists             *int64          `json:"public_gists,omitempty"`
	PublicRepos             *int64          `json:"public_repos,omitempty"`
	ReceivedEventsURL       *string         `json:"received_events_url,omitempty"`
	ReposURL                *string         `json:"repos_url,omitempty"`
	SiteAdmin               *bool           `json:"site_admin,omitempty"`
	StarredURL              *string         `json:"starred_url,omitempty"`
	SubscriptionsURL        *string         `json:"subscriptions_url,omitempty"`
	SuspendedAt             *Timestamp      `json:"suspended_at,omitempty"`
	TextMatches             []TextMatch     `json:"text_matches,omitempty"`
	TotalPrivateRepos       *int64          `json:"total_private_repos,omitempty"`
	TwitterUsername         *string         `json:"twitter_username,omitempty"`
	TwoFactorAuthentication *bool           `json:"two_factor_authentication,omitempty"`
	Type                    *string         `json:"type,omitempty"`
	UpdatedAt               *Timestamp      `json:"updated_at,omitempty"`
	URL                     *string         `json:"url,omitempty"`
}

type Timestamp struct {
}

type Plan struct {
	Collaborators *int64  `json:"collaborators,omitempty"`
	FilledSeats   *int64  `json:"filled_seats,omitempty"`
	Name          *string `json:"name,omitempty"`
	PrivateRepos  *int64  `json:"private_repos,omitempty"`
	Seats         *int64  `json:"seats,omitempty"`
	Space         *int64  `json:"space,omitempty"`
}

type TextMatch struct {
	Fragment   *string `json:"fragment,omitempty"`
	Matches    []Match `json:"matches,omitempty"`
	ObjectType *string `json:"object_type,omitempty"`
	ObjectURL  *string `json:"object_url,omitempty"`
	Property   *string `json:"property,omitempty"`
}

type Match struct {
	Indices []int64 `json:"indices,omitempty"`
	Text    *string `json:"text,omitempty"`
}
