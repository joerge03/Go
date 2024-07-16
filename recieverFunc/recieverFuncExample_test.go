package main

import "testing"

func NewPlayer() Player {
	return Player{Name: "Ursa", Health: PlayerHealth{64, 100}, Energy: PlayerEnergy{43, 100}}
}

func TestModifyHealth(t *testing.T) {
	player := NewPlayer()
	player.ModifyHealth(300)
	if player.Health.Health > player.Health.MaxPlayerHealth {
		t.Errorf("Incorrect player health, Expected: %v, Result: %v ", player.Health.Health, player.Health.MaxPlayerHealth)
	}
}

func TestModifyEnergy(t *testing.T) {
	player := NewPlayer()
	player.ModifyEnergy(200)
	if player.Energy.Energy > player.Energy.MaxEnergy {
		t.Errorf("Incorrect player energy expected: %v, but recieved: %v", player.Energy.MaxEnergy, player.Energy.Energy)
	}
}

func TestAttack(t *testing.T) {
	player := NewPlayer()
	player.AttackPlayer(300)
	if player.Health.Health >= 300 {
		t.Errorf("Incorrect health, expected %v, but recieved: %v", uint(player.Health.Health-300), player.Health.Health)
	}
}
