package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Server struct {
	ID            int    `json:"id"`
	Hostname      string `json:"hostname"`
	GlobalIP      string `json:"global_ip"`
	PrivateIP     string `json:"private_ip"`
	IsBenchTarget bool   `json:"is_bench_target"`
}

type Team struct {
	ID      int       `json:"id"`
	Owner   int       `json:"owner"`
	Name    string    `json:"name"`
	Servers []*Server `json:"servers"`
}

type Job struct {
	ID     int    `json:"id"`
	Team   *Team  `json:"team"`
	Status string `json:"status"`
	Score  int    `json:"score"`
	Reason string `json:"reason"`
	Stdout string `json:"stdout"`
	Stderr string `json:"stderr"`
}

type Survey struct {
	Language string `json:"language"`
}

type Result struct {
	ID     int    `json:"id"`
	Status string `json:"status"`
	Score  int    `json:"score"`
	Passed bool   `json:"passed"`
	Reason string `json:"reason"`
	Stdout string `json:"stdout"`
	Stderr string `json:"stderr"`
	Survey Survey `json:"survey"`
}

type BenchmarkResult struct {
	Stdout string
	Stderr string
	Status string
}

type BenchmarkResultStdout struct {
	Pass     bool     `json:"pass"`
	Score    int      `json:"score"`
	Messages []string `json:"messages"`
	Language string   `json:"language"`
}

const (
	apiEndpointDev         = "http://localhost:3000"
	defaultInterval        = 3 * time.Second
	maxStderrLength        = 8 * 1024 * 1024
	maxNumMessage          = 20
	maxBenchmarkTime       = 180 * time.Second
	defaultBenchmarkerPath = "../bench/bench"
)

var (
	apiClient *http.Client
)

// errors
var (
	errorJobNotFound          = fmt.Errorf("job not found")
	errorJobDequeueFail       = fmt.Errorf("job dequeue failure")
	errorPortalAPIUnavailable = fmt.Errorf("portal api is unavailable")
)

func init() {
	apiClient = &http.Client{
		Timeout: 5 * time.Second,
	}
}

func dequeue(ep string) (*Job, error) {
	return &Job{
		ID: 0,
		Team: &Team{
			Servers: []*Server{
				{
					GlobalIP:      "localhost:1323",
					IsBenchTarget: true,
				},
			},
		},
	}, nil

	// uri := fmt.Sprintf("%s/internal/job/dequeue/", ep)
	// req, err := http.NewRequest(http.MethodPost, uri, nil)
	// if err != nil {
	// 	return nil, err
	// }

	// res, err := apiClient.Do(req)
	// if err != nil {
	// 	return nil, err
	// }
	// defer res.Body.Close()
	// defer io.Copy(ioutil.Discard, res.Body)

	// switch res.StatusCode {
	// case http.StatusOK:
	// 	job := Job{}
	// 	if err := json.NewDecoder(res.Body).Decode(&job); err != nil {
	// 		return nil, err
	// 	}
	// 	return &job, nil

	// // Job not found
	// case http.StatusNoContent:
	// 	return nil, errorJobNotFound

	// // 5XX
	// case http.StatusInternalServerError, http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout:
	// 	return nil, errorPortalAPIUnavailable

	// default:
	// 	return nil, errorJobDequeueFail
	// }
}

func joinN(messages []string, n int) string {
	if len(messages) > n {
		strings.Join(messages[:n], ",\n")
	}
	return strings.Join(messages, ",\n")
}

func createResult(job *Job, benchmarkResult *BenchmarkResult) *Result {
	status := "done"
	var benchmarkResultStdout BenchmarkResultStdout
	if err := json.NewDecoder(strings.NewReader(benchmarkResult.Stdout)).Decode(&benchmarkResultStdout); err != nil {
		msg := ""
		if benchmarkResult.Status == "timeout" {
			msg = "ベンチマーク実行を指定時間内に完了することができませんでした"
		}
		if benchmarkResult.Status == "fail" {
			msg = "運営に連絡してください"
		}
		benchmarkResultStdout = BenchmarkResultStdout{
			Pass:     false,
			Score:    0,
			Messages: []string{msg},
		}
		status = "aborted"
	}

	return &Result{
		ID:     job.ID,
		Status: status,
		Score:  benchmarkResultStdout.Score,
		Passed: benchmarkResultStdout.Pass,
		Reason: joinN(benchmarkResultStdout.Messages, maxNumMessage),
		Stdout: benchmarkResult.Stdout,
		Stderr: benchmarkResult.Stderr,
		Survey: Survey{
			Language: benchmarkResultStdout.Language,
		},
	}
}

