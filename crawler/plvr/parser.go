package plvr

import (
	"strconv"
	"time"
)

type RealEstateItem struct {
	SerialNumber              string   `csv:"編號"`
	City                      string   `csv:"-"`
	District                  string   `csv:"鄉鎮市區"`
	TransactionType           string   `csv:"交易標的"`
	Address                   string   `csv:"土地位置建物門牌"`
	LandShiftingArea          string   `csv:"土地移轉總面積平方公尺"`
	UrbanLandUse              string   `csv:"都市土地使用分區"`
	NonUrbanLandUse           string   `csv:"非都市土地使用分區"`
	NonUrbanLandDesignation   string   `csv:"非都市土地使用編定"`
	TransactionDate           DateTime `csv:"交易年月日"`
	TransactionPenNumber      string   `csv:"交易筆棟數"`
	Floor                     string   `csv:"移轉層次"`
	TotalFloor                string   `csv:"總樓層數"`
	BuildingType              string   `csv:"建物型態"`
	PrimaryUse                string   `csv:"主要用途"`
	PrimaryMaterial           string   `csv:"主要建材"`
	ConstructionCompleteDate  DateTime `csv:"建築完成年月"`
	BuildingAreaSqm           string   `csv:"建物移轉總面積平方公尺"`
	NumberOfRooms             uint     `csv:"建物現況格局-房"`
	NumberOfLivingRooms       uint     `csv:"建物現況格局-廳"`
	NumberOfBathrooms         uint     `csv:"建物現況格局-衛"`
	Partitioned               string   `csv:"建物現況格局-隔間"`
	HasManagementOrganization string   `csv:"有無管理組織"`
	TotalPrice                uint64   `csv:"總價元"`
	UnitPrice                 uint     `csv:"單價元平方公尺"`
	ParkingType               string   `csv:"車位類別"`
	ParkingArea               string   `csv:"車位移轉總面積(平方公尺)"`
	ParkingPrice              uint64   `csv:"車位總價元"`
	Notes                     string   `csv:"備註"`
	MainBuildingArea          string   `csv:"主建物面積"`
	SubsidiaryBuildingArea    string   `csv:"附屬建物面積"`
	BalconyArea               string   `csv:"陽台面積"`
	Elevator                  string   `csv:"電梯"`
	TransactionIdentifier     string   `csv:"移轉編號"`
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
