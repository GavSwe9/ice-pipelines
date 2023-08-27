package main

type GameResponse struct {
	GameData GameData `json:"gameData"`
	LiveData LiveData `json:"liveData"`
}

type GameData struct {
	Game     Game     `json:"game"`
	DateTime DateTime `json:"datetime"`
	Teams    Teams    `json:"teams"`
}

type Game struct {
	GamePk int    `json:"pk"`
	Season string `json:"season"`
	Type   string `json:"type"`
}

type DateTime struct {
	DateTime string `json:"dateTime"`
}

type Teams struct {
	AwayTeam Team `json:"away"`
	HomeTeam Team `json:"home"`
}

type Team struct {
	Id int `json:"id"`
}

type LiveData struct {
	Plays    Plays    `json:"plays"`
	Boxscore Boxscore `json:"boxscore"`
}

type Plays struct {
	AllPlays []Play `json:"allPlays"`
}

type Boxscore struct {
	BoxscoreTeams BoxscoreTeams `json:"teams"`
}

type BoxscoreTeams struct {
	BoxscoreTeamAway BoxscoreTeam `json:"away"`
	BoxscoreTeamHome BoxscoreTeam `json:"home"`
}

type BoxscoreTeam struct {
	Team      Team              `json:"team"`
	Goalies   []int             `json:"goalies"`
	OnIcePlus []OnIcePlusPlayer `json:"onIcePlus"`
}

type OnIcePlusPlayer struct {
	PlayerId      int `json:"playerId"`
	ShiftDuration int `json:"shiftDuration"`
}

type Play struct {
	Players     []Players   `json:"players"`
	Result      Result      `json:"result"`
	About       About       `json:"about"`
	Coordinates Coordinates `json:"coordinates"`
	Team        Team        `json:"team"`
}

type Players struct {
	Player     Player `json:"player"`
	PlayerType string `json:"playerType"`
}

type Player struct {
	PlayerId int `json:"id"`
}

type Result struct {
	Event         string `json:"event"`
	EventCode     string `json:"eventCode"`
	EventTypeId   string `json:"eventTypeId"`
	Description   string `json:"description"`
	SecondaryType string `json:"string"`
}

type About struct {
	EventIdx            int    `json:"eventIdx"`
	EvendId             int    `json:"eventId"`
	Period              int    `json:"period"`
	PeriodType          string `json:"periodType"`
	OrdinalNum          string `json:"ordinalNum"`
	PeriodTime          string `json:"periodTime"`
	PeriodTimeRemaining string `json:"periodTimeRemaining"`
	DateTime            string `json:"dateTime"`
	Goals               Goals  `json:"goals"`
}

type Coordinates struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
}

type Goals struct {
	Away int `json:"away"`
	Home int `json:"home"`
}

type OnIceRecord struct {
	gamePk    int
	teamId    int
	eventIdx  int
	lineHash  string
	skaterId1 int
	skaterId2 int
	skaterId3 int
	skaterId4 int
	skaterId5 int
	skaterId6 int
	goalieId  int
}
