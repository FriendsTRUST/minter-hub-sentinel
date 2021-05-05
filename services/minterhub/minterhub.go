package minterhub

import (
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
)

const (
	getLatestBlock = "/blocks/latest"
	getBlock       = "/blocks/{height}"
)

type Service struct {
	apiUrls             []string
	currentNodeApiIndex int
	testnet             bool
	logger              *logrus.Logger
	http                *resty.Client
}

func New(apiUrls []string, logger *logrus.Logger) (*Service, error) {
	http := resty.New().
		SetRetryCount(1).
		SetHostURL(apiUrls[0])

	s := &Service{
		apiUrls:             apiUrls,
		currentNodeApiIndex: 0,
		logger:              logger,
		http:                http,
	}

	return s, nil
}

func (svc *Service) GetLatestBlock() (*GetBlockResponse, error) {
	var resp *resty.Response
	var err error
	var res GetBlockResponse

	svc.try(func() error {
		resp, err = svc.http.R().
			SetResult(&res).
			SetError(&res).
			Get(getLatestBlock)

		if err != nil {
			return err
		}

		if resp.StatusCode() == 404 {
			return nil
		}

		return err
	})

	if resp.StatusCode() == 404 {
		return &res, NewBlockNotFoundError(&res)
	}

	return &res, err
}

func (svc *Service) GetBlock(height int) (*GetBlockResponse, error) {
	var resp *resty.Response
	var err error
	var res GetBlockResponse

	svc.try(func() error {
		resp, err = svc.http.R().
			SetPathParam("height", strconv.Itoa(height)).
			SetResult(&res).
			SetError(&res).
			Get(getBlock)

		if err != nil {
			return err
		}

		if resp.StatusCode() == 404 {
			return nil
		}

		return err
	})

	if resp.StatusCode() == 404 {
		return &res, NewBlockNotFoundError(&res)
	}

	return &res, err
}

func (svc *Service) try(callback func() error) {
	lastIndex := len(svc.apiUrls) - 1

	for i, url := range svc.apiUrls {
		svc.http.SetHostURL(url)

		err := callback()

		if err != nil {
			if i == lastIndex {
				return
			}

			continue
		}

		return
	}

	return
}
