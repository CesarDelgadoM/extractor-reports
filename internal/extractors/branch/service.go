package branch

import (
	"github.com/CesarDelgadoM/extractor-reports/internal/requests"
	"github.com/CesarDelgadoM/extractor-reports/internal/workerpool"
	"github.com/CesarDelgadoM/extractor-reports/pkg/httperrors"
	"github.com/CesarDelgadoM/extractor-reports/pkg/logger/zap"
)

type IService interface {
	ProducerReport(params requests.RestaurantRequest) error
}

type BranchService struct {
	dispatcher *workerpool.WorkerPool
	store      requests.ISet
	extractor  IExtractor
}

func NewBranchService(dispatcher *workerpool.WorkerPool, store requests.ISet, extractor IExtractor) IService {
	return &BranchService{
		dispatcher: dispatcher,
		store:      store,
		extractor:  extractor,
	}
}

func (s *BranchService) ProducerReport(params requests.RestaurantRequest) error {
	// Validates that the same request not executed many times
	if s.store.Exist(params.String()) {
		zap.Log.Warn("Request is already processing")
		return httperrors.RequestAlreadyGenerating
	}

	s.dispatcher.Submit(func() {
		s.extractor.ExtractData(params)
	})

	return nil
}
