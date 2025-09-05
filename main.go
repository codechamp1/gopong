package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
)

type Player struct {
	name  string
	score int
}

type ShotType int

const (
	Serve ShotType = iota
	Shoot
	Miss
)

type Shot struct {
	source   *Player
	shotType ShotType
}

type Game struct {
	players []*Player
}

const MinPlayers = 2
const MaxPoints = 11

func play(player *Player, in <-chan Shot, done chan struct{}, players []*Player, pongs map[*Player]chan Shot) {
	for {
		select {
		case <-done:
			return
		case pong := <-in:
			if pong.shotType == Miss {
				pong.source.score++
				fmt.Printf("Player %s shot at player %s and missed, and increased its score to %d\n", pong.source.name, player.name, pong.source.score)
				if pong.source.score == MaxPoints {
					close(done)
					return
				}
			}

			target := pickWeightedTarget(players, player)

			stroke := Shoot
			if rand.Float64() <= 0.1 {
				stroke = Miss
			}

			pongs[target] <- Shot{player, stroke}
		}
	}
}

func pickWeightedTarget(players []*Player, exclude *Player) *Player {
	minScore, maxScore := players[0].score, players[0].score
	for _, p := range players {
		if p == exclude {
			continue
		}
		if p.score < minScore {
			minScore = p.score
		}
		if p.score > maxScore {
			maxScore = p.score
		}
	}

	var targets []*Player
	var weights []int
	totalWeight := 0

	for _, p := range players {
		if p == exclude {
			continue
		}

		w := maxScore - p.score + 1
		targets = append(targets, p)
		weights = append(weights, w)
		totalWeight += w
	}

	// use Weighted Random Algorithm
	r := rand.Intn(totalWeight)
	for i, w := range weights {
		if r < w {
			return targets[i]
		}
		r -= w
	}

	// fallback
	return targets[len(targets)-1]
}

func (g *Game) ReportResults() {
	var winner *Player
	for _, p := range g.players {
		if p.score >= MaxPoints {
			winner = p
			break
		}
	}

	if winner != nil {
		fmt.Printf("\nPlayer %s reached %d points and lost!\n", winner.name, MaxPoints)
	}

	fmt.Println("Final scores:")
	for _, p := range g.players {
		if p != winner {
			fmt.Printf("Player %s: %d points\n", p.name, p.score)
		}
	}
}

func main() {
	var players []*Player
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Enter player names. Type 'done' when finished. You need at least 2 players to start the game.")

	for {
		fmt.Print("Enter player name: ")
		input, _ := reader.ReadString('\n')
		name := strings.TrimSpace(input)

		if strings.ToLower(name) == "done" {
			if len(players) < MinPlayers {
				fmt.Println(fmt.Sprintf("You need atleast %d players!", MinPlayers))
				continue
			}
			break
		}

		if name == "" {
			fmt.Println("Name cannot be empty. Try again.")
			continue
		}

		players = append(players, &Player{name, 0})
	}

	var server *Player
	for {
		fmt.Print("Enter the name of the player who will serve first: ")
		input, _ := reader.ReadString('\n')
		serverName := strings.TrimSpace(input)

		found := false
		for _, p := range players {
			if strings.EqualFold(p.name, serverName) {
				server = p
				found = true
				break
			}
		}

		if !found {
			fmt.Println("Name not in the player list. Try again.")
			continue
		}
		break
	}

	game := Game{
		players,
	}

	gameOver := make(chan struct{})
	pongs := make(map[*Player]chan Shot)

	for _, player := range game.players {
		ch := make(chan Shot)
		pongs[player] = ch
		go play(player, pongs[player], gameOver, game.players, pongs)
	}

	firstTarget := pickWeightedTarget(game.players, server)
	pongs[firstTarget] <- Shot{server, Serve}
	<-gameOver
	game.ReportResults()
}
