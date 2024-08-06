package main

import (
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"golang.org/x/crypto/ssh"

	"crypto/sha256"
	"database/sql"
	"encoding/hex"

	_ "github.com/go-sql-driver/mysql"
)

var (
	errBadPassword = errors.New("permission denied")
	serverVersions = []string{
		"SSH-2.0-OpenSSH_8.9p1 Ubuntu-3ubuntu0.10",
		"SSH-2.0-OpenSSH_6.6.1p1 Ubuntu-2ubuntu2.3",
		"SSH-2.0-OpenSSH_6.7p1 Debian-5+deb8u3",
		"SSH-2.0-OpenSSH_7.2p2 Ubuntu-4ubuntu2.10",
		"SSH-2.0-OpenSSH_7.4",
		"SSH-2.0-OpenSSH_8.0",
		"SSH-2.0-OpenSSH_8.4p1 Debian-2~bpo10+1",
		"SSH-2.0-OpenSSH_8.4p1 Debian-5+deb11u1",
		"SSH-2.0-OpenSSH_8.9p1 Ubuntu-3ubuntu0.6",
	}
	db *sql.DB
)

func main() {
	// Read database connection details from environment variables
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	if dbUser == "" || dbPassword == "" || dbHost == "" || dbPort == "" || dbName == "" {
		log.Fatal("Database connection details are not fully set in environment variables")
	}

	// Build Data Source Name
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbName)

	// Connect to MySQL
	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	if len(os.Args) > 1 {
		logPath := fmt.Sprintf("%s/fakessh-%s.log", os.Args[1], time.Now().Format("2006-01-02-15-04-05-000"))
		logFile, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
		if err != nil {
			log.Println("Failed to open log file:", logPath, err)
			return
		}
		defer logFile.Close()
		log.SetOutput(logFile)
	}

	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	serverConfig := &ssh.ServerConfig{
		MaxAuthTries:      6,
		PasswordCallback:  passwordCallback,
		PublicKeyCallback: publicKeyCallback,
		ServerVersion:     serverVersions[0],
	}

	privateKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	signer, _ := ssh.NewSignerFromSigner(privateKey)
	serverConfig.AddHostKey(signer)

	listener, err := net.Listen("tcp", ":22")
	if err != nil {
		log.Println("Failed to listen:", err)
		return
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Failed to accept:", err)
			break
		}
		go handleConn(conn, serverConfig)
	}
}

func passwordCallback(conn ssh.ConnMetadata, password []byte) (*ssh.Permissions, error) {
	passwordHash := sha256.Sum256(password)
	passwordHashHex := hex.EncodeToString(passwordHash[:])
	log.Println("Password login attempt:", conn.RemoteAddr(), string(conn.ClientVersion()), conn.User(), string(password))

	// Save to database
	_, err := db.Exec("INSERT IGNORE INTO ssh (user, password, sha256) VALUES (?, ?, ?)", conn.User(), string(password), passwordHashHex)
	if err != nil {
		log.Println("Failed to save login attempt to database:", err)
	}

	time.Sleep(100 * time.Millisecond)
	return nil, errBadPassword
}

func publicKeyCallback(conn ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {
	log.Println("Public key login attempt:", conn.RemoteAddr(), string(conn.ClientVersion()), conn.User())
	time.Sleep(100 * time.Millisecond)
	return nil, errBadPassword
}

func handleConn(conn net.Conn, serverConfig *ssh.ServerConfig) {
	defer conn.Close()
	log.Println(conn.RemoteAddr())
	ssh.NewServerConn(conn, serverConfig)
}
