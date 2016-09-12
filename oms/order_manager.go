package oms

import (
	"fmt"
	"sync"
)

type ClOrdIDGenerator interface {
	Next() string
}

type OrderManager struct {
	sync.RWMutex
	orderID int
	clOrdID ClOrdIDGenerator

	orders        map[int]*Order
	clOrdIDLookup map[string]*Order
}

func NewOrderManager(idGen ClOrdIDGenerator) *OrderManager {
	return &OrderManager{
		clOrdIDLookup: make(map[string]*Order),
		orders:        make(map[int]*Order),
		clOrdID:       idGen,
	}
}

func (om *OrderManager) GetAll() []*Order {
	orders := make([]*Order, 0, len(om.clOrdIDLookup))
	for _, v := range om.clOrdIDLookup {
		orders = append(orders, v)
	}

	return orders
}

func (om *OrderManager) Get(id int) (*Order, error) {
	var err error
	order, ok := om.orders[id]
	if !ok {
		err = fmt.Errorf("could not find order with id %v", id)
	}

	return order, err
}

func (om *OrderManager) GetByClOrdID(clOrdID string) (*Order, error) {
	var err error
	order, ok := om.clOrdIDLookup[clOrdID]
	if !ok {
		err = fmt.Errorf("could not find order with clordid %v", clOrdID)
	}

	return order, err
}

func (om *OrderManager) Save(order *Order) error {
	order.ID = om.nextOrderID()
	order.ClOrdID = om.clOrdID.Next()

	om.orders[order.ID] = order
	om.clOrdIDLookup[order.ClOrdID] = order

	return nil
}

func (om *OrderManager) AssignNextClOrdID(order *Order) string {
	clOrdID := om.clOrdID.Next()
	om.clOrdIDLookup[clOrdID] = order
	return clOrdID
}

func (om *OrderManager) nextOrderID() int {
	om.orderID++
	return om.orderID
}
