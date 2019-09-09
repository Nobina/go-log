package log

import (
	nlog "log"
	"os"
)

var (
	defaultClient = &Client{
		modules: []Module{},
		loggers: map[Level]*nlog.Logger{
			LevelDebug:    nlog.New(os.Stderr, "DEBUG: ", nlog.Ldate|nlog.Ltime|nlog.Lshortfile),
			LevelInfo:     nlog.New(os.Stderr, "INFO: ", nlog.Ldate|nlog.Ltime|nlog.Lshortfile),
			LevelWarning:  nlog.New(os.Stderr, "WARNING: ", nlog.Ldate|nlog.Ltime|nlog.Lshortfile),
			LevelError:    nlog.New(os.Stderr, "ERROR: ", nlog.Ldate|nlog.Ltime|nlog.Lshortfile),
			LevelCritical: nlog.New(os.Stderr, "CRITICAL: ", nlog.Ldate|nlog.Ltime|nlog.Lshortfile),
		},
	}
)

type Level string

var (
	LevelDebug    = Level("debug")
	LevelInfo     = Level("info")
	LevelWarning  = Level("warning")
	LevelError    = Level("error")
	LevelCritical = Level("critical")
)

type Module interface {
	Logf(level Level, fmt string, v ...interface{})
	Push() error
}

type Client struct {
	modules []Module
	loggers map[Level]*nlog.Logger
}

func (c *Client) logf(level Level, fmt string, v ...interface{}) *Client {
	c.loggers[level].Printf(fmt, v...)

	for _, module := range c.modules {
		module.Logf(level, fmt, v...)
	}

	return c
}

func (c *Client) Use(module Module)                      { c.modules = append(c.modules, module) }
func (c *Client) Debug(v string)                         { c.logf(LevelDebug, v) }
func (c *Client) Debugf(fmt string, v ...interface{})    { c.logf(LevelDebug, fmt, v...) }
func (c *Client) Info(v string)                          { c.logf(LevelInfo, v) }
func (c *Client) Infof(fmt string, v ...interface{})     { c.logf(LevelInfo, fmt, v...) }
func (c *Client) Warning(v string)                       { c.logf(LevelWarning, v) }
func (c *Client) Warningf(fmt string, v ...interface{})  { c.logf(LevelWarning, fmt, v...) }
func (c *Client) Error(v string)                         { c.logf(LevelError, v) }
func (c *Client) Errorf(fmt string, v ...interface{})    { c.logf(LevelError, fmt, v...) }
func (c *Client) Critical(v string)                      { c.logf(LevelCritical, v) }
func (c *Client) Criticalf(fmt string, v ...interface{}) { c.logf(LevelCritical, fmt, v...) }

func (c *Client) Push() error {
	for _, module := range c.modules {
		if err := module.Push(); err != nil {
			return err
		}
	}

	return nil
}

func Debug(v string)                         { defaultClient.Debug(v) }
func Debugf(fmt string, v ...interface{})    { defaultClient.Debugf(fmt, v...) }
func Info(v string)                          { defaultClient.Info(v) }
func Infof(fmt string, v ...interface{})     { defaultClient.Infof(fmt, v...) }
func Warning(v string)                       { defaultClient.Warning(v) }
func Warningf(fmt string, v ...interface{})  { defaultClient.Warningf(fmt, v...) }
func Error(v string)                         { defaultClient.Error(v) }
func Errorf(fmt string, v ...interface{})    { defaultClient.Errorf(fmt, v...) }
func Critical(v string)                      { defaultClient.Critical(v) }
func Criticalf(fmt string, v ...interface{}) { defaultClient.Criticalf(fmt, v...) }

func New() *Client {
	c := &Client{
		modules: []Module{},
		loggers: map[Level]*nlog.Logger{
			LevelDebug:    nlog.New(os.Stdout, "DEBUG: ", nlog.Ldate|nlog.Ltime|nlog.Lshortfile),
			LevelInfo:     nlog.New(os.Stdout, "INFO: ", nlog.Ldate|nlog.Ltime|nlog.Lshortfile),
			LevelWarning:  nlog.New(os.Stderr, "WARNING: ", nlog.Ldate|nlog.Ltime|nlog.Lshortfile),
			LevelError:    nlog.New(os.Stderr, "ERROR: ", nlog.Ldate|nlog.Ltime|nlog.Lshortfile),
			LevelCritical: nlog.New(os.Stderr, "CRITICAL: ", nlog.Ldate|nlog.Ltime|nlog.Lshortfile),
		},
	}

	return c
}
