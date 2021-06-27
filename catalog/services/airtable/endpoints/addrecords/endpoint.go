package addrecords

import (
	"context"

	at "github.com/mehanizm/airtable"
	"github.com/open-integration/oi/catalog/services/airtable/types"
	"github.com/open-integration/oi/catalog/shared/airtable"
	api "github.com/open-integration/oi/pkg/api/v1"
	"github.com/open-integration/oi/pkg/logger"
	"github.com/open-integration/oi/pkg/service"
)

func Handle(ctx context.Context, lgr logger.Logger, svc *service.Service, req *api.CallRequest) (*api.CallResponse, error) {
	args := &types.AddRecordsArguments{}
	if err := service.UnmarshalRequestArgumentsInto(req, args); err != nil {
		return service.BuildErrorResponse(err)
	}
	records := []*at.Record{}
	for _, r := range args.Records {
		fields := map[string]interface{}{}
		for k, v := range r.Fields {
			fields[k] = v
		}
		records = append(records, &at.Record{
			Fields: fields,
		})
	}
	chunks := chunkSlice(records, 10)
	table := airtable.GetTable(args.Auth.APIKey, args.Auth.DatabaseID, args.Auth.TableName)
	lgr.Info("Adding records", "total", len(records), "chunks", len(chunks))
	success := []types.Record{}
	for _, c := range chunks {
		res, err := table.AddRecords(&at.Records{
			Records: c,
		})
		if err != nil {
			lgr.Info("failed to add chunk", "err", err.Error())
		}
		for _, r := range res.Records {
			success = append(success, types.Record{
				Fields:      r.Fields,
				CreatedTime: &r.CreatedTime,
				Deleted:     &r.Deleted,
				ID:          &r.ID,
				Typecast:    &r.Typecast,
			})
		}

	}
	return service.BuildSuccessfullResponse(types.AddRecordsReturns{
		Records: success,
	})
}

func chunkSlice(slice []*at.Record, chunkSize int) [][]*at.Record {
	var chunks [][]*at.Record
	for {
		if len(slice) == 0 {
			break
		}

		// necessary check to avoid slicing beyond
		// slice capacity
		if len(slice) < chunkSize {
			chunkSize = len(slice)
		}

		chunks = append(chunks, slice[0:chunkSize])
		slice = slice[chunkSize:]
	}

	return chunks
}
