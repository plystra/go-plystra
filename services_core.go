package plystra

import "context"

func (s UsersService) List(ctx context.Context, query Query) ([]Map, error) {
	return s.client.getList(ctx, "/api/v1/users", query)
}
func (s UsersService) Create(ctx context.Context, input Map) (Map, error) {
	return s.client.postMap(ctx, "/api/v1/users", input)
}
func (s UsersService) Get(ctx context.Context, userID string) (Map, error) {
	return s.client.getMap(ctx, "/api/v1/users/"+esc(userID), nil)
}
func (s UsersService) Update(ctx context.Context, userID string, input Map) (Map, error) {
	return s.client.patchMap(ctx, "/api/v1/users/"+esc(userID), input)
}
func (s UsersService) Disable(ctx context.Context, userID string, input Map) (Map, error) {
	return s.client.postMap(ctx, "/api/v1/users/"+esc(userID)+"/disable", input)
}
func (s UsersService) Restore(ctx context.Context, userID string, input Map) (Map, error) {
	return s.client.postMap(ctx, "/api/v1/users/"+esc(userID)+"/restore", input)
}

func (s SpacesService) List(ctx context.Context, query Query) ([]Map, error) {
	return s.client.getList(ctx, "/api/v1/spaces", query)
}
func (s SpacesService) Create(ctx context.Context, input Map) (Map, error) {
	return s.client.postMap(ctx, "/api/v1/spaces", input)
}
func (s SpacesService) Get(ctx context.Context, spaceID string) (Map, error) {
	return s.client.getMap(ctx, "/api/v1/spaces/"+esc(spaceID), nil)
}
func (s SpacesService) Update(ctx context.Context, spaceID string, input Map) (Map, error) {
	return s.client.patchMap(ctx, "/api/v1/spaces/"+esc(spaceID), input)
}
func (s SpacesService) Disable(ctx context.Context, spaceID string, input Map) (Map, error) {
	return s.client.postMap(ctx, "/api/v1/spaces/"+esc(spaceID)+"/disable", input)
}
func (s SpacesService) Restore(ctx context.Context, spaceID string, input Map) (Map, error) {
	return s.client.postMap(ctx, "/api/v1/spaces/"+esc(spaceID)+"/restore", input)
}
func (s SpacesService) Groups(ctx context.Context, spaceID string, query Query) ([]Map, error) {
	return s.client.getList(ctx, "/api/v1/spaces/"+esc(spaceID)+"/groups", query)
}
func (s SpacesService) GroupTree(ctx context.Context, spaceID string) ([]Map, error) {
	return s.client.getList(ctx, "/api/v1/spaces/"+esc(spaceID)+"/groups/tree", nil)
}
func (s SpacesService) Members(ctx context.Context, spaceID string, query Query) ([]Map, error) {
	return s.client.getList(ctx, "/api/v1/spaces/"+esc(spaceID)+"/members", query)
}
func (s SpacesService) UserMembers(ctx context.Context, spaceID string, query Query) ([]Map, error) {
	return s.client.getList(ctx, "/api/v1/spaces/"+esc(spaceID)+"/user-members", query)
}
func (s SpacesService) Roles(ctx context.Context, spaceID string, query Query) ([]Map, error) {
	return s.client.getList(ctx, "/api/v1/spaces/"+esc(spaceID)+"/roles", query)
}
func (s SpacesService) MemberRoleGrants(ctx context.Context, spaceID string, query Query) ([]Map, error) {
	return s.client.getList(ctx, "/api/v1/spaces/"+esc(spaceID)+"/member-role-grants", query)
}
func (s SpacesService) Resources(ctx context.Context, spaceID string, query Query) ([]Map, error) {
	return s.client.getList(ctx, "/api/v1/spaces/"+esc(spaceID)+"/resources", query)
}
func (s SpacesService) AuditLogs(ctx context.Context, spaceID string, query Query) ([]Map, error) {
	return s.client.getList(ctx, "/api/v1/spaces/"+esc(spaceID)+"/audit-logs", query)
}

func (s GroupsService) Create(ctx context.Context, spaceID string, input Map) (Map, error) {
	return s.client.postMap(ctx, "/api/v1/spaces/"+esc(spaceID)+"/groups", input)
}
func (s GroupsService) Get(ctx context.Context, groupID string) (Map, error) {
	return s.client.getMap(ctx, "/api/v1/groups/"+esc(groupID), nil)
}
func (s GroupsService) GetInSpace(ctx context.Context, spaceID, groupID string) (Map, error) {
	return s.client.getMap(ctx, "/api/v1/spaces/"+esc(spaceID)+"/groups/"+esc(groupID), nil)
}
func (s GroupsService) Update(ctx context.Context, spaceID, groupID string, input Map) (Map, error) {
	return s.client.patchMap(ctx, "/api/v1/spaces/"+esc(spaceID)+"/groups/"+esc(groupID), input)
}
func (s GroupsService) Disable(ctx context.Context, spaceID, groupID string, input Map) (Map, error) {
	return s.client.postMap(ctx, "/api/v1/spaces/"+esc(spaceID)+"/groups/"+esc(groupID)+"/disable", input)
}

