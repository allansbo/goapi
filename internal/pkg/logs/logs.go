package logs

import (
	"io"
	"log"
	"log/slog"
	"os"
	"path"

	"gopkg.in/natefinch/lumberjack.v2"
)

// ConfigLog initializes the logging configuration for the application.
// The log file is rotated using lumberjack to manage size and backups.
// The log level is set to Info, and the logs are formatted in JSON.
// The log file is stored in the "logs" directory under the current working directory.
func ConfigLog(outputs ...io.Writer) {
	output, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	level := slog.LevelInfo
	bkpSize := 5
	bkpPath := path.Join(output, "logs", "app.log")

	rotate := lumberjack.Logger{
		Filename:   bkpPath,
		MaxSize:    bkpSize, // megabytes
		MaxBackups: 2,
		Compress:   true,
		MaxAge:     28, // days
	}

	outputs = append(outputs, &rotate)
	multi := io.MultiWriter(outputs...)

	logger := slog.NewJSONHandler(multi, &slog.HandlerOptions{
		Level:     level,
		AddSource: true,
	})

	slog.SetDefault(slog.New(logger))
}
