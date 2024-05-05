// package main

// import (
// 	"database/sql"
// 	"fmt"
// 	"net/http"

// 	_ "github.com/go-sql-driver/mysql"
// )

// type Products struct {
// 	//product_id	name	description	price	category_id	 brand_id	image_url
// 	Product_id          uint16
// 	Product_name        string
// 	Product_description string
// 	Product_price       float64
// 	Product_categoryId  uint16
// 	Product_brandId     uint16
// 	Product_img         string
// }

// func main() {
// 	handleFunc()
// }
// func save_products(w http.ResponseWriter, r *http.Request) {
// 	product_name := r.FormValue("product_name")
// 	product_description := r.FormValue("product_description")
// 	product_price := r.FormValue("product_price")
// 	product_categoryId := r.FormValue("product_categoryId")
// 	product_brandId := r.FormValue("product_brandId")
// 	product_img := r.FormValue("product_img")

// 	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/clothes_store")
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer db.Close()

// 	insert, err := db.Query(fmt.Sprintf("INSERT INTO `products` (`name`, `description`, `price`, `category_id`, `brand_id`, `image_url`) VALUES('%s', '%s', '%g', '%d', '%d', '%s')", product_name, product_description, product_price, product_categoryId, product_brandId, product_img))
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer insert.Close()
// }

// func save_article(w http.ResponseWriter, r *http.Request) {

// 	username := r.FormValue("username")
// 	usersurname := r.FormValue("usersurname")
// 	iin := r.FormValue("iin")
// 	pass := r.FormValue("pass")
// 	phone := r.FormValue("phone")

// 	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/golang")
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer db.Close()

// 	insert, err := db.Query(fmt.Sprintf("INSERT INTO `articles` (`username`, `usersurname`, `iin`, `pass`, `phone`) VALUES('%s', '%s', '%d', '%s', '%d')", username, usersurname, iin, pass, phone))
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer insert.Close()

// 	http.Redirect(w, r, "/", http.StatusSeeOther)
// }
// func handleFunc() {
// 	http.HandleFunc("/products", products)
// 	http.ListenAndServe(":8080", nil)
// }

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

// User struct to represent a user
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// Secret key for signing JWT tokens
var jwtKey = []byte("secret_key")

// Credentials struct to receive username and password from the client
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Claims struct to define the JWT claims
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// Hardcoded user for demonstration purposes (replace with database lookup)
var user = User{
	ID:       1,
	Username: "user",
	Password: "password",
}

// Handler for user login
func loginHandler(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Validate user credentials (replace with database lookup)
	if creds.Username == user.Username && creds.Password == user.Password {
		expirationTime := time.Now().Add(5 * time.Minute)

		// Create JWT token
		claims := &Claims{
			Username: creds.Username,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: expirationTime.Unix(),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(jwtKey)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Send the token back to the client
		w.Write([]byte(tokenString))
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
}

// Middleware to validate JWT token
func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		tokenString := authHeader[len("Bearer "):]
		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(*Claims); ok && token.Valid {
			fmt.Printf("User %s authenticated\n", claims.Username)
			next.ServeHTTP(w, r)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	})
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/login", loginHandler).Methods("POST")
	r.HandleFunc("/protected", authMiddleware(protectedHandler)).Methods("GET")
	http.Handle("/", r)
	fmt.Println("Server started at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

// Handler for protected route
func protectedHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("You have accessed the protected route!"))
}
