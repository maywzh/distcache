# DistCache
Distributed Cache system based on Go

## Run
```bash
chmod +x ./run.sh && ./run.sh
```

## BenchMark
You can bench the system with `wrk`
```bash
wrk -t12 -c4000 -d1m http://localhost:9999/api?key=Tom
```