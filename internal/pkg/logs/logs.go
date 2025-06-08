package logs

import (
	"io"
	"log"
	"log/slog"
	"os"
	"path"

	"gopkg.in/natefinch/lumberjack.v2"
)

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
		Level: level,
	})

	slog.SetDefault(slog.New(logger))
}
