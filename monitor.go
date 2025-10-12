package pie

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// Monitor monitoring system
type Monitor struct {
	engine         *Engine
	commandMonitor *CommandMonitor
	poolMonitor    *PoolMonitor
	serverMonitor  *ServerMonitor
}

// CommandMonitor command monitor
type CommandMonitor struct {
	enabled     bool
	logger      *log.Logger
	commandChan chan *CommandEvent
}

// PoolMonitor connection pool monitor
type PoolMonitor struct {
	enabled  bool
	logger   *log.Logger
	poolChan chan *PoolEvent
}

// ServerMonitor server monitor
type ServerMonitor struct {
	enabled      bool
	logger       *log.Logger
	serverChan   chan *ServerEvent
	topologyChan chan *TopologyEvent
}

// CommandEvent command event
type CommandEvent struct {
	CommandName  string        `json:"commandName"`
	RequestID    int64         `json:"requestId"`
	ConnectionID string        `json:"connectionId"`
	Duration     time.Duration `json:"duration"`
	Command      bson.M        `json:"command"`
	Reply        bson.M        `json:"reply"`
	Failure      error         `json:"failure,omitempty"`
	Timestamp    time.Time     `json:"timestamp"`
}

// PoolEvent connection pool event
type PoolEvent struct {
	Type         string    `json:"type"`
	Address      string    `json:"address"`
	ConnectionID string    `json:"connectionId"`
	PoolSize     int       `json:"poolSize"`
	Available    int       `json:"available"`
	Timestamp    time.Time `json:"timestamp"`
}

// ServerEvent server event
type ServerEvent struct {
	Type      string    `json:"type"`
	Address   string    `json:"address"`
	Topology  string    `json:"topology"`
	Timestamp time.Time `json:"timestamp"`
}

// TopologyEvent topology event
type TopologyEvent struct {
	Type      string    `json:"type"`
	Topology  string    `json:"topology"`
	Servers   []string  `json:"servers"`
	Timestamp time.Time `json:"timestamp"`
}

// NewMonitor create monitoring system
func NewMonitor(engine *Engine) *Monitor {
	return &Monitor{
		engine:         engine,
		commandMonitor: NewCommandMonitor(),
		poolMonitor:    NewPoolMonitor(),
		serverMonitor:  NewServerMonitor(),
	}
}

// NewCommandMonitor create command monitor
func NewCommandMonitor() *CommandMonitor {
	return &CommandMonitor{
		enabled:     false,
		logger:      log.New(log.Writer(), "[CMD] ", log.LstdFlags),
		commandChan: make(chan *CommandEvent, 1000),
	}
}

// NewPoolMonitor create connection pool monitor
func NewPoolMonitor() *PoolMonitor {
	return &PoolMonitor{
		enabled:  false,
		logger:   log.New(log.Writer(), "[POOL] ", log.LstdFlags),
		poolChan: make(chan *PoolEvent, 1000),
	}
}

// NewServerMonitor create server monitor
func NewServerMonitor() *ServerMonitor {
	return &ServerMonitor{
		enabled:      false,
		logger:       log.New(log.Writer(), "[SERVER] ", log.LstdFlags),
		serverChan:   make(chan *ServerEvent, 1000),
		topologyChan: make(chan *TopologyEvent, 1000),
	}
}

// Enable enable monitoring
func (m *Monitor) Enable() *Monitor {
	m.commandMonitor.Enable()
	m.poolMonitor.Enable()
	m.serverMonitor.Enable()
	return m
}

// Disable disable monitoring
func (m *Monitor) Disable() *Monitor {
	m.commandMonitor.Disable()
	m.poolMonitor.Disable()
	m.serverMonitor.Disable()
	return m
}

// EnableCommandMonitor enable command monitoring
func (m *Monitor) EnableCommandMonitor() *Monitor {
	m.commandMonitor.Enable()
	return m
}

// DisableCommandMonitor disable command monitoring
func (m *Monitor) DisableCommandMonitor() *Monitor {
	m.commandMonitor.Disable()
	return m
}

// EnablePoolMonitor enable connection pool monitoring
func (m *Monitor) EnablePoolMonitor() *Monitor {
	m.poolMonitor.Enable()
	return m
}

// DisablePoolMonitor disable connection pool monitoring
func (m *Monitor) DisablePoolMonitor() *Monitor {
	m.poolMonitor.Disable()
	return m
}

