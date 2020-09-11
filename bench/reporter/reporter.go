package reporter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/isucon/isucon10-portal/bench-tool.go/benchrun"
	isuxportalResources "github.com/isucon/isucon10-portal/proto.go/isuxportal/resources"
	"github.com/isucon10-qualify/isucon10-qualify/bench/score"
)

type Stdout struct {
	Pass     bool      `json:"pass"`
	Score    int64     `json:"score"`
	Messages []Message `json:"messages"`
	Language string    `json:"language"`
}

type Logger struct {
	buf    *bytes.Buffer
	writer *io.Writer
}

func NewLogger() Logger {
	return Logger{
		buf: &bytes.Buffer{},
	}
}

func (w Logger) Write(p []byte) (n int, err error) {
	return w.buf.Write(p)
}

func (w *Logger) String() string {
	return w.buf.String()
}

var reporter benchrun.Reporter
var result *isuxportalResources.BenchmarkResult
var mu sync.RWMutex
var logger Logger
var writer io.Writer

func init() {
	var err error
	reporter, err = benchrun.NewReporter(false)
	if err != nil {
		log.Fatal(err)
	}

	result = &isuxportalResources.BenchmarkResult{
		Finished: false,
		Passed:   false,
		Score:    0,
		ScoreBreakdown: &isuxportalResources.BenchmarkResult_ScoreBreakdown{
			Raw:       0,
			Deduction: 0,
		},
		Execution: &isuxportalResources.BenchmarkResult_Execution{
			Reason: "",
			Stdout: "",
			Stderr: "",
		},
		SurveyResponse: &isuxportalResources.SurveyResponse{
			Language: "",
		},
	}

	logger = NewLogger()
	writer = io.MultiWriter(logger, os.Stderr)
}

func Report(msgs []string, critical, application, trivial int) error {
	err := update(msgs, critical, application, trivial)
	if err != nil {
		return err
	}

	mu.RLock()
	defer mu.RUnlock()
	err = reporter.Report(result)
	if err != nil {
		return err
	}

	if result.Finished {
		fmt.Println(result.Execution.Stdout)
	}
	return nil
}

func SetFinished(finished bool) {
	mu.Lock()
	defer mu.Unlock()
	result.Finished = finished
}

func SetPassed(passed bool) {
	mu.Lock()
	defer mu.Unlock()
	result.Passed = passed
}

func SetReason(reason string) {
	mu.Lock()
	defer mu.Unlock()
	result.Execution.Reason = reason
}

func update(msgs []string, critical, application, trivial int) error {
	mu.Lock()
	defer mu.Unlock()

	row := score.GetScore()
	deducation := int64(application * 50)
	score := row - deducation
	if score < 0 {
		score = 0
	}

	result.ScoreBreakdown.Raw = row
	result.ScoreBreakdown.Deduction = deducation
	result.Score = score

	if result.Score < 0 {
		result.Execution.Reason = "スコアが0点を下回りました"
	}

	output := Stdout{
		Pass:     result.Passed && result.Score > 0,
		Score:    score,
		Messages: UniqMsgs(msgs),
		Language: result.SurveyResponse.Language,
	}
	bytes, err := json.Marshal(output)
	if err != nil {
		return err
	}
	result.Execution.Stdout = string(bytes)
	result.Execution.Stderr = logger.String()

	return nil
}

func Logf(format string, v ...interface{}) {
	t := time.Now()
	year, month, day := t.Date()
	hour, min, sec := t.Clock()
	prefix := fmt.Sprintf("%04d/%02d/%02d %02d:%02d:%02d: ", year, month, day, hour, min, sec)

	mu.Lock()
	defer mu.Unlock()
	fmt.Fprintf(writer, prefix+format+"\n", v...)
}

func SetLanguage(language string) {
	mu.Lock()
	defer mu.Unlock()
	result.SurveyResponse.Language = language
}

type Message struct {
	Text  string `json:"text"`
	Count int    `json:"count"`
}

func UniqMsgs(allMsgs []string) []Message {
	if len(allMsgs) == 0 {
		return []Message{}
	}

	sort.Strings(allMsgs)
	msgs := make([]Message, 0, len(allMsgs))

	preMsg := allMsgs[0]
	cnt := 0

	// 適当にuniqする
	for _, msg := range allMsgs {
		if preMsg != msg {
			msgs = append(msgs, Message{
				Text:  preMsg,
				Count: cnt,
			})
			preMsg = msg
			cnt = 1
		} else {
			cnt++
		}
	}
	msgs = append(msgs, Message{
		Text:  preMsg,
		Count: cnt,
	})

	return msgs
}
