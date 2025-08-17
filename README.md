# Some info about this project


### Makefile instructions

✅ How to use this Makefile on Windows

Set GOARCH=arm or GOARCH=arm64 depending on your Pi OS.

Run make or make build → produces a Linux binary like:

gps-microservice-v0.1.2-0-g297ad9c-arm64


Copy the binary to your Raspberry Pi using scp or similar.

On the Pi:

chmod +x gps-microservice-v0.1.2-0-g297ad9c-arm64
./gps-microservice-v0.1.2-0-g297ad9c-arm64


Run make test, make vet, make lint on Windows before cross-compiling to ensure code quality.