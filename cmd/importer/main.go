package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize"

	"github.com/alexwbaule/give-help/v2/generated/models"
)

func main() {
	Execute(os.Args[1])
}

func Execute(filename string) {
	/*
		app, err := app.New("give-help-service")

		rt, err := runtimeApp.NewRuntime(app)
		if err != nil {
			log.Fatal(err.Error())
		}

		defer rt.CloseDatabase()
	*/
	//userSvc := userHandler.New(rt.GetDatabase())

	f, err := excelize.OpenFile(filename)

	if err != nil {
		log.Println(err)
		return
	}

	importUsers("Users", f)

	log.Printf("Done!")
}

func importUsers(sheetName string, f *excelize.File) map[int]error {
	h := map[string]int{
		"UserID":         0,
		"Name":           0,
		"Description":    0,
		"Tags":           0,
		"Images":         0,
		"CreatedAt":      0,
		"LastUpdate":     0,
		"URL":            0,
		"Email":          0,
		"Facebook":       0,
		"Instagram":      0,
		"Google":         0,
		"Address":        0,
		"City":           0,
		"State":          0,
		"ZipCode":        0,
		"Country":        0,
		"Lat":            0,
		"Long":           0,
		"RegisterFrom":   0,
		"PCountry":       0,
		"PRegion":        0,
		"PNumbers":       0,
		"DeviceID":       0,
		"AllowShareData": 0,
	}

	rows := f.GetRows(sheetName)

	if len(rows) < 2 {
		err := fmt.Errorf("no lines to read")
		log.Println(err)

		return map[int]error{0: err}
	}

	ret := map[int]error{}

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
		user := models.User{
			Contact: &models.Contact{
				Phones: []*models.Phone{},
			},
			Location:   &models.Location{},
			Images:     []string{},
			Reputation: &models.Reputation{},
			Tags:       models.Tags{},
		}

		user.UserID = models.UserID(row[h["UserID"]])
		user.Name = row[h["Name"]]
		user.Description = row[h["Description"]]
		user.Tags = strings.Split(row[h["Tags"]], ",")
		user.Images = strings.Split(row[h["Images"]], ",")
		user.Contact.URL = row[h["URL"]]
		user.Contact.Email = row[h["Email"]]
		user.Contact.Facebook = row[h["Facebook"]]
		user.Contact.Instagram = row[h["Instagram"]]
		user.Contact.Google = row[h["Google"]]
		user.Location.Address = row[h["Address"]]
		user.Location.City = row[h["City"]]
		user.Location.State = row[h["State"]]
		user.Location.ZipCode = getInt(row[h["ZipCode"]])
		user.Location.Country = row[h["Country"]]
		user.Location.Lat = getFloat(row[h["Lat"]])
		user.Location.Long = getFloat(row[h["Long"]])
		user.RegisterFrom = row[h["RegisterFrom"]]
		user.AllowShareData = getBool(row[h["AllowShareData"]])

		for _, phone := range strings.Split(row[h["PNumbers"]], ",") {
			user.Contact.Phones = append(user.Contact.Phones, &models.Phone{
				CountryCode: row[h["PCountry"]],
				Region:      row[h["PRegion"]],
				PhoneNumber: phone,
			})
		}

		if len(user.Contact.Phones) > 0 {
			user.Contact.Phones[0].IsDefault = true
		}

		/*
			//import user
			if err != nil {
				ret[i] = err
			}
		*/
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
