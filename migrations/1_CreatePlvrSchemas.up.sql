
CREATE TABLE IF NOT EXISTS ref_plvr_land_city (
  city_code VARCHAR(1) PRIMARY KEY,
  city_name_en TEXT,
  city_name_zh TEXT
);
TRUNCATE ref_plvr_land_city;
INSERT INTO ref_plvr_land_city
VALUES
('a', 'Taipei City', '台北市'),
('b', 'Taichung City', '台中市'),
('c', 'Keelung City', '基隆市'),
('d', 'Tainan City', '台南市'),
('e', 'Kaohsiung City', '高雄市'),
('f', 'New Taipei City', '新北市'),
('g', 'Yilan Country', '宜蘭縣'),
('h', 'Taoyuan City', '桃園市'),
('i', 'Chiayi City', '嘉義市'),
('j', 'Hsinchu Country', '新竹縣'),
('l', 'Taichung Country', '台中縣'),
('k', 'Miaoli Country', '苗栗縣'),
('m', 'Nantou Country', '南投縣'),
('n', 'Changhua Country', '彰化縣'),
('o', 'Hsinchu City', '新竹市'),
('p', 'Yunlin Country', '雲林縣'),
('q', 'Chiayi Country', '嘉義縣'),
('r', 'Tainan Country', '台南縣'),
('s', 'Kaohsiung Country', '高雄縣'),
('t', 'Pingtung Country', '屏東縣'),
('u', 'Hualien Country', '花蓮縣'),
('v', 'Taitung Country', '台東縣'),
('w', 'Kinmen Country', '金門縣'),
('x', 'Penghu Country', '澎湖縣'),
('y', 'Yangmingshan', '陽明山'),
('z', 'Lianjiang Country', '連江縣');

CREATE TABLE IF NOT EXISTS plvr_land_house_sale (
  serial_number TEXT,
  city TEXT,
  district TEXT,
  transaction_type TEXT,
  address TEXT,
  land_shifting_area_sqm DECIMAL(10,4),
  urban_land_use TEXT,
  non_urban_land_use TEXT,
  non_urban_land_designation TEXT,
  transaction_date_raw TEXT,
  transaction_date DATE,
  transaction_pen_number TEXT,
  floor TEXT,
  total_floor TEXT,
  building_type TEXT,
  primary_use TEXT,
  primary_material TEXT,
  construction_complete_date_raw TEXT,
  construction_complete_date DATE,
  building_area_sqm DECIMAL(10,4),
  number_of_rooms INT2,
  number_of_living_rooms INT2,
  number_of_bathrooms INT2,
  partitioned TEXT,
  has_management_organization TEXT,
  total_price INT8,
  unit_price_per_sqm INT4,
  parking_type TEXT,
  parking_area_sqm DECIMAL(10,4),
  parking_price INT8,
  notes TEXT,
  main_building_area_sqm DECIMAL(10,4),
  subsidiary_building_area_sqm DECIMAL(10,4),
  balcony_area_sqm DECIMAL(10,4),
  elevator TEXT,
  transaction_identifier TEXT
);
COMMENT ON TABLE plvr_land_house_sale IS '實價登錄 - 房屋買賣交易';
COMMENT ON COLUMN plvr_land_house_sale.transaction_type IS '交易標的';
COMMENT ON COLUMN plvr_land_house_sale.address IS '土地位置建物門牌';
COMMENT ON COLUMN plvr_land_house_sale.land_shifting_area_sqm IS '土地移轉總面積平方公尺';
COMMENT ON COLUMN plvr_land_house_sale.urban_land_use IS '都市土地使用分區';
COMMENT ON COLUMN plvr_land_house_sale.non_urban_land_use IS '非都市土地使用分區';
COMMENT ON COLUMN plvr_land_house_sale.non_urban_land_designation IS '非都市土地使用編定';
COMMENT ON COLUMN plvr_land_house_sale.transaction_pen_number IS '交易筆棟數';
COMMENT ON COLUMN plvr_land_house_sale.floor IS '移轉層次';
COMMENT ON COLUMN plvr_land_house_sale.total_floor IS '總樓層數';
COMMENT ON COLUMN plvr_land_house_sale.building_type IS '建物型態';
COMMENT ON COLUMN plvr_land_house_sale.primary_use IS '主要用途';
COMMENT ON COLUMN plvr_land_house_sale.primary_material IS '主要建材';
COMMENT ON COLUMN plvr_land_house_sale.building_area_sqm IS '建物移轉總面積平方公尺';
COMMENT ON COLUMN plvr_land_house_sale.partitioned IS '建物現況格局-隔間';
COMMENT ON COLUMN plvr_land_house_sale.parking_type IS '車位類別';
COMMENT ON COLUMN plvr_land_house_sale.main_building_area_sqm IS '主建物面積';
COMMENT ON COLUMN plvr_land_house_sale.subsidiary_building_area_sqm IS '附屬建物面積';
COMMENT ON COLUMN plvr_land_house_sale.transaction_identifier IS '移轉編號';

