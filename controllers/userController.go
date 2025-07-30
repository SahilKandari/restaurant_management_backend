package controllers

import (
	"context"
	"fmt"
	"net/http"
	"restaurant-management/helpers"
	"restaurant-management/models"
	"restaurant-management/utils"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		var users []models.User

		rows, err := Db.QueryContext(ctx, "SELECT * FROM users")
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users from database", "details": err.Error()})
			return
		}
		defer rows.Close()

		for rows.Next() {
			var user models.User
			if err := rows.Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.Phone, &user.Role, &user.Token, &user.AvatarURL, &user.CreatedAt, &user.UpdatedAt); err != nil {
				c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan user data", "details": err.Error()})
				return
			}
			user.Password = "" // Clear password before sending response
			user.Token = ""    // Clear token before sending response
			users = append(users, user)
		}
		c.IndentedJSON(http.StatusOK, gin.H{"users": users})
	}
}

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		id := c.Param("user_id")
		if id == "" {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Provide user ID"})
			return
		}

		var user models.User

		if err := Db.QueryRowContext(ctx, "SELECT * FROM users WHERE id = $1", id).Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.Phone, &user.Role, &user.Token, &user.AvatarURL, &user.CreatedAt, &user.UpdatedAt); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Unable to find the user with the given indentifier", "details": err.Error()})
			return
		}

		user.Password = ""
		user.Token = ""

		c.IndentedJSON(http.StatusOK, gin.H{"message": "User fetched successfully", "user": user})
	}
}

func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		var user models.User
		if err := c.BindJSON(&user); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Provide correct details to create user", "details": err.Error()})
			return
		}

		// Validate the user struct
		if err := utils.ValidateStruct(user); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": err.Error()})
			return
		}

		// Check if username, email or phone already exists
		rows, err := Db.QueryContext(ctx, "SELECT * FROM users WHERE username = $1 OR email = $2 OR phone = $3", user.Username, user.Email, user.Phone)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to check existing user", "details": err.Error()})
			return
		}
		defer rows.Close()
		if rows.Next() {
			c.IndentedJSON(http.StatusConflict, gin.H{"error": "User with this username, email or phone already exists"})
			return
		}

		// Create hash password
		hashedPassword, err := HashPassword(user.Password)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Failed to securely hash user password", "details": err.Error()})
			return
		}

		// Generate token
		token, err := helpers.CreateToken(user.Username)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to token for the user", "details": err.Error()})
			return
		}
		// Create user in database
		if err := Db.QueryRowContext(ctx, "INSERT INTO users (username, password, email, phone, role, token, avatar_url) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, username, email, phone, role, token, avatar_url, created_at, updated_at", user.Username, hashedPassword, user.Email, user.Phone, user.Role, token, user.AvatarURL).Scan(&user.ID, &user.Username, &user.Email, &user.Phone, &user.Role, &user.Token, &user.AvatarURL, &user.CreatedAt, &user.UpdatedAt); err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user in database", "details": err.Error()})
			return
		}

		user.Password = "" // Clear password before sending response
		user.Token = token // Set the generated token

		c.IndentedJSON(http.StatusCreated, gin.H{"message": "User created successfully", "user": user})
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		var input struct {
			Identifier string `json:"identifier" validate:"required"`
			Password   string `json:"password" validate:"required,min=6,max=100"`
		}
		var user models.User

		if err := c.BindJSON(&input); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Provide identifier(username, email or phone) and password(min=6,max=100)", "details": err.Error()})
			return
		}

		if err := Db.QueryRowContext(ctx, "SELECT * FROM users WHERE username = $1 OR email = $2 OR phone = $3", input.Identifier, input.Identifier, input.Identifier).Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.Phone, &user.Role, &user.Token, &user.AvatarURL, &user.CreatedAt, &user.UpdatedAt); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Unable to find the user with the given indentifier", "identifier": input.Identifier, "details": err.Error()})
			return
		}

		//Verify password
		if isPasswordCorrect := VerifyPassword(input.Password, user.Password); isPasswordCorrect == false {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Password is incorrect"})
			return
		}

		//Generate Token
		token, err := helpers.CreateToken(user.Username)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to create token for the user", "details": err.Error()})
			return
		}

		if err := Db.QueryRowContext(ctx, "UPDATE users SET token = $1 WHERE id = $2 RETURNING id, username, email, phone, role, token, avatar_url, created_at, updated_at", token, user.ID).Scan(&user.ID, &user.Username, &user.Email, &user.Phone, &user.Role, &user.Token, &user.AvatarURL, &user.CreatedAt, &user.UpdatedAt); err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to update token for the user in database", "username": user.Username, "details": err.Error()})
			return
		}

		user.Password = ""

		c.IndentedJSON(http.StatusOK, gin.H{"message": "User login successfully", "user": user})
	}
}

func UpdateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		id := c.Param("user_id")
		if id == "" {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Provide correct details"})
			return
		}

		var user models.UserWithOldPassword
		var previousUser models.User

		if err := c.BindJSON(&user); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Provide correct details", "details": err.Error()})
			return
		}

		if err := Db.QueryRowContext(ctx, "SELECT * FROM users WHERE id = $1", id).Scan(&previousUser.ID, &previousUser.Username, &previousUser.Password, &previousUser.Email, &previousUser.Phone, &previousUser.Role, &previousUser.Token, &previousUser.AvatarURL, &previousUser.CreatedAt, &previousUser.UpdatedAt); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Failed to fetch user from database", "details": err.Error()})
			return
		}

		//Verify password
		if isPasswordCorrect := VerifyPassword(user.OldPassword, previousUser.Password); isPasswordCorrect == false {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Password is incorrect"})
			return
		}

		// Create hash password
		hashedPassword, err := HashPassword(user.Password)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Failed to securely hash user password", "details": err.Error()})
			return
		}

		//Generate Token
		token, err := helpers.CreateToken(user.Username)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to create token for the user", "details": err.Error()})
			return
		}

		if err := Db.QueryRowContext(ctx, "UPDATE users SET username = $1, password = $2, email = $3, phone = $4, role = $5, token = $6, avatar_url = $7, updated_at = CURRENT_TIMESTAMP WHERE id = $8 RETURNING id, username, email, phone, role, token, avatar_url, created_at, updated_at", user.Username, hashedPassword, user.Email, user.Phone, user.Role, token, user.AvatarURL, id).Scan(&user.ID, &user.Username, &user.Email, &user.Phone, &user.Role, &user.Token, &user.AvatarURL, &user.CreatedAt, &user.UpdatedAt); err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to update token for the user in database", "username": user.Username, "details": err.Error()})
			return
		}

		user.Password = ""

		c.IndentedJSON(http.StatusOK, gin.H{"message": "User updated successfully", "user": user})
	}
}

func DeleteUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		id := c.Param("user_id")
		if id == "" {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Provide User ID"})
			return
		}

		var confirmUserPassword models.ConfirmPassword
		var user models.User

		if err := c.BindJSON(&confirmUserPassword); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Please provide password", "details": err.Error()})
			return
		}

		if err := Db.QueryRow("SELECT * FROM users WHERE id = $1", id).Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.Phone, &user.Role, &user.Token, &user.AvatarURL, &user.CreatedAt, &user.UpdatedAt); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Failed to fetch user from database", "details": err.Error()})
			return
		}

		if isPasswordCorrect := VerifyPassword(confirmUserPassword.Password, user.Password); isPasswordCorrect == false {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Password is incorrect"})
			return
		}

		result, err := Db.ExecContext(ctx, "DELETE FROM users WHERE id = $1", id)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete order from database", "details": err.Error()})
			return
		}
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve affected rows for user", "details": err.Error()})
			return
		}
		if rowsAffected == 0 {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": "No user found with given ID"})
			return
		}

		c.IndentedJSON(http.StatusNoContent, gin.H{"message": "User delete successfully"})
	}
}

func SendPasswordResetEmail() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		var input models.PasswordResetEmail
		if err := c.BindJSON(&input); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Email is required to reset the password"})
			return
		}

		var user models.User
		if err := Db.QueryRowContext(ctx, "SELECT * FROM users WHERE email = $1", input.Email).Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.Phone, &user.Role, &user.Token, &user.AvatarURL, &user.CreatedAt, &user.UpdatedAt); err != nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": "No user found with the provided email"})
			return
		}

		// generate otp and add it to redis
		otp := utils.GenerateOTP()
		if err := utils.AddOTPtoRedis(otp, input.Email, c); err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to add OTP to Redis", "details": err.Error()})
			return
		}

		// send the otp to user through email
		err := utils.SendOTP(otp, user.Email)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to send OTP to user: " + user.Email, "details": err.Error()})
			return
		}

		c.IndentedJSON(http.StatusNoContent, gin.H{"message": "OTP send successfully"})
	}
}

func PasswordReset() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		var input models.ResetPassword

		if err := c.BindJSON(&input); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Provide email, otp and new password to reset the password", "details": err.Error()})
			return
		}
		fmt.Println("Input:", input)
		var user models.User
		if err := Db.QueryRowContext(ctx, "SELECT * FROM users WHERE email = $1", input.Email).Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.Phone, &user.Role, &user.Token, &user.AvatarURL, &user.CreatedAt, &user.UpdatedAt); err != nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": "No user found with the provided email"})
			return
		}

		// Verify OTP from Redis
		err, isInternalErr := utils.VerifyOTP(input.OTP, input.Email, c)
		if err != nil {
			var code int
			if isInternalErr {
				code = http.StatusInternalServerError
			} else {
				code = http.StatusUnauthorized
			}

			c.IndentedJSON(code, gin.H{"error": "Failed to verify OTP", "details": err.Error()})
			return
		}

		hashedPassword, err := HashPassword(input.NewPassword)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash new password", "details": err.Error()})
			return
		}

		if _, err := Db.ExecContext(ctx, "UPDATE users SET password = $1 WHERE email = $2", hashedPassword, input.Email); err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password in database", "details": err.Error()})
			return
		}

		c.IndentedJSON(http.StatusOK, gin.H{"message": "Password reset successfully"})
	}
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func VerifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
