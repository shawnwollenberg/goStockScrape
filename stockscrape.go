package StockScrape

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

//StockData will be built into JSON for output of scrapeStockHistory
type StockData struct {
	Symbol       string
	CompanyName  string
	DtStockPrice map[string]float64
}

//ScrapeStockHistory retrieves the latest historical close prices for a symbol.  For Crypto add on "-USD" (ex: ETH-USD)
func ScrapeStockHistory(symbol string) []byte {
	//for crypto use ETH-USD as symbol ... if error try adding "-USD"
	fmt.Println("Starting " + symbol + "....")
	var coName string
	dtClosePrice := make(map[string]float64)
	response, err := http.Get("https://finance.yahoo.com/quote/" + symbol + "/history/")
	if err != nil {
		///LOG + RETURN ERRORS
		fmt.Println(err.Error())
	}
	defer response.Body.Close()
	dataInBytes, err := ioutil.ReadAll(response.Body)
	pageContent := string(dataInBytes)
	coNameStart := strings.Index(pageContent, "<h1 class=") //only one instance of 'h1 class=' and </h1
	if coNameStart != -1 {
		coNameEnd := strings.Index(pageContent, "</h1")
		coName = pageContent[coNameStart+1+strings.Index(pageContent[coNameStart:coNameEnd], ">") : coNameEnd]
	}
	if len(coName) > 0 {
		if coName[0:1] == " " {
			coName = ""
		}
	}
	if coName == "" { //not a valid company name... try running with "-USD" to see if it's a crypto
		if !strings.Contains(symbol, "-USD") {
			return ScrapeStockHistory(symbol + "-USD")
		}
		return []byte(`{"ERROR":"No Name For This"}`)
	}
	tblStart := strings.Index(pageContent, "<table")
	if tblStart == -1 {
		///LOG + RETURN ERRORS
		return []byte(`{"ERROR":"No opening element found"}`)
	}

	tblEnd := strings.Index(pageContent, "</tbody")
	if tblEnd == -1 {
		///LOG + RETURN ERRORS
		return []byte(`{"ERROR":"No closing tag found"}`)
	}

	dataTbl := pageContent[tblStart:tblEnd]
	sRows := strings.Split(dataTbl, "<tr")
	for i := 0; i < len(sRows); i++ {
		sCols := strings.Split(sRows[i], "<span")
		holdStr := ""
		for j := 0; j < len(sCols); j++ {
			if strings.Index(sCols[j], "<table") < 0 {
				colStart := strings.Index(sCols[j], ">") + 1
				colEnd := strings.Index(sCols[j], "</span")
				if colEnd > 0 {
					colOutput := sCols[j][colStart:colEnd]
					holdStr += colOutput + "|"
				}
			}
		}
		if i > 1 && len(holdStr) > 3 && strings.Index(holdStr, "Dividend") < 1 {
			strings.Split(holdStr, "|")
			//fmt.Println(len(holdStr), "-", holdStr)
			sColsOutput := strings.Split(holdStr, "|")
			if len(sColsOutput) > 5 {
				price, err := strconv.ParseFloat(strings.Replace(sColsOutput[5], ",", "", -1), 64)
				if err != nil {
					price = 0
				}
				dtClosePrice[sColsOutput[0]] = price
			}

			//insert += "Insert into " + dbconnect.DbName + ".dbo.[00T_Stock_Price_Load] (Symbol,[Date],[Close],Volume) Values('" + indexLookup[a][1] + "','" + sColsOutput[0] + "'," + strings.Replace(sColsOutput[5], ",", "", -1) + "," + strings.Replace(sColsOutput[6], ",", "", -1) + ");"
		}
	}
	m := StockData{symbol, coName, dtClosePrice}
	b, err := json.Marshal(m)
	return b
}
