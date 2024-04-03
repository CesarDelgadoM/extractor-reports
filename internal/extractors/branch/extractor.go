package branch

import (
	"strings"

	"github.com/CesarDelgadoM/extractor-reports/internal/producer"
	"github.com/CesarDelgadoM/extractor-reports/internal/requests"
	"github.com/CesarDelgadoM/extractor-reports/pkg/httperrors"
	"github.com/CesarDelgadoM/extractor-reports/pkg/logger/zap"
	"github.com/CesarDelgadoM/extractor-reports/pkg/stream"
)

const (
	queueSuffix = "-restaurant-queue"
	bindSuffix  = "-restaurant-bind"
	batch       = 10
)

type IExtractor interface {
	ExtractData(params requests.RestaurantRequest)
}

type BranchExtractor struct {
	store      requests.ISet
	producer   producer.IChannel
	repository IBranchRepository
}

func NewBranchExtractor(store requests.ISet, producer producer.IChannel, repository IBranchRepository) IExtractor {
	return &BranchExtractor{
		store:      store,
		producer:   producer,
		repository: repository,
	}
}

func (e *BranchExtractor) ExtractData(params requests.RestaurantRequest) {
	// Store request in md5 hash
	e.store.Set(params.String())
	defer e.store.Delete(params.String())

	message := producer.Message{
		Userid: params.Userid,
		Type:   params.Type,
		Format: params.Format,
		Status: 1,
	}

	zap.Log.Info("Extracting restaurant data: ", params)

	restaurant, err := e.repository.Find(params.Userid, params.Name)
	if err != nil {
		zap.Log.Error(httperrors.ErrRestaurantNotFound, err)
		return
	}

	queuename := strings.ToLower(restaurant.Name) + queueSuffix
	bindkey := strings.ToLower(restaurant.Name) + bindSuffix

	// Publish queuename
	producer.PublishQueueName(e.producer, queuename)

	// Set restaurant data to message
	message.Data = *restaurant

	queue := e.producer.Queue(&stream.QueueOpts{
		Name: queuename,
	})

	e.producer.BindQueue(&stream.BindOpts{
		Name: queue.Name,
		Key:  bindkey,
	})

	// Publish restaurant data
	e.producer.Publish(&stream.PublishOpts{
		RoutingKey: bindkey,
		Body:       message.ToBytes(),
	})

	zap.Log.Info("Extracting branches: ", params)

	size := e.repository.Size(params.Userid, params.Name)
	if size == -1 {
		zap.Log.Error("Branches size not found: ", err)
		return
	}

	// Extraction of branches by batches
	var skip int64
	for size > 0 {
		branches, err := e.repository.GetPage(params.Userid, params.Name, skip, batch)
		if err != nil {
			zap.Log.Error("branches extraction failure: ", err)
			return
		}

		// Set branches data to message
		message.Data = *branches

		// Validate if extraction finished
		if size-batch <= 0 {
			message.Status = 0
		}

		// Publish branches data
		e.producer.Publish(&stream.PublishOpts{
			RoutingKey: bindkey,
			Body:       message.ToBytes(),
		})

		skip = skip + batch
		size = size - batch
	}
	zap.Log.Info("Extraction restaurant finished: ", params)
}
