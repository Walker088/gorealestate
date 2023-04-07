package plvr

import (
	"context"
	"fmt"
	"strings"

	"github.com/Walker088/gorealestate/common"
	e "github.com/Walker088/gorealestate/error"
	"github.com/gocarina/gocsv"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	DbInsertionError      = "PS00001"
	RocEraFormattingError = "PS00002"

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

type HouseSaleItem struct {
	SerialNumber                string `csv:"編號"`
	District                    string `csv:"鄉鎮市區" db:"district"`
	TransactionType             string `csv:"交易標的" db:"transaction_type"`
	Address                     string `csv:"土地位置建物門牌" db:"address"`
	LandShiftingArea            string `csv:"土地移轉總面積平方公尺" db:"land_shifting_area_sqm"`
	UrbanLandUse                string `csv:"都市土地使用分區" db:"urban_land_use"`
	NonUrbanLandUse             string `csv:"非都市土地使用分區" db:"non_urban_land_use"`
	NonUrbanLandDesignation     string `csv:"非都市土地使用編定" db:"non_urban_land_designation"`
	TransactionDateRaw          string `csv:"交易年月日" db:"transaction_date"`
	TransactionPenNumber        string `csv:"交易筆棟數" db:"transaction_pen_number"`
	Floor                       string `csv:"移轉層次" db:"floor"`
	TotalFloor                  string `csv:"總樓層數" db:"total_floor"`
	BuildingType                string `csv:"建物型態" db:"building_type"`
	PrimaryUse                  string `csv:"主要用途" db:"primary_use"`
	PrimaryMaterial             string `csv:"主要建材" db:"primary_material"`
	ConstructionCompleteDateRaw string `csv:"建築完成年月" db:"construction_complete_date"`
	BuildingAreaSqm             string `csv:"建物移轉總面積平方公尺" db:"building_area_sqm"`
	NumberOfRooms               int    `csv:"建物現況格局-房" db:"number_of_rooms"`
	NumberOfLivingRooms         int    `csv:"建物現況格局-廳" db:"number_of_living_rooms"`
	NumberOfBathrooms           int    `csv:"建物現況格局-衛" db:"number_of_bathrooms"`
	Partitioned                 string `csv:"建物現況格局-隔間" db:"partitioned"`
	HasManagementOrganization   string `csv:"有無管理組織" db:"has_management_organization"`
	TotalPrice                  int    `csv:"總價元" db:"total_price"`
	UnitPrice                   int    `csv:"單價元平方公尺" db:"unit_price_per_sqm"`
	ParkingType                 string `csv:"車位類別" db:"parking_type"`
	ParkingArea                 string `csv:"車位移轉總面積(平方公尺)" db:"parking_area_sqm"`
	ParkingPrice                int    `csv:"車位總價元" db:"parking_price"`
	Notes                       string `csv:"備註" db:"notes"`
	MainBuildingArea            string `csv:"主建物面積" db:"main_building_area_sqm"`
	SubsidiaryBuildingArea      string `csv:"附屬建物面積" db:"subsidiary_building_area_sqm"`
	BalconyArea                 string `csv:"陽台面積" db:"balcony_area_sqm"`
	Elevator                    string `csv:"電梯" db:"elevator"`
	TransactionIdentifier       string `csv:"移轉編號" db:"transaction_identifier"`
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
func (h HouseSaleItem) save(pool *pgxpool.Pool, city string) *e.ErrorData {
	query := `
	INSERT INTO plvr_land_house_sale (
		serial_number, city, district, transaction_type, address, land_shifting_area_sqm,
		urban_land_use, non_urban_land_use, non_urban_land_designation, transaction_date_raw, transaction_date, transaction_pen_number,
		floor, total_floor, building_type, primary_use, primary_material, 
		construction_complete_date_raw, construction_complete_date, building_area_sqm, number_of_rooms, number_of_living_rooms, 
		number_of_bathrooms, partitioned, has_management_organization, total_price, unit_price_per_sqm, 
		parking_type, parking_area_sqm, parking_price, notes, main_building_area_sqm, 
		subsidiary_building_area_sqm, balcony_area_sqm, elevator, transaction_identifier
	)
	VALUES (
		$1, $2, $3, $4, $5, $6, 
		$7, $8, $9, $10, $11, $12,
		$13, $14, $15, $16, $17, 
		$18, $19, $20, $21, $22, 
		$23, $24, $25, $26, $27,
		$28, $29, $30, $31, $32,
		$33, $34, $35, $36
	);
	`
	transacDate, _ := common.RocEraToCommonEra(h.TransactionDateRaw)
	constructDate, _ := common.RocEraToCommonEra(h.ConstructionCompleteDateRaw)
	if _, err := pool.Exec(
		context.Background(),
		query,
		h.SerialNumber, city, h.District, h.TransactionType, h.Address, h.LandShiftingArea,
		h.UrbanLandUse, h.NonUrbanLandUse, h.NonUrbanLandDesignation, h.TransactionDateRaw, transacDate, h.TransactionPenNumber,
		h.Floor, h.TotalFloor, h.BuildingType, h.PrimaryUse, h.PrimaryMaterial,
		h.ConstructionCompleteDateRaw, constructDate, h.BuildingAreaSqm, h.NumberOfRooms, h.NumberOfLivingRooms,
		h.NumberOfBathrooms, h.Partitioned, h.HasManagementOrganization, h.TotalPrice, h.UnitPrice,
		h.ParkingType, h.ParkingArea, h.ParkingPrice, h.Notes, h.MainBuildingArea,
		h.SubsidiaryBuildingArea, h.BalconyArea, h.Elevator, h.TransactionIdentifier,
	); err != nil {
		return e.NewErrorData(
			DbInsertionError,
			fmt.Sprintf("Error: %s on %s", err.Error(), h.toString(city)),
			fmt.Sprintf("%s.HouseSaleItem.save", currentPackage),
			nil,
			nil,
		)
	}
	return nil
}
func (h HouseSaleItem) toString(cityCode string) string {
	city := cityMap[cityCode]
	return fmt.Sprintf(`
	[HouseSaleItem] [%s] City=%s District=%s TransacType=%s TransacDate=%s
	`, h.SerialNumber, city, h.District, h.TransactionType, h.TransactionDateRaw)
}

type NewHouseItem struct {
	SerialNumber                string `csv:"編號"`
	District                    string `csv:"鄉鎮市區"`
	TransactionType             string `csv:"交易標的"`
	Address                     string `csv:"土地位置建物門牌"`
	LandShiftingArea            string `csv:"土地移轉總面積平方公尺"`
	UrbanLandUse                string `csv:"都市土地使用分區"`
	NonUrbanLandUse             string `csv:"非都市土地使用分區"`
	NonUrbanLandDesignation     string `csv:"非都市土地使用編定"`
	TransactionDateRaw          string `csv:"交易年月日"`
	TransactionPenNumber        string `csv:"交易筆棟數"`
	Floor                       string `csv:"移轉層次"`
	TotalFloor                  string `csv:"總樓層數"`
	BuildingType                string `csv:"建物型態"`
	PrimaryUse                  string `csv:"主要用途"`
	PrimaryMaterial             string `csv:"主要建材"`
	ConstructionCompleteDateRaw string `csv:"建築完成年月"`
	BuildingAreaSqm             string `csv:"建物移轉總面積平方公尺"`
	NumberOfRooms               int    `csv:"建物現況格局-房"`
	NumberOfLivingRooms         int    `csv:"建物現況格局-廳"`
	NumberOfBathrooms           int    `csv:"建物現況格局-衛"`
	Partitioned                 string `csv:"建物現況格局-隔間"`
	HasManagementOrganization   string `csv:"有無管理組織"`
	TotalPrice                  int    `csv:"總價元"`
	UnitPrice                   int    `csv:"單價元平方公尺"`
	ParkingType                 string `csv:"車位類別"`
	ParkingArea                 string `csv:"車位移轉總面積平方公尺"`
	ParkingPrice                int    `csv:"車位總價元"`
	Notes                       string `csv:"備註"`
}

func NewNewHouseItems(csvBytes []byte) ([]NewHouseItem, *e.ErrorData) {
	head := strings.Split(string(csvBytes), "\n")[0:1]
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
func (n NewHouseItem) save(pool *pgxpool.Pool, city string) *e.ErrorData {
	query := `
	INSERT INTO plvr_land_new_house (
		serial_number, city, district, transaction_type, address, land_shifting_area_sqm,
		urban_land_use, non_urban_land_use, non_urban_land_designation, transaction_date_raw, transaction_date, transaction_pen_number,
		floor, total_floor, building_type, primary_use, primary_material, 
		construction_complete_date_raw, construction_complete_date, building_area_sqm, number_of_rooms, number_of_living_rooms, 
		number_of_bathrooms, partitioned, has_management_organization, total_price, unit_price_per_sqm, 
		parking_type, parking_area_sqm, parking_price, notes
	)
	VALUES (
		$1, $2, $3, $4, $5, $6, 
		$7, $8, $9, $10, $11, $12, 
		$13, $14, $15, $16, $17,
		$18, $19, $20, $21, $22,
		$23, $24, $25, $26, $27,
		$28, $29, $30, $31
	);
	`
	transacDate, _ := common.RocEraToCommonEra(n.TransactionDateRaw)
	constructDate, _ := common.RocEraToCommonEra(n.ConstructionCompleteDateRaw)
	if _, err := pool.Exec(
		context.Background(),
		query,
		n.SerialNumber, city, n.District, n.TransactionType, n.Address, n.LandShiftingArea,
		n.UrbanLandUse, n.NonUrbanLandUse, n.NonUrbanLandDesignation, n.TransactionDateRaw, transacDate, n.TransactionPenNumber,
		n.Floor, n.TotalFloor, n.BuildingType, n.PrimaryUse, n.PrimaryMaterial,
		n.ConstructionCompleteDateRaw, constructDate, n.BuildingAreaSqm, n.NumberOfRooms, n.NumberOfLivingRooms,
		n.NumberOfBathrooms, n.Partitioned, n.HasManagementOrganization, n.TotalPrice, n.UnitPrice,
		n.ParkingType, n.ParkingArea, n.ParkingPrice, n.Notes,
	); err != nil {
		return e.NewErrorData(
			DbInsertionError,
			fmt.Sprintf("Error: %s on %s", err.Error(), n.toString(city)),
			fmt.Sprintf("%s.NewHouseItem.save", currentPackage),
			nil,
			nil,
		)
	}
	return nil
}
func (n NewHouseItem) toString(cityCode string) string {
	city := cityMap[cityCode]
	return fmt.Sprintf(`
	[NewHouseItem] [%s] City=%s District=%s TransacType=%s TransacDate=%s
	`, n.SerialNumber, city, n.District, n.TransactionType, n.TransactionDateRaw)
}

type RentalItem struct {
	SerialNumber                string `csv:"編號"`
	District                    string `csv:"鄉鎮市區"`
	TransactionType             string `csv:"交易標的"`
	Address                     string `csv:"土地位置建物門牌"`
	LandShiftingArea            string `csv:"土地面積平方公尺"`
	UrbanLandUse                string `csv:"都市土地使用分區"`
	NonUrbanLandUse             string `csv:"非都市土地使用分區"`
	NonUrbanLandDesignation     string `csv:"非都市土地使用編定"`
	TransactionDateRaw          string `csv:"租賃年月日"`
	TransactionPenNumber        string `csv:"租賃筆棟數"`
	Floor                       string `csv:"租賃層次"`
	TotalFloor                  string `csv:"總樓層數"`
	BuildingType                string `csv:"建物型態"`
	PrimaryUse                  string `csv:"主要用途"`
	PrimaryMaterial             string `csv:"主要建材"`
	ConstructionCompleteDateRaw string `csv:"建築完成年月"`
	BuildingAreaSqm             string `csv:"建物總面積平方公尺"`
	NumberOfRooms               int    `csv:"建物現況格局-房"`
	NumberOfLivingRooms         int    `csv:"建物現況格局-廳"`
	NumberOfBathrooms           int    `csv:"建物現況格局-衛"`
	Partitioned                 string `csv:"建物現況格局-隔間"`
	HasManagementOrganization   string `csv:"有無管理組織"`
	HasFurniture                string `csv:"有無附傢俱"`
	TotalPrice                  int    `csv:"總額元"`
	UnitPrice                   int    `csv:"單價元平方公尺"`
	ParkingType                 string `csv:"車位類別"`
	ParkingArea                 string `csv:"車位面積平方公尺"`
	ParkingPrice                int    `csv:"車位總額元"`
	Notes                       string `csv:"備註"`
}

func NewRentalItems(csvBytes []byte) ([]RentalItem, *e.ErrorData) {
	head := strings.Split(string(csvBytes), "\n")[0:1]
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
func (r RentalItem) save(pool *pgxpool.Pool, city string) *e.ErrorData {
	query := `
	INSERT INTO plvr_land_rental (
		serial_number, city, district, transaction_type, address, land_shifting_area_sqm,
		urban_land_use, non_urban_land_use, non_urban_land_designation, transaction_date_raw, transaction_date, transaction_pen_number,
		floor, total_floor, building_type, primary_use, primary_material, 
		construction_complete_date_raw, construction_complete_date, building_area_sqm, number_of_rooms, number_of_living_rooms, 
		number_of_bathrooms, partitioned, has_management_organization, total_price, unit_price_per_sqm, 
		parking_type, parking_area_sqm, parking_price, notes
	)
	VALUES (
		$1, $2, $3, $4, $5, $6, 
		$7, $8, $9, $10, $11, $12,
		$13, $14, $15, $16, $17,
		$18, $19, $20, $21, $22,
		$23, $24, $25, $26, $27,
		$28, $29, $30, $31
	);
	`
	transacDate, _ := common.RocEraToCommonEra(r.TransactionDateRaw)
	constructDate, _ := common.RocEraToCommonEra(r.ConstructionCompleteDateRaw)
	if _, err := pool.Exec(
		context.Background(),
		query,
		r.SerialNumber, city, r.District, r.TransactionType, r.Address, r.LandShiftingArea,
		r.UrbanLandUse, r.NonUrbanLandUse, r.NonUrbanLandDesignation, r.TransactionDateRaw, transacDate, r.TransactionPenNumber,
		r.Floor, r.TotalFloor, r.BuildingType, r.PrimaryUse, r.PrimaryMaterial,
		r.ConstructionCompleteDateRaw, constructDate, r.BuildingAreaSqm, r.NumberOfRooms, r.NumberOfLivingRooms,
		r.NumberOfBathrooms, r.Partitioned, r.HasManagementOrganization, r.TotalPrice, r.UnitPrice,
		r.ParkingType, r.ParkingArea, r.ParkingPrice, r.Notes,
	); err != nil {
		return e.NewErrorData(
			DbInsertionError,
			fmt.Sprintf("Error: %s on %s", err.Error(), r.toString(city)),
			fmt.Sprintf("%s.NewHouseItem.save", currentPackage),
			nil,
			nil,
		)
	}
	return nil
}
func (r RentalItem) toString(cityCode string) string {
	city := cityMap[cityCode]
	return fmt.Sprintf(`
	[RentalItem] [%s] City=%s District=%s TransacType=%s TransacDate=%s
	`, r.SerialNumber, city, r.District, r.TransactionType, r.TransactionDateRaw)
}
