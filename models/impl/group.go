package impl

type GroupModel struct {
	TableGroups       string
	ColumnGroupID     string
	ColumnOwnerID     string
	ColumnScopeID     string
	ColumnGroupName   string
	ColumnDescription string
	ColumnIcon        string
	ColumnStatus      string
	ColumnCreatedAt   string
	ColumnUpdatedAt   string
}

func NewGroupModel() *GroupModel {
	return &GroupModel{
		TableGroups:       "groups",
		ColumnGroupID:     "group_id",
		ColumnOwnerID:     "owner_id",
		ColumnScopeID:     "scope_id",
		ColumnGroupName:   "group_name",
		ColumnDescription: "description",
		ColumnIcon:        "icon",
		ColumnStatus:      "status",
		ColumnCreatedAt:   "created_at",
		ColumnUpdatedAt:   "updated_at",
	}
}

// Additional methods for InsertGroup, UpdateGroup, DeleteGroup, GetGroupByID, GetAllGroups, etc.
