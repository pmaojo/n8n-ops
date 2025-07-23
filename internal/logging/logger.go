package logging

import (
        "fmt"
        "io"
        "log"
        "os"
        "path/filepath"
        "runtime"
        "strings"
        "time"
)

// LogLevel representa los niveles de logging
type LogLevel int

const (
        DEBUG LogLevel = iota
        INFO
        WARN
        ERROR
        FATAL
)

var levelNames = map[LogLevel]string{
        DEBUG: "DEBUG",
        INFO:  "INFO",
        WARN:  "WARN",
        ERROR: "ERROR",
        FATAL: "FATAL",
}

var levelColors = map[LogLevel]string{
        DEBUG: "\033[36m", // Cyan
        INFO:  "\033[32m", // Green
        WARN:  "\033[33m", // Yellow
        ERROR: "\033[31m", // Red
        FATAL: "\033[35m", // Magenta
}

const colorReset = "\033[0m"

// Logger estructura principal del logger
type Logger struct {
        level      LogLevel
        output     io.Writer
        fileOutput io.Writer
        prefix     string
        colorized  bool
}

// Config configuración del logger
type Config struct {
        Level     LogLevel
        Output    io.Writer
        LogFile   string
        Prefix    string
        Colorized bool
}

var defaultLogger *Logger

func init() {
        // Configurar logger por defecto
        defaultLogger = NewLogger(&Config{
                Level:     INFO,
                Output:    os.Stdout,
                Prefix:    "[n8n-ops]",
                Colorized: true,
        })
}

// NewLogger crea un nuevo logger con la configuración especificada
func NewLogger(config *Config) *Logger {
        logger := &Logger{
                level:     config.Level,
                output:    config.Output,
                prefix:    config.Prefix,
                colorized: config.Colorized,
        }

        // Configurar logging a archivo si se especifica
        if config.LogFile != "" {
                if err := logger.setupFileLogging(config.LogFile); err != nil {
                        log.Printf("Error setting up file logging: %v", err)
                }
        }

        return logger
}

// setupFileLogging configura el logging a archivo
func (l *Logger) setupFileLogging(filename string) error {
        // Crear directorio si no existe
        dir := filepath.Dir(filename)
        if err := os.MkdirAll(dir, 0755); err != nil {
                return fmt.Errorf("failed to create log directory: %w", err)
        }

        // Abrir archivo de log
        file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
        if err != nil {
                return fmt.Errorf("failed to open log file: %w", err)
        }

        l.fileOutput = file
        return nil
}

// SetLevel establece el nivel de logging
func (l *Logger) SetLevel(level LogLevel) {
        l.level = level
}

// GetLevel obtiene el nivel actual de logging
func (l *Logger) GetLevel() LogLevel {
        return l.level
}

// log es el método interno para escribir logs
func (l *Logger) log(level LogLevel, format string, args ...interface{}) {
        if level < l.level {
                return
        }

        // Obtener información del caller
        _, file, line, ok := runtime.Caller(2)
        caller := "unknown"
        if ok {
                caller = fmt.Sprintf("%s:%d", filepath.Base(file), line)
        }

        // Formatear timestamp
        timestamp := time.Now().Format("2006-01-02 15:04:05.000")

        // Crear mensaje
        levelName := levelNames[level]
        message := fmt.Sprintf(format, args...)

        // Formatear línea de log
        logLine := fmt.Sprintf("%s %s %s [%s] %s\n",
                timestamp, l.prefix, levelName, caller, message)

        // Output colorizado para consola
        if l.output != nil {
                colorizedLine := logLine
                if l.colorized {
                        color := levelColors[level]
                        colorizedLine = fmt.Sprintf("%s%s%s %s %s%s%s [%s] %s\n",
                                color, timestamp, colorReset,
                                l.prefix,
                                color, levelName, colorReset,
                                caller, message)
                }
                fmt.Fprint(l.output, colorizedLine)
        }

        // Output sin color para archivo
        if l.fileOutput != nil {
                fmt.Fprint(l.fileOutput, logLine)
        }
}

// Debug registra un mensaje de debug
func (l *Logger) Debug(format string, args ...interface{}) {
        l.log(DEBUG, format, args...)
}

// Info registra un mensaje informativo
func (l *Logger) Info(format string, args ...interface{}) {
        l.log(INFO, format, args...)
}

