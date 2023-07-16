package main

type GameResponse struct {
	LiveData LiveData `json:"liveData"`
}

type LiveData struct {
	Plays Plays `json:"plays"`
	Boxscore Boxscore `json:"boxscore"`
}

type Plays struct {
	AllPlays []Play `json:"allPlays"`
}

type Boxscore struct {
	BoxscoreTeams BoxscoreTeams `json:"teams"`
}

type BoxscoreTeams struct {
	BoxscoreTeamAway BoxscoreTeamAway `json:"away"`
	BoxscoreTeamHome BoxscoreTeamHome `json:"home"`
}

type BoxscoreTeamAway struct {
	OnIcePlus []OnIcePlusPlayer `json:"onIcePlus"`
}

type BoxscoreTeamHome struct {
	OnIcePlus []OnIcePlusPlayer `json:"onIcePlus"`
}

type OnIcePlusPlayer struct {
	PlayerId int `json:"playerId"`
	ShiftDuration int `json:"shiftDuration"`
}

type Play struct {
	Players []Players `json:"players"`
	Result Result `json:"result"`
	About About `json:"about"`
	Coordinates Coordinates `json:"coordinates"`
	Team Team `json:"team"`
}

type Players struct {
	Player Player `json:"player"`
	PlayerType string `json:"playerType"`
}

type Player struct {
	PlayerId int `json:"id"`
}

type Result struct {
	Event string `json:"event"`
	EventCode string `json:"eventCode"`
	EventTypeId string `json:"eventTypeId"`
	Description string `json:"description"`
	SecondaryType string `json:"string"`
}

type About struct {
	EventIdx int `json:"eventIdx"`
	EvendId int `json:"eventId"`
	Period int `json:"period"`
	PeriodType string `json:"periodType"`
	OrdinalNum string `json:"ordinalNum"`
	PeriodTime string `json:"periodTime"`
	PeriodTimeRemaining string `json:"periodTimeRemaining"`
	DateTime string `json:"dateTime"`
	Goals Goals `json:"goals"`
}

type Coordinates struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
}

type Team struct {
	Id float32 `json:"id"`
}

type Goals struct {
	Away int `json:"away"`
	Home int `json:"home"`
}