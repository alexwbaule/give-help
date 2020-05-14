package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"firebase.google.com/go/auth"
	"github.com/360EntSecGroup-Skylar/excelize"

	"github.com/alexwbaule/give-help/v2/generated/models"
	proposalHandler "github.com/alexwbaule/give-help/v2/handlers/proposal"
	userHandler "github.com/alexwbaule/give-help/v2/handlers/user"
	"github.com/alexwbaule/give-help/v2/internal/common"
	runtimeApp "github.com/alexwbaule/give-help/v2/runtime"
	"github.com/alexwbaule/go-app"
)

type FirebaseUser struct {
	Email       string
	PhoneNumber string
	Password    string
	DisplayName string
	PhotoURL    string
	Line        int
}

func main() {
	Execute(os.Args[1])
}

func Execute(filename string) {
	app, err := app.New("give-help-service")

	rt, err := runtimeApp.NewRuntime(app)
	if err != nil {
		log.Fatal(err.Error())
	}

	defer rt.CloseDatabase()

	f, err := excelize.OpenFile(filename)

	if err != nil {
		log.Println(err)
		return
	}

	users, errs := importUsers(f, rt)

	for line, err := range errs {
		log.Printf("[line %d] Error to try import user: %s\n", line, err)
	}

	if len(users) == 0 {
		log.Println("[line %d] No users to import")
	}

	errs = importProps(users, f, rt)

	for line, err := range errs {
		log.Printf("[line %d] Error to try import proposal: %s\n", line, err)
	}

	log.Printf("Finish!")
}

func importUsers(f *excelize.File, rt *runtimeApp.Runtime) (map[string]*models.User, map[int]error) {
	userIds := map[string]*models.User{}
	errs := map[int]error{}

	h := map[string]int{
		"UserID":         -1,
		"Name":           -1,
		"Description":    -1,
		"Tags":           -1,
		"Images":         -1,
		"CreatedAt":      -1,
		"LastUpdate":     -1,
		"URL":            -1,
		"Email":          -1,
		"Facebook":       -1,
		"Instagram":      -1,
		"Twitter":        -1,
		"Address":        -1,
		"City":           -1,
		"State":          -1,
		"ZipCode":        -1,
		"Country":        -1,
		"Lat":            -1,
		"Long":           -1,
		"RegisterFrom":   -1,
		"PCountry":       -1,
		"PRegion":        -1,
		"PNumbers":       -1,
		"DeviceID":       -1,
		"AllowShareData": -1,
	}

	rows := f.GetRows("Users")

	if len(rows) < 2 {
		err := fmt.Errorf("[Users] no lines to read")
		log.Println(err)

		errs[0] = err

		return userIds, errs
	}

	svc := userHandler.New(rt.GetDatabase())

	for line, row := range rows {
		//header
		if line == 0 {
			for pos, name := range rows[0] {
				if _, found := h[name]; found {
					h[name] = pos
				}
			}
			continue
		}

		//row
		user := &models.User{
			Contact: &models.Contact{
				Phones: []*models.Phone{},
			},
			Location:   &models.Location{},
			Images:     []string{},
			Reputation: &models.Reputation{},
			Tags:       models.Tags{},
		}

		sheetId := getData("UserID", h, row)

		if len(sheetId) == 0 {
			sheetId = common.GetULID()
		}

		user.Name = getData("Name", h, row)
		user.Description = getData("Description", h, row)
		user.Tags = strings.Split(getData("Tags", h, row), ",")
		user.Images = strings.Split(getData("Images", h, row), ",")
		user.Contact.URL = getData("URL", h, row)
		user.Contact.Email = getData("Email", h, row)
		user.Contact.Facebook = getData("Facebook", h, row)
		user.Contact.Instagram = getData("Instagram", h, row)
		user.Contact.Twitter = getData("Twitter", h, row)
		user.Location.Address = getData("Address", h, row)
		user.Location.City = getData("City", h, row)
		user.Location.State = getData("State", h, row)
		*user.Location.ZipCode = getInt(getData("ZipCode", h, row))
		user.Location.Country = getData("Country", h, row)
		*user.Location.Lat = getFloat(getData("Lat", h, row))
		*user.Location.Long = getFloat(getData("Long", h, row))
		user.RegisterFrom = getData("RegisterFrom", h, row)
		user.AllowShareData = getBool(getData("AllowShareData", h, row))

		for _, phone := range strings.Split(getData("PNumbers", h, row), ",") {
			user.Contact.Phones = append(user.Contact.Phones, &models.Phone{
				CountryCode: getData("PCountry", h, row),
				Region:      getData("PRegion", h, row),
				PhoneNumber: phone,
			})
		}

		if len(user.Contact.Phones) > 0 {
			user.Contact.Phones[0].IsDefault = true
		}

		//get firebase user id
		id, err := upsertFirebase(line, user, rt)
		if err != nil {
			errs[line] = err
			continue
		}

		user.UserID = models.UserID(id)

		//insert user
		_, err = svc.Insert(user, id)
		if err != nil {
			errs[line] = err
			continue
		}

		userIds[sheetId] = user

		log.Printf("[line %d] [API] User added\n", line)
	}

	return userIds, errs
}

