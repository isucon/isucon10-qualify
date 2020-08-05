package client

import (
	"context"
	"strconv"

	"golang.org/x/sync/errgroup"

	"github.com/isucon10-qualify/isucon10-qualify/bench/asset"
)

func (c *Client) AccessTopPage(ctx context.Context) error {
	eg, childCtx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		_, err := c.GetRecommendedChair(childCtx)
		return err
	})

	eg.Go(func() error {
		_, err := c.GetRecommendedEstate(childCtx)
		return err
	})

	if err := eg.Wait(); err != nil {
		return err
	}

	return nil
}

func (c *Client) AccessChairDetailPage(ctx context.Context, id int64) (*asset.Chair, *EstatesResponse, error) {
	eg, childCtx := errgroup.WithContext(ctx)

	chairCh := make(chan *asset.Chair, 1)
	estatesCh := make(chan *EstatesResponse, 1)

	eg.Go(func() error {
		chair, err := c.GetChairDetailFromID(childCtx, strconv.FormatInt(id, 10))
		if err != nil {
			chairCh <- nil
			return err
		}
		if chair == nil {
			chairCh <- nil
			return nil
		}

		chairCh <- chair
		return nil
	})

	eg.Go(func() error {
		estates, err := c.GetRecommendedEstatesFromChair(childCtx, id)
		if err != nil {
			estatesCh <- nil
			return err
		}

		estatesCh <- estates
		return nil
	})

	if err := eg.Wait(); err != nil {
		return nil, nil, err
	}

	return <-chairCh, <-estatesCh, nil
}

func (c *Client) AccessEstateDetailPage(ctx context.Context, id int64) (*asset.Estate, error) {
	estate, err := c.GetEstateDetailFromID(ctx, strconv.FormatInt(id, 10))
	if err != nil {
		return nil, err
	}

	return estate, nil
}

func (c *Client) AccessChairSearchPage(ctx context.Context) error {
	_, err := c.GetChairSearchCondition(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) AccessEstateSearchPage(ctx context.Context) error {
	_, err := c.GetEstateSearchCondition(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) AccessEstateNazottePage(ctx context.Context) error {
	return nil
}
