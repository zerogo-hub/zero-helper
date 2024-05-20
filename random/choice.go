package random

import (
	"errors"
	"math/rand"
)

type integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

// Choice 每一个选项
type Choice[T any, W integer] struct {
	Item   T
	Weight W
}

// Chooser 选择器
type Chooser[T any, W integer] struct {
	choices []Choice[T, W]
	totals  []int64
	max     int64
}

type WeightHandler[T any, W integer] func(item T) W

func NewChooser[T any, W integer](items []T, getWeight WeightHandler[T, W]) (*Chooser[T, W], error) {
	if len(items) == 0 {
		return nil, errors.New("items can not be empty")
	}

	choices := make([]Choice[T, W], 0, len(items))
	totals := make([]int64, 0, len(items))
	max := int64(0)

	for _, item := range items {
		weight := getWeight(item)
		max += int64(weight)

		choice := Choice[T, W]{
			Item:   item,
			Weight: weight,
		}
		choices = append(choices, choice)
		totals = append(totals, max)
	}

	return &Chooser[T, W]{choices: choices, totals: totals, max: max}, nil
}

func (c *Chooser[T, W]) Pick() T {
	r := rand.Int63n(c.max) + 1

	// 实际应用中，数量较少，直接遍历
	i := 0
	for j, v := range c.totals {
		if r <= v {
			i = j
			break
		}
	}
	return c.choices[i].Item
}

func PickOnce[T any, W integer](items []T, getWeight WeightHandler[T, W]) T {
	c, _ := NewChooser(items, getWeight)
	return c.Pick()
}
