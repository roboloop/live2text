package btt

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

type event string

const (
	eventMessage event = "message"
	eventFailed  event = "failed"
)

func (s *Server) TextStream(w http.ResponseWriter, r *http.Request) {
	// http://localhost:8080/btt/floating-page
	defer r.Body.Close()
	w.Header().Set("Content-Type", "text/event-stream; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// Always 200 code for SSE
	w.WriteHeader(http.StatusOK)
	_, ok := w.(http.Flusher)
	if !ok {
		fmt.Fprintf(w, "event: %s\ndata: %s\n\n", eventFailed, "Streaming unsupported")
		return
	}

	textCh, errCh, err := s.services.Btt().StreamText(r.Context())
	if err != nil {
		sendSSE(w, eventFailed, err.Error())
		return
	}

	var text string
	for {
		select {
		case text, ok = <-textCh:
			if !ok {
				s.logger.InfoContext(r.Context(), "TextStream channel closed")
				return
			}

			sendSSE(w, eventMessage, strings.ReplaceAll(text, "\n", " "))
		case err, ok = <-errCh:
			if ok {
				s.logger.ErrorContext(r.Context(), "Error during stream", "error", err)
				sendSSE(w, eventFailed, err.Error())
				// time.Sleep(30 * time.Second)
				return
			}
		case <-r.Context().Done():
			s.logger.ErrorContext(r.Context(), "Request context cancelled", "error", r.Context().Err())
			return
		}
	}
}

func sendSSE(w io.Writer, event event, msg string) {
	fmt.Fprintf(w, "event: %s\ndata: %s\n\n", event, msg)
	w.(http.Flusher).Flush()
}