// EnableServerMonitor enable server monitoring
func (m *Monitor) EnableServerMonitor() *Monitor {
	m.serverMonitor.Enable()
	return m
}

// DisableServerMonitor disable server monitoring
func (m *Monitor) DisableServerMonitor() *Monitor {
	m.serverMonitor.Disable()
	return m
}

// SetLogger set logger
func (m *Monitor) SetLogger(logger *log.Logger) *Monitor {
	m.commandMonitor.SetLogger(logger)
	m.poolMonitor.SetLogger(logger)
	m.serverMonitor.SetLogger(logger)
	return m
}

// GetCommandEvents get command event channel
func (m *Monitor) GetCommandEvents() <-chan *CommandEvent {
	return m.commandMonitor.GetEvents()
}

// GetPoolEvents get connection pool event channel
func (m *Monitor) GetPoolEvents() <-chan *PoolEvent {
	return m.poolMonitor.GetEvents()
}

// GetServerEvents get server event channel
func (m *Monitor) GetServerEvents() <-chan *ServerEvent {
	return m.serverMonitor.GetEvents()
}

// GetTopologyEvents get topology event channel
func (m *Monitor) GetTopologyEvents() <-chan *TopologyEvent {
	return m.serverMonitor.GetTopologyEvents()
}

// CommandMonitor methods

// Enable enable command monitoring
func (cm *CommandMonitor) Enable() {
	cm.enabled = true
}

// Disable disable command monitoring
func (cm *CommandMonitor) Disable() {
	cm.enabled = false
}

// IsEnabled check if enabled
func (cm *CommandMonitor) IsEnabled() bool {
	return cm.enabled
}

// SetLogger set logger
func (cm *CommandMonitor) SetLogger(logger *log.Logger) {
	cm.logger = logger
}

// GetEvents get event channel
func (cm *CommandMonitor) GetEvents() <-chan *CommandEvent {
	return cm.commandChan
}

// LogCommand log command event
func (cm *CommandMonitor) LogCommand(event *CommandEvent) {
	if !cm.enabled {
		return
	}

	if event.Failure != nil {
		cm.logger.Printf("❌ ERROR [%v] %s: %v - %v", event.Duration, event.CommandName, event.Command, event.Failure)
	} else {
		cm.logger.Printf("✅ [%v] %s: %v", event.Duration, event.CommandName, event.Command)
	}

	// send to event channel
	select {
	case cm.commandChan <- event:
	default:
		// channel full, drop event
	}
}

// PoolMonitor methods

// Enable enable connection pool monitoring
func (pm *PoolMonitor) Enable() {
	pm.enabled = true
}

// Disable disable connection pool monitoring
func (pm *PoolMonitor) Disable() {
	pm.enabled = false
}

// IsEnabled check if enabled
func (pm *PoolMonitor) IsEnabled() bool {
	return pm.enabled
}

// SetLogger set logger
func (pm *PoolMonitor) SetLogger(logger *log.Logger) {
	pm.logger = logger
}

// GetEvents get event channel
func (pm *PoolMonitor) GetEvents() <-chan *PoolEvent {
	return pm.poolChan
}

// LogPoolEvent log connection pool event
func (pm *PoolMonitor) LogPoolEvent(event *PoolEvent) {
	if !pm.enabled {
		return
	}

	pm.logger.Printf("🏊 %s [%s] PoolSize: %d, Available: %d",
		event.Type, event.Address, event.PoolSize, event.Available)

	// send to event channel
	select {
	case pm.poolChan <- event:
	default:
		// channel full, drop event
	}
}

// ServerMonitor methods

// Enable enable server monitoring
func (sm *ServerMonitor) Enable() {
	sm.enabled = true
}

// Disable disable server monitoring
func (sm *ServerMonitor) Disable() {
	sm.enabled = false
}

// IsEnabled check if enabled
func (sm *ServerMonitor) IsEnabled() bool {
	return sm.enabled
}

// SetLogger set logger
func (sm *ServerMonitor) SetLogger(logger *log.Logger) {
	sm.logger = logger
}

// GetEvents get event channel
func (sm *ServerMonitor) GetEvents() <-chan *ServerEvent {
	return sm.serverChan
}

// GetTopologyEvents get topology event channel
func (sm *ServerMonitor) GetTopologyEvents() <-chan *TopologyEvent {
	return sm.topologyChan
}

