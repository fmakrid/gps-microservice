GPS Microservice

This project provides a GPS microservice written in Go. You can build it for Raspberry Pi (Linux ARM/ARM64) from Windows or Linux, run tests on Windows, and deploy to production.

Makefile Instructions
Development & Testing (Windows)

Run all tests and checks locally before building for production:

# Run unit tests only
make test

# Run go vet for static code checks
make vet

# Run golint for style/linting checks
make lint

# Run all checks (vet + lint + tests)
make check


Notes:

Test files are _test.go files, using helpers from helpers_test.go.

go test may show paths like:

ok      github.com/fmakrid/gps-microservice     (cached)


(cached) means Go reused previous results. Use go test -count=1 ./... to force rerun.

Production Build (Linux ARM/ARM64)

Set the target architecture depending on your Raspberry Pi OS:

# For 32-bit OS
set GOARCH=arm

# For 64-bit OS
set GOARCH=arm64


Build the Linux binary (from Windows):

make build


This produces a binary like:

gps-microservice-v0.1.2-0-g297ad9c-arm64


Copy the binary to your Raspberry Pi:

scp gps-microservice-v0.1.2-0-g297ad9c-arm64 pi@raspberrypi.local:/home/pi/


On the Pi:

chmod +x gps-microservice-v0.1.2-0-g297ad9c-arm64
./gps-microservice-v0.1.2-0-g297ad9c-arm64

Notes

The Makefile supports cross-compilation from Windows.

All development tests must run on Windows to ensure code correctness before producing the Linux binary.

Make sure Go modules are up-to-date:

go mod tidy