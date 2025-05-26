package main

type TeamsResponse struct {
	Teams []Team `json:"standings"`
}

type Team struct {
	Id           int        `json:"id"`
	Name         string     `json:"name"`
	Abbreviation string     `json:"abbreviation"`
	Division     Division   `json:"division"`
	Conference   Conference `json:"conference"`
	FranchiseId  int        `json:"franchiseId"`
}

type Division struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Conference struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}
