package stockscrape

/*
import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestStockScrapeBad(t *testing.T) {
	want := `{"ERROR":"No Name For This"}`
	if got := string(ScrapeStockHistory("Shawn")); got != want {
		t.Errorf("ScrapeStockHistory(\"Shawn\") = %q, want %q", got, want)
	}
}

func TestStockScrapeCompanyName(t *testing.T) {
	want := "Mettler-Toledo International Inc. (MTD)"
	var s StockData
	x := ScrapeStockHistory("MTD")
	err := json.Unmarshal(x, &s)
	if err != nil {
		fmt.Println(err.Error())
	}
	if got := s.CompanyName; got != want {
		t.Errorf("ScrapeStockHistory(\"MTD\") Company Name = %q, want %q", got, want)
	}

}

func TestDate(t *testing.T) {
	want := "2020-04-19"
	if got := convertYDate("Apr 19, 2020"); got != want {
		t.Errorf("convertYDate(\"Apr 19, 2020\") = %q, want %q", got, want)
	}
	convertYDate("Apr 19, 2020")
}
*/
