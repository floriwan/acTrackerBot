package types

type AdsbExchData struct {
	Ac    []ac `json:"ac"`
	Total int  `json:"total"`
	Ctime int  `json:"ctime"`
	Now   int  `json:"now"`
	Ptime int  `json:"ptime"`
}

type ac struct {
	Hex              string   `json:"hex"`
	Actype           string   `json:"actype"`
	Flight           string   `json:"flight"`
	R                string   `json:"r"`
	T                string   `json:"t"`
	Alt_baro         int      `json:"alt_baro"`
	Alt_geom         int      `json:"alt_geom"`
	Gs               float32  `json:"gs"`
	Track            float32  `json:"track"`
	Baro_rate        int      `json:"baro_rate"`
	Squawk           string   `json:"squawk"`
	Emergency        string   `json:"emergency"`
	Category         string   `json:"category"`
	Nav_qnh          float32  `json:"nav_qnh"`
	Nav_altitude_mcp int      `json:"nav_altitude_mcp"`
	Nav_heading      float32  `json:"nav_heading"`
	Lat              float32  `json:"lat"`
	Lon              float32  `json:"lon"`
	Nic              int      `json:"nic"`
	Rc               int      `json:"rc"`
	Seen_pos         float32  `json:"seen_pos"`
	Version          int      `json:"version"`
	Nic_baro         int      `json:"nic_baro"`
	Nac_p            int      `json:"nac_p"`
	Nac_v            int      `json:"nac_v"`
	Sil              int      `json:"sil"`
	Sil_type         string   `json:"sil_type"`
	Gva              int      `json:"gva"`
	Sda              int      `json:"sda"`
	Alert            int      `json:"alert"`
	Spi              int      `json:"spi"`
	Mlat             []string `json:"mlat"`
	Tisb             []string `json:"tisb"`
	Messages         int      `json:"messages"`
	Seen             float32  `json:"seen"`
	Rssi             float32  `json:"rssi"`
}
