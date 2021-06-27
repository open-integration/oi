package call

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	urlpkg "net/url"

	"github.com/open-integration/oi/catalog/services/http/types"
	api "github.com/open-integration/oi/pkg/api/v1"
	"github.com/open-integration/oi/pkg/logger"
	"github.com/open-integration/oi/pkg/service"
)

func Handle(ctx context.Context, lgr logger.Logger, svc *service.Service, req *api.CallRequest) (*api.CallResponse, error) {
	args := &types.Arguments{}
	if err := service.UnmarshalRequestArgumentsInto(req, args); err != nil {
		return nil, err
	}
	if args.URL == nil {
		return service.BuildErrorResponse(fmt.Errorf("url is required"))
	}
	if args.Verb == nil {
		return service.BuildErrorResponse(fmt.Errorf("http verb is required"))
	}
	u, err := urlpkg.Parse(*args.URL)
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
		Method:     *args.Verb,
		URL:        u,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header:     headers,
		Body:       body,
		Host:       u.Host,
	}

	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	lgr.Info("Request returns", "status", resp.Status)
	respheaders := []types.ReturnsHeader{}
	for name, value := range resp.Header {
		respheaders = append(respheaders, types.ReturnsHeader{
			Name:  &name,
			Value: &value[0],
		})
	}
	defer resp.Body.Close()
	bodyData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	payload := &types.Returns{
		Status:  float64(resp.StatusCode),
		Headers: respheaders,
		Body:    string(bodyData),
	}
	return service.BuildSuccessfullResponse(payload)

}
