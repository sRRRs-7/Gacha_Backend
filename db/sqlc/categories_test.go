package db

import (
	"context"
	"testing"

	"github.com/sRRRs-7/GachaPon/utils"
	"github.com/stretchr/testify/require"
)

func RandomCategory(t *testing.T) Category {
	arg := utils.RandomCategory()

	category, err := testQueries.CreateCategory(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, category)

	require.Equal(t, arg, category.Category)
	require.NotZero(t, category.CreatedAt)

	return category
}

func TestCreateCategory(t *testing.T) {
	RandomCategory(t)
}

func TestGetCategory(t *testing.T) {
	category1 := RandomCategory(t)

	category2, err := testQueries.GetCategory(context.Background(), category1.Category)
	require.NoError(t, err)
	require.NotEmpty(t, category2)

	require.Equal(t, category1.Category, category2.Category)
	require.NotZero(t, category2.CreatedAt)
}

func TestListCategories(t *testing.T) {
	var categories1 = [5]Category{}
	for i := 0; i < 5; i++ {
		categories1[i] = RandomCategory(t)
	}

	arg := ListCategoriesParams{
		Limit: 100,
		Offset: 0,
	}

	categories2, err := testQueries.ListCategories(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, categories2)
}

