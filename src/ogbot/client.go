package main

import (
	"flag"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"ogbot/helpers"
	"ogbot/ogdata"
	"ogbot/oghttp"
	"os"
	"runtime"
	"strings"
	"time"
)

type OGBot struct {
	meta         ogdata.MetaData
	data         ogdata.GameData
	current      ogdata.Page
	cookieHeader []string
	client       *http.Client
	logger       helpers.Logger
	dump         bool
	fleets       []*ogdata.FleetMovement
}

func (b *OGBot) printFleets() {
	b.logger.Printf("%v fleet(s) movement(s)", len(b.fleets))
	for _, v := range b.fleets {
		v.Print(b.logger)
	}
}

func (b *OGBot) goToPage(label, args string) {
	// make the overview request
	req := oghttp.MakePageRequest(label, args, b.meta, b.cookieHeader)
	helpers.LogMark(label+" page request...", b.logger)
	resp, _ := b.client.Do(req)
	if b.dump {
		helpers.DumpResponse(resp, b.logger)
	}

	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		b.logger.Printf("%s\n", err.Error())
		return
	}
	b.current.PageType = label
	b.current.Content = string(contents)
}

func (b *OGBot) ChecksFleetMovements() {
	b.sleep(20)
	b.goToPage("eventList", "&ajax=1")
	split := strings.Split(b.current.Content, "eventFleet")
	if len(split) == 0 {
		b.fleets = nil
		return
	}
	b.fleets = make([]*ogdata.FleetMovement, 0)
	for _, v := range split[1:] {
		b.fleets = append(b.fleets, ogdata.GetFleetMovementInfo(v))
	}
}

// get planet data: planet info, ressources and docked ships
func (b *OGBot) goToPlanet(id string) {
	b.goToPage("fleet1", "&cp="+id)
}

func (b *OGBot) login() {
	// make the login request
	req := oghttp.MakeLoginRequest(b.meta, b.logger, b.dump)
	helpers.LogMark("Login request...", b.logger)
	resp, _ := b.client.Do(req)
	if b.dump {
		helpers.DumpResponse(resp, b.logger)
	}
	b.cookieHeader = resp.Header["Set-Cookie"]

	b.goToPage("overview", "")
	if helpers.IsLogginPage(b.current.Content) {
		b.logger.Printf("not logged in\n")
	}
}

func (b *OGBot) UpdatePlanetData() {
	id, name, coord := ogdata.GetCurrentPlanet(b.current.Content)
	b.data.Planets[id] = ogdata.Planet{
		Name:        name,
		Coordinates: coord,
		Resources: ogdata.Resources{
			Metal:     ogdata.GetResourceValue("resources_metal", b.current.Content),
			Crystal:   ogdata.GetResourceValue("resources_crystal", b.current.Content),
			Deuterium: ogdata.GetResourceValue("resources_deuterium", b.current.Content),
		},
		DockedFleet: ogdata.GetDockedFleet(b.current.Content),
	}
}

func (b *OGBot) UpdatePlanetsData() {
	planets := ogdata.ListAvailablePlanetIds(b.current.Content)
	for _, v := range planets {
		b.sleep(20)
		b.goToPlanet(v)
		b.UpdatePlanetData()
	}
}

func (b *OGBot) ChecksReconnect() {
	b.logger.Printf("checks if still logged in...")
	b.goToPage("overview", "")
	if helpers.IsLogginPage(b.current.Content) {
		b.logger.Printf("no more logged in... trying to reconnect...")
		b.login()
	}
}

func (b *OGBot) sleep(amount int) {
	random := rand.Intn(amount)
	b.logger.Printf("Random sleep: %+v seconds", random)
	time.Sleep(time.Second * time.Duration(random))
}

func (b *OGBot) Run() {
	b.login()
	b.UpdatePlanetsData()

	// checks frequently we are logged in
	ticker := time.NewTicker(1789 * time.Second)
	defer ticker.Stop()
	for _ = range ticker.C {
		b.sleep(1000)
		b.ChecksReconnect()
		b.UpdatePlanetsData()
		b.ChecksFleetMovements()
		b.data.Print(b.logger)
		b.printFleets()
	}
}

func makeOGBot(login, pass, uni, lang string, logger helpers.Logger, dump bool) *OGBot {
	return &OGBot{
		meta: ogdata.MetaData{
			Login: login,
			Pass:  pass,
			Uni:   uni,
			Lang:  lang,
		},
		data: ogdata.GameData{
			Planets: make(map[string]ogdata.Planet),
		},
		client: oghttp.MakeHttpClient(logger, dump),
		logger: logger,
		dump:   dump,
	}
}

func main() {
	login := flag.String("login", "", "login")
	pass := flag.String("pass", "", "password")
	universe := flag.String("uni", "", "universe. Ex: s131, s132...")
	lang := flag.String("lang", "", "language. Ex: en, fr...")
	logFile := flag.String("log", "", "optional log filename")
	dump := flag.Bool("dump", false, "dump http requests and responses")
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())

	logger := log.New(os.Stderr, "", log.LstdFlags)
	if len(*logFile) > 0 {
		file, err := os.Create(*logFile)
		if err != nil {
			log.Fatalln("unable to create file", *logFile)
		}
		defer file.Close()
		logger = log.New(io.MultiWriter(file, os.Stderr), "", log.LstdFlags)
	}
	logger.Println("login", *login)
	logger.Println("pass", "*n/a*")
	logger.Println("uni", *universe)
	logger.Println("lang", *lang)
	if len(*logFile) > 0 {
		logger.Println("log", *logFile)
	}

	b := makeOGBot(*login, *pass, *universe, *lang, logger, *dump)
	b.Run()
}
