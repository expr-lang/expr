package real_world

//go:generate sh -c "go run ./real_docs/generate.go > README.md"

type DirectFlightsDays struct {
	Start string
	Days  string
}

type RouteSegment struct {
	Origin                string
	OriginName            string
	Destination           string
	DestinationName       string
	Date                  string
	OriginCountry         string
	DestinationCountry    string
	TranslatedOrigin      string
	TranslatedDestination string
	UserOrigin            string
	UserDestination       string
	DirectFlightsDays     *DirectFlightsDays
}

type Passengers struct {
	Adults   uint32
	Children uint32
	Infants  uint32
}

type UserAgentFeatures struct {
	Assisted     bool
	TopPlacement bool
	TourTickets  bool
}

type SearchParamsEnv struct {
	Segments           []*RouteSegment
	OriginCountry      string
	DestinationCountry string
	SearchDepth        int
	Passengers         *Passengers
	TripClass          string
	UserIP             string
	KnowEnglish        bool
	Market             string
	Marker             string
	CleanMarker        string
	Locale             string
	ReferrerHost       string
	CountryCode        string
	CurrencyCode       string
	IsOpenJaw          bool
	Os                 string
	OsVersion          string
	AppVersion         string
	IsAffiliate        bool
	InitializedAt      int64
	Random             float32
	TravelPayoutsAPI   bool
	Features           *UserAgentFeatures
	GateID             int32
	UserAgentDevice    string
	UserAgentType      string
	IsDesktop          bool
	IsMobile           bool
}

type Env struct {
	SearchParamsEnv
}

func NewEnv() Env {
	return Env{
		SearchParamsEnv: SearchParamsEnv{
			Segments: []*RouteSegment{
				{
					Origin:      "VOG",
					Destination: "SHJ",
				},
				{
					Origin:      "SHJ",
					Destination: "VOG",
				},
			},
			OriginCountry:      "RU",
			DestinationCountry: "RU",
			SearchDepth:        44,
			Passengers:         &Passengers{1, 0, 0},
			TripClass:          "Y",
			UserIP:             "::1",
			KnowEnglish:        true,
			Market:             "ru",
			Marker:             "123456.direct",
			CleanMarker:        "123456",
			Locale:             "ru",
			ReferrerHost:       "www.aviasales.ru",
			CountryCode:        "",
			CurrencyCode:       "usd",
			IsOpenJaw:          false,
			Os:                 "",
			OsVersion:          "",
			AppVersion:         "",
			IsAffiliate:        true,
			InitializedAt:      1570788719,
			Random:             0.13497187,
			TravelPayoutsAPI:   false,
			Features:           &UserAgentFeatures{},
			GateID:             421,
			UserAgentDevice:    "DESKTOP",
			UserAgentType:      "WEB",
			IsDesktop:          true,
			IsMobile:           false,
		},
	}
}