// LogServerEvent log server event
func (sm *ServerMonitor) LogServerEvent(event *ServerEvent) {
	if !sm.enabled {
		return
	}

	sm.logger.Printf("🖥️ %s [%s] Topology: %s",
		event.Type, event.Address, event.Topology)

	// send to event channel
	select {
	case sm.serverChan <- event:
	default:
		// channel full, drop event
	}
}

// LogTopologyEvent log topology event
func (sm *ServerMonitor) LogTopologyEvent(event *TopologyEvent) {
	if !sm.enabled {
		return
	}

	sm.logger.Printf("🌐 %s Topology: %s, Servers: %v",
		event.Type, event.Topology, event.Servers)

	// send to event channel
	select {
	case sm.topologyChan <- event:
	default:
		// channel full, drop event
	}
}

// MonitorOptions monitoring options
type MonitorOptions struct {
	CommandMonitor bool
	PoolMonitor    bool
	ServerMonitor  bool
	Logger         *log.Logger
}

// NewMonitorOptions create monitoring options
func NewMonitorOptions() *MonitorOptions {
	return &MonitorOptions{
		CommandMonitor: true,
		PoolMonitor:    false,
		ServerMonitor:  false,
		Logger:         log.New(log.Writer(), "[MONITOR] ", log.LstdFlags),
	}
}

// WithCommandMonitor enable command monitoring
func (mo *MonitorOptions) WithCommandMonitor() *MonitorOptions {
	mo.CommandMonitor = true
	return mo
}

// WithPoolMonitor enable connection pool monitoring
func (mo *MonitorOptions) WithPoolMonitor() *MonitorOptions {
	mo.PoolMonitor = true
	return mo
}

// WithServerMonitor enable server monitoring
func (mo *MonitorOptions) WithServerMonitor() *MonitorOptions {
	mo.ServerMonitor = true
	return mo
}

// WithLogger set logger
func (mo *MonitorOptions) WithLogger(logger *log.Logger) *MonitorOptions {
	mo.Logger = logger
	return mo
}

// Build build monitoring options
func (mo *MonitorOptions) Build() *MonitorOptions {
	return mo
}

// EventHandler event handler interface
type EventHandler interface {
	HandleCommandEvent(event *CommandEvent)
	HandlePoolEvent(event *PoolEvent)
	HandleServerEvent(event *ServerEvent)
	HandleTopologyEvent(event *TopologyEvent)
}

// DefaultEventHandler default event handler
type DefaultEventHandler struct {
	logger *log.Logger
}

// NewDefaultEventHandler create default event handler
func NewDefaultEventHandler() *DefaultEventHandler {
	return &DefaultEventHandler{
		logger: log.New(log.Writer(), "[EVENT] ", log.LstdFlags),
	}
}

// HandleCommandEvent handle command event
func (deh *DefaultEventHandler) HandleCommandEvent(event *CommandEvent) {
	deh.logger.Printf("Command: %s [%v] %v", event.CommandName, event.Duration, event.Command)
}

// HandlePoolEvent handle connection pool event
func (deh *DefaultEventHandler) HandlePoolEvent(event *PoolEvent) {
	deh.logger.Printf("Pool: %s [%s] Size: %d, Available: %d",
		event.Type, event.Address, event.PoolSize, event.Available)
}

// HandleServerEvent handle server event
func (deh *DefaultEventHandler) HandleServerEvent(event *ServerEvent) {
	deh.logger.Printf("Server: %s [%s] Topology: %s",
		event.Type, event.Address, event.Topology)
}

// HandleTopologyEvent handle topology event
func (deh *DefaultEventHandler) HandleTopologyEvent(event *TopologyEvent) {
	deh.logger.Printf("Topology: %s %s Servers: %v",
		event.Type, event.Topology, event.Servers)
}

// EventProcessor event processor
type EventProcessor struct {
	monitor  *Monitor
	handler  EventHandler
	stopChan chan struct{}
}

// NewEventProcessor create event processor
func NewEventProcessor(monitor *Monitor, handler EventHandler) *EventProcessor {
	return &EventProcessor{
		monitor:  monitor,
		handler:  handler,
		stopChan: make(chan struct{}),
	}
}

// Start start event processing
func (ep *EventProcessor) Start(ctx context.Context) {
	go ep.processCommandEvents(ctx)
	go ep.processPoolEvents(ctx)
	go ep.processServerEvents(ctx)
	go ep.processTopologyEvents(ctx)
}

// Stop stop event processing
func (ep *EventProcessor) Stop() {
	close(ep.stopChan)
}