CREATE TABLE IF NOT EXISTS plvr_land_new_house (
  serial_number TEXT,
  city TEXT,
  district TEXT,
  transaction_type TEXT,
  address TEXT,
  land_shifting_area_sqm DECIMAL(10,4),
  urban_land_use TEXT,
  non_urban_land_use TEXT,
  non_urban_land_designation TEXT,
  transaction_date_raw TEXT,
  transaction_date DATE,
  transaction_pen_number TEXT,
  floor TEXT,
  total_floor TEXT,
  building_type TEXT,
  primary_use TEXT,
  primary_material TEXT,
  construction_complete_date_raw TEXT,
  construction_complete_date DATE,
  building_area_sqm DECIMAL(10,4),
  number_of_rooms INTEGER,
  number_of_living_rooms INTEGER,
  number_of_bathrooms INTEGER,
  partitioned TEXT,
  has_management_organization TEXT,
  total_price INTEGER,
  unit_price_per_sqm INTEGER,
  parking_type TEXT,
  parking_area_sqm DECIMAL(10,4),
  parking_price INTEGER,
  notes TEXT
);
COMMENT ON TABLE plvr_land_new_house IS '實價登錄 - 新成屋交易';
COMMENT ON COLUMN plvr_land_new_house.transaction_type IS '交易標的';
COMMENT ON COLUMN plvr_land_new_house.address IS '土地位置建物門牌';
COMMENT ON COLUMN plvr_land_new_house.land_shifting_area_sqm IS '土地移轉總面積平方公尺';
COMMENT ON COLUMN plvr_land_new_house.urban_land_use IS '都市土地使用分區';
COMMENT ON COLUMN plvr_land_new_house.non_urban_land_use IS '非都市土地使用分區';
COMMENT ON COLUMN plvr_land_new_house.non_urban_land_designation IS '非都市土地使用編定';
COMMENT ON COLUMN plvr_land_new_house.transaction_pen_number IS '交易筆棟數';
COMMENT ON COLUMN plvr_land_new_house.floor IS '移轉層次';
COMMENT ON COLUMN plvr_land_new_house.total_floor IS '總樓層數';
COMMENT ON COLUMN plvr_land_new_house.building_type IS '建物型態';
COMMENT ON COLUMN plvr_land_new_house.primary_use IS '主要用途';
COMMENT ON COLUMN plvr_land_new_house.primary_material IS '主要建材';
COMMENT ON COLUMN plvr_land_new_house.building_area_sqm IS '建物移轉總面積平方公尺';
COMMENT ON COLUMN plvr_land_new_house.partitioned IS '建物現況格局-隔間';
COMMENT ON COLUMN plvr_land_new_house.parking_type IS '車位類別';

