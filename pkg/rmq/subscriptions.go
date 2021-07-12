package rmq

import (
	"context"
)

func (c *Client) SubscribeForDetailedRoutes(handler handlerFunc) (func(context.Context), error) {
	f, err := c.subscribe(MetaDetailedRoutesAccept, MetaDetailedRoutesTransmit, handler)
	if err != nil {
		return nil, err
	}

	c.logger.Infof("Subscribed for detailed routes")

	return f, nil
}

func (c *Client) ListenRPCForDetailedRoutes(handler handlerFunc) (func(context.Context), error) {
	f, err := c.listenRPC(MetaDetailedRoutesRPC, handler)
	if err != nil {
		return nil, err
	}

	c.logger.Infof("Listening RPC for detailed routes")

	return f, nil
}

func (c *Client) subscribe(
	metaAccept *Meta,
	metaTransmit *Meta,
	handler handlerFunc,
) (func(context.Context), error) {
	c.logger.Infof("meta to accept: %v", metaAccept)
	c.logger.Infof("meta to transmit: %v", metaTransmit)

	delivery, err := c.Subscribe(metaAccept)
	if err != nil {
		return nil, err
	}

	return c.wrapHandler(metaTransmit, delivery, handler), nil
}

func (c *Client) listenRPC(meta *Meta, handler handlerFunc) (func(context.Context), error) {
	c.logger.Infof("meta: %v", meta)

	delivery, err := c.WorkOnTask(meta.QName, 1)
	if err != nil {
		return nil, err
	}

	return c.wrapHandler(meta, delivery, handler), nil
}