func (s MembersService) Create(ctx context.Context, spaceID string, input Map) (Map, error) {
	return s.client.postMap(ctx, "/api/v1/spaces/"+esc(spaceID)+"/members", input)
}
func (s MembersService) Get(ctx context.Context, memberID string) (Map, error) {
	return s.client.getMap(ctx, "/api/v1/members/"+esc(memberID), nil)
}
func (s MembersService) GetInSpace(ctx context.Context, spaceID, memberID string) (Map, error) {
	return s.client.getMap(ctx, "/api/v1/spaces/"+esc(spaceID)+"/members/"+esc(memberID), nil)
}
func (s MembersService) Update(ctx context.Context, spaceID, memberID string, input Map) (Map, error) {
	return s.client.patchMap(ctx, "/api/v1/spaces/"+esc(spaceID)+"/members/"+esc(memberID), input)
}
func (s MembersService) Disable(ctx context.Context, spaceID, memberID string, input Map) (Map, error) {
	return s.client.postMap(ctx, "/api/v1/spaces/"+esc(spaceID)+"/members/"+esc(memberID)+"/disable", input)
}

func (s UserMembersService) Create(ctx context.Context, spaceID string, input Map) (Map, error) {
	return s.client.postMap(ctx, "/api/v1/spaces/"+esc(spaceID)+"/user-members", input)
}
func (s UserMembersService) Get(ctx context.Context, userMemberID string) (Map, error) {
	return s.client.getMap(ctx, "/api/v1/user-members/"+esc(userMemberID), nil)
}
func (s UserMembersService) GetInSpace(ctx context.Context, spaceID, userMemberID string) (Map, error) {
	return s.client.getMap(ctx, "/api/v1/spaces/"+esc(spaceID)+"/user-members/"+esc(userMemberID), nil)
}
func (s UserMembersService) Update(ctx context.Context, spaceID, userMemberID string, input Map) (Map, error) {
	return s.client.patchMap(ctx, "/api/v1/spaces/"+esc(spaceID)+"/user-members/"+esc(userMemberID), input)
}
func (s UserMembersService) Revoke(ctx context.Context, spaceID, userMemberID string, input Map) (Map, error) {
	return s.client.postMap(ctx, "/api/v1/spaces/"+esc(spaceID)+"/user-members/"+esc(userMemberID)+"/revoke", input)
}

func (s RolesService) Create(ctx context.Context, spaceID string, input Map) (Map, error) {
	return s.client.postMap(ctx, "/api/v1/spaces/"+esc(spaceID)+"/roles", input)
}
func (s RolesService) Get(ctx context.Context, roleID string) (Map, error) {
	return s.client.getMap(ctx, "/api/v1/roles/"+esc(roleID), nil)
}
func (s RolesService) GetInSpace(ctx context.Context, spaceID, roleID string) (Map, error) {
	return s.client.getMap(ctx, "/api/v1/spaces/"+esc(spaceID)+"/roles/"+esc(roleID), nil)
}
func (s RolesService) Update(ctx context.Context, spaceID, roleID string, input Map) (Map, error) {
	return s.client.patchMap(ctx, "/api/v1/spaces/"+esc(spaceID)+"/roles/"+esc(roleID), input)
}
func (s RolesService) Disable(ctx context.Context, spaceID, roleID string, input Map) (Map, error) {
	return s.client.postMap(ctx, "/api/v1/spaces/"+esc(spaceID)+"/roles/"+esc(roleID)+"/disable", input)
}

func (s MemberRolesService) List(ctx context.Context, spaceID string, query Query) ([]Map, error) {
	return s.client.getList(ctx, "/api/v1/spaces/"+esc(spaceID)+"/member-role-grants", query)
}
func (s MemberRolesService) Create(ctx context.Context, spaceID string, input Map) (Map, error) {
	return s.client.postMap(ctx, "/api/v1/spaces/"+esc(spaceID)+"/member-role-grants", input)
}
func (s MemberRolesService) Get(ctx context.Context, spaceID, memberRoleID string) (Map, error) {
	return s.client.getMap(ctx, "/api/v1/spaces/"+esc(spaceID)+"/member-role-grants/"+esc(memberRoleID), nil)
}
func (s MemberRolesService) Revoke(ctx context.Context, spaceID, memberRoleID string, input Map) (Map, error) {
	return s.client.postMap(ctx, "/api/v1/spaces/"+esc(spaceID)+"/member-role-grants/"+esc(memberRoleID)+"/revoke", input)
}

