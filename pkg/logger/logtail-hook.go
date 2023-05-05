package logger

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
)

type LogtailHook struct {
	SourceToken string
	MinLevel    logrus.Level
	Formatter   logrus.Formatter
	MaxRetry    int
}

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

func (hook *LogtailHook) Fire(entry *logrus.Entry) error {

	formatted, err := hook.Formatter.Format(entry)
	if err != nil {
		return err
	}

	// Batching is disabled, just send the single log now
	hook.send([][]byte{formatted})

	return nil
}

func (hook *LogtailHook) Levels() []logrus.Level {
	return logrus.AllLevels[:hook.MinLevel+1]
}

func (hook *LogtailHook) send(batch [][]byte) {
	fmt.Println("Sending batch of", len(batch), "logs")

	buf := make([]byte, 0)
	for i, line := range batch {
		if i == 0 {
			buf = append(buf, '[')
		}
		buf = append(buf, line...)

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

	i := 0

	for {
		resp, err := http.DefaultClient.Do(req)
		if err != nil || (resp != nil && resp.StatusCode >= 400) {
			fmt.Println(resp)
			// print body
			buf := new(bytes.Buffer)
			buf.ReadFrom(resp.Body)
			i++
			if hook.MaxRetry < 0 || i >= hook.MaxRetry {
				fmt.Println(err.Error())
				// err := fmt.Errorf("failed to send after %d retries", hook.MaxRetry)
				return
			}
			continue
		}
		return
	}
}

// Close closes the hook
func (hook *LogtailHook) Close() {
	// Nothing to do
}
