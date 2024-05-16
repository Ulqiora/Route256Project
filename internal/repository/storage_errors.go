package repository

import "errors"

var ErrorNotFoundedStateOrder = errors.New("this state was not founded")
var ErrorOrderNotCreated = errors.New("this order was not created")
var ErrorOrderNotFounded = errors.New("this order was not founded")
var ErrorDataBase = errors.New("internal error of DockerfileDatabase")
var ErrorIssuingOrderForCustomer = errors.New("this orders can't be issuing")
var ErrorTimeOutForReturnOrder = errors.New("it is not possible to return the product because the time for a refund has ended")
var ErrorOrderIsNotReceived = errors.New("this order was not received")
var ErrorOrderCantReturnedToCourier = errors.New("this order was not returned to order_courier")

var ErrorObjectNotFounded = errors.New("not found")
var ErrorObjectNotCreated = errors.New("not created")