func (s PermissionsService) List(ctx context.Context, query Query) ([]Map, error) {
	return s.client.getList(ctx, "/api/v1/permissions", query)
}
func (s PermissionsService) Create(ctx context.Context, input Map) (Map, error) {
	return s.client.postMap(ctx, "/api/v1/permissions", input)
}
func (s PermissionsService) Get(ctx context.Context, permissionID string) (Map, error) {
	return s.client.getMap(ctx, "/api/v1/permissions/"+esc(permissionID), nil)
}
func (s PermissionsService) Update(ctx context.Context, permissionID string, input Map) (Map, error) {
	return s.client.patchMap(ctx, "/api/v1/permissions/"+esc(permissionID), input)
}
func (s PermissionsService) Disable(ctx context.Context, permissionID string, input Map) (Map, error) {
	return s.client.postMap(ctx, "/api/v1/permissions/"+esc(permissionID)+"/disable", input)
}

func (s RolePermissionsService) List(ctx context.Context, query Query) ([]Map, error) {
	return s.client.getList(ctx, "/api/v1/role-permissions", query)
}
func (s RolePermissionsService) Create(ctx context.Context, input Map) (Map, error) {
	return s.client.postMap(ctx, "/api/v1/role-permissions", input)
}
func (s RolePermissionsService) Get(ctx context.Context, rolePermissionID string) (Map, error) {
	return s.client.getMap(ctx, "/api/v1/role-permissions/"+esc(rolePermissionID), nil)
}
func (s RolePermissionsService) Revoke(ctx context.Context, rolePermissionID string, input Map) (Map, error) {
	return s.client.deleteMap(ctx, "/api/v1/role-permissions/"+esc(rolePermissionID), input)
}

func (s ResourceTypesService) List(ctx context.Context, query Query) ([]Map, error) {
	return s.client.getList(ctx, "/api/v1/resource-types", query)
}
func (s ResourceTypesService) Upsert(ctx context.Context, input Map) (Map, error) {
	return s.client.postMap(ctx, "/api/v1/resource-types", input)
}
func (s ResourceTypesService) Get(ctx context.Context, key string) (Map, error) {
	return s.client.getMap(ctx, "/api/v1/resource-types/"+esc(key), nil)
}
func (s ResourceTypesService) Actions(ctx context.Context, key string, query Query) ([]Map, error) {
	return s.client.getList(ctx, "/api/v1/resource-types/"+esc(key)+"/actions", query)
}
func (s ResourceTypesService) UpsertAction(ctx context.Context, key string, input Map) (Map, error) {
	return s.client.postMap(ctx, "/api/v1/resource-types/"+esc(key)+"/actions", input)
}
func (s ResourceTypesService) Mapping(ctx context.Context, key string) (Map, error) {
	return s.client.getMap(ctx, "/api/v1/resource-types/"+esc(key)+"/mapping", nil)
}
func (s ResourceTypesService) UpsertMapping(ctx context.Context, key string, input Map) (Map, error) {
	return s.client.postMap(ctx, "/api/v1/resource-types/"+esc(key)+"/mapping", input)
}

func (s ResourcesService) List(ctx context.Context, query Query) ([]Map, error) {
	return s.client.getList(ctx, "/api/v1/resources", query)
}
func (s ResourcesService) Create(ctx context.Context, input Map) (Map, error) {
	return s.client.postMap(ctx, "/api/v1/resources", input)
}
func (s ResourcesService) CreateInSpace(ctx context.Context, spaceID string, input Map) (Map, error) {
	return s.client.postMap(ctx, "/api/v1/spaces/"+esc(spaceID)+"/resources", input)
}
func (s ResourcesService) Get(ctx context.Context, resourceType, resourceID string) (Map, error) {
	return s.client.getMap(ctx, "/api/v1/resources/"+esc(resourceType)+"/"+esc(resourceID), nil)
}
func (s ResourcesService) GetInSpace(ctx context.Context, spaceID, resourceID string) (Map, error) {
	return s.client.getMap(ctx, "/api/v1/spaces/"+esc(spaceID)+"/resources/"+esc(resourceID), nil)
}
func (s ResourcesService) Update(ctx context.Context, spaceID, resourceID string, input Map) (Map, error) {
	return s.client.patchMap(ctx, "/api/v1/spaces/"+esc(spaceID)+"/resources/"+esc(resourceID), input)
}
func (s ResourcesService) Archive(ctx context.Context, spaceID, resourceID string, input Map) (Map, error) {
	return s.client.postMap(ctx, "/api/v1/spaces/"+esc(spaceID)+"/resources/"+esc(resourceID)+"/archive", input)
}