func report(ep string, job *Job, result *Result) error {
	buf := &bytes.Buffer{}
	if err := json.NewEncoder(buf).Encode(result); err != nil {
		return err
	}

	uri := fmt.Sprintf("%s/internal/job/%d/report/", ep, job.ID)
	req, err := http.NewRequest(http.MethodPost, uri, buf)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := apiClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	defer io.Copy(ioutil.Discard, res.Body)

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("%v: %v", res.StatusCode, res.Status)
	}

	return nil
}

func findBenchmarkTargetServer(job *Job) (*Server, error) {
	for _, server := range job.Team.Servers {
		if server.IsBenchTarget {
			return server, nil
		}
	}
	return nil, fmt.Errorf("benchmark target server not found")
}

func runBenchmarker(benchmarkerPath string, job *Job) (*BenchmarkResult, error) {
	target, err := findBenchmarkTargetServer(job)
	if err != nil {
		return &BenchmarkResult{}, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), maxBenchmarkTime)
	defer cancel()
	cmd := exec.CommandContext(
		ctx,
		benchmarkerPath,
		fmt.Sprintf("-target-url=http://%s", target.GlobalIP),
		fmt.Sprintf("-data-dir=../initial-data"),
	)

	var (
		stdout bytes.Buffer
		stderr bytes.Buffer
	)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	status := "success"
	done := make(chan error, 1)
	go func() {
		done <- cmd.Run()
	}()

	select {
	case e := <-done:
		err = e
		if err != nil {
			status = "fail"
		}
	case <-ctx.Done():
		status = "timeout"
		err = fmt.Errorf("benchmarking timeout")
	}

	// triming too long stderr
	stderrStr := stderr.String()
	if len(stderrStr) > maxStderrLength {
		stderrStr = stderrStr[:maxStderrLength]
	}

	return &BenchmarkResult{
		Stdout: stdout.String(),
		Stderr: stderrStr,
		Status: status,
	}, err
}

func printPrettyResult(result *Result) {
	tmpResult := &Result{
		ID:     result.ID,
		Status: result.Status,
		Score:  result.Score,
		Passed: result.Passed,
		Reason: result.Reason,
		Stdout: "see above",
		Stderr: "see above",
		Survey: result.Survey,
	}
	log.Println("============Benchmark stderr start====================")
	log.Println(result.Stderr)
	log.Println("============Benchmark stderr end======================")
	log.Println("============Benchmark stdout start====================")
	log.Println(result.Stdout)
	log.Println("============Benchmark stdout end======================")
	log.Println("============Result start==============================")
	json.NewEncoder(os.Stderr).Encode(tmpResult)
	log.Println("============Result end================================")
}

func main() {
	var (
		apiEndpoint     string
		interval        time.Duration
		benchmarkerPath string
	)

	flag.StringVar(&apiEndpoint, "ep", apiEndpointDev, "API Endpoint")
	flag.DurationVar(&interval, "interval", defaultInterval, "Dequeuing interval second")
	flag.StringVar(&benchmarkerPath, "benchmarker", defaultBenchmarkerPath, "Benchmarker path")
	flag.Parse()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for range ticker.C {
		job, err := dequeue(apiEndpoint)
		if err != nil {
			log.Println(err)
			continue
		}

		log.Println("Dequeued benchmark job")
		log.Println("============Benchmark job start====================")
		json.NewEncoder(os.Stderr).Encode(job)
		log.Println("============Benchmark job end======================")

		log.Printf("Run benchmark")
		benchmarkResult, err := runBenchmarker(benchmarkerPath, job)
		if err != nil {
			log.Println("Run benchmark fail: ", err)
		}

		log.Printf("Report benchmark result start")
		result := createResult(job, benchmarkResult)
		printPrettyResult(result)
		if err := report(apiEndpoint, job, result); err != nil {
			log.Println("Report benchmark result fail: ", err)
		} else {
			log.Printf("Report benchmark result done")
		}
	}
}
