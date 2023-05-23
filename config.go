package opentelemetry_stacktrace_processor // import "github.com/joostlek/opentelemetry-stacktrace-processor"

type Config struct {
	SourceMapDirs []string `mapstructure:"source_map_dirs"`
}
