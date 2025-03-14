package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)



func getWasabiStatsURL(pageNum int) string {
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	today := time.Now().Format("2006-01-02")
	return fmt.Sprintf("https://stats.wasabisys.com/v1/standalone/utilizations/bucket?from=%s&to=%s&pageNum=%d", yesterday, today, pageNum)
}

type BucketStats struct {
	BucketUtilizationNum      int     `json:"BucketUtilizationNum"`
	AcctNum                  int     `json:"AcctNum"`
	AcctPlanNum              int     `json:"AcctPlanNum"`
	BucketNum                int     `json:"BucketNum"`
	StartTime                string  `json:"StartTime"`
	EndTime                  string  `json:"EndTime"`
	CreateTime               string  `json:"CreateTime"`
	NumBillableObjects       float64 `json:"NumBillableObjects"`
	NumBillableDeletedObjects float64 `json:"NumBillableDeletedObjects"`
	RawStorageSizeBytes      float64 `json:"RawStorageSizeBytes"`
	DeletedStorageSizeBytes  float64 `json:"DeletedStorageSizeBytes"`
	PaddedStorageSizeBytes   float64 `json:"PaddedStorageSizeBytes"`
	MetadataStorageSizeBytes float64 `json:"MetadataStorageSizeBytes"`
	OrphanedStorageSizeBytes float64 `json:"OrphanedStorageSizeBytes"`
	StorageWroteBytes        float64 `json:"StorageWroteBytes"`
	StorageReadBytes         float64 `json:"StorageReadBytes"`
	NumAPICalls             float64 `json:"NumAPICalls"`
	UploadBytes             float64 `json:"UploadBytes"`
	DownloadBytes           float64 `json:"DownloadBytes"`
	NumGETCalls             float64 `json:"NumGETCalls"`
	NumPUTCalls             float64 `json:"NumPUTCalls"`
	NumDELETECalls          float64 `json:"NumDELETECalls"`
	NumLISTCalls            float64 `json:"NumLISTCalls"`
	NumHEADCalls            float64 `json:"NumHEADCalls"`
	Bucket                  string  `json:"Bucket"`
	Region                  string  `json:"Region"`
}

type PageInfo struct {
	RecordCount int `json:"RecordCount"`
	PageCount   int `json:"PageCount"`
	PageSize    int `json:"PageSize"`
	PageNum     int `json:"PageNum"`
}

type WasabiResponse struct {
	PageInfo PageInfo       `json:"PageInfo"`
	Records  []BucketStats `json:"Records"`
}
var metrics = map[string]*prometheus.GaugeVec{}


func initMetrics() {
	fields := map[string]string{
		"wasabi_raw_storage_bytes": "Total raw storage bytes used per bucket",
    	"wasabi_deleted_storage_bytes": "Total deleted storage bytes per bucket",
    	"wasabi_padded_storage_bytes": "Total padded storage bytes per bucket",
    	"wasabi_metadata_storage_bytes": "Total metadata storage bytes per bucket",
    	//"wasabi_orphaned_storage_bytes": "Total orphaned storage bytes per bucket",
    	"wasabi_storage_wrote_bytes": "Total storage wrote bytes per bucket",
    	"wasabi_storage_read_bytes": "Total storage read bytes per bucket",
    	"wasabi_api_calls_total": "Total API calls per bucket",
    	"wasabi_upload_bytes": "Total uploaded bytes per bucket",
    	"wasabi_download_bytes": "Total downloaded bytes per bucket",
    	"wasabi_num_billable_objects": "Total number of billable objects per bucket",
    	"wasabi_num_billable_deleted_objects": "Total number of billable deleted objects per bucket",
		//"wasabi_num_get_calls": "Total GET API calls per bucket",
		//"wasabi_num_put_calls": "Total PUT API calls per bucket",
		//"wasabi_num_delete_calls": "Total DELETE API calls per bucket",
		//"wasabi_num_list_calls": "Total LIST API calls per bucket",
		//"wasabi_num_head_calls": "Total HEAD API calls per bucket",
	}

	for name, help := range fields {
		metrics[name] = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: name,
				Help: help,
			},
			[]string{"bucket", "region", "account"},
		)
		prometheus.MustRegister(metrics[name])
	}
}

