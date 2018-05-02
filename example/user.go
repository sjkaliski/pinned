package main

// User is a user.
type User struct {
	ID        uint64
	Email     string
	Name      string
	CreatedAt int64
}

// Data implements the pinned.Versionable interface.
func (u *User) Data() map[string]interface{} {
	return map[string]interface{}{
		"id":         u.ID,
		"email":      u.Email,
		"name":       u.Name,
		"created_at": u.CreatedAt,
	}
}

func userNameFieldChange(mapping map[string]interface{}) map[string]interface{} {
	mapping["full_name"] = mapping["name"]
	delete(mapping, "name")
	return mapping
}

func userCreatedAtChange(mapping map[string]interface{}) map[string]interface{} {
	delete(mapping, "created_at")
	return mapping
}
