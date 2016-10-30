package ogdata

import (
	"fmt"
	"ogbot/helpers"
	"strconv"
	"strings"
)

type Page struct {
	PageType string
	Content  string
}

type MetaData struct {
	Login string
	Pass  string
	Uni   string
	Lang  string
}

type Resources struct {
	Metal     int
	Crystal   int
	Deuterium int
}

type Planet struct {
	Name        string
	Coordinates string
	Resources   Resources
	//204 -> 10 means 10 ships of type 204
	DockedFleet map[string]int
}

type GameData struct {
	Planets map[string]Planet
}

func (g *GameData) Print(logger helpers.Logger) {
	for k, v := range g.Planets {
		numShips := 0
		for _, num := range v.DockedFleet {
			numShips += num
		}
		logger.Printf("(%s) Planet [%s]: %dM, %dC, %dD. %d ships docked", k, v.Coordinates, v.Resources.Metal,
			v.Resources.Crystal, v.Resources.Deuterium, numShips)
	}
}

const (
	AttackMovement    = 1
	TransportMovement = 3
	SpyMovement       = 6
)

func getMovementType(move int) string {
	switch {
	case move == 1:
		return "Attack"
	case move == 2:
		return "Grouped Attack"
	case move == 3:
		return "Transport"
	case move == 4:
		return "Deployment"
	case move == 5:
		return "Defend"
	case move == 6:
		return "Espionage"
	case move == 7:
		return "Colonisation"
	case move == 8:
		return "Exploit"
	case move == 9:
		return "Moon destruction"
	case move == 15:
		return "Expedition"
	default:
		return "Unknown"
	}
}

type FleetMovement struct {
	id      string
	from    string
	to      string
	move    int // MovementType
	arrival string
	hostile bool
}

func (f *FleetMovement) Print(logger helpers.Logger) {
	var hostile string
	if f.hostile {
		hostile = "/!\\"
	}
	logger.Printf("(%s) at %s [%s] -> [%s] (%s) %s", f.id, f.arrival,
		f.from, f.to, getMovementType(f.move), hostile)
}

func getAttributeValue(att, content string) string {
	split := strings.Split(content, att+"=\"")
	if len(split) < 2 {
		return ""
	}
	final := strings.Split(split[1], "\"")
	return final[0]
}

func getFleetDataValue(att, content, sep1, sep2 string) string {
	split := strings.Split(content, att)
	second := strings.Split(split[1], sep1)
	final := strings.Split(second[1], sep2)
	return final[0]
}

func GetFleetMovementInfo(content string) *FleetMovement {
	move := &FleetMovement{}
	move.hostile = strings.Contains(content, "countDown hostile")

	eventRowId := getAttributeValue("id", content)
	split := strings.Split(eventRowId, "-")
	move.id = split[1]
	mission, err := strconv.Atoi(getAttributeValue("data-mission-type", content))
	if err != nil {
		return nil
	}
	move.move = mission
	move.from = getFleetDataValue("coordsOrigin", content, "[", "]")
	move.to = getFleetDataValue("destCoords", content, "[", "]")
	move.arrival = getFleetDataValue("arrivalTime", content, ">", " ")
	return move
}

func getTagValue(tag, page string) (string, error) {
	split := strings.Split(page, tag)
	if len(split) != 2 {
		return "", fmt.Errorf("too many %s tags (%d)", tag, len(split)-1)
	}
	splitValue := strings.Split(split[1], ">")
	lastSplit := strings.Split(splitValue[1], "<")
	return helpers.RemoveNoise(lastSplit[0]), nil
}

func GetResourceValue(resource, page string) int {
	value, err := getTagValue(resource, page)
	if err != nil {
		fmt.Println(err.Error())
		return -1
	}
	val, _ := strconv.Atoi(value)
	return val
}

func getMetaOGame(page, suffix string) string {
	split := strings.Split(page, "ogame-"+suffix)
	split = strings.Split(split[1], "\"")
	return split[2]
}

func ListAvailablePlanetIds(page string) []string {
	split := strings.Split(page, "id=\"planet-")
	var list []string
	for _, value := range split[1:] {
		temp := strings.Split(value, "\"")
		list = append(list, temp[0])
	}
	return list
}

func GetCurrentPlanet(page string) (string, string, string) {
	id := getMetaOGame(page, "planet-id")
	name := getMetaOGame(page, "planet-name")
	coord := getMetaOGame(page, "planet-coordinates")
	return id, name, coord
}

func GetDockedFleet(content string) map[string]int {
	dockedFleet := make(map[string]int)
	split := strings.Split(content, "#shipsChosen")
	for _, v := range split[1:] {
		shipData := strings.Split(v, ",")
		shipId := helpers.RemoveNoise(shipData[1])
		splitnumber := strings.Split(shipData[2], ")")
		number, err := strconv.Atoi(splitnumber[0])
		if err != nil {
			return nil
		}
		dockedFleet[shipId] = number
	}
	return dockedFleet
}
