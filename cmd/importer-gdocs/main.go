package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/alexwbaule/give-help/v2/generated/models"
	"github.com/alexwbaule/give-help/v2/internal/common"
	"github.com/go-openapi/strfmt"
	"google.golang.org/api/option"

	app "github.com/alexwbaule/go-app"

	proposalHandler "github.com/alexwbaule/give-help/v2/handlers/proposal"
	tagsHandler "github.com/alexwbaule/give-help/v2/handlers/tags"
	userHandler "github.com/alexwbaule/give-help/v2/handlers/user"
	runtimeApp "github.com/alexwbaule/give-help/v2/runtime"
)

type Proposal struct {
	ProposaldID    string
	DeviceID       string
	Timestamp      time.Time
	Name           string
	Email          string
	Lat            float64
	Long           float64
	SheetType      string
	Type           string // volunteer|taker|job|local_business
	Side           string
	Tags           []string // categoria
	Description    string
	URL            string
	Address        string
	PhoneNumbers   []string
	PhoneRegion    string
	PhoneCountry   string
	AllowShareData bool
	Images         []string
	Facebook       string
	Instagram      string
	Line           int
	Ranking        float64
}

type FirebaseUser struct {
	Email       string
	PhoneNumber string
	Password    string
	DisplayName string
	PhotoURL    string
	Line        int
}

var rt *runtimeApp.Runtime
var fbase *firebase.App
var userSvc *userHandler.User
var proposalSvc *proposalHandler.Proposal
var tagsSvc *tagsHandler.Tags

func main() {
	log.Printf("Import from tsf gdocs - start\n")
	for i, a := range os.Args {
		log.Printf("Arg %d: %s\n", i, a)
	}

	app, err := app.New("give-help-service")
	cfg := app.Config()

	cfg.SetDefault("service.Host", "127.0.0.1")
	cfg.SetDefault("service.Port", "8081")
	cfg.SetDefault("service.TLSWriteTimeout", "15m")
	cfg.SetDefault("service.WriteTimeout", "15m")

	rt, err = runtimeApp.NewRuntime(app)
	if err != nil {
		log.Fatal(err.Error())
	}

	defer rt.CloseDatabase()

	props, err := loadFromFile(os.Args[1], 4)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("read %d lines - Done!\n", len(props))

	j, err := json.MarshalIndent(props, "", "\t")

	err = ioutil.WriteFile(os.Args[1]+".output.json", j, 0644)

	log.Printf("output in %s\n", os.Args[1]+".output.json")

	initFirebase()

	userSvc = userHandler.New(rt.GetDatabase())
	proposalSvc = proposalHandler.New(rt.GetDatabase())
	tagsSvc = tagsHandler.New(rt.GetDatabase())

	for _, p := range props {
		u := parserToFirebaseUser(p)
		id, err := insertFirebase(u)

		if err != nil {
			log.Printf("[ERROR] fail to try add user on firebase: %s\n", err)
			continue
		}

		err = insertDbTags(p.Tags)
		if err != nil {
			log.Printf("[ERROR] [id=%s] fail: %s\n", id, err)
		}

		insertDbUser(p, id)
		if err != nil {
			log.Printf("[ERROR] [id=%s] fail: %s\n", id, err)
		}

		insertDbProposal(p, id)
		if err != nil {
			log.Printf("[ERROR] [id=%s] fail: %s\n", id, err)
		}

		log.Printf("[line=%d] [id=%s] Import ok!\n", p.Line, id)
	}
	log.Printf("Import from tsf gdocs - done!\n")
}

func showError(prop Proposal) {
	j, _ := json.MarshalIndent(prop, "", "\t")
	log.Printf("Error on: \n%s", string(j))
}

func loadFromFile(path string, offset int) ([]Proposal, error) {
	ret := []Proposal{}

	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
		return ret, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	line := 0
	for scanner.Scan() {
		if line > offset {
			data := scanner.Text()
			u, err := parser(data, line)

			if err != nil {
				log.Printf("[line: %d] - Fail to parser data: %s", line, data)
			} else {
				log.Printf("[line: %d] - line parsed!", line)
			}

			if len(u.Name) > 0 {
				ret = append(ret, u)
			}
		}
		line++
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
		return ret, err
	}

	return ret, err
}

