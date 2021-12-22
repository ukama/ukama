package integration

const GetConnectedUsers string = `{
		getConnectedUsers(filter:WEEK){
			totalUser
		}
	}`

type GetConnectedUsersResponse struct {
	ConnectedUser struct {
		TotalUsers    int `json:"totalUser"`
	} `json:"getConnectedUsers"`
}
