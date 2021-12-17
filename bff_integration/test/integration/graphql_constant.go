package integration

const GetConnectedUsers string = `{
		getConnectedUsers(filter:WEEK){
			totalUser
			residentUsers
			guestUsers
		}
	}`

const GetNodesByOrg string = `{
		getNodesByOrg(orgId: "%s"){
			orgName
    		nodes {
      			nodeId
      			state
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
		ResidentUsers int `json:"residentUsers"`
		GuestUsers    int `json:"guestUsers"`
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
	NodedId     string `json:"nodeId"`
	State       string `json:"state"`
	Description string `json:"description"`
	Title       string `json:"title"`
	TotalUser   int    `json:"totalUser"`
}
