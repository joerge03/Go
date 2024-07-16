package main

import "fmt"

const (
	Small = iota
	Medium
	Large
)

const (
	Ground = iota
	Air
)

type AutomateList interface {
	conveyour
	shipper
}

type (
	Beltsize int
	Shipping int
)

func (b *Beltsize) String() string {
	return []string{"Small", "Medium", "Large"}[*b]
}

func (s *Shipping) String() string {
	return []string{"Ground", "Air"}[*s]
}

type conveyour interface {
	Convey() Beltsize
}
type shipper interface {
	Shipper() Shipping
}

type SpamMail struct {
	Amount int
}

func (s *SpamMail) String() string {
	return "Spam mail"
}

func (s *SpamMail) Convey() Beltsize {
	return Medium
}

func (s *SpamMail) Shipper() Shipping {
	return Ground
}

func automate(item AutomateList) {
	fmt.Println(item.Convey(), item.Shipper())
}

type toxicWaster struct {
	weight int
}

func (t *toxicWaster) Convey() Beltsize {
	return Small
}

func main3() {
	test := Shipping(Ground)

	mail := SpamMail{300}

	// toxicWaste := toxicWaster{2}

	automate(&mail)

	// automate(&toxicWaste)

	fmt.Println(test.String())
}
