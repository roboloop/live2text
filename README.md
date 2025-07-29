# Live2Text

Converts speech to text using Google Cloud Speech-to-Text and shows the live transcription on the MacBook Touch Bar with BetterTouchTool.

## Preview

![preview.gif](./preview.gif)

## Features

- 🎙️ Real-time speech recognition from any audio source.
- 💻 Touch Bar support to display transcribed text.
- 🎧 Selectable audio input device for flexible source control.
- 🌐 Multi-language recognition with configurable language settings.
- 🖼️ Clean and embed modes to match your preferred UI experience.
- 📋 Clipboard integration to quickly copy recognized text.
- 📊 Built-in metrics to monitor traffic and app performance.
- 🔄 macOS background integration for seamless system-level operation.

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

## Known Issues

Device list is no updated.

## TODO

- Improve README with explanations and demonstrations
- Improve CI (remove duplicates, check the cache)
- Dockerfile?
- Repository badges
