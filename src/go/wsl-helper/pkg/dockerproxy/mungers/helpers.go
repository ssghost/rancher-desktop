package mungers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// readRequestBodyJSON reads the incoming HTTP request body as if it was JSON,
// unmarshalled into the provided object.  A copy of the data is placed in the
// request body, so that it can be used by downstream consumers as necessary.
func readRequestBodyJSON(req *http.Request, data interface{}) error {
	buf, err := io.ReadAll(req.Body)
	if err != nil {
		return fmt.Errorf("could not read request body: %w", err)
	}

	err = json.Unmarshal(buf, data)
	req.Body = io.NopCloser(bytes.NewBuffer(buf))
	if err != nil {
		return fmt.Errorf("could not unmarshal request body: %w", err)
	}

	return nil
}

// readResponseBodyJSON reads the outgoing HTTP response body as if it was JSON,
// unmarshalled into the provided object.  A copy of the data is placed in the
// response body, so that it can be used directly if no modification neeeded to
// occur.
func readResponseBodyJSON(resp *http.Response, data interface{}) error {
	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("could not read response body: %w", err)
	}

	err = json.Unmarshal(buf, data)
	resp.Body = io.NopCloser(bytes.NewBuffer(buf))
	if err != nil {
		return fmt.Errorf("could not unmarshal response body: %w", err)
	}

	return nil
}
