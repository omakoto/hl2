- How to capture a prof:
./scripts/build.sh -r -cpuprofile hl.prof -r ./samples/highlighter-logcat.toml <./samples/sample.log | wc -l

- How to view:
echo "web" | go tool pprof hl.prof
