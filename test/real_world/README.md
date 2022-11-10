### Variables
| Name | Type |
|------|------|
| AppVersion | `string` |
| CleanMarker | `string` |
| CountryCode | `string` |
| CurrencyCode | `string` |
| DestinationCountry | `string` |
| Features | [UserAgentFeatures](#UserAgentFeatures) |
| GateID | `int` |
| InitializedAt | `int` |
| IsAffiliate | `bool` |
| IsDesktop | `bool` |
| IsMobile | `bool` |
| IsOpenJaw | `bool` |
| KnowEnglish | `bool` |
| Locale | `string` |
| Marker | `string` |
| Market | `string` |
| OriginCountry | `string` |
| Os | `string` |
| OsVersion | `string` |
| Passengers | [Passengers](#Passengers) |
| Random | `float` |
| ReferrerHost | `string` |
| SearchDepth | `int` |
| SearchParamsEnv | [SearchParamsEnv](#SearchParamsEnv) |
| Segments | array([RouteSegment](#RouteSegment)) |
| TravelPayoutsAPI | `bool` |
| TripClass | `string` |
| UserAgentDevice | `string` |
| UserAgentType | `string` |
| UserIP | `string` |
| false | `bool` |
| true | `bool` |

### Functions
| Name | Return type |
|------|-------------|
| all(array(`any`), `func`) | `bool` |
| any(array(`any`), `func`) | `bool` |
| count(array(`any`), `func`) | `int` |
| filter(array(`any`), `func`) | array(`any`) |
| len(array(`any`)) | `int` |
| map(array(`any`), `func`) | array(`any`) |
| none(array(`any`), `func`) | `bool` |
| one(array(`any`), `func`) | `bool` |

### Types
#### DirectFlightsDays
| Field | Type |
|---|---|
| Days | `string` |
| Start | `string` |

#### Passengers
| Field | Type |
|---|---|
| Adults | `int` |
| Children | `int` |
| Infants | `int` |

#### RouteSegment
| Field | Type |
|---|---|
| Date | `string` |
| Destination | `string` |
| DestinationCountry | `string` |
| DestinationName | `string` |
| DirectFlightsDays | [DirectFlightsDays](#DirectFlightsDays) |
| Origin | `string` |
| OriginCountry | `string` |
| OriginName | `string` |
| TranslatedDestination | `string` |
| TranslatedOrigin | `string` |
| UserDestination | `string` |
| UserOrigin | `string` |

#### SearchParamsEnv
| Field | Type |
|---|---|
| AppVersion | `string` |
| CleanMarker | `string` |
| CountryCode | `string` |
| CurrencyCode | `string` |
| DestinationCountry | `string` |
| Features | [UserAgentFeatures](#UserAgentFeatures) |
| GateID | `int` |
| InitializedAt | `int` |
| IsAffiliate | `bool` |
| IsDesktop | `bool` |
| IsMobile | `bool` |
| IsOpenJaw | `bool` |
| KnowEnglish | `bool` |
| Locale | `string` |
| Marker | `string` |
| Market | `string` |
| OriginCountry | `string` |
| Os | `string` |
| OsVersion | `string` |
| Passengers | [Passengers](#Passengers) |
| Random | `float` |
| ReferrerHost | `string` |
| SearchDepth | `int` |
| Segments | array([RouteSegment](#RouteSegment)) |
| TravelPayoutsAPI | `bool` |
| TripClass | `string` |
| UserAgentDevice | `string` |
| UserAgentType | `string` |
| UserIP | `string` |

#### UserAgentFeatures
| Field | Type |
|---|---|
| Assisted | `bool` |
| TopPlacement | `bool` |
| TourTickets | `bool` |