CREATE TABLE IF NOT EXISTS plvr_land_rental (
  serial_number TEXT,
  city TEXT,
  district TEXT,
  transaction_type TEXT,
  address TEXT,
  land_shifting_area_sqm DECIMAL(10,4),
  urban_land_use TEXT,
  non_urban_land_use TEXT,
  non_urban_land_designation TEXT,
  transaction_date_raw TEXT,
  transaction_date DATE,
  transaction_pen_number TEXT,
  floor TEXT,
  total_floor TEXT,
  building_type TEXT,
  primary_use TEXT,
  primary_material TEXT,
  construction_complete_date_raw TEXT,
  construction_complete_date DATE,
  building_area_sqm DECIMAL(10,4),
  number_of_rooms INTEGER,
  number_of_living_rooms INTEGER,
  number_of_bathrooms INTEGER,
  partitioned TEXT,
  has_management_organization TEXT,
  total_price INTEGER,
  unit_price_per_sqm INTEGER,
  parking_type TEXT,
  parking_area_sqm DECIMAL(10,4),
  parking_price INTEGER,
  notes TEXT
);
COMMENT ON TABLE plvr_land_rental IS '實價登錄 - 租房交易';
COMMENT ON COLUMN plvr_land_rental.transaction_type IS '交易標的';
COMMENT ON COLUMN plvr_land_rental.address IS '土地位置建物門牌';
COMMENT ON COLUMN plvr_land_rental.land_shifting_area_sqm IS '土地移轉總面積平方公尺';
COMMENT ON COLUMN plvr_land_rental.urban_land_use IS '都市土地使用分區';
COMMENT ON COLUMN plvr_land_rental.non_urban_land_use IS '非都市土地使用分區';
COMMENT ON COLUMN plvr_land_rental.non_urban_land_designation IS '非都市土地使用編定';
COMMENT ON COLUMN plvr_land_rental.transaction_pen_number IS '交易筆棟數';
COMMENT ON COLUMN plvr_land_rental.floor IS '移轉層次';
COMMENT ON COLUMN plvr_land_rental.total_floor IS '總樓層數';
COMMENT ON COLUMN plvr_land_rental.building_type IS '建物型態';
COMMENT ON COLUMN plvr_land_rental.primary_use IS '主要用途';
COMMENT ON COLUMN plvr_land_rental.primary_material IS '主要建材';
COMMENT ON COLUMN plvr_land_rental.building_area_sqm IS '建物移轉總面積平方公尺';
COMMENT ON COLUMN plvr_land_rental.partitioned IS '建物現況格局-隔間';
COMMENT ON COLUMN plvr_land_rental.parking_type IS '車位類別';

CREATE TABLE IF NOT EXISTS plvr_land_parse_failed (
  serial_number TEXT,
  city TEXT,
  district TEXT,
  transaction_type TEXT,
  address TEXT,
  land_shifting_area_sqm TEXT,
  urban_land_use TEXT,
  non_urban_land_use TEXT,
  non_urban_land_designation TEXT,
  transaction_date_raw TEXT,
  transaction_date DATE,
  transaction_pen_number TEXT,
  floor TEXT,
  total_floor TEXT,
  building_type TEXT,
  primary_use TEXT,
  primary_material TEXT,
  construction_complete_date_raw TEXT,
  construction_complete_date DATE,
  building_area_sqm TEXT,
  number_of_rooms TEXT,
  number_of_living_rooms TEXT,
  number_of_bathrooms TEXT,
  partitioned TEXT,
  has_management_organization TEXT,
  total_price TEXT,
  unit_price_per_sqm TEXT,
  parking_type TEXT,
  parking_area_sqm TEXT,
  parking_price TEXT,
  notes TEXT,
  main_building_area_sqm TEXT,
  subsidiary_building_area_sqm TEXT,
  balcony_area_sqm TEXT,
  elevator TEXT,
  transaction_identifier TEXT
);

CREATE TABLE IF NOT EXISTS plvr_download_history (
  remote_addr     TEXT PRIMARY KEY,
  downloaded_time TIMESTAMP WITH TIME ZONE
);
