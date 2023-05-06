package schema

type NewGroupReq struct {
	Public  bool     `json:"public"`
	Name    string   `json:"name" binding:"required"`
	Members []string `json:"members" binding:"required"`
}

type GetGroupByIDReq struct {
	ID string `uri:"id" binding:"required"`
}
