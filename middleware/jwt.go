package middleware

import (
	"booking-klinik/utils"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: implement authentication middleware
		tokenString := c.GetHeader("Authorization")
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		claims, err := utils.VerifyJWT(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid token"})
			return
		}

		fmt.Println(claims.Email, claims.UserID, claims.Role)
		c.Set("email", claims.Email)
		c.Set("userID", claims.UserID)
		c.Set("role", claims.Role)
		c.Next()
	}
}

func RoleCheckMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		fmt.Println(role)
		if !exists {
			c.AbortWithStatusJSON(403, gin.H{"error": "Role not found"})
		}

		roleStr := role.(string)
		for _, allowedRole := range allowedRoles {
			if strings.EqualFold(roleStr, allowedRole) {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(403, gin.H{"error": "You don't have permission to access this resource"})
	}
}
