package logger

import "io"

type Format string

const (
	FormatText Format = "text"
	FormatJson Format = "json"
)

type Config struct {
	Path   string
	Format Format
	Output io.Writer
}
