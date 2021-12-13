package integration

const GetConnectedUsers string = `{
		getConnectedUsers(filter:WEEK){
			totalUser
			residentUsers
			guestUsers
		}
	}`

type GetConnectedUsersResponse struct {
	ConnectedUser struct {
		TotalUsers    int `json:"totalUser"`
		ResidentUsers int `json:"residentUsers"`
		GuestUsers    int `json:"guestUsers"`
	} `json:"getConnectedUsers"`
}
