package logic

import (
	"giogo/ui/pages/minesweeper/model"
	"image"
)

func RevealedCells(startingPoint image.Point, matrix [][]*model.MineElement) (revealedPos []image.Point, countOfFloodedCells uint16) {
	revealedPos = make([]image.Point, 0, 8)
	revealedPos = append(revealedPos, startingPoint)
	floodedPos := make([]image.Point, 0, 8)
	floodedPos = append(floodedPos, startingPoint)
	height := len(matrix)
	width := len(matrix[0])
	countOfFloodedCells = 1

	matrix[startingPoint.Y][startingPoint.X].PropOff(model.HiddenBits)

	for iterator := 0; iterator < len(floodedPos); iterator++ {
		// Venni a jelenlegi elem rejtett környezetét és azokat hozzáadni egy listához
		for i := -1; i <= 1; i++ {
			rowIndex := floodedPos[iterator].Y + i

			// Kilóg felül, vagy alul
			if rowIndex < 0 || rowIndex > int(height-1) {
				continue
			}

			// Column loop
			for j := -1; j <= 1; j++ {
				if j == 0 && i == 0 {
					continue
				}

				colIndex := floodedPos[iterator].X + j

				// Kilóg bal, vagy jobb oldalt
				if colIndex < 0 || colIndex > int(width-1) {
					continue
				}

				// Az adott elem értékét megvizsgálni
				element := matrix[rowIndex][colIndex]
				if !element.IsHidden() || element.IsMarked() {
					continue
				}

				// felfedni és növelni a felfedettek számát
				element.PropOff(model.HiddenBits)
				countOfFloodedCells++

				revealedPos = append(revealedPos, element.Pos)

				// Ha 0, akkor tovább vizsgálni abból a pontból a cellákat
				if element.Value == 0 {
					floodedPos = append(floodedPos, element.Pos)
				}
			}
		}
	}

	return
}
