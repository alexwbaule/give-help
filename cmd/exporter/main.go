package main

import (
	"fmt"
	"log"
	"strings"

	app "github.com/alexwbaule/go-app"

	"github.com/360EntSecGroup-Skylar/excelize"

	"github.com/alexwbaule/give-help/v2/generated/models"

	proposalHandler "github.com/alexwbaule/give-help/v2/handlers/proposal"
	tagsHandler "github.com/alexwbaule/give-help/v2/handlers/tags"
	userHandler "github.com/alexwbaule/give-help/v2/handlers/user"
	runtimeApp "github.com/alexwbaule/give-help/v2/runtime"
)

const outputFileName string = "output.xlsx"

func main() {
	app, err := app.New("give-help-service")

	rt, err := runtimeApp.NewRuntime(app)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Printf("Database config:\n\tHost:   %s\n\tDBName: %s\n\tUser:   %s\n",
		app.Config().GetString("database.Host"),
		app.Config().GetString("database.DBName"),
		app.Config().GetString("database.User"),
	)

	defer rt.CloseDatabase()

	f := createOutputFile()

	exportUsers(rt, f)
	exportProps(rt, f)
	exportTags(rt, f)

	f.SetActiveSheet(0)

	if err := f.SaveAs(outputFileName); err != nil {
		log.Println(err)
	}

	log.Printf("Done! see: %s", outputFileName)
}

func createOutputFile() *excelize.File {
	f := excelize.NewFile()

	f.NewSheet("Users")
	f.NewSheet("Proposals")
	f.NewSheet("Tags")
	f.DeleteSheet("Sheet1")

	return f
}

func exportUsers(rt *runtimeApp.Runtime, f *excelize.File) {
	userSvc := userHandler.New(rt.GetDatabase())
	users, err := userSvc.LoadAll()
	if err != nil {
		log.Fatal(err.Error())
	}

	userHeader := map[string]string{
		"UserID":         "A%d",
		"Name":           "B%d",
		"Description":    "C%d",
		"Tags":           "D%d",
		"Images":         "E%d",
		"CreatedAt":      "F%d",
		"LastUpdate":     "G%d",
		"URL":            "H%d",
		"Email":          "I%d",
		"Facebook":       "J%d",
		"Instagram":      "K%d",
		"Google":         "L%d",
		"Address":        "M%d",
		"City":           "N%d",
		"State":          "O%d",
		"ZipCode":        "P%d",
		"Country":        "Q%d",
		"Lat":            "R%d",
		"Long":           "S%d",
		"RegisterFrom":   "T%d",
		"PCountry":       "U%d",
		"PRegion":        "V%d",
		"PNumbers":       "W%d",
		"DeviceID":       "X%d",
		"AllowShareData": "Y%d",
	}

	currentLine := 1

	for k, v := range userHeader {
		f.SetCellValue("Users", fmt.Sprintf(v, currentLine), k)
	}

	for _, u := range users {
		if u.Contact == nil {
			u.Contact = &models.Contact{
				Phones: []*models.Phone{},
			}
		}

		if u.Location == nil {
			u.Location = &models.Location{}
		}

		currentLine++

		f.SetCellValue("Users", fmt.Sprintf(userHeader["UserID"], currentLine), u.UserID)
		f.SetCellValue("Users", fmt.Sprintf(userHeader["Name"], currentLine), u.Name)
		f.SetCellValue("Users", fmt.Sprintf(userHeader["Description"], currentLine), u.Description)
		f.SetCellValue("Users", fmt.Sprintf(userHeader["Tags"], currentLine), strings.Join(u.Tags, ","))
		f.SetCellValue("Users", fmt.Sprintf(userHeader["Images"], currentLine), strings.Join(u.Images, ","))
		f.SetCellValue("Users", fmt.Sprintf(userHeader["CreatedAt"], currentLine), u.CreatedAt)
		f.SetCellValue("Users", fmt.Sprintf(userHeader["LastUpdate"], currentLine), u.LastUpdate)
		f.SetCellValue("Users", fmt.Sprintf(userHeader["URL"], currentLine), u.Contact.URL)
		f.SetCellValue("Users", fmt.Sprintf(userHeader["Email"], currentLine), u.Contact.Email)
		f.SetCellValue("Users", fmt.Sprintf(userHeader["Facebook"], currentLine), u.Contact.Facebook)
		f.SetCellValue("Users", fmt.Sprintf(userHeader["Instagram"], currentLine), u.Contact.Instagram)
		f.SetCellValue("Users", fmt.Sprintf(userHeader["Google"], currentLine), u.Contact.Google)
		f.SetCellValue("Users", fmt.Sprintf(userHeader["Address"], currentLine), u.Location.Address)
		f.SetCellValue("Users", fmt.Sprintf(userHeader["City"], currentLine), u.Location.City)
		f.SetCellValue("Users", fmt.Sprintf(userHeader["State"], currentLine), u.Location.State)
		f.SetCellValue("Users", fmt.Sprintf(userHeader["ZipCode"], currentLine), u.Location.ZipCode)
		f.SetCellValue("Users", fmt.Sprintf(userHeader["Country"], currentLine), u.Location.Country)
		f.SetCellValue("Users", fmt.Sprintf(userHeader["Lat"], currentLine), u.Location.Lat)
		f.SetCellValue("Users", fmt.Sprintf(userHeader["Long"], currentLine), u.Location.Long)
		f.SetCellValue("Users", fmt.Sprintf(userHeader["RegisterFrom"], currentLine), u.RegisterFrom)
		f.SetCellValue("Users", fmt.Sprintf(userHeader["PCountry"], currentLine), getPhoneCountryCode(u))
		f.SetCellValue("Users", fmt.Sprintf(userHeader["PRegion"], currentLine), getPhoneRegion(u))
		f.SetCellValue("Users", fmt.Sprintf(userHeader["PNumbers"], currentLine), getPhones(u))
		f.SetCellValue("Users", fmt.Sprintf(userHeader["DeviceID"], currentLine), u.DeviceID)
		f.SetCellValue("Users", fmt.Sprintf(userHeader["AllowShareData"], currentLine), u.AllowShareData)
	}

	log.Printf("Exported %d users", len(users))
}

