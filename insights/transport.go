package insights

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/Microsoft/ApplicationInsights-Go/appinsights/contracts"
	"github.com/docker/docker/pkg/urlutil"
	"github.com/sirupsen/logrus"

	"gitlab.com/michael.golfi/appinsights/constants"
)

func parseURL(endpoint string) (*url.URL, error) {
	if !urlutil.IsURL(endpoint) {
		return nil, fmt.Errorf("expected endpoint format https://dns_name_or_ip:port, received: %s", endpoint)
	}

	insightsURL, err := url.Parse(endpoint)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to parse %s as url value in %s", constants.DriverName, endpoint, constants.EndpointKey)
	}

	if !insightsURL.IsAbs() || insightsURL.Path == "" || insightsURL.Path == "/" || insightsURL.RawQuery != "" || insightsURL.Fragment != "" {
		return nil, fmt.Errorf("expected endpoint format scheme://dns_name_or_ip:port, received: %s", endpoint)
	}

	return insightsURL, nil
}

func verifyInsightsConnection(uri string) error {
	client := http.Client{}
	req, err := http.NewRequest(http.MethodOptions, uri, nil)
	if err != nil {
		return err
	}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	if res.Body != nil {
		defer res.Body.Close()
	}
	if res.StatusCode != http.StatusOK {
		var body []byte
		body, err = ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("%s: failed to verify connection - %s - %s", constants.DriverName, res.Status, body)
	}
	return nil
}

// REVIEW
func (l *insightsLogger) queueMessageAsync(message *contracts.Envelope) error {
	l.lock.RLock()
	defer l.lock.RUnlock()
	if l.closedCond != nil {
		return fmt.Errorf("%s: driver is closed", constants.DriverName)
	}
	l.stream <- message
	return nil
}

// REVIEW
func (l *insightsLogger) postMessages(messages []*contracts.Envelope, lastChance bool) []*contracts.Envelope {
	messagesLen := len(messages)

	ctx, cancel := context.WithTimeout(context.Background(), l.sendTimeout)
	defer cancel()

	for i := 0; i < messagesLen; i += l.postMessagesBatchSize {
		upperBound := i + l.postMessagesBatchSize
		if upperBound > messagesLen {
			upperBound = messagesLen
		}

		if err := l.tryPostMessages(ctx, messages[i:upperBound]); err != nil {
			logrus.WithError(err).WithField("module", "logger/appinsights").Warn("Error while sending logs")
			if messagesLen-i >= l.bufferMaximum || lastChance {
				// If this is last chance - print them all to the daemon log
				if lastChance {
					upperBound = messagesLen
				}
				// Not all sent, but buffer has got to its maximum, let's log all messages
				// we could not send and return buffer minus one batch size
				for j := i; j < upperBound; j++ {
					if jsonEvent, err := json.Marshal(messages[j]); err != nil {
						logrus.Error(err)
					} else {
						logrus.Error(fmt.Errorf("failed to send a message '%s'", string(jsonEvent)))
					}
				}
				return messages[upperBound:messagesLen]
			}
			// Not all sent, returning buffer from where we have not sent messages
			return messages[i:messagesLen]
		}
	}
	// All sent, return empty buffer
	return messages[:0]
}

// REVIEW
func (l *insightsLogger) tryPostMessages(ctx context.Context, messages []*contracts.Envelope) error {
	if len(messages) == 0 {
		return nil
	}
	var buffer bytes.Buffer
	var writer io.Writer
	var gzipWriter *gzip.Writer
	var err error
	// If gzip compression is enabled - create gzip writer with specified compression
	// level. If gzip compression is disabled, use standard buffer as a writer
	if l.gzipCompression {
		gzipWriter, err = gzip.NewWriterLevel(&buffer, l.gzipCompressionLevel)
		if err != nil {
			return err
		}
		writer = gzipWriter
	} else {
		writer = &buffer
	}
	for _, message := range messages {
		jsonEvent, err := json.Marshal(message)
		if err != nil {
			return err
		}
		if _, err := writer.Write(jsonEvent); err != nil {
			return err
		}
	}
	// If gzip compression is enabled, tell it, that we are done
	if l.gzipCompression {
		err = gzipWriter.Close()
		if err != nil {
			return err
		}
	}
	req, err := http.NewRequest("POST", l.url, bytes.NewBuffer(buffer.Bytes()))
	if err != nil {
		return err
	}
	req = req.WithContext(ctx)
	// Tell if we are sending gzip compressed body
	if l.gzipCompression {
		req.Header.Set("Content-Encoding", "gzip")
	}
	res, err := l.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		var body []byte
		body, err = ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("%s: failed to send event - %s - %s", constants.DriverName, res.Status, body)
	}
	io.Copy(ioutil.Discard, res.Body)
	return nil
}
