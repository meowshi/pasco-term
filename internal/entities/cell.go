package entities

import (
	"fmt"
	"strconv"
)

type Cell struct {
	Column string
	Row    int
}

func NewCell(column string, row int) Cell {
	return Cell{
		Column: column,
		Row:    row,
	}
}

func NewCellFromString(cell string) Cell {
	row, _ := strconv.Atoi(cell[1:])
	return NewCell(string(cell[0]), row)
}

func (c *Cell) AddColumn(n byte) *Cell {
	if c.Column[0] <= 'Z' && c.Column[0]+n >= 'a' || c.Column[0]+n > 'z' {
		panic("таблица оказалась очень большая я не предусмотрел что колонки заняли столько места. попробуйте перенести таблицу ниже, так чтобы она не заходила за колонку Z")
	}
	c.Column = string(c.Column[0] + n)
	return c
}

func (c *Cell) AddRow(n int) *Cell {
	c.Row += n

	return c
}

func (c *Cell) ToString() string {
	return fmt.Sprintf("%s%d", c.Column, c.Row)
}