func exportProps(rt *runtimeApp.Runtime, f *excelize.File) {
	proposalSvc := proposalHandler.New(rt.GetDatabase())
	props, err := proposalSvc.LoadFromFilter(&models.Filter{
		PageSize: 50000,
	})

	if err != nil {
		log.Fatal(err.Error())
	}

	currentLine := 1

	propsHeader := map[string]string{
		"ProposalID":       "A%d",
		"UserID":           "B%d",
		"Title":            "C%d",
		"Description":      "D%d",
		"Side":             "E%d",
		"ProposalType":     "F%d",
		"Tags":             "G%d",
		"IsActive":         "H%d",
		"CreatedAt":        "I%d",
		"LastUpdate":       "J%d",
		"ProposalValidate": "K%d",
		"AreaTags":         "L%d",
		"Lat":              "M%d",
		"Long":             "N%d",
		"Range":            "O%d",
		"Images":           "P%d",
		"EstimatedValue":   "Q%d",
		"ExposeUserData":   "R%d",
		"DataToShare":      "S%d",
		"Ranking":          "T%d",
	}

	for k, v := range propsHeader {
		f.SetCellValue("Proposals", fmt.Sprintf(v, currentLine), k)
	}

	for _, p := range props.Result {
		currentLine++

		if p.TargetArea == nil {
			p.TargetArea = &models.Area{}
		}

		dts := []string{}

		for _, d := range p.DataToShare {
			dts = append(dts, string(d))
		}

		f.SetCellValue("Proposals", fmt.Sprintf(propsHeader["ProposalID"], currentLine), p.ProposalID)
		f.SetCellValue("Proposals", fmt.Sprintf(propsHeader["UserID"], currentLine), p.UserID)
		f.SetCellValue("Proposals", fmt.Sprintf(propsHeader["Title"], currentLine), p.Title)
		f.SetCellValue("Proposals", fmt.Sprintf(propsHeader["Description"], currentLine), p.Description)
		f.SetCellValue("Proposals", fmt.Sprintf(propsHeader["Side"], currentLine), p.Side)
		f.SetCellValue("Proposals", fmt.Sprintf(propsHeader["ProposalType"], currentLine), p.ProposalType)
		f.SetCellValue("Proposals", fmt.Sprintf(propsHeader["Tags"], currentLine), strings.Join(p.Tags, ","))
		f.SetCellValue("Proposals", fmt.Sprintf(propsHeader["IsActive"], currentLine), p.IsActive)
		f.SetCellValue("Proposals", fmt.Sprintf(propsHeader["CreatedAt"], currentLine), p.CreatedAt)
		f.SetCellValue("Proposals", fmt.Sprintf(propsHeader["LastUpdate"], currentLine), p.LastUpdate)
		f.SetCellValue("Proposals", fmt.Sprintf(propsHeader["ProposalValidate"], currentLine), p.ProposalValidate)
		f.SetCellValue("Proposals", fmt.Sprintf(propsHeader["AreaTags"], currentLine), strings.Join(p.TargetArea.AreaTags, ","))
		f.SetCellValue("Proposals", fmt.Sprintf(propsHeader["Lat"], currentLine), p.TargetArea.Lat)
		f.SetCellValue("Proposals", fmt.Sprintf(propsHeader["Long"], currentLine), p.TargetArea.Long)
		f.SetCellValue("Proposals", fmt.Sprintf(propsHeader["Range"], currentLine), p.TargetArea.Range)
		f.SetCellValue("Proposals", fmt.Sprintf(propsHeader["Images"], currentLine), strings.Join(p.Images, ","))
		f.SetCellValue("Proposals", fmt.Sprintf(propsHeader["EstimatedValue"], currentLine), p.EstimatedValue)
		f.SetCellValue("Proposals", fmt.Sprintf(propsHeader["ExposeUserData"], currentLine), p.ExposeUserData)
		f.SetCellValue("Proposals", fmt.Sprintf(propsHeader["DataToShare"], currentLine), strings.Join(dts, ","))
		f.SetCellValue("Proposals", fmt.Sprintf(propsHeader["Ranking"], currentLine), p.Ranking)
	}

	log.Printf("Exported %d proposals", len(props.Result))
}

