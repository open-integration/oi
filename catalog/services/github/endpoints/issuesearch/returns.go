// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    issueSearchReturns, err := UnmarshalIssueSearchReturns(bytes)
//    bytes, err = issueSearchReturns.Marshal()

package issuesearch

import "encoding/json"

func UnmarshalIssueSearchReturns(data []byte) (IssueSearchReturns, error) {
	var r IssueSearchReturns
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *IssueSearchReturns) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type IssueSearchReturns struct {
	Issues []Issue `json:"issues,omitempty"`
}

type Issue struct {
	ActiveLockReason  *string           `json:"active_lock_reason,omitempty"`
	Assignee          *User             `json:"assignee,omitempty"`          
	Assignees         []User            `json:"assignees,omitempty"`         
	AuthorAssociation *string           `json:"author_association,omitempty"`
	Body              *string           `json:"body,omitempty"`              
	ClosedAt          *string           `json:"closed_at,omitempty"`         
	ClosedBy          *User             `json:"closed_by,omitempty"`         
	Comments          *int64            `json:"comments,omitempty"`          
	CommentsURL       *string           `json:"comments_url,omitempty"`      
	CreatedAt         *string           `json:"created_at,omitempty"`        
	EventsURL         *string           `json:"events_url,omitempty"`        
	HTMLURL           *string           `json:"html_url,omitempty"`          
	ID                *int64            `json:"id,omitempty"`                
	Labels            []Label           `json:"labels,omitempty"`            
	LabelsURL         *string           `json:"labels_url,omitempty"`        
	Locked            *bool             `json:"locked,omitempty"`            
	Milestone         *Milestone        `json:"milestone,omitempty"`         
	NodeID            *string           `json:"node_id,omitempty"`           
	Number            *int64            `json:"number,omitempty"`            
	PullRequest       *PullRequestLinks `json:"pull_request,omitempty"`      
	Reactions         *Reactions        `json:"reactions,omitempty"`         
	Repository        *Repository       `json:"repository,omitempty"`        
	RepositoryURL     *string           `json:"repository_url,omitempty"`    
	State             *string           `json:"state,omitempty"`             
	TextMatches       []TextMatch       `json:"text_matches,omitempty"`      
	Title             *string           `json:"title,omitempty"`             
	UpdatedAt         *string           `json:"updated_at,omitempty"`        
	URL               *string           `json:"url,omitempty"`               
	User              *User             `json:"user,omitempty"`              
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

type Label struct {
	Color       *string `json:"color,omitempty"`      
	Default     *bool   `json:"default,omitempty"`    
	Description *string `json:"description,omitempty"`
	ID          *int64  `json:"id,omitempty"`         
	Name        *string `json:"name,omitempty"`       
	NodeID      *string `json:"node_id,omitempty"`    
	URL         *string `json:"url,omitempty"`        
}

type Milestone struct {
	ClosedAt     *string `json:"closed_at,omitempty"`    
	ClosedIssues *int64  `json:"closed_issues,omitempty"`
	CreatedAt    *string `json:"created_at,omitempty"`   
	Creator      *User   `json:"creator,omitempty"`      
	Description  *string `json:"description,omitempty"`  
	DueOn        *string `json:"due_on,omitempty"`       
	HTMLURL      *string `json:"html_url,omitempty"`     
	ID           *int64  `json:"id,omitempty"`           
	LabelsURL    *string `json:"labels_url,omitempty"`   
	NodeID       *string `json:"node_id,omitempty"`      
	Number       *int64  `json:"number,omitempty"`       
	OpenIssues   *int64  `json:"open_issues,omitempty"`  
	State        *string `json:"state,omitempty"`        
	Title        *string `json:"title,omitempty"`        
	UpdatedAt    *string `json:"updated_at,omitempty"`   
	URL          *string `json:"url,omitempty"`          
}

type PullRequestLinks struct {
	DiffURL  *string `json:"diff_url,omitempty"` 
	HTMLURL  *string `json:"html_url,omitempty"` 
	PatchURL *string `json:"patch_url,omitempty"`
	URL      *string `json:"url,omitempty"`      
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

type Repository struct {
	AllowMergeCommit    *bool           `json:"allow_merge_commit,omitempty"`    
	AllowRebaseMerge    *bool           `json:"allow_rebase_merge,omitempty"`    
	AllowSquashMerge    *bool           `json:"allow_squash_merge,omitempty"`    
	ArchiveURL          *string         `json:"archive_url,omitempty"`           
	Archived            *bool           `json:"archived,omitempty"`              
	AssigneesURL        *string         `json:"assignees_url,omitempty"`         
	AutoInit            *bool           `json:"auto_init,omitempty"`             
	BlobsURL            *string         `json:"blobs_url,omitempty"`             
	BranchesURL         *string         `json:"branches_url,omitempty"`          
	CloneURL            *string         `json:"clone_url,omitempty"`             
	CodeOfConduct       *CodeOfConduct  `json:"code_of_conduct,omitempty"`       
	CollaboratorsURL    *string         `json:"collaborators_url,omitempty"`     
	CommentsURL         *string         `json:"comments_url,omitempty"`          
	CommitsURL          *string         `json:"commits_url,omitempty"`           
	CompareURL          *string         `json:"compare_url,omitempty"`           
	ContentsURL         *string         `json:"contents_url,omitempty"`          
	ContributorsURL     *string         `json:"contributors_url,omitempty"`      
	CreatedAt           *Timestamp      `json:"created_at,omitempty"`            
	DefaultBranch       *string         `json:"default_branch,omitempty"`        
	DeleteBranchOnMerge *bool           `json:"delete_branch_on_merge,omitempty"`
	DeploymentsURL      *string         `json:"deployments_url,omitempty"`       
	Description         *string         `json:"description,omitempty"`           
	Disabled            *bool           `json:"disabled,omitempty"`              
	DownloadsURL        *string         `json:"downloads_url,omitempty"`         
	EventsURL           *string         `json:"events_url,omitempty"`            
	Fork                *bool           `json:"fork,omitempty"`                  
	ForksCount          *int64          `json:"forks_count,omitempty"`           
	ForksURL            *string         `json:"forks_url,omitempty"`             
	FullName            *string         `json:"full_name,omitempty"`             
	GitCommitsURL       *string         `json:"git_commits_url,omitempty"`       
	GitRefsURL          *string         `json:"git_refs_url,omitempty"`          
	GitTagsURL          *string         `json:"git_tags_url,omitempty"`          
	GitURL              *string         `json:"git_url,omitempty"`               
	GitignoreTemplate   *string         `json:"gitignore_template,omitempty"`    
	HasDownloads        *bool           `json:"has_downloads,omitempty"`         
	HasIssues           *bool           `json:"has_issues,omitempty"`            
	HasPages            *bool           `json:"has_pages,omitempty"`             
	HasProjects         *bool           `json:"has_projects,omitempty"`          
	HasWiki             *bool           `json:"has_wiki,omitempty"`              
	Homepage            *string         `json:"homepage,omitempty"`              
	HooksURL            *string         `json:"hooks_url,omitempty"`             
	HTMLURL             *string         `json:"html_url,omitempty"`              
	ID                  *int64          `json:"id,omitempty"`                    
	IsTemplate          *bool           `json:"is_template,omitempty"`           
	IssueCommentURL     *string         `json:"issue_comment_url,omitempty"`     
	IssueEventsURL      *string         `json:"issue_events_url,omitempty"`      
	IssuesURL           *string         `json:"issues_url,omitempty"`            
	KeysURL             *string         `json:"keys_url,omitempty"`              
	LabelsURL           *string         `json:"labels_url,omitempty"`            
	Language            *string         `json:"language,omitempty"`              
	LanguagesURL        *string         `json:"languages_url,omitempty"`         
	License             *License        `json:"license,omitempty"`               
	LicenseTemplate     *string         `json:"license_template,omitempty"`      
	MasterBranch        *string         `json:"master_branch,omitempty"`         
	MergesURL           *string         `json:"merges_url,omitempty"`            
	MilestonesURL       *string         `json:"milestones_url,omitempty"`        
	MirrorURL           *string         `json:"mirror_url,omitempty"`            
	Name                *string         `json:"name,omitempty"`                  
	NetworkCount        *int64          `json:"network_count,omitempty"`         
	NodeID              *string         `json:"node_id,omitempty"`               
	NotificationsURL    *string         `json:"notifications_url,omitempty"`     
	OpenIssuesCount     *int64          `json:"open_issues_count,omitempty"`     
	Organization        *Organization   `json:"organization,omitempty"`          
	Owner               *User           `json:"owner,omitempty"`                 
	Parent              *Repository     `json:"parent,omitempty"`                
	Permissions         map[string]bool `json:"permissions,omitempty"`           
	Private             *bool           `json:"private,omitempty"`               
	PullsURL            *string         `json:"pulls_url,omitempty"`             
	PushedAt            *Timestamp      `json:"pushed_at,omitempty"`             
	ReleasesURL         *string         `json:"releases_url,omitempty"`          
	Size                *int64          `json:"size,omitempty"`                  
	Source              *Repository     `json:"source,omitempty"`                
	SSHURL              *string         `json:"ssh_url,omitempty"`               
	StargazersCount     *int64          `json:"stargazers_count,omitempty"`      
	StargazersURL       *string         `json:"stargazers_url,omitempty"`        
	StatusesURL         *string         `json:"statuses_url,omitempty"`          
	SubscribersCount    *int64          `json:"subscribers_count,omitempty"`     
	SubscribersURL      *string         `json:"subscribers_url,omitempty"`       
	SubscriptionURL     *string         `json:"subscription_url,omitempty"`      
	SvnURL              *string         `json:"svn_url,omitempty"`               
	TagsURL             *string         `json:"tags_url,omitempty"`              
	TeamID              *int64          `json:"team_id,omitempty"`               
	TeamsURL            *string         `json:"teams_url,omitempty"`             
	TemplateRepository  *Repository     `json:"template_repository,omitempty"`   
	TextMatches         []TextMatch     `json:"text_matches,omitempty"`          
	Topics              []string        `json:"topics,omitempty"`                
	TreesURL            *string         `json:"trees_url,omitempty"`             
	UpdatedAt           *Timestamp      `json:"updated_at,omitempty"`            
	URL                 *string         `json:"url,omitempty"`                   
	Visibility          *string         `json:"visibility,omitempty"`            
	WatchersCount       *int64          `json:"watchers_count,omitempty"`        
}

type CodeOfConduct struct {
	Body *string `json:"body,omitempty"`
	Key  *string `json:"key,omitempty"` 
	Name *string `json:"name,omitempty"`
	URL  *string `json:"url,omitempty"` 
}

type License struct {
	Body           *string  `json:"body,omitempty"`          
	Conditions     []string `json:"conditions,omitempty"`    
	Description    *string  `json:"description,omitempty"`   
	Featured       *bool    `json:"featured,omitempty"`      
	HTMLURL        *string  `json:"html_url,omitempty"`      
	Implementation *string  `json:"implementation,omitempty"`
	Key            *string  `json:"key,omitempty"`           
	Limitations    []string `json:"limitations,omitempty"`   
	Name           *string  `json:"name,omitempty"`          
	Permissions    []string `json:"permissions,omitempty"`   
	SpdxID         *string  `json:"spdx_id,omitempty"`       
	URL            *string  `json:"url,omitempty"`           
}

type Organization struct {
	AvatarURL                            *string `json:"avatar_url,omitempty"`                              
	BillingEmail                         *string `json:"billing_email,omitempty"`                           
	Blog                                 *string `json:"blog,omitempty"`                                    
	Collaborators                        *int64  `json:"collaborators,omitempty"`                           
	Company                              *string `json:"company,omitempty"`                                 
	CreatedAt                            *string `json:"created_at,omitempty"`                              
	DefaultRepositoryPermission          *string `json:"default_repository_permission,omitempty"`           
	DefaultRepositorySettings            *string `json:"default_repository_settings,omitempty"`             
	Description                          *string `json:"description,omitempty"`                             
	DiskUsage                            *int64  `json:"disk_usage,omitempty"`                              
	Email                                *string `json:"email,omitempty"`                                   
	EventsURL                            *string `json:"events_url,omitempty"`                              
	Followers                            *int64  `json:"followers,omitempty"`                               
	Following                            *int64  `json:"following,omitempty"`                               
	HasOrganizationProjects              *bool   `json:"has_organization_projects,omitempty"`               
	HasRepositoryProjects                *bool   `json:"has_repository_projects,omitempty"`                 
	HooksURL                             *string `json:"hooks_url,omitempty"`                               
	HTMLURL                              *string `json:"html_url,omitempty"`                                
	ID                                   *int64  `json:"id,omitempty"`                                      
	IsVerified                           *bool   `json:"is_verified,omitempty"`                             
	IssuesURL                            *string `json:"issues_url,omitempty"`                              
	Location                             *string `json:"location,omitempty"`                                
	Login                                *string `json:"login,omitempty"`                                   
	MembersAllowedRepositoryCreationType *string `json:"members_allowed_repository_creation_type,omitempty"`
	MembersCanCreateInternalRepositories *bool   `json:"members_can_create_internal_repositories,omitempty"`
	MembersCanCreatePages                *bool   `json:"members_can_create_pages,omitempty"`                
	MembersCanCreatePrivatePages         *bool   `json:"members_can_create_private_pages,omitempty"`        
	MembersCanCreatePrivateRepositories  *bool   `json:"members_can_create_private_repositories,omitempty"` 
	MembersCanCreatePublicPages          *bool   `json:"members_can_create_public_pages,omitempty"`         
	MembersCanCreatePublicRepositories   *bool   `json:"members_can_create_public_repositories,omitempty"`  
	MembersCanCreateRepositories         *bool   `json:"members_can_create_repositories,omitempty"`         
	MembersURL                           *string `json:"members_url,omitempty"`                             
	Name                                 *string `json:"name,omitempty"`                                    
	NodeID                               *string `json:"node_id,omitempty"`                                 
	OwnedPrivateRepos                    *int64  `json:"owned_private_repos,omitempty"`                     
	Plan                                 *Plan   `json:"plan,omitempty"`                                    
	PrivateGists                         *int64  `json:"private_gists,omitempty"`                           
	PublicGists                          *int64  `json:"public_gists,omitempty"`                            
	PublicMembersURL                     *string `json:"public_members_url,omitempty"`                      
	PublicRepos                          *int64  `json:"public_repos,omitempty"`                            
	ReposURL                             *string `json:"repos_url,omitempty"`                               
	TotalPrivateRepos                    *int64  `json:"total_private_repos,omitempty"`                     
	TwitterUsername                      *string `json:"twitter_username,omitempty"`                        
	TwoFactorRequirementEnabled          *bool   `json:"two_factor_requirement_enabled,omitempty"`          
	Type                                 *string `json:"type,omitempty"`                                    
	UpdatedAt                            *string `json:"updated_at,omitempty"`                              
	URL                                  *string `json:"url,omitempty"`                                     
}