// processCommandEvents process command events
func (ep *EventProcessor) processCommandEvents(ctx context.Context) {
	for {
		select {
		case event := <-ep.monitor.GetCommandEvents():
			ep.handler.HandleCommandEvent(event)
		case <-ctx.Done():
			return
		case <-ep.stopChan:
			return
		}
	}
}

// processPoolEvents process connection pool events
func (ep *EventProcessor) processPoolEvents(ctx context.Context) {
	for {
		select {
		case event := <-ep.monitor.GetPoolEvents():
			ep.handler.HandlePoolEvent(event)
		case <-ctx.Done():
			return
		case <-ep.stopChan:
			return
		}
	}
}

// processServerEvents process server events
func (ep *EventProcessor) processServerEvents(ctx context.Context) {
	for {
		select {
		case event := <-ep.monitor.GetServerEvents():
			ep.handler.HandleServerEvent(event)
		case <-ctx.Done():
			return
		case <-ep.stopChan:
			return
		}
	}
}

// processTopologyEvents process topology events
func (ep *EventProcessor) processTopologyEvents(ctx context.Context) {
	for {
		select {
		case event := <-ep.monitor.GetTopologyEvents():
			ep.handler.HandleTopologyEvent(event)
		case <-ctx.Done():
			return
		case <-ep.stopChan:
			return
		}
	}
}

// MetricsCollector metrics collector
type MetricsCollector struct {
	commandCounts map[string]int64
	commandTimes  map[string]time.Duration
	poolStats     map[string]*PoolStats
	serverStats   map[string]*ServerStats
}

// PoolStats connection pool stats
type PoolStats struct {
	TotalConnections  int64
	ActiveConnections int64
	IdleConnections   int64
	CreatedAt         time.Time
	LastUpdated       time.Time
}

// ServerStats server stats
type ServerStats struct {
	IsMaster      bool
	LastPing      time.Time
	RoundTripTime time.Duration
	CreatedAt     time.Time
	LastUpdated   time.Time
}

// NewMetricsCollector create metrics collector
func NewMetricsCollector() *MetricsCollector {
	return &MetricsCollector{
		commandCounts: make(map[string]int64),
		commandTimes:  make(map[string]time.Duration),
		poolStats:     make(map[string]*PoolStats),
		serverStats:   make(map[string]*ServerStats),
	}
}

// CollectCommand collect command metrics
func (mc *MetricsCollector) CollectCommand(event *CommandEvent) {
	mc.commandCounts[event.CommandName]++
	mc.commandTimes[event.CommandName] += event.Duration
}

// CollectPool collect connection pool metrics
func (mc *MetricsCollector) CollectPool(event *PoolEvent) {
	if mc.poolStats[event.Address] == nil {
		mc.poolStats[event.Address] = &PoolStats{
			CreatedAt: time.Now(),
		}
	}

	stats := mc.poolStats[event.Address]
	stats.TotalConnections = int64(event.PoolSize)
	stats.ActiveConnections = int64(event.PoolSize - event.Available)
	stats.IdleConnections = int64(event.Available)
	stats.LastUpdated = time.Now()
}

// CollectServer collect server metrics
func (mc *MetricsCollector) CollectServer(event *ServerEvent) {
	if mc.serverStats[event.Address] == nil {
		mc.serverStats[event.Address] = &ServerStats{
			CreatedAt: time.Now(),
		}
	}

	stats := mc.serverStats[event.Address]
	stats.LastUpdated = time.Now()
}

// GetCommandMetrics get command metrics
func (mc *MetricsCollector) GetCommandMetrics() map[string]any {
	result := make(map[string]any)
	for cmd, count := range mc.commandCounts {
		result[cmd] = map[string]any{
			"count":     count,
			"totalTime": mc.commandTimes[cmd],
			"avgTime":   mc.commandTimes[cmd] / time.Duration(count),
		}
	}
	return result
}

// GetPoolMetrics get connection pool metrics
func (mc *MetricsCollector) GetPoolMetrics() map[string]*PoolStats {
	return mc.poolStats
}

// GetServerMetrics get server metrics
func (mc *MetricsCollector) GetServerMetrics() map[string]*ServerStats {
	return mc.serverStats
}

// Reset reset metrics
func (mc *MetricsCollector) Reset() {
	mc.commandCounts = make(map[string]int64)
	mc.commandTimes = make(map[string]time.Duration)
	mc.poolStats = make(map[string]*PoolStats)
	mc.serverStats = make(map[string]*ServerStats)
}