func exportTags(rt *runtimeApp.Runtime, f *excelize.File) {
	tagsSvc := tagsHandler.New(rt.GetDatabase())
	tags, err := tagsSvc.Load()
	if err != nil {
		log.Fatal(err.Error())
	}

	currentLine := 1

	tagsHeader := map[string]string{
		"Tag": "A%d",
	}

	f.SetCellValue("Tags", fmt.Sprintf(tagsHeader["Tag"], currentLine), "Tags")

	for _, t := range tags {
		currentLine++
		f.SetCellValue("Tags", fmt.Sprintf(tagsHeader["Tag"], currentLine), t)
	}

	log.Printf("Exported %d tags", len(tags))
}

func getPhones(user *models.User) string {
	ret := ""
	if user.Contact != nil {
		if len(user.Contact.Phones) > 0 {
			arr := []string{}
			for _, p := range user.Contact.Phones {
				arr = append(arr, p.PhoneNumber)
			}

			ret = strings.Join(arr, ",")
		}
	}

	return ret
}

func getPhoneCountryCode(user *models.User) string {
	ret := "+55"
	if user.Contact != nil {
		if len(user.Contact.Phones) > 0 {
			for _, p := range user.Contact.Phones {
				if p.IsDefault {
					ret = p.CountryCode
				}
			}
		}
	}

	return ret
}

func getPhoneRegion(user *models.User) string {
	ret := "11"
	if user.Contact != nil {
		if len(user.Contact.Phones) > 0 {
			for _, p := range user.Contact.Phones {
				if p.IsDefault {
					ret = p.Region
				}
			}
		}
	}

	return ret
}
