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
)

type OnIcePlayResult struct {
	awayLineHash OnIceRecord
	homeLineHash OnIceRecord
}

type OnIceRecord struct {
	gamePk int
	teamId int
	eventIdx int
	lineHash string
	skaterId1 int
	skaterId2 int
	skaterId3 int
	skaterId4 int
	skaterId5 int
	skaterId6 int
}

func GetPlayersOnIce(gamePk int, records []Play) {
	outChannel := make(chan OnIceRecord)

	for _, play := range records {
		timeStamp := formatTimeStamp(play.About.DateTime)
		go processOnIce(gamePk, play.About.EventIdx, timeStamp, outChannel)
	}
	
	for i:=0;i<len(records)*2;i++ {
		fmt.Println(<-outChannel)
	}
	
	close(outChannel)
}


func processOnIce(gamePk int, eventIdx int, timeStamp string, outChannel chan<- OnIceRecord) {
	response, err := http.Get(fmt.Sprintf("https://statsapi.web.nhl.com/api/v1/game/%s/feed/live?timecode=%s", strconv.Itoa(gamePk), timeStamp));

	if err != nil {
		log.Fatal(err)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	
	if err != nil {
		log.Fatal(err)
	}
	
	var responseObject GameResponse
	json.Unmarshal(responseData, &responseObject)

	awayLineHash := getTeamLineHash(gamePk, eventIdx, responseObject.LiveData.Boxscore.BoxscoreTeams.BoxscoreTeamAway.OnIcePlus)
	homeLineHash := getTeamLineHash(gamePk, eventIdx, responseObject.LiveData.Boxscore.BoxscoreTeams.BoxscoreTeamHome.OnIcePlus)
	
	// outChannel <- OnIcePlayResult{
	// 	awayLineHash: awayLineHash,
	// 	homeLineHash: homeLineHash,
	// }

	outChannel <- awayLineHash
	outChannel <- homeLineHash
}

func getTeamLineHash(gamePk int, eventIdx int, onIcePlus []OnIcePlusPlayer) (OnIceRecord) {
	onIceRecord := OnIceRecord {
		gamePk: gamePk,
		eventIdx: eventIdx,
	}

	playerIdList := make([]int, 0, 6)
	for _, onIcePlayer := range onIcePlus { 
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

func formatTimeStamp(timeStamp string) (string) {
	timeStamp = strings.ReplaceAll(timeStamp, "-", "")
	timeStamp = strings.ReplaceAll(timeStamp, ":", "")
	timeStamp = strings.ReplaceAll(timeStamp, "T", "_")
	timeStamp = strings.ReplaceAll(timeStamp, "Z", "")
	return timeStamp
}
