package orderbook

import (
	"encoding/json"
	"strconv"
	"sync"
	"time"
)

// https://docs.deribit.com/#book-instrument_name-group-depth-interval
type Orderbook struct {
	ask map[json.Number]json.Number
	bid map[json.Number]json.Number
	mut sync.Mutex
}

func New() *Orderbook {
	return &Orderbook{}
}

func (o *Orderbook) Empty() bool {
	{
		o.mut.Lock()
		defer o.mut.Unlock()
	}

	return len(o.ask) == 0 && len(o.bid) == 0
}

func (o *Orderbook) MarshalJSON() ([]byte, error) {
	{
		o.mut.Lock()
		defer o.mut.Unlock()
	}

	return json.Marshal(&struct {
		Ask map[json.Number]json.Number `json:"ask"`
		Bid map[json.Number]json.Number `json:"bid"`
		Tim time.Time                   `json:"tim"`
	}{
		Ask: o.ask,
		Bid: o.bid,
		Tim: time.Now().UTC().Round(time.Second),
	})
}

func (o *Orderbook) Middleware(upd Response) error {
	{
		o.mut.Lock()
		defer o.mut.Unlock()
	}

	{
		o.ask = map[json.Number]json.Number{}
		o.bid = map[json.Number]json.Number{}
	}

	for _, x := range upd.Asks {
		o.ask[musnum(x[0])] = musnum(x[1])
	}

	for _, x := range upd.Bids {
		o.bid[musnum(x[0])] = musnum(x[1])
	}

	return nil
}

func musnum(flo float64) json.Number {
	return json.Number(strconv.FormatFloat(flo, 'f', -1, 64))
}
