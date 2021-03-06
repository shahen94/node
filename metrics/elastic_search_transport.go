/*
 * Copyright (C) 2019 The "MysteriumNetwork/node" Authors.
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package metrics

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/mysteriumnetwork/node/market/mysterium"
	"github.com/mysteriumnetwork/node/requests"
)

// NewElasticSearchTransport creates transport allowing to send events to ElasticSearch through HTTP
func NewElasticSearchTransport(url string, timeout time.Duration) Transport {
	return &elasticSearchTransport{http: newMysteriumHTTPTransport(timeout), url: url}
}

type elasticSearchTransport struct {
	http mysterium.HTTPTransport
	url  string
}

func (transport *elasticSearchTransport) sendEvent(event event) error {
	req, err := requests.NewPostRequest(transport.url, "/", event)
	if err != nil {
		return err
	}

	response, err := transport.http.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("error while reading response body: %v", err)
	}
	body := string(bodyBytes)

	if response.StatusCode != 200 {
		return fmt.Errorf("unexpected response status: %v, body: %v", response.Status, body)
	}

	if strings.ToUpper(body) != "OK" {
		return fmt.Errorf("unexpected response body: %v", body)
	}

	return nil
}

func newMysteriumHTTPTransport(timeout time.Duration) mysterium.HTTPTransport {
	return &http.Client{
		Transport: &http.Transport{
			//Don't reuse tcp connections for request - see ip/rest_resolver.go for details
			DisableKeepAlives: true,
		},
		Timeout: timeout,
	}
}
