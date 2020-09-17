package internal

import (
	"context"
)

const (
	// plugin types
	PluginTypeKafka = "kafka" // generic kafka plugin
	PluginTypeKDS   = "kds"   // generic kds plugin
	PluginTypeSQS   = "sqs"   // generic sqs plugin
	PluginTypeHTTP  = "http"  // generic http / webhook plugin
	PluginTypeEars  = "ears"  // loopback plugin, may be useful for event splitting
	PluginTypeGears = "gears" // deliver to kafka in masheens envelope
	PluginTypeNull  = "null"  // black whole for events for testing
	PluginTypeDebug = "debug" // black whole for events for testing
)

//TODO: should pattern be an interface?
//TODO: what about pluggable routers? do they need pluggable routing table extensions too?
//TODO: json serialization conventions?

type (

	// A RoutingEntry represents an entry in the EARS routing table
	RoutingTableEntry struct {
		PartnerId       string      `json: "partner_id"` // partner ID for quota and rate limiting
		AppId           string      `json: "app_id"`     // app ID for quota and rate limiting
		SrcType         string      `json: "src_type"`   // source plugin type, e.g. kafka, kds, sqs, webhook
		SrcParams       interface{} `json: "src_params"` // plugin specific configuration parameters
		SrcHash         string      `json: "src_hash"`   // hash over all plugin configurations
		srcRef          *Plugin     // pointer to plugin instance
		DstType         string      `json: "dst_type"`   // destination plugin type
		DstParams       interface{} `json: "dst_params"` // plugin specific configuration parameters
		DstHash         string      `json: "dst_hash"`   // hash over all plugin configurations
		dstRef          *Plugin     // pointer to plugin instance
		RoutingData     interface{} `json: "routing_data"`       // destination specific routing parameters, may contain dynamic elements pulled from incoming event
		MatchPattern    interface{} `json: "match_pattern"`      // json pattern that must be matched for route to be taken
		FilterPattern   interface{} `json: "filter_pattern"`     // json pattern that must not match for route to be taken
		Transformation  interface{} `json: "transformation"`     // simple structural transformation (otpional)
		EventTsPath     string      `json: "event_ts_path"`      // jq path to extract timestamp from event (optional) - maybe this should be a pluggable router feature
		EventTsPeriodMs int         `json: "event_ts_period_ms"` // optional event timeout - maybe this should be a pluggable router feature
		EventSplitPath  string      `json: "event_split_path"`   // optional path to array to be split in event payload - maybe this should be a pluggable router feature
		DeliveryMode    string      `json: "delivery_mode"`      // possible values: fire_and_forget, at_least_once, exactly_once
		Debug           bool        `json: "debug"`              // if true generate debug logs and metrics for events taking this route
		Hash            string      `json: "hash"`               // hash over all route entry configurations
		Ts              int         `json: "ts"`                 // timestamp when route was created or updated
	}

	Pattern interface{}

	// A RoutingTable is a slice of routing entries and reprrsents the EARS routing table
	RoutingTable []*RoutingTableEntry

	// A RoutingTableIndex is a hashmap mapping a routing entry hash to a routing entry pointer
	RoutingTableIndex map[string]*RoutingTableEntry

	// An EarsPlugin represents an input or output plugin instance
	Plugin struct {
		Hash         string               `json: "hash"`      // hash over all plugin configurations
		Type         string               `json: "type"`      // source plugin type, e.g. kafka, kds, sqs, webhook
		Params       interface{}          `json: "params"`    // plugin specific configuration parameters
		IsInput      bool                 `json: "is_input"`  // if true plugin is input plugin
		IsOutput     bool                 `json: "is_output"` // if true plugin is output plugin
		State        string               `json: "state"`     // plugin state
		inputRoutes  []*RoutingTableEntry // list of routes using this plugin instance as source plugin
		outputRoutes []*RoutingTableEntry // list of routes using this plugin instance as output plugin
	}

	// A PluginIndex is a hashmap mapping a plugin instance hash to a plugin instance
	PluginIndex map[string]*Plugin

	// An EarsEvent bundles even payload and metadata
	Event struct {
		Payload  interface{}    `json:"payload"`  // event payload
		Metadata *EventMetadata `json:"metadata"` // event metadata
		Source   *Plugin        `json:"source"`   // pointer to source plugin instance
	}

	// EarsEventMetadata bundles event meta data
	EventMetadata struct {
		Ts int `json: "ts"` // timestamp when event was received
		// tbd
	}

	EventQueue struct { //tbd
	}

	Worker struct { //tbd
	}

	WorkerPool struct { // tbd
	}
)

