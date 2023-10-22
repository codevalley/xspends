package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"`
}

var db *sql.DB
var jwtKey = []byte(os.Getenv("JWT_KEY"))

type Claims struct {
	UserID string `json:"user_id"`
	jwt.StandardClaims
}

func init_db() {
	var err error
	db, err = sql.Open("mysql", "root:@tcp(tidb-cluster-tidb.tidb-cluster.svc.cluster.local:4000)/xpends")
	if err != nil {
		log.Fatal(err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
}

func generateToken(user User) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID: user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func register(c *gin.Context) {
	var newUser User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), 14)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not hash password"})
		return
	}
	newUser.Password = string(hashedPassword)
	newUser.ID = uuid.New().String()

	_, err = db.Exec("INSERT INTO users (id, username, password) VALUES (?, ?, ?)", newUser.ID, newUser.Username, newUser.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not insert user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": newUser.ID, "username": newUser.Username})
}

func login(c *gin.Context) {
	var creds User
	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	row := db.QueryRow("SELECT id, password FROM users WHERE username=?", creds.Username)
	var foundUser User
	err := row.Scan(&foundUser.ID, &foundUser.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not retrieve user"})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(creds.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	token, err := generateToken(foundUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}
