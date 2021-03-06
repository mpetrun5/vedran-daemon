package lb

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"

	"github.com/NodeFactoryIo/vedran-daemon/internal/node"
	mocks "github.com/NodeFactoryIo/vedran-daemon/mocks/node"
	"github.com/stretchr/testify/assert"
)

func Test_metricsService_Send(t *testing.T) {
	setup()
	defer teardown()

	nodeClient := &mocks.Client{}

	peerCount := float64(19)
	bestBlockHeight := float64(432933)
	finalizedBlockHeight := float64(432640)
	readyTransactionCount := float64(0)
	expectedMetrics := node.Metrics{
		PeerCount:             &peerCount,
		BestBlockHeight:       &bestBlockHeight,
		FinalizedBlockHeight:  &finalizedBlockHeight,
		ReadyTransactionCount: &readyTransactionCount}

	tests := []struct {
		name                    string
		wantErr                 bool
		lbHandleFunc            handleFnMock
		want                    int
		fetchMetricsMockMetrics *node.Metrics
		fetchMetricsMockError   error
	}{
		{
			name:    "Returns error if fetching metrics fails",
			wantErr: true,
			want:    0,
			lbHandleFunc: func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "Not Found", 404)
			},
			fetchMetricsMockError:   fmt.Errorf("Failed fetching metrics"),
			fetchMetricsMockMetrics: nil},
		{
			name:    "Returns error if sending metrics fails",
			wantErr: true,
			want:    0,
			lbHandleFunc: func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "Not Found", 404)
			},
			fetchMetricsMockError:   nil,
			fetchMetricsMockMetrics: &expectedMetrics},
		{
			name:    "Returns resp if sending metrics succedes",
			wantErr: false,
			want:    200,
			lbHandleFunc: func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodPut, r.Method)

				var expectedMetrics node.Metrics
				defer r.Body.Close()
				body, _ := ioutil.ReadAll((r.Body))
				_ = json.Unmarshal(body, &expectedMetrics)

				assert.Equal(t, expectedMetrics, node.Metrics{
					PeerCount:             &peerCount,
					BestBlockHeight:       &bestBlockHeight,
					FinalizedBlockHeight:  &finalizedBlockHeight,
					ReadyTransactionCount: &readyTransactionCount})

				_, _ = io.WriteString(w, `{"status": "ok"}`)
			},
			fetchMetricsMockError:   nil,
			fetchMetricsMockMetrics: &expectedMetrics},
	}
	for _, tt := range tests {
		setup()

		t.Run(tt.name, func(t *testing.T) {
			mux.HandleFunc("/api/v1/nodes/metrics", tt.lbHandleFunc)
			mockURL, _ := url.Parse(server.URL)
			client := NewClient(mockURL)
			nodeClient.On("GetMetrics").Once().Return(
				tt.fetchMetricsMockMetrics,
				tt.fetchMetricsMockError)

			got, err := client.Metrics.Send(nodeClient)

			if (err != nil) != tt.wantErr {
				t.Errorf("metricsService.Send() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != nil && got.StatusCode != 200 {
				t.Errorf("metricsService.Send() statusCode = %d, want %d", got.StatusCode, tt.want)
				return
			}
		})

		teardown()
	}
}
