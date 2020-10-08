# SAP Receiver

**Status: Not fully implemented**
This receiver can not be used yet.

## Overview
SAP receiver skeleton which just wakes up at the given interval and gets the size of the given website.

## Configuration

Example:

```yaml
receivers:
  sap:
    url: "http://www.splunk.com"
    collection_interval: 30s
```

### url (Optional)
The address of the the endpoint to get the metrics from.

Default: `http://www.splunk.com`

### collection_interval (Optional)
How often to collect the metrics.

Default: `30s`

