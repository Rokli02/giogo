package logic

import (
	"image"
	"math/rand"

	"giogo/ui/pages/minesweeper/model"
)

func GenerateMines(clickPos image.Point, matrix [][]*model.MineElement, maxMines uint16) uint16 {
	minePositions := make([]image.Point, 0, maxMines)
	height := len(matrix)
	width := len(matrix[0])

	// Calculate them

	for i, calcTries := 0, 0; i < cap(minePositions) || calcTries > 5; calcTries++ {
		validPos := true

		// Random num
		minePos := image.Point{int(rand.Int31n(int32(width))), int(rand.Int31n(int32(height)))}

		// Check if is it stored already or at 'clickPos'
		if minePos == clickPos {
			continue
		}

		for i := range minePositions {
			if minePositions[i] == minePos {
				validPos = false
				break
			}
		}

		if validPos {
			minePositions = append(minePositions, minePos)
			i++
			calcTries = 0
		}
	}

	// Clear minefield
	for rowIndex := range matrix {
		for colIndex := range matrix[rowIndex] {
			matrix[rowIndex][colIndex].Value = 0
		}
	}

	// Plant mines
	for i := range minePositions {
		matrix[minePositions[i].Y][minePositions[i].X].Value = -1
	}
	mines := uint16(len(minePositions))

	// Find mines in the neighborhood, MAAAN!
	for rowIndex := range matrix {
		for colIndex := range matrix[rowIndex] {
			element := matrix[rowIndex][colIndex]

			if element.Value == -1 {
				continue
			}

			element.Value = neighboringMines(rowIndex, colIndex, matrix)
		}
	}

	return mines
}

func neighboringMines(rowIndexParam, colIndexParam int, matrix [][]*model.MineElement) int8 {
	height := len(matrix)
	width := len(matrix[0])
	var sum int8 = 0
	// Row loop
	for i := -1; i <= 1; i++ {
		rowIndex := rowIndexParam + i
		if rowIndex < 0 || rowIndex > int(height-1) {
			continue
		}

		// Column loop
		for j := -1; j <= 1; j++ {
			if j == 0 && i == 0 {
				continue
			}

			colIndex := colIndexParam + j
			if colIndex < 0 || colIndex > int(width-1) {
				continue
			}

			if matrix[rowIndex][colIndex].Value == -1 {
				sum++
			}
		}
	}

	return sum
}
