package xml

// ScorePartwise Основная структура партитуры
type ScorePartwise struct {
	Version string `xml:"version,attr"`
	Parts   []Part `xml:"part"`
}

// Part Часть партитуры
type Part struct {
	ID              string           `xml:"id,attr"`
	Measures        []Measure        `xml:"measure"`
	Instrument      *Instrument      `xml:"instrument"`
	ScoreInstrument *ScoreInstrument `xml:"score-instrument"`
}

// Measure Мера в музыке (тактовый размер)
type Measure struct {
	Number     string      `xml:"number,attr"`
	Attributes *Attributes `xml:"attributes"`
	Notes      []Note      `xml:"note"`
	Legato     []Legato    `xml:"legato"`
	Tempo      *Tempo      `xml:"direction>sound"`
}

// Attributes Атрибуты, такие как ключ, время, нотный стан
type Attributes struct {
	Key   *Key   `xml:"key"`
	Time  *Time  `xml:"time"`
	Clef  *Clef  `xml:"clef"`
	Stave *Stave `xml:"stave"`
}

type Key struct {
	Fifths int `xml:"fifths"`
}

type Time struct {
	Beats    int `xml:"beats"`
	BeatType int `xml:"beat-type"`
}

type Clef struct {
	Sign string `xml:"sign"`
	Line int    `xml:"line"`
}

type Stave struct {
	StaveLine int `xml:"stave-line"`
}

// Note Нота в музыке
type Note struct {
	Pitch         *Pitch         `xml:"pitch"`
	Duration      int            `xml:"duration"`
	Type          string         `xml:"type"`
	Voice         int            `xml:"voice"`
	Chord         bool           `xml:"chord,omitempty"`
	Tie           *Tie           `xml:"tie"`
	Rest          *Rest          `xml:"rest"`
	Articulations []Articulation `xml:"articulations"`
}

type Pitch struct {
	Step   string `xml:"step"`
	Octave int    `xml:"octave"`
	Alter  int    `xml:"alter,omitempty"`
}

// Tie Связь для нот
type Tie struct {
	Type string `xml:"type,attr"`
}

// Rest Отдых (пауза)
type Rest struct {
	Duration int `xml:"duration"`
}

// Articulation Артикуляции для нот
type Articulation struct {
	Staccato bool `xml:"staccato,omitempty"`
	Accent   bool `xml:"accent,omitempty"`
}

// Legato Легато
type Legato struct {
	Slur bool `xml:"slur"`
}

// Tempo Дирижирование
type Tempo struct {
	Tempo int `xml:"tempo,attr"`
}

// Instrument Инструмент
type Instrument struct {
	ID   string `xml:"id,attr"`
	Name string `xml:"name"`
}

type ScoreInstrument struct {
	ID   string `xml:"id,attr"`
	Name string `xml:"name"`
}

// Tablature Табулатура
type Tablature struct {
	Measure []TabNote `xml:"measure"`
}

type TabNote struct {
	Fret          int            `xml:"fret"`
	String        int            `xml:"string"`
	Duration      int            `xml:"duration"`
	Articulations []Articulation `xml:"articulations"`
}

// InstrumentPart Инструментальная партия
type InstrumentPart struct {
	ID        string     `xml:"id,attr"`
	PartName  string     `xml:"part-name"`
	Tablature *Tablature `xml:"tablature,omitempty"`
}