type (
	// A Hasher can provide a unique deterministic string hash based on its configuration parameters
	Hasher interface {
		Hash() (string, error)
	}

	// A matcher can match a pattern against an object
	Matcher interface {
		Match(message interface{}, pattern interface{}) (bool, error) // if pattern is contained in message thhe function returns true
	}

	// A Doer does things - this is the interface for an EARS worker
	Doer interface {
		Start() error
		Stop() error
	}

	// A RouteModifier allows modifications to a routing table
	RouteModifier interface {
		AddRoute(ctx *context.Context, entry *RoutingTableEntry) error             // idempotent operation to add a routing entry to a local routing table
		RemoveRoute(ctx *context.Context, entry *RoutingTableEntry) error          // idempotent operation to remove a routing entry from a local routing table
		ReplaceAllRoutes(ctx *context.Context, entries []*RoutingTableEntry) error // replace complete local routing table
	}

	// A RouteNavigator
	RouteNavigator interface {
		GetAllRoutes(ctx *context.Context) ([]*RoutingTableEntry, error)                                 // obtain complete local routing table
		GetRoutesBySourcePlugin(ctx *context.Context, plugin *Plugin) ([]*RoutingTableEntry, error)      // get all routes for a specifc source plugin
		GetRoutesByDestinationPlugin(ctx *context.Context, plugin *Plugin) ([]*RoutingTableEntry, error) // get all routes for a specific destination plugin
		GetRoutesForEvent(ctx *context.Context, event *Event) ([]*RoutingTableEntry, error)              // get all routes for a given event (and source plugin)
	}

	// A RoutingTableManager supports CRUD operations on an EARS routing table
	// note: the local in memory cache and the database backed source of truth both implement this interface!
	RoutingTableManager interface {
		RouteNavigator
		RouteModifier
		Hasher
	}

	// An EventRouter represents and EARS worker
	EventRouter interface {
		RouteEvent(ctx *context.Context, event *Event) error
	}

	// An EventQueuer represents an ears event queue
	EventQueuer interface {
		AddEvent(ctx *context.Context, event *Event) error // add event to queue
		NextEvent(ctx *context.Context) (*Event, error)    // blocking call to get next event and remove it from queue
		GetEventCount(ctx *context.Context) (int error)    // get maximum number of elements in queue
		GetMaxEventCount(ctx *context.Context) int         // get capacity of event queue
	}

	// An EventSourceManager manages all event source plugins for a live ears instance
	EventSourceManager interface {
		GetAllEventSources(ctx *context.Context) ([]*Plugin, error)                            // get all event sourced
		GetEventSourcesByType(ctx *context.Context, sourceType string) ([]*Plugin, error)      // get event sources by plugin type
		GetEventSourcesByState(ctx *context.Context, sourceState string) ([]*Plugin, error)    // get event sources by plugin state
		GetEventSourceByRoute(ctx *context.Context, route *RoutingTableEntry) (*Plugin, error) // get event source for route entry
		AddEventSource(ctx *context.Context, source *Plugin) (*Plugin, error)                  // adds event source and starts listening for events if event source doesn't already exist, otherwise increments counter
		RemoveEventSource(ctx *context.Context, source *Plugin) error                          // stops listening for events and removes event source if event route counter is down to zero
	}

	// An EventDestinationManager manages all event destination plugins for a live ears instance
	EventDestinationManager interface {
		GetAllDestinations(ctx *context.Context) ([]*Plugin, error)                                  // get all event sourced
		GetEventDestinationsByType(ctx *context.Context, sourceType string) ([]*Plugin, error)       // get event sources by plugin type
		GetEventDestinationsByState(ctx *context.Context, sourceState string) ([]*Plugin, error)     // get event sources by plugin state
		GetEventDestinationsByRoute(ctx *context.Context, route *RoutingTableEntry) (*Plugin, error) // get event source for route entry
		AddEventDestination(ctx *context.Context, source *Plugin) (*Plugin, error)                   // adds event source and starts listening for events if event source doesn't already exist, otherwise increments counter
		RemoveEventDestination(ctx *context.Context, source *Plugin) error                           // stops listening for events and removes event source if event route counter is down to zero
	}
)
