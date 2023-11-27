package collector

import (
	"encoding/json"
	"fmt"
	"github.com/go-kit/log/level"
	"io"
	"net/http"
	"os"
)

func (c *DigicertCollector) FetchDigicertData() (*OrderList, error) {

	var orderList OrderList

	if c.digicertMock {
		// Open data file
		jsonFile, err := os.Open("mock.json")
		if err != nil {
			level.Error(c.logger).Log("err", fmt.Sprintf("cannot open json file : %s", err))
			return nil, err
		}
		data, err := io.ReadAll(jsonFile)
		defer jsonFile.Close()
		if err != nil {
			level.Error(c.logger).Log("err", fmt.Sprintf("error while ready the file : %s", err))
			return nil, err
		}

		err = json.Unmarshal(data, &orderList)
		if err != nil {
			level.Error(c.logger).Log("err", fmt.Sprintf("error while unmarshalling the body : %s", err))
			return nil, err
		}

		return &orderList, nil
	} else {
		req, err := http.NewRequest("GET", c.digicertAPIEndpoint, nil)
		if err != nil {
			level.Error(c.logger).Log("msg", err)
			return nil, err
		}

		// This one line implements the authentication required for the task.
		req.Header.Add("X-DC-DEVKEY", c.digicertAPIKey)
		req.Header.Add("Content-Type", "application/json")

		// Make request and show output.
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			level.Error(c.logger).Log("err", fmt.Sprintf("error while calling Digicert API : %s", err))
			return nil, err
		}

		body, err := io.ReadAll(resp.Body)
		defer resp.Body.Close()
		if err != nil {
			level.Error(c.logger).Log("err", fmt.Sprintf("error while ready the response body : %s", err))
			return nil, err
		}

		err = json.Unmarshal(body, &orderList)
		if err != nil {
			level.Error(c.logger).Log("err", fmt.Sprintf("error while unmarshalling the body : %s", err))
			return nil, err
		}

		return &orderList, nil
	}

}
