package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func RandomCreateGallery(t *testing.T) Gallery {
	account, err := testQueries.GetAccount(context.Background(), 1)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	item, err := testQueries.GetItem(context.Background(), 1)
	require.NoError(t, err)
	require.NotEmpty(t, item)

	arg := CreateGalleryParams{
		OwnerID: account.ID,
		ItemID: item.ID,
	}

	gallery, err := testQueries.CreateGallery(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, gallery)

	require.Equal(t, arg.OwnerID, gallery.OwnerID)
	require.Equal(t, arg.ItemID, gallery.ItemID)
	require.NotZero(t, gallery.CreatedAt)

	return gallery
}

func TestCreateGallery(t *testing.T) {
	RandomCreateGallery(t)
}

func TestGetGallery(t *testing.T) {
	gallery1 := RandomCreateGallery(t)

	gallery2, err := testQueries.GetGallery(context.Background(), gallery1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, gallery2)

	require.Equal(t, gallery1.OwnerID, gallery2.OwnerID)
	require.Equal(t, gallery1.ItemID, gallery2.ItemID)
	require.WithinDuration(t, gallery1.CreatedAt, gallery2.CreatedAt, time.Second)
}

func TestListGalleriesById(t *testing.T) {
	var gallery Gallery
	for i := 0; i < 5; i++ {
		gallery = RandomCreateGallery(t)
	}

	account, err := testQueries.GetAccount(context.Background(), 1)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	arg := ListGalleriesByIdParams{
		OwnerID: account.ID,
		Limit: 5,
		Offset: 0,
	}

	galleries, err := testQueries.ListGalleriesById(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, galleries)

	for _, g := range galleries {
		require.NotEmpty(t, g)
		require.Equal(t, gallery.OwnerID , g.OwnerID)
	}
}

func TestListGalleriesByItemId(t *testing.T) {
	var gallery Gallery
	for i := 0; i < 5; i++ {
		gallery = RandomCreateGallery(t)
	}

	item, err := testQueries.GetItem(context.Background(), 1)
	require.NoError(t, err)
	require.NotEmpty(t, item)

	arg := ListGalleriesByItemIdParams{
		ItemID: item.ID,
		Limit: 5,
		Offset: 0,
	}

	galleries, err := testQueries.ListGalleriesByItemId(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, galleries)

	for _, g := range galleries {
		require.NotEmpty(t, g)
		require.Equal(t, gallery.ItemID, g.ItemID)
	}
}
