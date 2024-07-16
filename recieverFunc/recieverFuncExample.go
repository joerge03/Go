package main

import "fmt"

type PlayerHealth struct {
	Health, MaxPlayerHealth uint
}

type PlayerEnergy struct {
	Energy, MaxEnergy uint
}

type Player struct {
	Name   string
	Health PlayerHealth
	Energy PlayerEnergy
}

func (player *Player) ModifyHealth(health uint) {
	if health <= player.Health.MaxPlayerHealth {
		player.Health.Health = health
		fmt.Printf("Health modified from %v to %v \n", player.Health.Health, health)
	} else {
		fmt.Println("Cannot exceed above max health")
	}
}

func (player *Player) ModifyEnergy(energy uint) {
	if energy > player.Energy.MaxEnergy {
		fmt.Println("Cannot exceed above max energy")
	} else {
		fmt.Printf("Energy modified from %v to %v \n", player.Energy.Energy, energy)
		player.Energy.Energy = energy
	}
}

func (player *Player) AttackPlayer(damage uint) {
	playerHealth := &player.Health.Health
	if *playerHealth-damage > *playerHealth {
		*playerHealth = 0
	} else {
		*playerHealth -= damage
	}
}

func ReciverFuncExample() {
	players := []Player{{
		Name:   "Yasuo",
		Health: PlayerHealth{Health: 15, MaxPlayerHealth: 50},
		Energy: PlayerEnergy{Energy: 50, MaxEnergy: 350},
	}}

	players[0].ModifyEnergy(2)
	players[0].ModifyHealth(50)

	players[0].AttackPlayer(500)

	fmt.Println(players)
}
