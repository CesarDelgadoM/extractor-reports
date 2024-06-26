package branch

import (
	"strings"

	"github.com/CesarDelgadoM/extractor-reports/internal/producer"
	"github.com/CesarDelgadoM/extractor-reports/internal/producer/databus"
	"github.com/CesarDelgadoM/extractor-reports/internal/requests"
	"github.com/CesarDelgadoM/extractor-reports/internal/utils"
	"github.com/CesarDelgadoM/extractor-reports/pkg/httperrors"
	"github.com/CesarDelgadoM/extractor-reports/pkg/logger/zap"
	"github.com/CesarDelgadoM/extractor-reports/pkg/stream"
)

const (
	restaurant_type = "restaurant"
	branch_type     = "branch"

	reportType  = "branch"
	queueSuffix = "-branches-queue-"
	bindSuffix  = "-branches-bind-"

	batch = 10
)

type IExtractor interface {
	ExtractData(params requests.RestaurantRequest)
}

type BranchExtractor struct {
	store      requests.ISet
	databus    databus.IDataBus
	producer   producer.IProducer
	repository IBranchRepository
}

func NewBranchExtractor(store requests.ISet, databus databus.IDataBus, producer producer.IProducer, repository IBranchRepository) IExtractor {
	return &BranchExtractor{
		store:      store,
		databus:    databus,
		producer:   producer,
		repository: repository,
	}
}

func (e *BranchExtractor) ExtractData(params requests.RestaurantRequest) {
	defer e.store.Delete(params.String())

	// Store request in md5 hash
	e.store.Set(params.String())

	message := producer.Message{
		Userid: params.Userid,
		Format: params.Format,
		Type:   restaurant_type,
		Status: 1,
	}

	zap.Log.Info("Extracting restaurant data: ", params)

	restaurant, err := e.repository.Find(params.Userid, params.Name)
	if err != nil {
		zap.Log.Error(httperrors.ErrRestaurantNotFound, err)
		return
	}

	timestamp := utils.TimestampID()
	queuename := strings.ToLower(restaurant.Name) + queueSuffix + timestamp
	bindkey := strings.ToLower(restaurant.Name) + bindSuffix + timestamp

	// Publish queuename to data bus
	e.databus.PublishQueueName(producer.MessageQueueNames{
		ReportType: reportType,
		QueueName:  queuename,
	})

	// Set restaurant data to message
	message.Data = utils.ToBytes(restaurant)

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
		Body:       utils.ToBytes(message),
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
		message.Data = utils.ToBytes(branches)
		message.Type = branch_type

		// Validate if extraction finished
		if size-batch <= 0 {
			message.Status = 0
		}

		// Publish branches data
		e.producer.Publish(&stream.PublishOpts{
			RoutingKey: bindkey,
			Body:       utils.ToBytes(message),
		})

		skip = skip + batch
		size = size - batch
	}
	zap.Log.Info("Extraction restaurant finished: ", params)
}
