package plvr

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	e "github.com/Walker088/gorealestate/error"
	"github.com/gocarina/gocsv"
)

var (
	cityMap = map[string]string{
		"a": "Taipei", "b": "Taichung", "c": "Keelung", "d": "Tainan",
		"e": "Kaohsiung", "f": "New Taipei", "g": "Yilan", "h": "Taoyuan",
		"j": "Hsinchu Country", "k": "Miaoli", "l": "Taichung Country", "m": "Nantou",
		"n": "Changhua", "p": "Yunlin", "q": "Chiayi Country", "r": "Tainan County",
		"s": "Kaohsiung County", "t": "Pingtung", "u": "Hualien", "v": "Taitung",
		"x": "Penghu", "y": "Yangmingshan", "w": "Kinmen", "z": "Lianjiang",
		"i": "Chiayi", "o": "Hsinchu",
	}
)

type ParsedItem interface {
	save(string) *e.ErrorData
	toString(string) string
}

type HouseSaleItem struct {
	SerialNumber string `csv:"編號"`
	//City                      string `csv:"-"`
	District                  string `csv:"鄉鎮市區"`
	TransactionType           string `csv:"交易標的"`
	Address                   string `csv:"土地位置建物門牌"`
	LandShiftingArea          string `csv:"土地移轉總面積平方公尺"`
	UrbanLandUse              string `csv:"都市土地使用分區"`
	NonUrbanLandUse           string `csv:"非都市土地使用分區"`
	NonUrbanLandDesignation   string `csv:"非都市土地使用編定"`
	TransactionDate           string `csv:"交易年月日"`
	TransactionPenNumber      string `csv:"交易筆棟數"`
	Floor                     string `csv:"移轉層次"`
	TotalFloor                string `csv:"總樓層數"`
	BuildingType              string `csv:"建物型態"`
	PrimaryUse                string `csv:"主要用途"`
	PrimaryMaterial           string `csv:"主要建材"`
	ConstructionCompleteDate  string `csv:"建築完成年月"`
	BuildingAreaSqm           string `csv:"建物移轉總面積平方公尺"`
	NumberOfRooms             string `csv:"建物現況格局-房"`
	NumberOfLivingRooms       string `csv:"建物現況格局-廳"`
	NumberOfBathrooms         string `csv:"建物現況格局-衛"`
	Partitioned               string `csv:"建物現況格局-隔間"`
	HasManagementOrganization string `csv:"有無管理組織"`
	TotalPrice                string `csv:"總價元"`
	UnitPrice                 string `csv:"單價元平方公尺"`
	ParkingType               string `csv:"車位類別"`
	ParkingArea               string `csv:"車位移轉總面積(平方公尺)"`
	ParkingPrice              string `csv:"車位總價元"`
	Notes                     string `csv:"備註"`
	MainBuildingArea          string `csv:"主建物面積"`
	SubsidiaryBuildingArea    string `csv:"附屬建物面積"`
	BalconyArea               string `csv:"陽台面積"`
	Elevator                  string `csv:"電梯"`
	TransactionIdentifier     string `csv:"移轉編號"`
}

func NewHouseSaleItems(csvBytes []byte) ([]HouseSaleItem, *e.ErrorData) {
	head := strings.Split(string(csvBytes), "\n")[0:1]
	body := strings.Split(string(csvBytes), "\n")[2:]
	cleaned := head[0] + "\n" + strings.Join(body, "\n")

	items := []HouseSaleItem{}

	if err := gocsv.UnmarshalBytes([]byte(cleaned), &items); err != nil {
		return nil, e.NewErrorData(
			UnmarshalCsvError,
			err.Error(),
			fmt.Sprintf("%s.parse", currentPackage),
			nil,
			nil,
		)
	}
	return items, nil
}
func (h *HouseSaleItem) save(city string) *e.ErrorData {
	return nil
}
func (h *HouseSaleItem) toString(cityCode string) string {
	city := cityMap[cityCode]
	return fmt.Sprintf(`
	[HouseSaleItem] [%s] City=%s District=%s TransacType=%s TransacDate=%s
	`, h.SerialNumber, city, h.District, h.TransactionType, h.TransactionDate)
}

type NewHouseItem struct {
	SerialNumber string `csv:"serial number"`
	//City                      string `csv:"-"`
	District                  string `csv:"The villages and towns urban district"`
	TransactionType           string `csv:"transaction sign"`
	Address                   string `csv:"land sector position building sector house number plate"`
	LandShiftingArea          string `csv:"land shifting total area square meter"`
	UrbanLandUse              string `csv:"the use zoning or compiles and checks"`
	NonUrbanLandUse           string `csv:"the non-metropolis land use district"`
	NonUrbanLandDesignation   string `csv:"non-metropolis land use"`
	TransactionDate           string `csv:"transaction year month and day"`
	TransactionPenNumber      string `csv:"transaction pen number"`
	Floor                     string `csv:"shifting level"`
	TotalFloor                string `csv:"total floor number"`
	BuildingType              string `csv:"building state"`
	PrimaryUse                string `csv:"main use"`
	PrimaryMaterial           string `csv:"main building materials"`
	ConstructionCompleteDate  string `csv:"construction to complete the years"`
	BuildingAreaSqm           string `csv:"building shifting total area"`
	NumberOfRooms             string `csv:"Building present situation pattern - room"`
	NumberOfLivingRooms       string `csv:"building present situation pattern - hall"`
	NumberOfBathrooms         string `csv:"building present situation pattern - health"`
	Partitioned               string `csv:"building present situation pattern - compartmented"`
	HasManagementOrganization string `csv:"Whether there is manages the organization"`
	TotalPrice                string `csv:"total price NTD"`
	UnitPrice                 string `csv:"the unit price (NTD / square meter)"`
	ParkingType               string `csv:"the berth category"`
	ParkingArea               string `csv:"berth shifting total area square meter"`
	ParkingPrice              string `csv:"the berth total price NTD"`
	Notes                     string `csv:"the note"`
}

