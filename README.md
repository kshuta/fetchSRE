# Fetch SRE challenge
This repo serves as the submission for the Fetch SRE challenge.

# Running the submission
## Requirements
- Go version 1.22

## Steps:
1. Clone the repo.

2. Run `go mod tidy` to install dependencies.

3. Build the app with `go build ./cmd/app`.

4. (Optional) Start the test server: `go run ./cmd/testserver`

This test server is designed to have a latency of 1 seconds every 3 requests. This can be used for testing as needed.

5. Run the app with `./app {yaml file with endpoints}`.

