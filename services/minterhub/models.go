package minterhub

import "time"

type GetBlockResponse struct {
	Block struct {
		Header struct {
			Height int `json:"height,string"`
		} `json:"header"`
		LastCommit struct {
			Height     string `json:"height"`
			Signatures []struct {
				ValidatorAddress string    `json:"validator_address"`
				Timestamp        time.Time `json:"timestamp"`
				Signature        string    `json:"signature"`
			} `json:"signatures"`
		} `json:"last_commit"`
	} `json:"block"`

	Error *string `json:"error"`
}