func parser(line string, index int) (Proposal, error) {
	fields := strings.Split(line, "\t")

	if len(fields) < 17 {
		return Proposal{}, fmt.Errorf("invalid input format (expected at least 17 fields)")
	}

	ret := Proposal{
		ProposaldID:    fields[0],
		DeviceID:       fields[1],
		Timestamp:      getTime(fields[2]),
		Name:           fields[3],
		Email:          fields[4],
		Lat:            getFloat(fields[5]),
		Long:           getFloat(fields[6]),
		SheetType:      fields[7],
		Tags:           getArray(fields[8]),
		Description:    fields[9],
		URL:            fields[10],
		Address:        fields[11],
		PhoneNumbers:   getPhoneNumbers(fields[12]),
		PhoneRegion:    getPhoneRegion(fields[13]),
		PhoneCountry:   getPhoneCountry(fields[14]),
		AllowShareData: getBool(fields[15]),
		Images:         getArray(fields[16]),
		Facebook:       parserURL(fields[10], "facebook"),
		Instagram:      parserURL(fields[10], "instagram"),
		Line:           index,
		Ranking:        getFloat(fields[17]),
	}

	t, s := getType(fields[7], ret)

	ret.Type = t
	ret.Side = s

	if len(ret.Facebook) > 0 || len(ret.Instagram) > 0 {
		ret.URL = ""
	}

	//animals
	found := 0
	for _, t := range common.NormalizeTagArray(ret.Tags) {
		switch t {
		case "animais":
		case "dogs":
		case "gatos":
			found++
		}
	}

	//Natal
	if strings.Contains(ret.Description, "Natal") {
		found++
	}

	if found == 0 {
		ret.Ranking = 1
		for _, t := range common.NormalizeTagArray(ret.Tags) {
			switch t {
			case "crianças":
				ret.Ranking += 3
			case "alimentação":
				ret.Ranking += 2
			case "saúde":
				ret.Ranking += 1
			}
		}
	}

	return ret, nil
}

func getPhoneRegion(input string) string {
	if len(input) == 0 {
		return "11"
	}

	return input
}

func getPhoneCountry(input string) string {
	if len(input) == 0 {
		return "+55"
	}

	if string(input[0]) != "+" {
		input = "+" + input
	}

	return input
}

func getPhoneNumbers(input string) []string {
	re := regexp.MustCompile(`[^0-9]`)

	ret := []string{}

	for _, t := range getArray(input) {
		ret = append(ret, string(re.ReplaceAll([]byte(t), []byte(""))))
	}

	return ret
}

func getTime(input string) time.Time {
	return time.Now()
}

func getFloat(input string) float64 {
	if len(input) == 0 {
		return 0.0
	}

	if ret, err := strconv.ParseFloat(input, 64); err == nil {
		return ret
	}

	return 0.0
}

func getArray(input string) []string {
	ret := []string{}
	for _, v := range strings.Split(input, ",") {
		item := strings.TrimSpace(v)

		if len(item) > 0 {
			ret = append(ret, item)
		}
	}

	return ret
}

func getBool(input string) bool {
	i := strings.ToUpper(input)

	switch i {
	case "NÃO":
	case "NAO":
	case "N":
		return false
	}

	return true
}

const (
	InputTypeVolunteer     string = "voluteer"
	InputTypeLocalBusiness string = "local_business"
	InputTypeJob           string = "job"
	InputTypeTaker         string = "taker"
)

func getType(input string, prop Proposal) (string, string) {
	t := models.TypeService
	s := models.SideRequest

	switch strings.ToLower(input) {
	case strings.ToLower(InputTypeVolunteer):
		t = models.TypeService
		s = models.SideRequest
	case strings.ToLower(InputTypeLocalBusiness):
		t = models.TypeService
		s = models.SideLocalBusiness
	case strings.ToLower(InputTypeJob):
		t = models.TypeJob
		s = models.SideRequest
	case strings.ToLower(InputTypeTaker):
		t = models.TypeJob
		s = models.SideOffer
	}

	return string(t), string(s)
}

func parserURL(input string, target string) string {
	if strings.Contains(strings.ToLower(input), strings.ToLower(target)) {
		return strings.ToLower(input)
	}

	return ""
}

func parserToFirebaseUser(prop Proposal) FirebaseUser {
	ret := FirebaseUser{Line: prop.Line}

	if len(prop.Email) > 0 {
		ret.Email = prop.Email
	} else {
		re := regexp.MustCompile(`[^A-Za-z0-9]`)
		replaced := re.ReplaceAll([]byte(prop.Name), []byte(""))

		ret.Email = fmt.Sprintf("%s@%s.com", string(replaced), string(re.ReplaceAll([]byte(prop.SheetType), []byte(""))))
	}

	ret.Password = ret.Email

	if len(prop.PhoneNumbers) > 0 {
		if len(prop.PhoneNumbers[0]) > 0 {
			ret.PhoneNumber = fmt.Sprintf("%s%s%s", prop.PhoneCountry, prop.PhoneRegion, prop.PhoneNumbers[0])
		}
	}

	ret.DisplayName = prop.Name

	if len(prop.Images) > 0 {
		ret.PhotoURL = prop.Images[0]
	}

	return ret
}

