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
	"Soozan-ws/interservice"
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
	userChannel := channels.GetOrCreateUserChannel(userID)
	userChannel.Add(conn)
	log.Printf("User %f joined their private channel", userID)

	// Ensure in DC user is removed from his group | closing connection
	defer func() {
		userChannel.Remove(conn)
		conn.Close()
		log.Printf("User %f disconnected", userID)
	}()

	for {
		msg, op, err := wsutil.ReadClientData(conn)
		if err != nil {
			log.Println("Read error:", err)
			break
		}

		response, users, err := interservice.Responde(msg)

		if err == nil {
			if users[0] == 0 {
				// if no user is set to the array send the response to current conn
				err = wsutil.WriteServerMessage(conn, op, response)
			} else {
				conns := channels.GetMultipleUserConnections(users)
				for _, conn := range conns {
					err = wsutil.WriteServerMessage(conn, op, response)
				}

			}
			if err != nil {
				log.Println("Write error:", err)
			}

		} else {
			errorAsByte := []byte(err.Error())
			if err := wsutil.WriteServerMessage(conn, op, errorAsByte); err != nil {
				log.Println("Write error:", err)
			}
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
