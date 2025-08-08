# Live2Text

**Live2Text** converts speech to text using Google Cloud Speech-to-Text and displays the live transcription on the MacBook Touch Bar using BetterTouchTool.

## 📺 Preview

![preview.gif](./preview.gif)

## ✨ Features

- Real-time speech recognition from any audio source
- Touch Bar display of transcribed text
- Selectable audio input devices
- Selectable recognition language
- Clean and embedded display modes
- Clipboard copy support
- Usage and performance metrics
- Background operation on macOS

## ⚙️ Requirements

- macOS
- [Google Cloud account](https://cloud.google.com/)
- [gcloud CLI](https://cloud.google.com/sdk/docs/install)
- [Better Touch Tool](https://folivora.ai/)

## 🚀 Installation

### Using Homebrew

```shell
brew tap roboloop/tap
brew install live2text
```

### Static build

Download from the [release page](https://github.com/roboloop/live2text/releases).

### Go package

[portaudio](https://www.portaudio.com/) required

```shell
go install github.com/roboloop/live2text/cmd/live2text@latest
```

## 🔧 Setup

1. Enable the Web Server in BetterTouchTool and note the listening port

2. Install the BetterTouchTool integration:

    ```shell
    live2text install [--args]
    ```

3. Start the background service:

    ```shell
    live2text serve [--args]
    ```

## 🐛 Known Issues

**Issue**: Audio device list does not update after changes  
**Solution**: Restart the application

## 📌 TODO

- Add configuration file support
- Refactor CI pipeline
- Fix `go test -race` failures
- Add repository badges


