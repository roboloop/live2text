# Live2Text

Converts speech to text using Google Cloud Speech-to-Text and shows the live transcription on the MacBook Touch Bar with BetterTouchTool.

## Requirements

- OS: MacOS
- [Google Cloud account](https://cloud.google.com/)
- [gcloud CLI](https://cloud.google.com/sdk/docs/install)
- [Better Touch Tool](https://folivora.ai/)
- [portaudio](https://formulae.brew.sh/formula/portaudio)

## Local publishing

```shell
make build
```

```shell
./bin/btt -btt-port 44444
```

```shell
./bin/live2text -btt-port 44444
```

## TODO

- Improve README with explanations and demonstrations
- Improve CI (remove duplicates, check the cache)
- Dockerfile?
- Repository badges
