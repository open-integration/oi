package call

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"net/http"
	urlpkg "net/url"

	api "github.com/open-integration/oi/pkg/api/v1"
	"github.com/open-integration/oi/pkg/logger"
	"github.com/open-integration/oi/pkg/service"
)

func Handle(ctx context.Context, lgr logger.Logger, svc *service.Service, req *api.CallRequest) (*api.CallResponse, error) {
	args := &CallArguments{}
	if err := service.UnmarshalRequestArgumentsInto(req, args); err != nil {
		return nil, err
	}
	u, err := urlpkg.Parse(args.URL)
	if err != nil {
		return nil, err
	}
	var body io.ReadCloser
	if args.Content != nil {
		body = ioutil.NopCloser(bytes.NewReader([]byte(*args.Content)))
	}

	client := http.Client{}
	headers := http.Header{}
	for _, h := range args.Headers {
		headers.Set(*h.Name, *h.Value)
	}

	request := &http.Request{
		Method:     args.Verb,
		URL:        u,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header:     headers,
		Body:       body,
		Host:       u.Host,
	}
	request.Method = args.Verb

	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	lgr.Info("Request returns", "status", resp.Status)
	respheaders := []Header{}
	for name, value := range resp.Header {
		respheaders = append(respheaders, Header{
			Name:  &name,
			Value: &value[0],
		})
	}
	defer resp.Body.Close()
	bodyData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	payload := &CallReturns{
		Status:  resp.StatusCode,
		Headers: respheaders,
		Body:    string(bodyData),
	}
	return service.BuildSuccessfullResponse(payload)

}
