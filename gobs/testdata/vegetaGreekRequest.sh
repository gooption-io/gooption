vegeta attack -duration=30s -rate=10 -targets=vegetaGreekRequest.txt > resultsGreekRequest.bin
cat resultsGreekRequest.bin | vegeta report -reporter=plot > plotGreekRequest.html