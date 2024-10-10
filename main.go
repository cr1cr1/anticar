package main

import (
	"database/sql"
	"flag"
	"fmt"

	_ "github.com/mattn/go-sqlite3"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/rs/zerolog/log"
)

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (u *User) Login() string {
	return "User login"
}

func main() {
	// Flags
	preloadDB := flag.Bool("preload", false, "Preload the database with mock users")
	flag.Parse()

	// Initialize a new Fiber app
	app := fiber.New()

	// Initialize SQLite database
	db, err := sql.Open("sqlite3", "./users.db")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to open database")
	}
	defer db.Close()

	// Create users table if it doesn't exist
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT,
        email TEXT,
        password TEXT
    )`)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create users table")
	}

	if *preloadDB {
		// Insert mock users
		insertMockUsers(db)
	}

	// Home page route
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to eBay-like application!")
	})

	// Product listing route
	app.Get("/products", func(c *fiber.Ctx) error {
		// In a real application, you would fetch products from a database
		products := []string{"Product 1", "Product 2", "Product 3"}
		return c.JSON(products)
	})

	// Product details route
	app.Get("/products/:id", func(c *fiber.Ctx) error {
		// In a real application, you would fetch product details from a database
		productID := c.Params("id")
		product := map[string]string{
			"id":    productID,
			"name":  "Product " + productID,
			"price": "$100",
		}
		return c.JSON(product)
	})

	// User login route
	app.Get("/login", func(c *fiber.Ctx) error {
		return adaptor.HTTPHandler(templ.Handler(LoginPage()))(c)
	})

	// POST /login route for handling form submission
	app.Post("/login", func(c *fiber.Ctx) error {
		user, err := getUserFromDb(db, c.FormValue("username"), c.FormValue("password"))
		if err != nil {
			return adaptor.HTTPHandler(templ.Handler(LoginResult(false, "Invalid credentials")))(c)
		}

		return adaptor.HTTPHandler(templ.Handler(LoginResult(true, user.Email)))(c)
	})

	// User registration route
	app.Get("/register", func(c *fiber.Ctx) error {
		// In a real application, you would handle user registration

		return c.SendString("Here will be a button to register")
	})

	// User registration route
	app.Post("/register", func(c *fiber.Ctx) error {
		// In a real application, you would handle user registration

		return c.SendString("User registration")
	})

	// Start the server on port 3000
	log.Fatal().Err(app.Listen(":3000")).Msg("failed to start server")
}

// getUserFromDB fetches a user from the database by ID
func getUserFromDB(db *sql.DB, userID int) (*User, error) {
	row := db.QueryRow("SELECT id, name, email, password FROM users WHERE id = ?", userID)

	var user User
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// insertMockUsers inserts 10 mock users into the database
func insertMockUsers(db *sql.DB) {
	log.Info().Msg("Inserting mock users into the database")
	for i := 1; i <= 10; i++ {
		id := fmt.Sprintf("%d", i)
		_, err := db.Exec("INSERT INTO users (name, email, password) VALUES (?, ?, ?)",
			"User"+id, "user"+id+"@example.com", "password"+id)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to insert mock users")
		}
	}
}

// getUserFromDb fetches a user from the database by email
func getUserFromDb(db *sql.DB, username, password string) (*User, error) {
	row := db.QueryRow("SELECT id, name, email, password FROM users WHERE name = ? AND password = ?", username, password)

	var user User
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
