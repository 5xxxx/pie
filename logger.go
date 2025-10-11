package pie

import (
	"fmt"
	"io"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// QueryLogger query logger
type QueryLogger struct {
	enabled   bool
	writer    io.Writer
	formatter LogFormatter
}

// LogEntry log entry
type LogEntry struct {
	Timestamp  time.Time
	Collection string
	Operation  string
	Filter     interface{}
	Update     interface{}
	Document   interface{}
	Pipeline   interface{}
	Options    interface{}
	Duration   time.Duration
	Error      error
}

// LogFormatter log formatter interface
type LogFormatter interface {
	Format(entry *LogEntry) string
}

// MongoShellFormatter MongoDB Shell formatter (default)
type MongoShellFormatter struct{}

// JSONFormatter JSON formatter
type JSONFormatter struct{}

// NewQueryLogger create query logger
func NewQueryLogger(writer io.Writer) *QueryLogger {
	return &QueryLogger{
		enabled:   false,
		writer:    writer,
		formatter: &MongoShellFormatter{},
	}
}

// Enable enable query log
func (ql *QueryLogger) Enable() {
	ql.enabled = true
}

// Disable disable query log
func (ql *QueryLogger) Disable() {
	ql.enabled = false
}

// IsEnabled check if enabled
func (ql *QueryLogger) IsEnabled() bool {
	return ql.enabled
}

// SetFormatter set log formatter
func (ql *QueryLogger) SetFormatter(formatter LogFormatter) {
	ql.formatter = formatter
}

// Log record log
func (ql *QueryLogger) Log(entry *LogEntry) {
	if !ql.enabled {
		return
	}

	logStr := ql.formatter.Format(entry)
	fmt.Fprintln(ql.writer, logStr)
}

// MongoShellFormatter.Format implement MongoDB Shell format output
func (f *MongoShellFormatter) Format(entry *LogEntry) string {
	var sb strings.Builder

	// Timestamp and duration
	sb.WriteString(fmt.Sprintf("[%s] [%v] ",
		entry.Timestamp.Format("2006-01-02 15:04:05"),
		entry.Duration.Round(time.Millisecond),
	))

	// Error mark
	if entry.Error != nil {
		sb.WriteString("❌ ERROR ")
	}

	// Database operation
	sb.WriteString(fmt.Sprintf("db.%s.%s", entry.Collection, entry.Operation))

	// Format parameters based on operation type
	switch entry.Operation {
	case "insertOne":
		if entry.Document != nil {
			sb.WriteString(formatDocument(entry.Document))
		}

	case "insertMany":
		if entry.Document != nil {
			sb.WriteString(formatDocuments(entry.Document))
		}

	case "find":
		if entry.Filter != nil {
			sb.WriteString(formatDocument(entry.Filter))
		}
		if entry.Options != nil {
			sb.WriteString(formatFindOptions(entry.Options))
		}

	case "findOne":
		if entry.Filter != nil {
			sb.WriteString(formatDocument(entry.Filter))
		}
		if entry.Options != nil {
			sb.WriteString(formatFindOneOptions(entry.Options))
		}

	case "updateOne", "updateMany":
		if entry.Filter != nil && entry.Update != nil {
			sb.WriteString(fmt.Sprintf("(%s, %s)",
				formatDocument(entry.Filter),
				formatDocument(entry.Update),
			))
		}
		if entry.Options != nil {
			sb.WriteString(formatUpdateOptions(entry.Options))
		}

	case "replaceOne":
		if entry.Filter != nil && entry.Document != nil {
			sb.WriteString(fmt.Sprintf("(%s, %s)",
				formatDocument(entry.Filter),
				formatDocument(entry.Document),
			))
		}

	case "deleteOne", "deleteMany":
		if entry.Filter != nil {
			sb.WriteString(formatDocument(entry.Filter))
		}

	case "aggregate":
		if entry.Pipeline != nil {
			sb.WriteString(formatPipeline(entry.Pipeline))
		}

	case "countDocuments", "estimatedDocumentCount":
		if entry.Filter != nil {
			sb.WriteString(formatDocument(entry.Filter))
		}
	}

	// Error information
	if entry.Error != nil {
		sb.WriteString(fmt.Sprintf(" - error: %v", entry.Error))
	}

	return sb.String()
}

// JSONFormatter.Format implement JSON format output
func (f *JSONFormatter) Format(entry *LogEntry) string {
	return fmt.Sprintf(`{"timestamp":"%s","collection":"%s","operation":"%s","duration":"%v","error":"%v"}`,
		entry.Timestamp.Format(time.RFC3339),
		entry.Collection,
		entry.Operation,
		entry.Duration,
		entry.Error,
	)
}

// formatDocument format document to MongoDB Shell format
func formatDocument(doc interface{}) string {
	if doc == nil {
		return "null"
	}

	switch v := doc.(type) {
	case bson.M:
		return formatBSONMap(v)
	case bson.D:
		return formatBSOND(v)
	case bson.A:
		return formatBSONArray(v)
	case map[string]interface{}:
		return formatMap(v)
	case []interface{}:
		return formatArray(v)
	default:
		// For struct, try to convert to bson.M
		if bsonDoc, err := bson.Marshal(doc); err == nil {
			var bsonMap bson.M
			if err := bson.Unmarshal(bsonDoc, &bsonMap); err == nil {
				return formatBSONMap(bsonMap)
			}
		}
		return fmt.Sprintf("%+v", doc)
	}
}

// formatDocuments format multiple documents
func formatDocuments(docs interface{}) string {
	if docs == nil {
		return "[]"
	}

	switch v := docs.(type) {
	case []interface{}:
		var parts []string
		for _, doc := range v {
			parts = append(parts, formatDocument(doc))
		}
		return fmt.Sprintf("[%s]", strings.Join(parts, ", "))
	default:
		return formatDocument(docs)
	}
}

// formatBSONMap format bson.M
func formatBSONMap(m bson.M) string {
	var parts []string
	for k, v := range m {
		parts = append(parts, fmt.Sprintf("%s: %s", k, formatValue(v)))
	}
	return fmt.Sprintf("{%s}", strings.Join(parts, ", "))
}

// formatBSOND format bson.D
func formatBSOND(d bson.D) string {
	var parts []string
	for _, elem := range d {
		parts = append(parts, fmt.Sprintf("%s: %s", elem.Key, formatValue(elem.Value)))
	}
	return fmt.Sprintf("{%s}", strings.Join(parts, ", "))
}

// formatBSONArray format bson.A
func formatBSONArray(a bson.A) string {
	var parts []string
	for _, v := range a {
		parts = append(parts, formatValue(v))
	}
	return fmt.Sprintf("[%s]", strings.Join(parts, ", "))
}

// formatMap format map[string]interface{}
func formatMap(m map[string]interface{}) string {
	var parts []string
	for k, v := range m {
		parts = append(parts, fmt.Sprintf("%s: %s", k, formatValue(v)))
	}
	return fmt.Sprintf("{%s}", strings.Join(parts, ", "))
}

// formatArray format []interface{}
func formatArray(a []interface{}) string {
	var parts []string
	for _, v := range a {
		parts = append(parts, formatValue(v))
	}
	return fmt.Sprintf("[%s]", strings.Join(parts, ", "))
}

// formatValue format value
func formatValue(v interface{}) string {
	if v == nil {
		return "null"
	}

	switch val := v.(type) {
	case string:
		return fmt.Sprintf(`"%s"`, val)
	case bson.ObjectID:
		return fmt.Sprintf(`ObjectId("%s")`, val.Hex())
	case time.Time:
		return fmt.Sprintf(`ISODate("%s")`, val.Format("2006-01-02T15:04:05.000Z"))
	case bson.M:
		return formatBSONMap(val)
	case bson.D:
		return formatBSOND(val)
	case bson.A:
		return formatBSONArray(val)
	case map[string]interface{}:
		return formatMap(val)
	case []interface{}:
		return formatArray(val)
	default:
		return fmt.Sprintf("%v", val)
	}
}

// formatFindOptions format FindOptions
func formatFindOptions(opts interface{}) string {
	if opts == nil {
		return ""
	}

	var parts []string
	// This needs to access builder's internal fields, but since they are private, we simplify the handling
	// In actual implementation, we might need to use reflection or other ways to get options
	return strings.Join(parts, "")
}

// formatFindOneOptions format FindOneOptions
func formatFindOneOptions(opts interface{}) string {
	if opts == nil {
		return ""
	}

	var parts []string
	// Similar to FindOptions handling
	return strings.Join(parts, "")
}

// formatUpdateOptions format UpdateOptions
func formatUpdateOptions(opts interface{}) string {
	if opts == nil {
		return ""
	}

	var parts []string
	// Similar to FindOptions handling
	return strings.Join(parts, "")
}

// formatPipeline format aggregation pipeline
func formatPipeline(pipeline interface{}) string {
	if pipeline == nil {
		return "[]"
	}

	switch v := pipeline.(type) {
	case []interface{}:
		var parts []string
		for _, stage := range v {
			parts = append(parts, formatDocument(stage))
		}
		return fmt.Sprintf("[%s]", strings.Join(parts, ", "))
	case bson.A:
		return formatBSONArray(v)
	default:
		return formatDocument(pipeline)
	}
}
