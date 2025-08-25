# GoPong 🏓

GoPong is a Go project that simulates a multiplayer ping-pong game in the terminal.
It was created as an exercise to explore **concurrency techniques** in Go, inspired by the book [*Gist of Go – Concurrency*](https://antonz.org/go-concurrency/).

---

## Overview

In GoPong, you can have as many players as you like.  
When you run the program, the terminal will ask you to enter the names of the players one by one.  
Type `done` when you have entered all the players.

After that, you will choose a player to start serving. The game will then simulate the rally between players.  
Points are scored when a player misses, and the first player to reach **11 points loses**.  
Randomness makes each match unique, creating unpredictable outcomes.

You can also use this program to decide something with your friends: just add your names, run the game, and see who is the winner. 🙂

---

## How to Play

1. Clone the repository:

```bash
git clone https://github.com/yourusername/GoPong.git
cd GoPong
```

2. Run the program
```bash
go run main.go
```
3. Follow the terminal instructions 🙂
