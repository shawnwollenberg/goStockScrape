package stockscrape

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

//StockData will be built into JSON for output of scrapeStockHistory
type StockData struct {
	Symbol       string
	CompanyName  string
	DtStockPrice map[string]float64
}

func convertYDate(dt string) string {
	t, _ := time.Parse("Jan 02, 2006", dt)
	return t.Format("2006-01-02")
	//return string(t.String())
}

//ScrapeStockHistory retrieves the latest historical close prices for a symbol.  For Crypto add on "-USD" (ex: ETH-USD)
func ScrapeStockHistory(symbol string) []byte {
	//for crypto use ETH-USD as symbol ... if error try adding "-USD"
	fmt.Println("Starting " + symbol + "....")
	var coName string
	client := &http.Client{}
	dtClosePrice := make(map[string]float64)
	//response, err := http.Get("https://finance.yahoo.com/quote/" + symbol + "/history/")
	request, err := http.NewRequest("GET", "https://finance.yahoo.com/quote/"+symbol+"/history/", nil)
	if err != nil {
		fmt.Println(err.Error())
	}
	request.Header.Set("User-Agent", "Not Firefox")

	//response, err := http.Get(indexLookup[a][0])
	response, err := client.Do(request)
	if err != nil {
		///LOG + RETURN ERRORS
		fmt.Println(err.Error())
	}
	defer response.Body.Close()
	dataInBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err.Error())
	}
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
			if !strings.Contains(sCols[j], "<table") {
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
				dtClosePrice[convertYDate(sColsOutput[0])] = price
			}
		}
	}
	m := StockData{symbol, coName, dtClosePrice}
	b, _ := json.Marshal(m)
	return b
}