func importProps(users map[string]*models.User, f *excelize.File, rt *runtimeApp.Runtime) map[int]error {
	h := map[string]int{
		"ProposalID":       -1,
		"UserID":           -1,
		"Title":            -1,
		"Description":      -1,
		"Side":             -1,
		"ProposalType":     -1,
		"Tags":             -1,
		"IsActive":         -1,
		"CreatedAt":        -1,
		"LastUpdate":       -1,
		"ProposalValidate": -1,
		"AreaTags":         -1,
		"Lat":              -1,
		"Long":             -1,
		"Range":            -1,
		"Images":           -1,
		"EstimatedValue":   -1,
		"ExposeUserData":   -1,
		"DataToShare":      -1,
		"Ranking":          -1,
	}

	rows := f.GetRows("Proposals")

	if len(rows) < 2 {
		err := fmt.Errorf("[proposals] no lines to read")
		log.Println(err)

		return map[int]error{0: err}
	}

	ret := map[int]error{}
	svc := proposalHandler.New(rt.GetDatabase())

	for line, row := range rows {
		//header
		if line == 0 {
			for pos, name := range rows[0] {
				if _, found := h[name]; found {
					h[name] = pos
				}
			}
			continue
		}

		//row
		prop := &models.Proposal{
			DataToShare: []models.DataToShare{},
			Images:      []string{},
			Tags:        models.Tags{},
			TargetArea:  &models.Area{},
			IsActive:    true,
		}

		sheetUserID := getData("UserID", h, row)

		if len(sheetUserID) == 0 {
			ret[line] = fmt.Errorf("[line %d] invalid user id", line)
			continue
		}

		user, found := users[sheetUserID]

		if !found {
			ret[line] = fmt.Errorf("[line %d] user id (%s) not found!", line, sheetUserID)
			continue
		}

		dts := []models.DataToShare{}

		if len(user.Contact.Facebook) > 0 {
			dts = append(dts, models.DataToShareFacebook)
		}

		if len(user.Contact.Instagram) > 0 {
			dts = append(dts, models.DataToShareInstagram)
		}

		if len(user.Contact.Email) > 0 {
			dts = append(dts, models.DataToShareEmail)
		}

		if len(user.Contact.URL) > 0 {
			dts = append(dts, models.DataToShareURL)
		}

		if len(user.Contact.Phones) > 0 {
			dts = append(dts, models.DataToSharePhone)
		}

		prop.DataToShare = dts

		prop.Title = getData("Title", h, row)
		prop.Description = getData("Description", h, row)

		prop.Side = models.Side(getData("Side", h, row))
		prop.ProposalType = models.Type(getData("ProposalType", h, row))

		prop.Tags = strings.Split(getData("Tags", h, row), ",")
		prop.Images = strings.Split(getData("Images", h, row), ",")

		*prop.TargetArea.Lat = getFloat(getData("Lat", h, row))
		*prop.TargetArea.Long = getFloat(getData("Long", h, row))
		prop.TargetArea.Range = getFloat(getData("Range", h, row))
		prop.TargetArea.AreaTags = strings.Split(getData("AreaTags", h, row), ",")

		*prop.EstimatedValue = getFloat(getData("EstimatedValue", h, row))
		prop.ExposeUserData = getBool(getData("ExposeUserData", h, row))

		*prop.Ranking = getFloat(getData("Ranking", h, row))

		//insert proposal
		propId, err := svc.Insert(prop)
		if err != nil {
			ret[line] = err
			continue
		}

		log.Printf("[line %d] [API] Proposal added (%s)\n", line, propId)
	}

	return ret
}

func getData(name string, h map[string]int, row []string) string {
	pos := h[name]
	ret := ""
	if pos >= 0 {
		ret = strings.TrimSpace(row[pos])
	}
	return ret
}

func getFloat(input string) float64 {
	if ret, err := strconv.ParseFloat(input, 64); err == nil {
		return ret
	}

	return 0
}

func getInt(input string) int64 {
	if ret, err := strconv.ParseInt(input, 10, 64); err == nil {
		return ret
	}

	return 0
}

func getBool(input string) bool {
	if ret, err := strconv.ParseBool(input); err == nil {
		return ret
	}

	return false
}

func parserToFirebaseUser(line int, user *models.User) FirebaseUser {
	ret := FirebaseUser{Line: line}

	if len(user.Contact.Email) > 0 {
		ret.Email = user.Contact.Email
	} else {
		re := regexp.MustCompile(`[^A-Za-z0-9]`)
		replaced := re.ReplaceAll([]byte(user.Name), []byte(""))

		ret.Email = fmt.Sprintf("%s@%s.com", string(replaced), "give-help-importer")
	}

	ret.Password = ret.Email

	if len(user.Contact.Phones) > 0 {
		ret.PhoneNumber = fmt.Sprintf("%s%s%s", user.Contact.Phones[0].CountryCode, user.Contact.Phones[0].Region, user.Contact.Phones[0].PhoneNumber)
	}

	ret.DisplayName = user.Name

	if len(user.Images) > 0 {
		ret.PhotoURL = user.Images[0]
	}

	return ret
}

func upsertFirebase(line int, input *models.User, rt *runtimeApp.Runtime) (string, error) {
	user := parserToFirebaseUser(line, input)

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

	client, err := rt.GetFirebase().Auth(ctx)
	if err != nil {
		log.Fatalf("[line %d] [firebase] Error getting Auth client: %v\n", line, err)
	}

	u, err := client.GetUserByEmail(ctx, user.Email)
	if err == nil {
		log.Printf("[line %d] [firebase] User (id=%s) already exists on firebase\n", line, u.UID)
		return u.UID, err
	}

	u, err = client.CreateUser(ctx, params)

	if err != nil {
		dt, _ := json.MarshalIndent(user, "", "\t")
		log.Fatalf("[line %d] [firebase] Error creating user: data%s\nerror: %s", line, string(dt), err)
		return "", err
	}

	log.Printf("[line %d] [firebase] Successfully created user (id=%s)\n", line, u.UID)

	return u.UID, err
}
