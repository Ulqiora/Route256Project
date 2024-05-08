package middleware

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/IBM/sarama"
	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

type SenderMessages interface {
	GetTopic() string
	SendAsyncMessage(message *sarama.ProducerMessage)
	Close() error
}
type RequestInfo struct {
	URL      string              `json:"URL"`
	Headers  map[string][]string `json:"headers"`
	BodyData string              `json:"body"`
}

func LogMiddleware(producer SenderMessages) func(_ http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			if request.Method == http.MethodGet ||
				request.Method == http.MethodPost ||
				request.Method == http.MethodDelete {

				message, err := configureMessage(request, producer.GetTopic())
				if err != nil {
					writer.WriteHeader(http.StatusInternalServerError)
					writer.Write([]byte(err.Error()))
					return
				}
				producer.SendAsyncMessage(message)
			}
			handler.ServeHTTP(writer, request)
		})
	}
}

func SimpleLogMiddleware() func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			if request.Method == http.MethodGet ||
				request.Method == http.MethodPost ||
				request.Method == http.MethodDelete {

				queryParams := mux.Vars(request)
				slog.Info(fmt.Sprintf("%s: %s", "MethodType", request.Method))
				for key, value := range queryParams {
					slog.Info(fmt.Sprintf("%s: %s", key, value))
				}
			}
			handler.ServeHTTP(writer, request)
		})
	}
}

func configureMessage(request *http.Request, topic string) (*sarama.ProducerMessage, error) {
	key, _ := uuid.NewV4()
	body := request.Body
	defer body.Close()

	// Чтение данных из тела запроса
	data := make([]byte, request.ContentLength)
	_, err := body.Read(data)
	if err != nil {
		return nil, errors.Wrap(err, "error reading body")
	}
	requestInfo := RequestInfo{
		URL:      request.URL.String(),
		Headers:  request.Header,
		BodyData: string(data),
	}
	bytesAll, err := json.Marshal(requestInfo)
	if err != nil {
		return nil, err
	}
	message := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key.String()),
		Value: sarama.ByteEncoder(bytesAll),
	}
	return message, nil
}
