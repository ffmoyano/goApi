package model

func UserToUserResponse(user User) UserResponse {

	var roles []RoleResponse

	for _, role := range user.Roles {
		roles = append(roles, RoleToRoleResponse(role))
	}

	return UserResponse{
		ID:       user.ID,
		Name:     user.Name,
		Username: user.Username,
		Email:    user.Email,
		Roles:    roles,
		Active:   user.Active,
	}
}

func RoleToRoleResponse(role Role) RoleResponse {
	return RoleResponse{
		ID:   role.ID,
		Name: role.Name,
	}
}
