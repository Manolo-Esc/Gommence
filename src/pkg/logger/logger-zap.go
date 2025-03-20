package logger

import (
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

/* Ejemplo de uso
func main() {
	logger := logger.NewLogger()
	defer logger.Sync()

	logger.Debug("Esto es un mensaje de debug")
	logger.Info("Esto es un mensaje de información", zap.String("context", "main"))
}
*/

type LoggerService interface {
	Info(string)
	Sync() error
}

type loggerServiceImpl struct {
	provider *zap.Logger
}

func (l *loggerServiceImpl) Info(key string) {
	if l.provider != nil {
		l.provider.Info(key)
	}
}

func (l *loggerServiceImpl) Sync() error {
	if l.provider != nil {
		return l.provider.Sync()
	}
	return nil
}

type LoggerConfig struct {
	UseConsole bool
	UseFile    bool
}

var (
	theLogger  *loggerServiceImpl
	createOnce sync.Once
	config     LoggerConfig = LoggerConfig{UseConsole: false, UseFile: true}
)

func GetLogger() LoggerService {
	createOnce.Do(func() {
		theLogger = &loggerServiceImpl{}
		theLogger.provider = newLogger()
	})
	return theLogger
}

func GetNopLogger() LoggerService {
	return &loggerServiceImpl{}
}

func newLogger() *zap.Logger {
	encoderConfig := zapcore.EncoderConfig{ // Encoder para formato legible (JSON en este caso)
		TimeKey:        "time",
		LevelKey:       "level",
		MessageKey:     "msg",
		CallerKey:      "caller",
		EncodeTime:     zapcore.ISO8601TimeEncoder, // Formato de tiempo legible
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
	}

	var core, consoleCore zapcore.Core
	if config.UseConsole {
		consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
		consoleLevel := zapcore.DebugLevel // Logs de debug y superiores a consola
		consoleCore = zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), consoleLevel)
		core = consoleCore
		//consoleCore := zapcore.NewCore(consoleEncoder, zapcore.Lock(zapcore.AddSync(zapcore.AddSync(fileWriter))), consoleLevel)
	}
	var fileCore zapcore.Core
	if config.UseFile {
		// Configuración de lumberjack para rotación de archivos
		fileWriter := zapcore.AddSync(&lumberjack.Logger{
			Filename:   "logs/app.log", // Archivo donde se escribirán los logs
			MaxSize:    10,             // Tamaño máximo del archivo en MB
			MaxBackups: 5,              // Número máximo de archivos de backup
			MaxAge:     30,             // Días máximos para mantener los archivos antiguos
			Compress:   true,           // Habilita la compresión de archivos rotados
		})
		fileEncoder := zapcore.NewJSONEncoder(encoderConfig)
		fileLevel := zapcore.InfoLevel // Logs de info y superiores a archivo
		fileCore = zapcore.NewCore(fileEncoder, fileWriter, fileLevel)
		core = fileCore
	}

	if config.UseConsole && config.UseFile {
		core = zapcore.NewTee(consoleCore, fileCore) // Constructor de logger con múltiples salidas
	}

	return zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
}
