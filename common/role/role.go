package role

import (
	"github.com/gofiber/fiber/v2"
)

func HasAny(c *fiber.Ctx, roles []string) bool {
	permissions := c.Locals("Permissions").([]string)
	for _, role := range roles {
		for _, permission := range permissions {
			if permission == role {
				return true
			}
		}
	}
	return false
}

func HasAll(c *fiber.Ctx, roles []string) bool {
	permissions := c.Locals("Permissions").([]string)
	roleSet := make(map[string]struct{}, len(roles))
	for _, role := range roles {
		roleSet[role] = struct{}{}
	}
	for _, permission := range permissions {
		delete(roleSet, permission)
	}
	return len(roleSet) == 0
}
