package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gupta29470/golang_sql_crud_without_orm/database"
	"github.com/gupta29470/golang_sql_crud_without_orm/models"
)

func CreateUser() gin.HandlerFunc {
	return func(context *gin.Context) {
		var user models.User

		bindJsonError := context.BindJSON(&user)
		if bindJsonError != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Data"})
			return
		}

		stmt, stmtError := database.DB.Prepare("INSERT INTO users (firstName, lastName, email, createdAt) VALUES (?, ?, ?, ?)")
		if stmtError != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": stmtError.Error()})
			return
		}

		defer stmt.Close()

		user.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		result, resultError := stmt.Exec(user.FirstName, user.LastName, user.Email, user.CreatedAt)
		if resultError != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": stmtError.Error()})
			return
		}

		id, _ := result.LastInsertId()
		user.ID = int(id)

		context.JSON(http.StatusOK, user)
	}
}

func GetAllUsers() gin.HandlerFunc {
	return func(context *gin.Context) {
		rows, queryError := database.DB.Query("SELECT id, firstName, lastName, email, createdAt, updatedAt from users WHERE deletedAt is NULL")
		if queryError != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": queryError.Error()})
			return
		}

		defer rows.Close()

		var users []models.User

		for rows.Next() {
			var user models.User
			var createdAtBytes, updatedAtBytes []uint8

			scanError := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &createdAtBytes, &updatedAtBytes)
			if scanError != nil {
				context.JSON(http.StatusInternalServerError, gin.H{"error": scanError.Error()})
				return
			}

			createdAtStr := string(createdAtBytes)
			updatedAtStr := string(updatedAtBytes)

			var timeParseError error

			user.CreatedAt, timeParseError = time.Parse("2006-01-02 15:04:05", createdAtStr)
			if timeParseError != nil {
				context.JSON(http.StatusInternalServerError, gin.H{"error": timeParseError.Error()})
				return
			}

			if updatedAtStr != "" {
				user.UpdatedAt, timeParseError = time.Parse("2006-01-02 15:04:05", updatedAtStr)
				if timeParseError != nil {
					context.JSON(http.StatusInternalServerError, gin.H{"error": timeParseError.Error()})
					return
				}
			}

			user.CreatedAt = user.CreatedAt.UTC()
			user.UpdatedAt = user.UpdatedAt.UTC()

			users = append(users, user)
		}

		context.JSON(http.StatusOK, users)
	}
}

func GetUser() gin.HandlerFunc {
	return func(context *gin.Context) {
		id := context.Param("id")

		var user models.User

		query := database.DB.QueryRow("SELECT id, firstName, lastName, email, createdAt, updatedAt FROM users WHERE id=? AND deletedAt IS NULL", id)

		var createdAtBytes []uint8
		var updatedAtBytes []uint8

		scanError := query.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &createdAtBytes, &updatedAtBytes)
		if scanError != nil {
			context.JSON(http.StatusNotFound, gin.H{"error": scanError.Error()})
			return
		}

		createdAtStr := string(createdAtBytes)
		updatedAtStr := string(updatedAtBytes)

		var timeParseError error

		user.CreatedAt, timeParseError = time.Parse("2006-01-02 15:04:05", createdAtStr)
		if timeParseError != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": timeParseError.Error()})
			return
		}

		if updatedAtStr != "" {
			user.UpdatedAt, timeParseError = time.Parse("2006-01-02 15:04:05", updatedAtStr)
			if timeParseError != nil {
				context.JSON(http.StatusInternalServerError, gin.H{"error": timeParseError.Error()})
				return
			}
		}

		user.CreatedAt = user.CreatedAt.UTC()
		user.UpdatedAt = user.UpdatedAt.UTC()

		context.JSON(http.StatusOK, user)
	}
}

func UpdateUser() gin.HandlerFunc {
	return func(context *gin.Context) {
		var updatedUser models.User

		bindJsonError := context.ShouldBindJSON(&updatedUser)
		if bindJsonError != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Data"})
			return
		}

		id := context.Param("id")
		var fetchedUser models.User

		query := database.DB.QueryRow("SELECT firstName, lastName, email, createdAt FROM users WHERE id=? and deletedAt IS NULL", id)

		var createdAtBytes []uint8

		scanError := query.Scan(&fetchedUser.FirstName, &fetchedUser.LastName, &fetchedUser.Email, &createdAtBytes)
		if scanError != nil {
			context.JSON(http.StatusNotFound, gin.H{"error": scanError.Error()})
			return
		}

		createdAtString := string(createdAtBytes)
		var timeParseError error
		fetchedUser.CreatedAt, timeParseError = time.Parse("2006-01-02 15:04:05", createdAtString)
		if timeParseError != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": timeParseError.Error()})
			return
		}

		if updatedUser.FirstName != "" {
			fetchedUser.FirstName = updatedUser.FirstName
		}

		if updatedUser.LastName != "" {
			fetchedUser.LastName = updatedUser.LastName
		}

		if updatedUser.Email != "" {
			fetchedUser.Email = updatedUser.Email
		}

		stmt, stmtError := database.DB.Prepare("UPDATE users SET firstName=?, lastName=?, email=?, createdAt=?, updatedAt=? WHERE id=?")
		if stmtError != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": stmtError.Error()})
			return
		}

		defer stmt.Close()

		fetchedUser.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		result, resultError := stmt.Exec(fetchedUser.FirstName, fetchedUser.LastName, fetchedUser.Email, fetchedUser.CreatedAt, fetchedUser.UpdatedAt, id)
		if resultError != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": stmtError.Error()})
			return
		}

		context.JSON(http.StatusOK, result)
	}
}

func DeleteUser() gin.HandlerFunc {
	return func(context *gin.Context) {
		id := context.Param("id")
		var createdAtBytes, updatedAtBytes []uint8

		query := database.DB.QueryRow("SELECT firstName, lastName, email, createdAt, updatedAt from users WHERE id=?", id)
		var fetchedUser models.User
		scanError := query.Scan(&fetchedUser.FirstName, &fetchedUser.LastName, &fetchedUser.Email, &createdAtBytes, &updatedAtBytes)
		if scanError != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": scanError.Error()})
			return
		}

		createdAtStr := string(createdAtBytes)
		updatedAtStr := string(updatedAtBytes)

		var timeParseError error
		fetchedUser.CreatedAt, timeParseError = time.Parse("2006-01-02 15:04:05", createdAtStr)
		if timeParseError != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": timeParseError.Error()})
			return
		}

		if updatedAtStr != "" {
			fetchedUser.UpdatedAt, timeParseError = time.Parse("2006-01-02 15:04:05", updatedAtStr)
			if timeParseError != nil {
				context.JSON(http.StatusInternalServerError, gin.H{"error": timeParseError.Error()})
				return
			}
		}

		fetchedUser.DeletedAt, timeParseError = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		if timeParseError != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": timeParseError.Error()})
			return
		}

		stmt, stmtError := database.DB.Prepare("UPDATE users SET firstName=?, lastName=?, email=?, createdAt=?, updatedAt=?, deletedAt=? WHERE id=?")
		if stmtError != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": stmtError.Error()})
			return
		}

		result, resultError := stmt.Exec(fetchedUser.FirstName, fetchedUser.LastName, fetchedUser.Email, fetchedUser.CreatedAt, fetchedUser.UpdatedAt, fetchedUser.DeletedAt, id)
		if resultError != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": resultError.Error()})
			return
		}

		context.JSON(http.StatusOK, result)
	}
}
