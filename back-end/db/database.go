package db

import (
	"database/sql"
	"log"

	"golang.org/x/crypto/bcrypt"
	_ "github.com/mattn/go-sqlite3"
)

func vumbleDB struct {
	db *sql.DB
}
var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
func generate(length int) string {
	// Generate a random string of the given length
	gen := make([]rune, length)
	for i := range gen {
		gen[i] = letters[rand.Intn(len(letters))]
	}
	return string(gen)
}
func newVumbleDB() *vumbleDB {
	// Open the database
	db, err := sql.Open("sqlite3", "./vumble.db")
	if err != nil {
		log.Fatal(err)
	}
	// Tables to create if they don't exist
	// USers table:
	// ID, Username, Password, Email, Firstname, Lastname, Bio, profilePic, oath_otken_id, oauth_token_secret
	// Likes table:
	// ID, FromID, toID
	// Matches table:
	// ID, FromID, toID
	// Messages table:
	// ID, FromID, toID, Message, Time

	// Oauth tokens table:
	// ID, userId, Token. Expiry


	// Create the tables if they don't exist
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, username TEXT, password TEXT, email TEXT, firstname TEXT, lastname TEXT, bio TEXT, profilepic TEXT, oauth_token_id TEXT, oauth_token_secret TEXT)")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS likes (id INTEGER PRIMARY KEY, fromid INTEGER, toid INTEGER)")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS matches (id INTEGER PRIMARY KEY, fromid INTEGER, toid INTEGER)")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS messages (id INTEGER PRIMARY KEY, fromid INTEGER, toid INTEGER, message TEXT, time TEXT)")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS oauth_tokens (id INTEGER PRIMARY KEY, userid INTEGER, token TEXT, expiry TEXT)")
	if err != nil {
		log.Fatal(err)
	}
	return &vumbleDB{db: db}
}

func (v *vumbleDB) getUser(username) (int, error) {
	// Check if the username exists
	// If it does, return the ID
	// If it doesn't, return an error

	conn := v.db
	rows, err := conn.Query("SELECT id FROM users WHERE username = ?", username)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var id int
	for rows.Next() {
		err := rows.Scan(&id)
		if err != nil {
			log.Fatal(err)
		}
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	return id, nil
}

func (v *vumbleDB) getOAuthToken(id) (string, error) {
	// Check if the user has an OAuth token
	// If they do, return it
	// If they don't, generate a new one and return it

	conn := v.db
	rows, err := conn.Query("SELECT token FROM oauth_tokens WHERE userid = ?", id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var token string
	for rows.Next() {
		err := rows.Scan(&token)
		if err != nil {
			log.Fatal(err)
		}
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	if token == "" {
		token, err = v.generateOAuthToken(id)
		if err != nil {
			log.Fatal(err)
		}
	}
	return token, nil
}

// Generates a new OAuth token and secret for the user with the given ID and returns them	
func (v *vumbleDB) generateOAuthToken(id) (string, error) {
	// Generate a token and secret
	// Store them in the database
	// Return them
    conn := v.db
	// Generate a 32 A-Za-z0-9 character token
	token := generate(32)
	// Save the token to the database
	_, err := conn.Exec("INSERT INTO oauth_tokens (userid, token) VALUES (?, ?)", id, token)
	if err != nil {
		log.Fatal(err)
	}
	return token, nil
}


// Returns the OAuth token and secret for the user with the given username
func (v *vumbleDB) createUser(username, password, email, firstname, lastname) (int, error) {
	// Check if the username is already taken
	// If it is, return an error
	// If it isn't, create the user and return the ID

	conn := v.db

	hashedPassword := bcrypt.GenerateFromPassword([]byte(password), 14);

	entry, err := conn.Exec("INSERT INTO users (username, password, email, firstname, lastname) VALUES (?, ?, ?, ?, ?)", username, hashedPassword, email, firstname, lastname)
	if err != nil {
		log.Fatal(err)
	}
	id, err := entry.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	//Generate an oauth token for the user
	token, err := v.generateOAuthToken(id)
	if err != nil {
		log.Fatal(err)
	}
	return id, nil
}

//Login
func (v *vumbleDB) login(username, password) (string, error) {
	// Check if the username exists
	// If it doesn't, return an error
	// If it does, check the password
	// If the password is correct, return the ID
	// If the password is incorrect, return an error

	conn := v.db
	rows, err := conn.Query("SELECT id, password FROM users WHERE username = ?", username)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var id int
	var hashedPassword string
	for rows.Next() {
		err := rows.Scan(&id, &hashedPassword)
		if err != nil {
			log.Fatal(err)
		}
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		// Passwords don't match
		return nil, err
	}
	token := v.getOAuthToken(id)

	return token, nil
}

// LikeUser
func (v *vumbleDB) likeUser(fromID, toID) (int, error) {
	// Check if the user has already liked the other user
	// If they have, return an error
	// If they haven't, create the like and return 1 for match and 0 for no match
	// If they are an AI user, then return always 0.
	conn := v.db
	
	// Check if the user has already been liked by the other user
	// If they have, create a match and return the ID
	rows, err := conn.Query("SELECT id FROM likes WHERE fromid = ? AND toid = ?", toID, fromID)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var id int
	for rows.Next() {
		err := rows.Scan(&id)
		if err != nil {
			log.Fatal(err)
		}
	}

	if id != 0 {
		// The user has already been liked by the other user
		// Create a match
		entry, err := conn.Exec("INSERT INTO matches (fromid, toid) VALUES (?, ?)", fromID, toID)
		if err != nil {
			log.Fatal(err)
		}
		id, err := entry.LastInsertId()
		if err != nil {
			log.Fatal(err)
		}
		return 1, nil
	}

	entry, err := conn.Exec("INSERT INTO likes (fromid, toid) VALUES (?, ?)", fromID, toID)
	if err != nil {
		log.Fatal(err)
	}
	id, err := entry.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}

	return 0, nil
}

// sendMessage
func (v *vumbleDB) sendMessage(fromID, toID, message) (int, error) {
	// Create the message and return the ID
	conn := v.db

	// Get epoch time
	epoch := time.Now().Unix()
	entry, err := conn.Exec("INSERT INTO messages (fromid, toid, message, time) VALUES (?, ?, ?, ?)", fromID, toID, message, epoch)
	if err != nil {
		log.Fatal(err)
	}
	id, err := entry.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	return id, nil
}