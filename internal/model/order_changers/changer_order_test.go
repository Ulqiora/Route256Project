package order_changers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"homework/internal/model"
	jtime "homework/pkg/wrapper/jsontime"
)

func TestChangerOrderBox_Change(t *testing.T) {
	mapChangers := map[model.TypePacking]ChangerOrder{
		model.TypeBox:     ChangerOrderBox{},
		model.TypeTape:    ChangerOrderTape{},
		model.TypePackage: ChangerOrderPackage{},
	}

	t.Run("ok", func(t *testing.T) {
		t.Parallel()
		obj := model.OrderInitData{
			CustomerID:  0,
			PickPointID: 0,
			ShelfLife:   jtime.TimeWrap{},
			Penny:       10000,
			Weight:      10,
			Type:        model.TypeBox,
		}
		obj, err := mapChangers[obj.Type].Change(obj)
		require.Nil(t, err)
		assert.Equal(t, obj.Penny, int64(12000))
	})
	t.Run("fail", func(t *testing.T) {
		t.Run("Incorrect weight", func(t *testing.T) {
			t.Parallel()
			obj := model.OrderInitData{
				CustomerID:  0,
				PickPointID: 0,
				ShelfLife:   jtime.TimeWrap{},
				Penny:       10000,
				Weight:      40,
				Type:        model.TypeBox,
			}
			obj, err := mapChangers[obj.Type].Change(obj)
			require.Equal(t, err, ErrorHeavyWeight)
			require.Equal(t, obj.Penny, int64(10000))
		})
		t.Run("fail type", func(t *testing.T) {
			t.Parallel()
			obj := model.OrderInitData{
				CustomerID:  0,
				PickPointID: 0,
				ShelfLife:   jtime.TimeWrap{},
				Penny:       10000,
				Weight:      40,
				Type:        "dasf",
			}
			_, ok := mapChangers[obj.Type]
			require.Equal(t, ok, false)
		})
	})
}

func TestChangerOrderTape_Change(t *testing.T) {
	mapChangers := map[model.TypePacking]ChangerOrder{
		model.TypeBox:     ChangerOrderBox{},
		model.TypeTape:    ChangerOrderTape{},
		model.TypePackage: ChangerOrderPackage{},
	}

	t.Run("ok", func(t *testing.T) {
		t.Parallel()
		obj := model.OrderInitData{
			CustomerID:  0,
			PickPointID: 0,
			ShelfLife:   jtime.TimeWrap{},
			Penny:       10000,
			Weight:      10,
			Type:        model.TypeTape,
		}
		obj, err := mapChangers[obj.Type].Change(obj)
		require.Nil(t, err)
		assert.Equal(t, obj.Penny, int64(10100))
	})
	t.Run("fail", func(t *testing.T) {
		t.Run("Incorrect weight", func(t *testing.T) {
			t.Parallel()
			obj := model.OrderInitData{
				CustomerID:  0,
				PickPointID: 0,
				ShelfLife:   jtime.TimeWrap{},
				Penny:       10000,
				Weight:      40,
				Type:        model.TypeTape,
			}
			obj, err := mapChangers[obj.Type].Change(obj)
			require.Nil(t, err)
			assert.Equal(t, obj.Penny, int64(10100))
		})
		t.Run("fail type", func(t *testing.T) {
			t.Parallel()
			obj := model.OrderInitData{
				CustomerID:  0,
				PickPointID: 0,
				ShelfLife:   jtime.TimeWrap{},
				Penny:       10000,
				Weight:      40,
				Type:        "dasf",
			}
			_, ok := mapChangers[obj.Type]
			require.Equal(t, ok, false)
		})
	})
}

func TestChangerOrderPackage_Change(t *testing.T) {
	mapChangers := map[model.TypePacking]ChangerOrder{
		model.TypeBox:     ChangerOrderBox{},
		model.TypeTape:    ChangerOrderTape{},
		model.TypePackage: ChangerOrderPackage{},
	}

	t.Run("ok", func(t *testing.T) {
		t.Parallel()
		obj := model.OrderInitData{
			CustomerID:  0,
			PickPointID: 0,
			ShelfLife:   jtime.TimeWrap{},
			Penny:       10000,
			Weight:      10,
			Type:        model.TypeTape,
		}
		obj, err := mapChangers[obj.Type].Change(obj)
		require.Nil(t, err)
		assert.Equal(t, obj.Penny, int64(10100))
	})
	t.Run("fail", func(t *testing.T) {
		t.Run("Incorrect weight", func(t *testing.T) {
			t.Parallel()
			obj := model.OrderInitData{
				CustomerID:  0,
				PickPointID: 0,
				ShelfLife:   jtime.TimeWrap{},
				Penny:       10000,
				Weight:      40,
				Type:        model.TypeTape,
			}
			obj, err := mapChangers[obj.Type].Change(obj)
			require.Nil(t, err)
			assert.Equal(t, obj.Penny, int64(10100))
		})
		t.Run("fail type", func(t *testing.T) {
			t.Parallel()
			obj := model.OrderInitData{
				CustomerID:  0,
				PickPointID: 0,
				ShelfLife:   jtime.TimeWrap{},
				Penny:       10000,
				Weight:      40,
				Type:        "dasf",
			}
			_, ok := mapChangers[obj.Type]
			require.Equal(t, ok, false)
		})
	})
}
