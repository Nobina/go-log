package log

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"time"
)

type SysLogInfluxConfig struct {
	HttpClient *http.Client
	Database   string
	Appname    string
	Host       string
	Token      string
	ProcID     string
	BaseURL    string
}

type SysLogInflux struct {
	httpClient *http.Client
	baseURL    string
	database   string
	appname    string
	host       string
	token      string
	procID     string
	stack      []string
}

func (m *SysLogInflux) Logf(level Level, format string, v ...interface{}) {
	row := "syslog"

	if m.database != "" {
		row += ",db=" + m.database
	}
	if m.appname != "" {
		row += ",appname=" + m.appname
	}
	row += ",facility=console"
	if m.host != "" {
		row += ",host=" + m.host
		row += ",hostname=" + m.host
	}
	if level == LevelDebug {
		row += ",severity=debug"
	} else if level == LevelInfo {
		row += ",severity=info"
	} else if level == LevelWarning {
		row += ",severity=warning"
	} else if level == LevelError {
		row += ",severity=error"
	} else if level == LevelCritical {
		row += ",severity=crit"
	}

	row += fmt.Sprintf(" message=\"%v\",timestamp=%v000000000i", fmt.Sprintf(format, v...), time.Now().Unix())

	if m.procID != "" {
		row += fmt.Sprintf(",procid=\"%v\"", m.procID)
	}

	m.stack = append(m.stack, row)
}

func (m *SysLogInflux) Push() error {
	if len(m.stack) == 0 {
		return nil
	}

	buf := new(bytes.Buffer)

	for i, row := range m.stack {
		if i != 0 {
			row = "\n" + row
		}

		buf.Write([]byte(row))
	}

	m.stack = []string{}

	req, err := http.NewRequest("POST", m.baseURL, buf)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "text/plain")
	req.Header.Set("Authorization", "Token "+m.token)

	resp, err := m.httpClient.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return err
	} else if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("bad status %v", resp.StatusCode)
	}

	return nil
}

func NewSysLogInflux(config *SysLogInfluxConfig) *SysLogInflux {
	logger := &SysLogInflux{
		httpClient: http.DefaultClient,
		stack:      []string{},
		database:   "system",
	}

	if config != nil {
		if config.HttpClient != nil {
			logger.httpClient = config.HttpClient
		}
		if config.BaseURL != "" {
			logger.baseURL = config.BaseURL
			if logger.baseURL == "" {
				logger.baseURL = os.Getenv("SYSLOG_INFLUX_BASE_URL")
			}
		}
		if config.Database != "" {
			logger.database = config.Database
			if logger.database == "" {
				logger.database = os.Getenv("SYSLOG_INFLUX_DATABASE")
			}
		}
		if config.Appname != "" {
			logger.appname = config.Appname
			if logger.appname == "" {
				logger.appname = os.Getenv("SYSLOG_INFLUX_APP_NAME")
			}
		}
		if config.Host != "" {
			logger.host = config.Host
			if logger.host == "" {
				logger.host = os.Getenv("SYSLOG_INFLUX_HOST")
			}
		}
		if config.ProcID != "" {
			logger.procID = config.ProcID
			if logger.procID == "" {
				logger.procID = os.Getenv("SYSLOG_INFLUX_PROCID")
			}
		}
		if config.Token != "" {
			logger.token = config.Token
			if logger.token == "" {
				logger.token = os.Getenv("SYSLOG_INFLUX_TOKEN")
			}
		}
	}

	if logger.host == "" {
		name, _ := os.Hostname()
		logger.host = name
	}
	if logger.appname == "" {
		logger.appname = logger.host
	}

	return logger
}
