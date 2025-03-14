# Wasabi Prometheus Exporter

This project is a Prometheus exporter for Wasabi bucket statistics. It fetches bucket utilization data from the Wasabi API and exposes it as Prometheus metrics.

## Features
- Periodically fetches Wasabi bucket statistics.
- Exposes bucket storage and API usage metrics via a Prometheus-compatible `/metrics` endpoint.
- Supports multiple Wasabi accounts.
- Handles API pagination for large datasets.
- Disables default Prometheus metrics to expose only relevant Wasabi statistics.
- Some API Call Informations are skipped as they are always empty at the moment

## Metrics Exported
The exporter exposes the following metrics for each bucket:
- `wasabi_raw_storage_bytes` - Total raw storage bytes used per bucket.
- `wasabi_deleted_storage_bytes` - Total deleted storage bytes per bucket.
- `wasabi_padded_storage_bytes` - Total padded storage bytes per bucket.
- `wasabi_metadata_storage_bytes` - Total metadata storage bytes per bucket.
- `wasabi_orphaned_storage_bytes` - Total orphaned storage bytes per bucket.
- `wasabi_storage_wrote_bytes` - Total bytes written to storage per bucket.
- `wasabi_storage_read_bytes` - Total bytes read from storage per bucket.
- `wasabi_api_calls_total` - Total number of API calls per bucket.
- `wasabi_upload_bytes` - Total bytes uploaded per bucket.
- `wasabi_download_bytes` - Total bytes downloaded per bucket.
- `wasabi_num_billable_objects` - Total number of billable objects per bucket.
- `wasabi_num_billable_deleted_objects` - Total number of billable deleted objects per bucket.

## Setup and Usage

### 1. Build and Run
#### Using Go
```sh
export WASABI_API_KEYS="account1=API_KEY_1,account2=API_KEY_2"
go build -o wasabi_exporter
./wasabi_exporter
```

#### Using Docker
```sh
docker build -t wasabi-prometheus-exporter .
docker run -e WASABI_API_KEYS="account1=API_KEY_1,account2=API_KEY_2" -p 8080:8080 wasabi-prometheus-exporter
```

### 2. Access Metrics
After running the service, Prometheus can scrape metrics from:
```
http://localhost:8080/metrics
```

## Environment Variables
- `WASABI_API_KEYS`: A comma-separated list of Wasabi API keys in the format `account_name=API_KEY`.

## License
This project is licensed under the MIT License.

