#!/usr/bin/env bash

curl -X GET "localhost:9200/_cluster/health?wait_for_status=yellow&timeout=50s&pretty"