func fetchStats(apiKey string, accountName string) {
	pageNum := 0
	for {
		wasabiStatsURL := getWasabiStatsURL(pageNum)
		req, err := http.NewRequest("GET", wasabiStatsURL, nil)
		if err != nil {
			log.Println("Error creating request:", err)
			return
		}
		req.Header.Set("Authorization", fmt.Sprintf("%s", apiKey))

		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			log.Println("Error making request:", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Println("Error: Non-200 response", resp.Status)
			return
		}

		var response WasabiResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		if err != nil {
			log.Println("Error decoding response:", err)
			return
		}

		for _, bucket := range response.Records {
			metrics["wasabi_raw_storage_bytes"].WithLabelValues(bucket.Bucket, bucket.Region, accountName).Set(bucket.RawStorageSizeBytes)
			metrics["wasabi_deleted_storage_bytes"].WithLabelValues(bucket.Bucket, bucket.Region, accountName).Set(bucket.DeletedStorageSizeBytes)
			metrics["wasabi_padded_storage_bytes"].WithLabelValues(bucket.Bucket, bucket.Region, accountName).Set(bucket.PaddedStorageSizeBytes)
			metrics["wasabi_metadata_storage_bytes"].WithLabelValues(bucket.Bucket, bucket.Region, accountName).Set(bucket.MetadataStorageSizeBytes)
			//metrics["wasabi_orphaned_storage_bytes"].WithLabelValues(bucket.Bucket, bucket.Region, accountName).Set(bucket.OrphanedStorageSizeBytes)
			metrics["wasabi_storage_wrote_bytes"].WithLabelValues(bucket.Bucket, bucket.Region, accountName).Set(bucket.StorageWroteBytes)
			metrics["wasabi_storage_read_bytes"].WithLabelValues(bucket.Bucket, bucket.Region, accountName).Set(bucket.StorageReadBytes)
			metrics["wasabi_api_calls_total"].WithLabelValues(bucket.Bucket, bucket.Region, accountName).Set(bucket.NumAPICalls)
			metrics["wasabi_upload_bytes"].WithLabelValues(bucket.Bucket, bucket.Region, accountName).Set(bucket.UploadBytes)
			metrics["wasabi_download_bytes"].WithLabelValues(bucket.Bucket, bucket.Region, accountName).Set(bucket.DownloadBytes)
			metrics["wasabi_num_billable_objects"].WithLabelValues(bucket.Bucket, bucket.Region, accountName).Set(bucket.NumBillableObjects)
			metrics["wasabi_num_billable_deleted_objects"].WithLabelValues(bucket.Bucket, bucket.Region, accountName).Set(bucket.NumBillableDeletedObjects)
			//metrics["wasabi_num_get_calls"].WithLabelValues(bucket.Bucket, bucket.Region, accountName).Set(bucket.NumGETCalls)
            //metrics["wasabi_num_put_calls"].WithLabelValues(bucket.Bucket, bucket.Region, accountName).Set(bucket.NumPUTCalls)
            //metrics["wasabi_num_delete_calls"].WithLabelValues(bucket.Bucket, bucket.Region, accountName).Set(bucket.NumDELETECalls)
            //metrics["wasabi_num_list_calls"].WithLabelValues(bucket.Bucket, bucket.Region, accountName).Set(bucket.NumLISTCalls)
            //metrics["wasabi_num_head_calls"].WithLabelValues(bucket.Bucket, bucket.Region, accountName).Set(bucket.NumHEADCalls)
		}

		if pageNum >= response.PageInfo.PageCount-1 {
			break
		}
		pageNum++
	}
}

func main() {
	initMetrics()

	apiKeys := os.Getenv("WASABI_API_KEYS")
	if apiKeys == "" {
		log.Fatal("WASABI_API_KEYS environment variable not set")
	}

	keyPairs := strings.Split(apiKeys, ",")
	apiKeyMap := make(map[string]string)

	for _, pair := range keyPairs {
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) == 2 {
			apiKeyMap[parts[0]] = parts[1]
		} else {
			log.Println("Invalid API key format, expected name=key but got:", pair)
		}
	}

	go func() {
		for {
			for accountName, apiKey := range apiKeyMap {
				fetchStats(apiKey, accountName)
			}
			time.Sleep(60 * time.Minute)
		}
	}()
	registry := prometheus.NewRegistry()
	for _, metric := range metrics {
		registry.MustRegister(metric)
	}
	http.Handle("/metrics",  promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
	log.Println("Starting Prometheus metrics server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
