package lwtest

// Search is the accumalation of Project/Repository search results
type Search struct {

	// Search results of the projects that matched the filter keywords.
	Projects []SearchProject `json:"project,omitempty"`

	// Search results of the repositories that matched the filter keywords.
	Repositories []SearchRepository `json:"repository,omitempty"`
}

// SearchProject maps search result JSON to object
type SearchProject struct {
	ID int64 `json:"id,omitempty"`

	Name string `json:"name,omitempty"`

	Public int32 `json:"public,omitempty"`
}

// SearchRepository maps search result JSON to object
type SearchRepository struct {
	ProjectID int32 `json:"project_id,omitempty"`

	ProjectName string `json:"project_name,omitempty"`

	ProjectPublic int32 `json:"project_public,omitempty"`

	RepositoryName string `json:"repository_name,omitempty"`
}

// ProjectReq contains parameters for Project creation
type ProjectReq struct {

	// The name of the project.
	ProjectName string `json:"project_name,omitempty"`

	// The public status of the project.
	Public int32 `json:"public,omitempty"`
}

type MemberReq struct {

	// The username of the user to be added as member
	Username string `json:"username"`

	// The role ids of the roles to be given to the member
	Roles []int `json:"roles"`
}
