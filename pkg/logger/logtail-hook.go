package logger

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
)

// LogtailHook demuxes logs to io.Writers based on
// severity. By default it uses the following outputs:
// error and higher -> os.Stderr
// warning and lower -> os.Stdout
type LogtailHook struct {
	SourceToken string
	MinLevel    logrus.Level
	Formatter   logrus.Formatter
	MaxRetry    int
}

// New returns a new LogtailHook by silencing the parent
// logger and configuring separate loggers for stderr and
// stdout with the parents loggers properties.
func NewLogtailHook(parent *logrus.Logger, sourceToken string, minLevel logrus.Level) *LogtailHook {

	return &LogtailHook{
		SourceToken: sourceToken,
		MinLevel:    minLevel,
		Formatter: &logrus.JSONFormatter{
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyMsg: "message",
			},
		},
		MaxRetry: 3,
	}
}

// Fire is triggered on new log entries
func (hook *LogtailHook) Fire(entry *logrus.Entry) error {

	formatted, err := hook.Formatter.Format(entry)
	if err != nil {
		return err
	}

	// Batching is disabled, just send the single log now
	hook.send([][]byte{formatted})

	return nil
}

// Levels returns all levels this hook should be registered to
func (hook *LogtailHook) Levels() []logrus.Level {
	return logrus.AllLevels[:hook.MinLevel+1]
}

func (hook *LogtailHook) send(batch [][]byte) {
	fmt.Println("Sending batch of", len(batch), "logs")

	buf := make([]byte, 0)
	for i, line := range batch {
		// First character is the opening bracket of the array
		if i == 0 {
			buf = append(buf, '[')
		}
		buf = append(buf, line...)
		// Last character is the closing bracket of the array and doesn't get a ','
		if i == len(batch)-1 {
			buf = append(buf, ']')
			continue
		}
		buf = append(buf, ',')
	}

	url := "https://in.logtail.com"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(buf))

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	header := http.Header{}
	header.Add("Authorization", fmt.Sprint("Bearer ", hook.SourceToken))
	header.Add("Content-Type", "application/json")

	req.Header = header
	fmt.Printf("req: %v", req)

	i := 0

	fmt.Println("Sending log to Logtail")
	for {
		resp, err := http.DefaultClient.Do(req)
		fmt.Println("error 1")
		if err != nil || (resp != nil && resp.StatusCode >= 400) {
			fmt.Println(resp)
			// print body
			buf := new(bytes.Buffer)
			buf.ReadFrom(resp.Body)
			fmt.Println(buf.String())
			i++
			if hook.MaxRetry < 0 || i >= hook.MaxRetry {
				fmt.Println(err.Error())
				// err := fmt.Errorf("failed to send after %d retries", hook.MaxRetry)
				return
			}
			continue
		}

		fmt.Println("error 3")
		return
	}
}

// Close closes the hook
func (hook *LogtailHook) Close() {
	// Nothing to do
}
