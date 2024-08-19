package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/seyhak/websocketgo/saper"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type ServedGame struct {
	// gameIdx uint8
	game *saper.Game
}

var userToConnection map[string]*websocket.Conn
var servedGames [1]ServedGame

// var playersToGame map[string]ServedGame

func convertGameToByteArr(game *saper.Game, user string) []byte {
	const END_PLAYERS_FLAG = 255
	const ADDITIONAL_FIELDS_LEN = 5
	moreAdditionalFields := 0
	// maxByteArr := [saper.AREA + ADDITIONAL_FIELDS_LEN]byte{}
	byteArr := make([]byte, saper.AREA+ADDITIONAL_FIELDS_LEN+len(game.Players)+100)
	byteArr[0] = saper.WIDTH
	player := game.Players[user]
	byteArr[1] = player.GameStatus
	byteArr[2] = byte(player.Points)
	byteArr[3] = byte(player.HueDeg)
	i := 4
	for playerName, playerValues := range game.Players {
		if playerName == user {
			continue
		} else {
			byteArr[i] = byte(playerValues.Points)
			i++
			byteArr[i] = playerValues.HueDeg
			i++
			moreAdditionalFields += 2
		}
	}
	byteArr[i] = byte(END_PLAYERS_FLAG)

	fmt.Printf("\n State %s, %v \n", user, game.State)
	for x, fieldYs := range game.State {
		for y, fieldYVal := range fieldYs {
			byteArr[(x*saper.WIDTH)+y+ADDITIONAL_FIELDS_LEN+moreAdditionalFields] = fieldYVal
		}
	}
	fmt.Printf("\nBYTE ARR %s, %v \n", user, byteArr)
	return byteArr[:saper.AREA+ADDITIONAL_FIELDS_LEN+len(game.Players)+moreAdditionalFields]
}
func getUser(r *http.Request) string {
	user, err := r.Cookie("user")
	if err != nil {
		log.Println("no user found", err)
		return ""
	}
	return user.Value
}

func getCurrentGame(user string) *ServedGame {
	if servedGames[0].game == nil {
		game := saper.GetGame()
		servedGames[0] = ServedGame{&game}
	}
	currentGame := &servedGames[0]
	currentGame.game.AddUserToGame(user)
	return currentGame
}

func reloadGame() *ServedGame {
	game := saper.GetGame()
	servedGames[0] = ServedGame{&game}
	return &servedGames[0]
}

func handleServeGame(user string) []byte {
	game := getCurrentGame(user)
	return convertGameToByteArr(game.game, user)
}

func handleFieldClick(p []byte, user string) error {
	game := getCurrentGame(user)
	game.game.HandleFieldClick(int(p[0]), int(p[1]), user)

	// fmt.Printf("\n conn  %v", conn)
	// fmt.Printf("\nFIELD IS %v", game.game.Field)
	// fmt.Printf("\n!!!!!!!!CHANGED FIELD 2 %v", game.game.State)
	var err error
	for otherP, con := range userToConnection {
		byteField := convertGameToByteArr(game.game, otherP)
		err = con.WriteMessage(websocket.BinaryMessage, byteField)

	}
	return err
}

func handler(w http.ResponseWriter, r *http.Request) {
	// println("header")
	// for name, values := range r.Header {
	// 	// Loop over all values for the name.
	// 	for _, value := range values {
	// 		println(name, value)
	// 	}
	// }
	user := getUser(r)
	println("user", user)
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	// w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Origin, Content-Length, Accept-Encoding, X-CSRF-Token")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Access-Control-Allow-Headers, Origin,Accept, X-Requested-With, Content-Type, Access-Control-Request-Method, Access-Control-Request-Headers")
	upgrader.CheckOrigin = func(r *http.Request) bool {
		// r.Header.Get("Origin")
		return true
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	// con := userToConnection[user]
	// if con == nil {
	userToConnection[user] = conn
	// }
	fmt.Println("START")
	for k, _ := range userToConnection {
		fmt.Printf("\n %v \n", k)
	}
	fmt.Println("END")
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("Reading messages from", r.Header.Get(("Origin")))
	for {
		_, p, err := conn.ReadMessage()
		log.Println("p", p)
		if err != nil {
			log.Println(err)
			return
		}

		clientHasNoData := len(p) == 0
		forceReload := len(p) == 1
		log.Println(p)
		if clientHasNoData {
			byteField := handleServeGame(user)
			if err := conn.WriteMessage(websocket.BinaryMessage, byteField); err != nil {
				log.Println(err)
				continue
			}
		} else if forceReload {
			reloadGame()
			byteField := handleServeGame(user)
			for _, con := range userToConnection {
				if err := con.WriteMessage(websocket.BinaryMessage, byteField); err != nil {
					log.Println(err)
					continue
				}
			}
		} else {
			log.Println("clientHas Data", p)
			if err := handleFieldClick(p, user); err != nil {
				log.Println(err)
				continue
			}
		}
		// if err := conn.WriteMessage(messageType, p); err != nil {
		// 	log.Println(err)
		// 	return
		// }
	}
}

func main() {
	log.Println("Starting server")
	userToConnection = make(map[string]*websocket.Conn)
	http.HandleFunc("/init", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
