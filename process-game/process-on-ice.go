package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/exp/slices"
)

type OnIcePlayResult struct {
	awayLineHash OnIceRecord
	homeLineHash OnIceRecord
}

func GetPlayersOnIce(gamePk int, records []Play) []OnIceRecord {
	outChannel := make(chan OnIceRecord)

	for _, play := range records {
		timeStamp := formatTimeStamp(play.About.DateTime)
		go processOnIce(gamePk, play.About.EventIdx, timeStamp, outChannel)
	}

	var onIceRecordList []OnIceRecord

	for i := 0; i < len(records)*2; i++ {
		onIceRecordList = append(onIceRecordList, <-outChannel)
		// fmt.Println(<-outChannel)
	}

	return onIceRecordList
}

func processOnIce(gamePk int, eventIdx int, timeStamp string, outChannel chan<- OnIceRecord) {
	response, err := http.Get(fmt.Sprintf("https://statsapi.web.nhl.com/api/v1/game/%s/feed/live?timecode=%s", strconv.Itoa(gamePk), timeStamp))

	if err != nil {
		log.Fatal(err)
	}

	responseData, err := ioutil.ReadAll(response.Body)

	if err != nil {
		log.Fatal(err)
	}

	var responseObject GameResponse
	json.Unmarshal(responseData, &responseObject)

	awayLineHash := getTeamLineHash(gamePk, eventIdx, responseObject.LiveData.Boxscore.BoxscoreTeams.BoxscoreTeamAway)
	homeLineHash := getTeamLineHash(gamePk, eventIdx, responseObject.LiveData.Boxscore.BoxscoreTeams.BoxscoreTeamHome)

	outChannel <- awayLineHash
	outChannel <- homeLineHash
}

func getTeamLineHash(gamePk int, eventIdx int, boxScoreTeam BoxscoreTeam) OnIceRecord {
	onIcePlus := boxScoreTeam.OnIcePlus

	onIceRecord := OnIceRecord{
		gamePk:   gamePk,
		teamId:   boxScoreTeam.Team.Id,
		eventIdx: eventIdx,
	}

	playerIdList := make([]int, 0, 6)
	for _, onIcePlayer := range onIcePlus {
		if slices.Contains(boxScoreTeam.Goalies, onIcePlayer.PlayerId) {
			onIceRecord.goalieId = onIcePlayer.PlayerId
			continue
		}
		playerIdList = append(playerIdList, onIcePlayer.PlayerId)
	}
	sort.Ints(playerIdList)

	for i, playerId := range playerIdList {
		switch i {
		case 0:
			onIceRecord.skaterId1 = playerId
		case 1:
			onIceRecord.skaterId2 = playerId
		case 2:
			onIceRecord.skaterId3 = playerId
		case 3:
			onIceRecord.skaterId4 = playerId
		case 4:
			onIceRecord.skaterId5 = playerId
		case 5:
			onIceRecord.skaterId6 = playerId
		}
	}

	lineStr := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(playerIdList)), "-"), "[]")
	lineHash := getLineHash(lineStr)
	onIceRecord.lineHash = lineHash

	return onIceRecord
}

func getLineHash(line string) string {
	hash := md5.Sum([]byte(line))
	return hex.EncodeToString(hash[:])
}

func formatTimeStamp(timeStamp string) string {
	timeStamp = strings.ReplaceAll(timeStamp, "-", "")
	timeStamp = strings.ReplaceAll(timeStamp, ":", "")
	timeStamp = strings.ReplaceAll(timeStamp, "T", "_")
	timeStamp = strings.ReplaceAll(timeStamp, "Z", "")
	return timeStamp
}
