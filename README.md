# Recognizer

Converts speech to text using Google Cloud Speech-to-Text and shows the live transcription on the MacBook Touch Bar with BetterTouchTool.

## Requirements

- OS: MacOS
- [Google Cloud account](https://cloud.google.com/)
- [gcloud CLI](https://cloud.google.com/sdk/docs/install)
- [Better Touch Tool](https://folivora.ai/)
- [jq](https://github.com/jqlang/jq)

## Local publish

```shell
go build -o ~/go/bin/live2text cmd/live2text/main.go
```

## TODO

- Refactor code (ensure consistent message formats in tests, errors, logs, and use tags instead of \[prefix\] in logs)
- Fully implement graceful shutdown
- Handle port occupation cases
- Properly manage context cancellation
- Integrate HTTP calls in btt.sh
- Fix goroutine leaks
- Update Go version
- Improve README with explanations and demonstrations
- Add tests, set up linter, and configure CI/CD
