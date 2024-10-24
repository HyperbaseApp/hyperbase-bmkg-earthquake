package model

import "time"

type AutoGempaModel struct {
	Infogempa InfoGempaModel `json:"Infogempa"`
}

func (a *AutoGempaModel) ToMap() map[string]any {
	return a.Infogempa.Gempa.ToMap()
}

type DataModel struct {
	Infogempa InfoGempaListModel `json:"Infogempa"`
}

func (d *DataModel) ToSliceOfMap() []map[string]any {
	data := make([]map[string]any, 0, len(d.Infogempa.Gempa))
	for _, g := range d.Infogempa.Gempa {
		data = append(data, g.ToMap())
	}
	return data
}

type InfoGempaModel struct {
	Gempa GempaModel `json:"gempa"`
}

type InfoGempaListModel struct {
	Gempa []GempaModel `json:"gempa"`
}

type GempaModel struct {
	Tanggal     string    `json:"Tanggal"`     // ex: 24 Agu 2024
	Jam         string    `json:"Jam"`         // ex: 23:54:34 WIB
	DateTime    time.Time `json:"DateTime"`    // ex: 2024-08-24T16:54:34+00:00
	Coordinates string    `json:"Coordinates"` // ex: -4.43,102.18
	Lintang     string    `json:"Lintang"`     // ex: 4.43 LS
	Bujur       string    `json:"Bujur"`       // ex: 102.18 BT
	Magnitude   string    `json:"Magnitude"`   // ex: 5.2
	Kedalaman   string    `json:"Kedalaman"`   // ex: 21 km
	Wilayah     string    `json:"Wilayah"`     // ex: Pusat gempa berada di laut 59km barat daya Seluma
	Potensi     string    `json:"Potensi"`     // ex: Gempa ini dirasakan untuk diteruskan pada masyarakat
	Dirasakan   string    `json:"Dirasakan"`   // ex: III-IV Kota Bengkulu, III-IV Bengkulu Utara, III Kaur, II - III Empat Lawang, II - III Pagar Alam
	Shakemap    string    `json:"Shakemap"`    // ex: 20240824235434.mmi.jpg
}

func (g *GempaModel) ToMap() map[string]any {
	return map[string]any{
		"tanggal":     g.Tanggal,
		"jam":         g.Jam,
		"datetime":    g.DateTime,
		"coordinates": g.Coordinates,
		"lintang":     g.Lintang,
		"bujur":       g.Bujur,
		"magnitude":   g.Magnitude,
		"kedalaman":   g.Kedalaman,
		"wilayah":     g.Wilayah,
		"potensi":     g.Potensi,
		"dirasakan":   g.Dirasakan,
		"shakemap":    g.Shakemap,
		"shakemap_url": func() string {
			if g.Shakemap != "" {
				return "https://data.bmkg.go.id/DataMKG/TEWS/" + g.Shakemap
			}
			return ""
		}(),
	}
}
