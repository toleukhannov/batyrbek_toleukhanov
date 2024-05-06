package controllers

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/toleukhannov/batyrbek_toleukhanov/models"
)

func HashPassword(password string) string {

}

func VerifyPassword(userPassword string, givenPassword string) (bool, string) {

}
func SignUp(db *sql.DB, validate *validator.Validate) gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		// Проверка наличия пользователя с таким же email
		var count int
		err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM users WHERE email = $1", user.Email).Scan(&count)
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists"})
			return
		}

		// Проверка наличия пользователя с таким же номером телефона
		err = db.QueryRowContext(ctx, "SELECT COUNT(*) FROM users WHERE phone = $1", user.Phone).Scan(&count)
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Phone is already in use"})
			return
		}

		// Хеширование пароля
		hashedPassword := HashPassword(*user.Password)
		user.Password = &hashedPassword

		// Установка времени создания и обновления пользователя
		currentTime := time.Now().Format(time.RFC3339)
		user.Created_At, _ = time.Parse(time.RFC3339, currentTime)
		user.Updated_At, _ = time.Parse(time.RFC3339, currentTime)

		// Генерация токенов
		token, refreshToken, _ := generate.TokenGenerator(*user.Email, *user.First_Name, *user.Last_Name, user.User_ID)
		user.Token = &token
		user.Refresh_Token = &refreshToken

		// Инициализация пустых полей
		user.UserCart = make([]models.ProductUser, 0)
		user.Address_Details = make([]models.Address, 0)
		user.Order_Status = make([]models.Order, 0)

		// Вставка пользователя в базу данных
		_, err = db.ExecContext(ctx, "INSERT INTO users (user_id, first_name, last_name, email, phone, password, created_at, updated_at, token, refresh_token) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)",
			user.User_ID, user.First_Name, user.Last_Name, user.Email, user.Phone, user.Password, user.Created_At, user.Updated_At, user.Token, user.Refresh_Token)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}

		c.JSON(http.StatusCreated, "Successfully Signed Up!!")
	}
}
func Login() gin.HandlerFunc {
    return func(c *gin.Context) {
        var user models.User
        var founduser models.User

        if err := c.BindJSON(&user); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        // Establish a connection to PostgreSQL
        db, err := sql.Open("postgres", "postgres://user:password@localhost/dbname?sslmode=disable")
        if err != nil {
            log.Fatal(err)
        }
        defer db.Close()

        // Query PostgreSQL for the user
        row := db.QueryRow("SELECT email, password, first_name, last_name, user_id FROM users WHERE email = $1", user.Email)
        err = row.Scan(&founduser.Email, &founduser.Password, &founduser.First_Name, &founduser.Last_Name, &founduser.User_ID)
        if err != nil {
            if err == sql.ErrNoRows {
                c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
                return
            }
            c.JSON(http.StatusInternalServerError, gin.H{"error": "An error occurred"})
            return
        }

        // Verify password
        passwordIsValid, msg := VerifyPassword(*user.Password, *founduser.Password)
        if !passwordIsValid {
            c.JSON(http.StatusUnauthorized, gin.H{"error": msg})
            return
        }

        // Generate tokens
        token, refreshToken, _ := generate.TokenGenerator(*founduser.Email, *founduser.First_Name, *founduser.Last_Name, founduser.User_ID)
        generate.UpdateAllTokens(token, refreshToken, founduser.User_ID)

        c.JSON(http.StatusOK, founduser)
    }
}


func ProductViewerAdmin() gin.HandlerFunc {

}

func SearchProduct() gin.HandlerFunc {

}

func SearchProductByQuery() gin.HandlerFunc {

}

func AddProduct() gin.HandlerFunc {

}
