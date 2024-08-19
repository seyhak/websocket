package saper

import (
	"fmt"
	"math/rand/v2"
)

type Player struct {
	Points     uint16
	GameStatus uint8
	HueDeg     uint8
}
type Field = [WIDTH][HEIGHT]uint8
type Game struct {
	Field                Field
	State                Field
	IsLost               bool
	IsWon                bool
	DiscoveredMinesCount uint8
	Players              map[string]Player
}

const WIDTH = 25
const HEIGHT = 25
const MINE_VALUE uint8 = 10
const CLEAN_VALUE uint8 = 11
const AREA = WIDTH * HEIGHT
const MAX_MINES = AREA / 15
const GAME_WON = 1
const GAME_LOST = 2
const GAME_IN_GAME = 3
const POINTS_TO_DECREASE_ON_MINE = 10

func generateField() Field {
	var field Field
	mineCount := 0
	x := 0
	for {
		for y := 0; y < HEIGHT; y++ {
			currentField := field[y][x]
			isMine := currentField == MINE_VALUE
			if isMine {
				continue
			}
			if mineCount >= MAX_MINES {
				break
			}
			generatedNumber := uint8(rand.IntN(200)) // 200/10 = 5%
			willBeMine := generatedNumber <= MINE_VALUE
			if willBeMine {
				field[y][x] = MINE_VALUE
				mineCount++
			} else {
				field[y][x] = 0
			}
		}
		x++
		if mineCount < MAX_MINES && x == WIDTH-1 {
			x = 0
		} else if x == WIDTH {
			break
		}
	}
	fmt.Printf("GENERATED FIELD %v", field)
	return field
}

func generateState() Field {
	var state Field
	return state
}

func generateGameState() Game {
	field := generateField()
	state := generateState()
	players := make(map[string]Player)
	return Game{field, state, false, false, 0, players}
}

func (game *Game) setFieldState(x int, y int, value uint8) {
	fmt.Printf("\n game.State %v", value)
	fmt.Printf("\n val %v", value)
	game.State[y][x] = value
	fmt.Printf("\n game.State2 %v", value)
}

func (game *Game) validateGameOver() bool {
	return game.DiscoveredMinesCount == MAX_MINES
}

func (game *Game) increasePlayerPoints(username string) {
	player := game.Players[username]
	points := player.Points
	points++
	fmt.Printf("ps %v \n", points)
	game.Players[username] = Player{points, player.GameStatus, player.HueDeg}
	fmt.Printf("after increase %v \n", game.Players[username])
}

func (game *Game) setGameOver(username string) {
	player := game.Players[username]
	player.Points = player.Points - POINTS_TO_DECREASE_ON_MINE
}
func (game *Game) setWinnerAndLosers(username string) {
	wonPlayer := username
	maxValue := game.Players[username].Points
	for player, info := range game.Players {
		if info.Points > maxValue {
			maxValue = info.Points
			wonPlayer = player
		}
	}
	wonPlayerStruct := game.Players[wonPlayer]

	game.Players[wonPlayer] = Player{wonPlayerStruct.Points, GAME_WON, wonPlayerStruct.HueDeg}
	for player, info := range game.Players {
		if player == wonPlayer {
			continue
		}
		game.Players[player] = Player{info.Points, GAME_LOST, wonPlayerStruct.HueDeg}
	}
}

func (game *Game) HandleFieldClick(x int, y int, username string) {
	fieldVal := game.Field[y][x]
	if fieldVal == MINE_VALUE {
		// gameOver
		game.setGameOver(username)
	} else {
		fieldVal = CLEAN_VALUE
		game.DiscoveredMinesCount++
		game.increasePlayerPoints(username)
	}
	game.setFieldState(x, y, uint8(fieldVal))

	if game.validateGameOver() {
		game.setWinnerAndLosers(username)
	}
	// fmt.Printf("CHANGED FIELD STATE %v", game.Field)
}

func (game *Game) AddUserToGame(user string) {
	_, exists := game.Players[user]
	if !exists {
		game.Players[user] = Player{0, GAME_IN_GAME, uint8(rand.IntN(360))}
	}
}

func GetGame() Game {
	return generateGameState()
}