func initFirebase() {
	opt := option.WithCredentialsFile("etc/serviceAccountKey.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)

	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	fbase = app
}

func insertFirebase(user FirebaseUser) (string, error) {
	params := &auth.UserToCreate{}
	params.Email(user.Email)
	params.EmailVerified(false)

	if len(user.PhoneNumber) > 0 {
		params.PhoneNumber(user.PhoneNumber)
	}

	params.Password(user.Password)
	params.DisplayName(user.DisplayName)

	if len(user.PhotoURL) > 0 {
		params.PhotoURL(user.PhotoURL)
	}

	params.Disabled(false)

	ctx := context.Background()

	client, err := fbase.Auth(ctx)
	if err != nil {
		log.Fatalf("error getting Auth client: %v\n", err)
	}

	u, err := client.GetUserByEmail(ctx, user.Email)
	if err == nil {
		log.Printf("[id=%s] User already exists on firebase\n", u.UID)
		return u.UID, err
	}

	u, err = client.CreateUser(ctx, params)

	if err != nil {
		dt, _ := json.MarshalIndent(user, "", "\t")
		log.Fatalf("error creating user: data%s\nerror: %s", string(dt), err)
		return "", err
	}

	log.Printf("[id=%s] Successfully created user\n", u.UID)

	return u.UID, err
}

func insertDbTags(tags []string) error {
	return tagsSvc.Insert(tags)
}

func insertDbUser(prop Proposal, userID string) error {
	phones := []*models.Phone{}

	for _, p := range prop.PhoneNumbers {
		phones = append(phones, &models.Phone{
			CountryCode: prop.PhoneCountry,
			Region:      prop.PhoneRegion,
			PhoneNumber: p,
		})
	}

	if len(phones) > 0 {
		phones[0].IsDefault = true
	}

	contact := &models.Contact{
		Email:     prop.Email,
		Facebook:  prop.Facebook,
		Instagram: prop.Instagram,
		URL:       prop.URL,
		Phones:    phones,
	}
	location := &models.Location{
		Address: prop.Address,
		City:    "São Paulo",
		Country: "Brasil",
		State:   "São Paulo",
	}

	data := &models.User{
		AllowShareData: prop.AllowShareData,
		Contact:        contact,
		Description:    prop.Description,
		DeviceID:       prop.DeviceID,
		Images:         prop.Images,
		Location:       location,
		Name:           prop.Name,
		RegisterFrom:   "admin",
		Tags:           prop.Tags,
		UserID:         models.UserID(userID),
	}

	err := userSvc.Update(data)

	if err != nil {
		log.Printf("[id=%s] error to try insert user: %s", userID, err)

		return err
	}

	return err
}

func insertDbProposal(prop Proposal, userID string) error {
	dts := []models.DataToShare{}

	if len(prop.Facebook) > 0 {
		dts = append(dts, models.DataToShareFacebook)
	}

	if len(prop.Instagram) > 0 {
		dts = append(dts, models.DataToShareInstagram)
	}

	if len(prop.Email) > 0 {
		dts = append(dts, models.DataToShareEmail)
	}

	if len(prop.URL) > 0 {
		dts = append(dts, models.DataToShareURL)
	}

	if len(prop.PhoneNumbers) > 0 {
		dts = append(dts, models.DataToSharePhone)
	}

	retID, err := proposalSvc.Insert(&models.Proposal{
		ProposalValidate: strfmt.DateTime(time.Now().AddDate(10, 0, 0)),
		DataToShare:      dts,
		Description:      prop.Description,
		ExposeUserData:   prop.AllowShareData,
		Images:           prop.Images,
		ProposalType:     models.Type(prop.Type),
		Side:             models.Side(prop.Side),
		Tags:             prop.Tags,
		Title:            prop.Name,
		UserID:           models.UserID(userID),
		IsActive:         true,
		Ranking:          &prop.Ranking,
	})

	if err != nil {
		log.Printf("[id=%s] error to try insert proposal: %s", retID, err)

		return err
	}

	return err
}
