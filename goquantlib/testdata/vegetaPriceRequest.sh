vegeta attack -duration=30s -rate=10 -targets=vegetaPriceRequest.txt > resultsPriceRequest.bin
cat resultsPriceRequest.bin | vegeta report -reporter=plot > plotPriceRequest.html