package main

type RosterResponse struct {
	Players []PlayerObj `json:"roster"`
}

type PlayerObj struct {
	Player   Person   `json:"person"`
	Position Position `json:"position"`
}

type Person struct {
	PlayerId int    `json:"id"`
	FullName string `json:"fullName"`
}

type Position struct {
	PositionCode string `json:"code"`
}
