package integration

const GetConnectedUsers string = `{
		getConnectedUsers(filter:WEEK){
			totalUser
		}
	}`

const GetNodesByOrg string = `{
		getNodesByOrg(orgId: "%s"){
			orgName
    		nodes {
      			id
      			description
      			title
    			totalUser
    		}
    		activeNodes
    		totalNodes
		}
	}`

type GetConnectedUsersResponse struct {
	ConnectedUser struct {
		TotalUsers    int `json:"totalUser"`
	} `json:"getConnectedUsers"`
}

type GetNodesByOrgResponse struct {
	GetNodesByOrg struct {
		OrgName     string  `json:"orgName"`
		ActiveNodes int     `json:"activeNodes"`
		TotalNodes  int     `json:"totalNodes"`
		Nodes       []Nodes `json:"nodes"`
	} `json:"getNodesByOrg"`
}

type Nodes struct {
	Id          string `json:"id"`
	Description string `json:"description"`
	Title       string `json:"title"`
	TotalUser   int    `json:"totalUser"`
}
