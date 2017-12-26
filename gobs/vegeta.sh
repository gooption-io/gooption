vegeta attack -duration=30s -rate=10 -targets=vegeta.txt > results.bin
cat results.bin | vegeta report -reporter=plot > plot.html