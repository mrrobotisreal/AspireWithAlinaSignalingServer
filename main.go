//package main
//
//import (
//	"fmt"
//	"github.com/gorilla/websocket"
//	"github.com/pion/webrtc/v3"
//	"log"
//	"net/http"
//)
//
//func getRtcOffer() {
//	// Create a new peer connection
//	peerConnection, err := webrtc.NewPeerConnection(webrtc.Configuration{})
//	if err != nil {
//		panic(err)
//	}
//	defer peerConnection.Close()
//
//	// Set up video track
//	videoTrack, err := webrtc.NewTrackLocalStaticSample(webrtc.RTPCodecCapability{MimeType: webrtc.MimeTypeVP8}, "video", "pion")
//	if err != nil {
//		panic(err)
//	}
//
//	_, err = peerConnection.AddTrack(videoTrack)
//	if err != nil {
//		panic(err)
//	}
//
//	// Generate an offer
//	offer, err := peerConnection.CreateOffer(nil)
//	if err != nil {
//		panic(err)
//	}
//	err = peerConnection.SetLocalDescription(offer)
//	if err != nil {
//		panic(err)
//	}
//
//	// This is where you'd send the SDP to the other peer
//	fmt.Println("SDP offer:\n", offer.SDP)
//}
//
//var upgrader = websocket.Upgrader{
//	CheckOrigin: func(r *http.Request) bool {
//		return true
//	},
//}
//var peers = make(map[*websocket.Conn]bool)
//
//func handleConnection(w http.ResponseWriter, r *http.Request) {
//	fmt.Println("Handling connection...")
//	fmt.Println(r.Body)
//	conn, err := upgrader.Upgrade(w, r, nil)
//	if err != nil {
//		fmt.Println("Upgrade error:", err)
//		return
//	}
//	defer conn.Close()
//	peers[conn] = true
//	for {
//		_, msg, err := conn.ReadMessage()
//		if err != nil {
//			delete(peers, conn)
//			break
//		}
//		for peer := range peers {
//			if peer != conn {
//				fmt.Println("msg: " + string(msg))
//				peer.WriteMessage(websocket.TextMessage, msg)
//			}
//		}
//	}
//}
//
//// CORS middleware
//func enableCors(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		w.Header().Set("Access-Control-Allow-Origin", "*")
//		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
//		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
//
//		if r.Method == "OPTIONS" {
//			w.WriteHeader(http.StatusOK)
//			return
//		}
//
//		next.ServeHTTP(w, r)
//	})
//}
//
//func main() {
//	http.HandleFunc("/video", handleConnection)
//
//	fmt.Println("Video Signaling Server running on port 9999...")
//	if err := http.ListenAndServe(":9999", enableCors(http.DefaultServeMux)); err != nil {
//		log.Fatal(err)
//	}
//}

package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var clients = make(map[*websocket.Conn]bool) // Track connected clients
var broadcast = make(chan Message)           // Broadcast channel

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow any origin
	},
}

// Message structure
type Message struct {
	Event   string `json:"event"`   // e.g., "offer", "answer", "ice-candidate"
	Payload string `json:"payload"` // SDP or ICE data
}

func main() {
	http.HandleFunc("/video", handleConnections)
	go handleMessages()

	log.Println("WebSocket server started on :9999")
	log.Fatal(http.ListenAndServe(":9999", nil))
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	clients[ws] = true

	for {
		var msg Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			delete(clients, ws)
			break
		}
		broadcast <- msg
	}
}

func handleMessages() {
	for {
		msg := <-broadcast
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				client.Close()
				delete(clients, client)
			}
		}
	}
}

//func getRtcOffer() {
//	// Create a new peer connection
//	peerConnection, err := webrtc.NewPeerConnection(webrtc.Configuration{})
//	if err != nil {
//		panic(err)
//	}
//	defer peerConnection.Close()
//
//	// Set up video track
//	videoTrack, err := webrtc.NewTrackLocalStaticSample(webrtc.RTPCodecCapability{MimeType: webrtc.MimeTypeVP8}, "video", "pion")
//	if err != nil {
//		panic(err)
//	}
//
//	_, err = peerConnection.AddTrack(videoTrack)
//	if err != nil {
//		panic(err)
//	}
//
//	// Generate an offer
//	offer, err := peerConnection.CreateOffer(nil)
//	if err != nil {
//		panic(err)
//	}
//	err = peerConnection.SetLocalDescription(offer)
//	if err != nil {
//		panic(err)
//	}
//
//	// This is where you'd send the SDP to the other peer
//	fmt.Println("SDP offer:\n", offer.SDP)
//}
