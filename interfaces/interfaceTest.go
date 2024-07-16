package main

import "fmt"

type PlayerReset interface {
	Reset()
}

type Player struct {
	Name   string
	Health uint
}

func (p *Player) Reset() {
	p.Health = 100
}

func execute(i PlayerReset) {
	i.Reset()
}

func resetWithDamage(i PlayerReset) {
	if player, ok := i.(*Player); ok {
		player.Health = 1
	}
}

func (p *Player) DamagePlayer(d uint) {
	p.Health -= d
}

func main2() {
	player := Player{Name: "test", Health: 100}
	player.DamagePlayer(1)

	execute(&player)

	fmt.Println(player)
}
