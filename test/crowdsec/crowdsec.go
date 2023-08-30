package crowdsec

import (
	"time"
)

type Event struct {
	Type            int               `yaml:"Type,omitempty" json:"Type,omitempty"`
	ExpectMode      int               `yaml:"ExpectMode,omitempty" json:"ExpectMode,omitempty"`
	Whitelisted     bool              `yaml:"Whitelisted,omitempty" json:"Whitelisted,omitempty"`
	WhitelistReason string            `yaml:"WhitelistReason,omitempty" json:"whitelist_reason,omitempty"`
	Stage           string            `yaml:"Stage,omitempty" json:"Stage,omitempty"`
	Line            Line              `yaml:"Line,omitempty" json:"Line,omitempty"`
	Parsed          map[string]string `yaml:"Parsed,omitempty" json:"Parsed,omitempty"`
	Enriched        map[string]string `yaml:"Enriched,omitempty" json:"Enriched,omitempty"`
	Unmarshaled     map[string]any    `yaml:"Unmarshaled,omitempty" json:"Unmarshaled,omitempty"`
	Overflow        RuntimeAlert      `yaml:"Overflow,omitempty" json:"Alert,omitempty"`
	Time            time.Time         `yaml:"Time,omitempty" json:"Time,omitempty"`
	StrTime         string            `yaml:"StrTime,omitempty" json:"StrTime,omitempty"`
	StrTimeFormat   string            `yaml:"StrTimeFormat,omitempty" json:"StrTimeFormat,omitempty"`
	MarshaledTime   string            `yaml:"MarshaledTime,omitempty" json:"MarshaledTime,omitempty"`
	Process         bool              `yaml:"Process,omitempty" json:"Process,omitempty"`
	Meta            map[string]string `yaml:"Meta,omitempty" json:"Meta,omitempty"`
}

func (e *Event) GetType() string {
	return ""
}

func (e *Event) GetMeta(key string) string {
	return ""
}

type Alert struct {
	Capacity        *int32      `json:"capacity"`
	CreatedAt       string      `json:"created_at,omitempty"`
	Decisions       []*Decision `json:"decisions"`
	Events          []*Event    `json:"events"`
	EventsCount     *int32      `json:"events_count"`
	ID              int64       `json:"id,omitempty"`
	Labels          []string    `json:"labels"`
	Leakspeed       *string     `json:"leakspeed"`
	MachineID       string      `json:"machine_id,omitempty"`
	Message         *string     `json:"message"`
	Meta            Meta        `json:"meta,omitempty"`
	Remediation     bool        `json:"remediation,omitempty"`
	Scenario        *string     `json:"scenario"`
	ScenarioHash    *string     `json:"scenario_hash"`
	ScenarioVersion *string     `json:"scenario_version"`
	Simulated       *bool       `json:"simulated"`
	Source          *Source     `json:"source"`
	StartAt         *string     `json:"start_at"`
	StopAt          *string     `json:"stop_at"`
	UUID            string      `json:"uuid,omitempty"`
	Edges           AlertEdges  `json:"edges"`
}

func (a *Alert) HasRemediation() bool {
	return true
}

func (a *Alert) GetScope() string {
	return ""
}

func (a *Alert) GetValue() string {
	return ""
}

func (a *Alert) GetScenario() string {
	return ""
}

func (a *Alert) GetEventsCount() int32 {
	return 0
}

func (a *Alert) GetMeta(_ string) string {
	return ""
}

func (s Source) GetValue() string {
	return *s.Value
}

func (s Source) GetScope() string {
	return *s.Scope
}

func (s Source) GetAsNumberName() string {
	return ""
}

type AlertEdges struct {
	Owner     *Machine    `json:"owner,omitempty"`
	Decisions []*Decision `json:"decisions,omitempty"`
	Events    []*Event    `json:"events,omitempty"`
	Metas     []*Meta     `json:"metas,omitempty"`
}

func (e AlertEdges) OwnerOrErr() (*Machine, error) {
	return nil, nil
}

func (e AlertEdges) DecisionsOrErr() ([]*Decision, error) {
	return nil, nil
}

func (e AlertEdges) EventsOrErr() ([]*Event, error) {
	return nil, nil
}

func (e AlertEdges) MetasOrErr() ([]*Meta, error) {
	return nil, nil
}

type Machine struct {
	ID            int          `json:"id,omitempty"`
	CreatedAt     *time.Time   `json:"created_at,omitempty"`
	UpdatedAt     *time.Time   `json:"updated_at,omitempty"`
	LastPush      *time.Time   `json:"last_push,omitempty"`
	LastHeartbeat *time.Time   `json:"last_heartbeat,omitempty"`
	MachineId     string       `json:"machineId,omitempty"`
	Password      string       `json:"-"`
	IpAddress     string       `json:"ipAddress,omitempty"`
	Scenarios     string       `json:"scenarios,omitempty"`
	Version       string       `json:"version,omitempty"`
	IsValidated   bool         `json:"isValidated,omitempty"`
	Status        string       `json:"status,omitempty"`
	AuthType      string       `json:"auth_type"`
	Edges         MachineEdges `json:"edges"`
}

type MachineEdges struct {
	Alerts []*Alert `json:"alerts,omitempty"`
}

type Decision struct {
	Duration  *string `json:"duration"`
	ID        int64   `json:"id,omitempty"`
	Origin    *string `json:"origin"`
	Scenario  *string `json:"scenario"`
	Scope     *string `json:"scope"`
	Simulated *bool   `json:"simulated,omitempty"`
	Type      *string `json:"type"`
	Until     string  `json:"until,omitempty"`
	UUID      string  `json:"uuid,omitempty"`
	Value     *string `json:"value"`
}

type Line struct {
	Raw     string `yaml:"Raw,omitempty"`
	Src     string `yaml:"Src,omitempty"`
	Time    time.Time
	Labels  map[string]string `yaml:"Labels,omitempty"`
	Process bool
	Module  string `yaml:"Module,omitempty"`
}

type ScopeType struct {
	Scope  string `yaml:"type"`
	Filter string `yaml:"expression"`
}

type RuntimeAlert struct {
	Mapkey      string            `yaml:"MapKey,omitempty" json:"MapKey,omitempty"`
	BucketId    string            `yaml:"BucketId,omitempty" json:"BucketId,omitempty"`
	Whitelisted bool              `yaml:"Whitelisted,omitempty" json:"Whitelisted,omitempty"`
	Reprocess   bool              `yaml:"Reprocess,omitempty" json:"Reprocess,omitempty"`
	Sources     map[string]Source `yaml:"Sources,omitempty" json:"Sources,omitempty"`
	Alert       *Alert            `yaml:"Alert,omitempty" json:"Alert,omitempty"`
	APIAlerts   []Alert           `yaml:"APIAlerts,omitempty" json:"APIAlerts,omitempty"`
}

func (r RuntimeAlert) GetSources() []string {
	return nil
}

type Source struct {
	AsName    string  `json:"as_name,omitempty"`
	AsNumber  string  `json:"as_number,omitempty"`
	Cn        string  `json:"cn,omitempty"`
	IP        string  `json:"ip,omitempty"`
	Latitude  float32 `json:"latitude,omitempty"`
	Longitude float32 `json:"longitude,omitempty"`
	Range     string  `json:"range,omitempty"`
	Scope     *string `json:"scope"`
	Value     *string `json:"value"`
}

type Meta []*MetaItems0

type MetaItems0 struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}