// Warn registra una advertencia
func (l *Logger) Warn(format string, args ...interface{}) {
        l.log(WARN, format, args...)
}

// Error registra un error
func (l *Logger) Error(format string, args ...interface{}) {
        l.log(ERROR, format, args...)
}

// Fatal registra un error fatal y termina el programa
func (l *Logger) Fatal(format string, args ...interface{}) {
        l.log(FATAL, format, args...)
        os.Exit(1)
}

// WithField añade un campo al contexto del log
func (l *Logger) WithField(key, value string) *ContextLogger {
        return &ContextLogger{
                logger: l,
                fields: map[string]string{key: value},
        }
}

// WithFields añade múltiples campos al contexto del log
func (l *Logger) WithFields(fields map[string]string) *ContextLogger {
        return &ContextLogger{
                logger: l,
                fields: fields,
        }
}

// ContextLogger logger con contexto (campos adicionales)
type ContextLogger struct {
        logger *Logger
        fields map[string]string
}

// formatWithContext formatea el mensaje incluyendo los campos del contexto
func (cl *ContextLogger) formatWithContext(format string, args ...interface{}) (string, []interface{}) {
        message := fmt.Sprintf(format, args...)
        
        if len(cl.fields) > 0 {
                var fieldParts []string
                for key, value := range cl.fields {
                        fieldParts = append(fieldParts, fmt.Sprintf("%s=%s", key, value))
                }
                message = fmt.Sprintf("%s [%s]", message, strings.Join(fieldParts, " "))
        }
        
        return "%s", []interface{}{message}
}

// Debug registra un mensaje de debug con contexto
func (cl *ContextLogger) Debug(format string, args ...interface{}) {
        newFormat, newArgs := cl.formatWithContext(format, args...)
        cl.logger.Debug(newFormat, newArgs...)
}

// Info registra un mensaje informativo con contexto
func (cl *ContextLogger) Info(format string, args ...interface{}) {
        newFormat, newArgs := cl.formatWithContext(format, args...)
        cl.logger.Info(newFormat, newArgs...)
}

// Warn registra una advertencia con contexto
func (cl *ContextLogger) Warn(format string, args ...interface{}) {
        newFormat, newArgs := cl.formatWithContext(format, args...)
        cl.logger.Warn(newFormat, newArgs...)
}

// Error registra un error con contexto
func (cl *ContextLogger) Error(format string, args ...interface{}) {
        newFormat, newArgs := cl.formatWithContext(format, args...)
        cl.logger.Error(newFormat, newArgs...)
}

// Fatal registra un error fatal con contexto
func (cl *ContextLogger) Fatal(format string, args ...interface{}) {
        newFormat, newArgs := cl.formatWithContext(format, args...)
        cl.logger.Fatal(newFormat, newArgs...)
}

// Funciones globales que usan el logger por defecto
func Debug(format string, args ...interface{}) {
        defaultLogger.Debug(format, args...)
}

func Info(format string, args ...interface{}) {
        defaultLogger.Info(format, args...)
}

func Warn(format string, args ...interface{}) {
        defaultLogger.Warn(format, args...)
}

func Error(format string, args ...interface{}) {
        defaultLogger.Error(format, args...)
}

func Fatal(format string, args ...interface{}) {
        defaultLogger.Fatal(format, args...)
}

func WithField(key, value string) *ContextLogger {
        return defaultLogger.WithField(key, value)
}

func WithFields(fields map[string]string) *ContextLogger {
        return defaultLogger.WithFields(fields)
}

func SetLevel(level LogLevel) {
        defaultLogger.SetLevel(level)
}

func SetupFileLogging(filename string) error {
        return defaultLogger.setupFileLogging(filename)
}

// Utilidades de configuración

// ParseLogLevel parsea un string a LogLevel
func ParseLogLevel(level string) LogLevel {
        switch strings.ToUpper(level) {
        case "DEBUG":
                return DEBUG
        case "INFO":
                return INFO
        case "WARN", "WARNING":
                return WARN
        case "ERROR":
                return ERROR
        case "FATAL":
                return FATAL
        default:
                return INFO
        }
}

// GetLogLevelName obtiene el nombre del nivel de log
func GetLogLevelName(level LogLevel) string {
        return levelNames[level]
}