package gitClient

type Company struct {
	Company       string `json:"company"`
	GitBranchName string `json:"git_branch_name"`
	Email         string `json:"email"`
}

type Environment struct {
	Production []Company `json:"production"`
	Test       []Company `json:"test"`
}
