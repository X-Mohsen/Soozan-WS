package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"

	"Soozan-ws/auth"
	"Soozan-ws/channels"
)

func handler(w http.ResponseWriter, r *http.Request) {
	// validate user JWT
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Unauthorized: missing token", http.StatusUnauthorized)
		return
	}

	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		http.Error(w, "Unauthorized: invalid token format", http.StatusUnauthorized)
		return
	}

	tokenString := tokenParts[1]

	userID, err := auth.ValidateJWT(tokenString)
	if err != nil || userID == 0 {
		http.Error(w, "Unauthorized: invalid token", http.StatusUnauthorized)
		return
	}

	// Upgrade to Web-Socket connection
	conn, _, _, err := ws.UpgradeHTTP(r, w)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}

	// Adding user to his own private group
	userChannel := channels.GetUserChannel(userID)
	userChannel.Add(conn)
	log.Printf("User %s joined their private channel", userID)

	// Ensure in DC user is removed from his group | closing connection
	defer func() {
		userChannel.Remove(conn)
		conn.Close()
		log.Printf("User %s disconnected", userID)
	}()

	for {
		msg, op, err := wsutil.ReadClientData(conn)
		if err != nil {
			log.Println("Read error:", err)
			break
		}
		log.Printf("Received: %s\n", msg)

		if err := wsutil.WriteServerMessage(conn, op, msg); err != nil {
			log.Println("Write error:", err)
			break
		}
	}
}

func main() {
	if err := auth.LoadPublicKey("public_key.pem"); err != nil {
		log.Fatalf("Error loading public key: %v", err)
	}

	http.HandleFunc("/ws", handler)
	fmt.Println("WebSocket server started at ws://localhost:8080/ws")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
