package alpha2

import (
	"encoding/json"
	"io/ioutil"
	"last_lesson/internal/config"
	"last_lesson/internal/mytypes"
	"last_lesson/internal/sub"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const (
	wikiURL string = "https://ru.wikipedia.org/wiki/ISO_3166-1"
)

func GetAlpha2Codes(sync bool) ([]string, error) {
	jsonData, err := GetAlpha2(sync)
	if err != nil {
		return []string{}, err
	}

	var codes []string
	for key := range jsonData {
		codes = append(codes, key)
	}

	return codes, err
}

func GetAlpha2(sync bool) (mytypes.Alpha2Codes, error) {
	config, err := config.Enviroment()
	if err != nil {
		return mytypes.Alpha2Codes{}, err
	}

	Alpha2CodesFile := config.Alpha2File

	ok, err := sub.FileExists(Alpha2CodesFile)
	if err != nil {
		return mytypes.Alpha2Codes{}, err
	}

	if sync || !ok {
		if err = Sync(Alpha2CodesFile); err != nil {
			return mytypes.Alpha2Codes{}, err
		}
	}

	byteData, err := ioutil.ReadFile(Alpha2CodesFile)
	if err != nil {
		return mytypes.Alpha2Codes{}, err
	}

	jsonData := make(mytypes.Alpha2Codes)
	if json.Unmarshal(byteData, &jsonData) != nil {
		return mytypes.Alpha2Codes{}, err
	}

	return jsonData, nil
}

func Sync(path string) error {
	data, err := json.Marshal(parseWiki())
	if err != nil {
		return err
	}

	ioutil.WriteFile(path, data, 0644)

	return nil
}

func Alpha2ToCountrySMS(data *[]mytypes.SMSData) error {
	alpha2, err := GetAlpha2(false)
	if err != nil {
		return err
	}

	for i := range *(data) {
		(*data)[i].Country = alpha2[(*data)[i].Country].Country
	}

	return nil
}

func Alpha2ToCountryMMS(data *[]mytypes.MMSData) error {
	alpha2, err := GetAlpha2(false)
	if err != nil {
		return err
	}

	for i := range *(data) {
		(*data)[i].Country = alpha2[(*data)[i].Country].Country
	}

	return nil
}

func parseWiki() mytypes.Alpha2Codes {
	resp, err := http.Get(wikiURL)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	codes := make(mytypes.Alpha2Codes)
	doc.Find("table.wikitable tr").Each(
		func(trId int, tr *goquery.Selection) {
			newRow := mytypes.Alpha2Row{
				Country: "",
				Alpha3:  "",
				ISO1:    "",
			}
			alpha2Id := ""

			tr.Find("td").Each(func(tdId int, td *goquery.Selection) {
				tdData := strings.TrimSpace(td.Text())

				switch tdId {
				case 0:
					newRow.Country = tdData
				case 1:
					alpha2Id = tdData
				case 2:
					newRow.Alpha3 = tdData
				case 3:
					newRow.ISO1 = tdData
				}
			})

			codes[alpha2Id] = newRow
		})

	return codes
}
