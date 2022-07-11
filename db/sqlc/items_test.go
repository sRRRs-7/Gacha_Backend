package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/sRRRs-7/GachaPon/utils"
	"github.com/stretchr/testify/require"
)

func RandomCreateItem(t *testing.T) Item {
	itemName := utils.RandomItemName()
	itemUrl := utils.RandomItemUrl()

	category, err := testQueries.GetCategory(context.Background(), "vuqlrn")
	require.NoError(t, err)
	require.NotEmpty(t, category)

	arg := CreateItemParams{
		ItemName: itemName,
		Rating: 3,
		ItemUrl: itemUrl,
		CategoryID: int32(category.ID),
	}

	item, err := testQueries.CreateItem(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, item)

	require.Equal(t, arg.ItemName, item.ItemName)
	require.Equal(t, arg.Rating, item.Rating)
	require.Equal(t, arg.ItemUrl, item.ItemUrl)
	require.Equal(t, arg.CategoryID, item.CategoryID)

	require.NotZero(t, item.CreatedAt)

	return item
}

func TestCreateItem(t *testing.T) {
	RandomCreateItem(t)
}

func TestGetItem(t *testing.T) {
	item1 := RandomCreateItem(t)

	item2, err := testQueries.GetItem(context.Background(), item1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, item2)

	require.Equal(t, item1.ItemName, item2.ItemName)
	require.Equal(t, item1.Rating, item2.Rating)
	require.Equal(t, item1.ItemUrl, item2.ItemUrl)
	require.Equal(t, item1.CategoryID, item2.CategoryID)

	require.WithinDuration(t, item1.CreatedAt, item2.CreatedAt, time.Second)
}

func TestListItemByCategoryId(t *testing.T) {
	item1, err := testQueries.GetItem(context.Background(), 1)
	require.NoError(t, err)
	require.NotEmpty(t, item1)

	arg := ListItemByCategoryIdParams{
		CategoryID: item1.CategoryID,
		Limit: 5,
		Offset: 0,
	}

	item2, err := testQueries.ListItemByCategoryId(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, item2)
}

func TestListItemsByCategoryId(t *testing.T) {
	arg:= ListItemsByCategoryIdParams{
		Limit: 5,
		Offset: 0,
	}

	items, err := testQueries.ListItemsByCategoryId(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, items)
}

func TestListItemsById(t *testing.T) {
	arg:= ListItemsByIdParams{
		Limit: 5,
		Offset: 0,
	}

	items, err := testQueries.ListItemsById(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, items)
}

func TestListItemsByItemName(t *testing.T) {
	arg:= ListItemsByItemNameParams{
		Limit: 5,
		Offset: 0,
	}

	items, err := testQueries.ListItemsByItemName(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, items)
}

func TestListItemsByRating(t *testing.T) {
	arg:= ListItemsByRatingParams{
		Limit: 5,
		Offset: 0,
	}

	items, err := testQueries.ListItemsByRating(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, items)
}

func TestUpdateItem(t *testing.T) {
	item1, err := testQueries.GetItem(context.Background(), 1)
	require.NoError(t, err)
	require.NotEmpty(t, item1)

	arg := UpdateItemParams{
		ID: item1.ID,
		ItemName: "sss",
		Rating: item1.Rating,
		ItemUrl: item1.ItemUrl,
		CategoryID: item1.CategoryID,
	}

	item2, err := testQueries.UpdateItem(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, item2)

	require.Equal(t, arg.ItemName, item2.ItemName)
	require.Equal(t, arg.Rating, item2.Rating)
	require.Equal(t, arg.ItemUrl, item2.ItemUrl)
	require.Equal(t, arg.CategoryID, item2.CategoryID)

	require.WithinDuration(t, item1.CreatedAt, item2.CreatedAt, time.Second)
}

func TestDeleteItem(t *testing.T) {
	item1 := RandomCreateItem(t)
	err := testQueries.DeleteItem(context.Background(), item1.ID)
	require.NoError(t, err)

	item2, err := testQueries.GetItem(context.Background(), item1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, item2)

}