package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"

	// This import is just an example, replace sw.WatchEvent with your actual struct
	sw "sigs.k8s.io/kube-scheduler-simulator/simulator/resourcewatcher/streamwriter"
)

type Client struct{}

func (c *Client) Connect(address string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", address, nil)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("performing request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("non-200 response code: %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)

	for scanner.Scan() {
		data := scanner.Bytes()

		var watchEvent sw.WatchEvent
		err := json.Unmarshal(data, &watchEvent)
		if err != nil {
			logrus.Errorf("Error unmarshalling JSON: %s", err)
		}

		if watchEvent.EventType == "MODIFIED" {
			//logrus.Infof("Received WatchEvent: %+v", watchEvent)
			parseEvent(watchEvent)
		}
	}

	// Check if the scanner ended with an error
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scanner error: %w", err)
	}

	return nil
}

func parseEvent(in sw.WatchEvent) {
	if m, ok := in.Obj.(map[string]interface{}); !ok {
		logrus.Errorf("Error parsing object: %v", in.Obj)
	} else {
		if annotations, ok := m["metadata"].(map[string]interface{})["annotations"].(map[string]interface{}); ok {
			logrus.Infof("Received event for pod %+v", annotations)
			// TODO extract scores
		}
	}

}