func NewNewHouseItems(csvBytes []byte) ([]NewHouseItem, *e.ErrorData) {
	head := strings.Split(string(csvBytes), "\n")[1:2]
	body := strings.Split(string(csvBytes), "\n")[2:]
	cleaned := head[0] + "\n" + strings.Join(body, "\n")

	items := []NewHouseItem{}

	if err := gocsv.UnmarshalBytes([]byte(cleaned), &items); err != nil {
		return nil, e.NewErrorData(
			UnmarshalCsvError,
			err.Error(),
			fmt.Sprintf("%s.parse", currentPackage),
			nil,
			nil,
		)
	}
	return items, nil
}
func (n *NewHouseItem) save(city string) *e.ErrorData {
	return nil
}
func (n *NewHouseItem) toString(cityCode string) string {
	city := cityMap[cityCode]
	return fmt.Sprintf(`
	[NewHouseItem] [%s] City=%s District=%s TransacType=%s TransacDate=%s
	`, n.SerialNumber, city, n.District, n.TransactionType, n.TransactionDate)
}

type RentalItem struct {
	SerialNumber string `csv:"serial number"`
	//City                      string `csv:"-"`
	District                  string `csv:"The villages and towns urban district"`
	TransactionType           string `csv:"transaction sign"`
	Address                   string `csv:"land sector position building sector house number plate"`
	LandShiftingArea          string `csv:"land shifting total area square meter"`
	UrbanLandUse              string `csv:"the use zoning or compiles and checks"`
	NonUrbanLandUse           string `csv:"the non-metropolis land use district"`
	NonUrbanLandDesignation   string `csv:"non-metropolis land use"`
	TransactionDate           string `csv:"transaction year month and day"`
	TransactionPenNumber      string `csv:"transaction pen number"`
	Floor                     string `csv:"shifting level"`
	TotalFloor                string `csv:"total floor number"`
	BuildingType              string `csv:"building state"`
	PrimaryUse                string `csv:"main use"`
	PrimaryMaterial           string `csv:"main building materials"`
	ConstructionCompleteDate  string `csv:"construction to complete the years"`
	BuildingAreaSqm           string `csv:"building shifting total area"`
	NumberOfRooms             string `csv:"Building present situation pattern - room"`
	NumberOfLivingRooms       string `csv:"building present situation pattern - hall"`
	NumberOfBathrooms         string `csv:"building present situation pattern - health"`
	Partitioned               string `csv:"building present situation pattern - compartmented"`
	HasManagementOrganization string `csv:"Whether there is manages the organization"`
	HasFurniture              string `csv:"Whether there is attaches the furniture"`
	TotalPrice                string `csv:"total price NTD"`
	UnitPrice                 string `csv:"the unit price (NTD / square meter)"`
	ParkingType               string `csv:"the berth category"`
	ParkingArea               string `csv:"berth shifting total area square meter"`
	ParkingPrice              string `csv:"the berth total price NTD"`
	Notes                     string `csv:"the note"`
}

func NewRentalItems(csvBytes []byte) ([]RentalItem, *e.ErrorData) {
	head := strings.Split(string(csvBytes), "\n")[1:2]
	body := strings.Split(string(csvBytes), "\n")[2:]
	cleaned := head[0] + "\n" + strings.Join(body, "\n")

	items := []RentalItem{}

	if err := gocsv.UnmarshalBytes([]byte(cleaned), &items); err != nil {
		return nil, e.NewErrorData(
			UnmarshalCsvError,
			err.Error(),
			fmt.Sprintf("%s.parse", currentPackage),
			nil,
			nil,
		)
	}
	return items, nil
}
func (n *RentalItem) save(city string) *e.ErrorData {
	return nil
}
func (r *RentalItem) toString(cityCode string) string {
	city := cityMap[cityCode]
	return fmt.Sprintf(`
	[RentalItem] [%s] City=%s District=%s TransacType=%s TransacDate=%s
	`, r.SerialNumber, city, r.District, r.TransactionType, r.TransactionDate)
}

type DateTime struct {
	time.Time
}

// Input: ROC Era, e.g., 1011019 = 2012/10/19
func (date *DateTime) UnmarshalCSV(csv string) (err error) {
	if len(csv) == 0 {
		return nil
	}

	rocToCommonEra := func(era int) int {
		return era + 1911
	}

	rocYear, err := strconv.Atoi(csv[:len(csv)-4])
	if err != nil {
		return err
	}
	month, err := strconv.Atoi(csv[3 : len(csv)-2])
	if err != nil {
		return err
	}

	day, err := strconv.Atoi(csv[5:])
	if err != nil {
		return err
	}

	date.Time = time.Date(rocToCommonEra(int(rocYear)), time.Month(month), day, 0, 0, 0, 0, time.Local)
	return nil
}
