package minterhub

type BlockNotFound struct {
	resp *GetBlockResponse
}

func NewBlockNotFoundError(resp *GetBlockResponse) *BlockNotFound {
	return &BlockNotFound{resp: resp}
}

func (e *BlockNotFound) Error() string {
	return *e.resp.Error
}
