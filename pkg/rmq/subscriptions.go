package rmq

import (
	"context"
)

func (c *Consumer) SubscribeForDetailedRoutes(handler handlerFunc, publisher *Publisher) (func(context.Context), error) {
	f, err := c.subscribe(GetMetaDetailedRoutesAccept(), GetMetaDetailedRoutesTransmit(), handler, publisher)
	if err != nil {
		return nil, err
	}

	c.logger.Infof("Subscribed for detailed routes")

	return f, nil
}

func (c *Consumer) ListenRPCForDetailedRoutes(handler handlerFunc, publisher *Publisher) (func(context.Context), error) {
	f, err := c.listenRPC(GetMetaDetailedRoutesRPC(), handler, publisher)
	if err != nil {
		return nil, err
	}

	c.logger.Infof("Listening RPC for detailed routes")

	return f, nil
}

func (c *Consumer) subscribe(
	metaAccept *Meta,
	metaTransmit *Meta,
	handler handlerFunc,
	publisher *Publisher,
) (func(context.Context), error) {
	c.logger.Infof("meta to accept: %v", *metaAccept)
	c.logger.Infof("meta to transmit: %v", *metaTransmit)

	delivery, err := c.Subscribe(metaAccept)
	if err != nil {
		return nil, err
	}

	return publisher.wrapHandler(metaTransmit, delivery, false, handler), nil
}

func (c *Consumer) listenRPC(
	meta *Meta,
	handler handlerFunc,
	publisher *Publisher,
) (func(context.Context), error) {
	c.logger.Infof("meta: %v", *meta)

	delivery, err := c.WorkOnTask(meta.Key, meta.PrefetchCount)
	if err != nil {
		return nil, err
	}

	return publisher.wrapHandler(meta, delivery, true, handler), nil
}
