package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Point struct {
	x, y int
}

func input() (int, int, [][]int, Point, Point, error) {
	reader := bufio.NewReader(os.Stdin)

	var n, m int
	if _, err := fmt.Fscanf(reader, "%d %d\n", &n, &m); err != nil {
		return 0, 0, nil, Point{}, Point{}, fmt.Errorf("неверный формат размеров лабиринта")
	}
	if n <= 0 || m <= 0 {
		return 0, 0, nil, Point{}, Point{}, fmt.Errorf("размеры лабиринта должны быть положительными")
	}

	maze := make([][]int, n)
	for i := 0; i < n; i++ {
		line, err := reader.ReadString('\n')
		if err != nil {
			return 0, 0, nil, Point{}, Point{}, fmt.Errorf("ошибка чтения строки лабиринта на строке %d", i+1)
		}
		line = strings.TrimSpace(line)
		cols := strings.Fields(line)
		if len(cols) != m {
			return 0, 0, nil, Point{}, Point{}, fmt.Errorf("неверное количество элементов в строке %d: ожидалось %d, получено %d", i+1, m, len(cols))
		}
		maze[i] = make([]int, m)
		for j, col := range cols {
			maze[i][j], err = strconv.Atoi(col)
			if err != nil {
				return 0, 0, nil, Point{}, Point{}, fmt.Errorf("неверный символ '%s' в строке %d", col, i+1)
			}
			if maze[i][j] < 0 || maze[i][j] > 9 {
				return 0, 0, nil, Point{}, Point{}, fmt.Errorf("значение '%d' вне допустимого диапазона (0-9) в строке %d", maze[i][j], i+1)
			}
		}
	}

	var startX, startY int
	if _, err := fmt.Fscanf(reader, "%d %d\n", &startX, &startY); err != nil {
		return 0, 0, nil, Point{}, Point{}, fmt.Errorf("неверный формат начальной точки")
	}

	var endX, endY int
	if _, err := fmt.Fscanf(reader, "%d %d\n", &endX, &endY); err != nil {
		return 0, 0, nil, Point{}, Point{}, fmt.Errorf("неверный формат конечной точки")
	}

	if !isValidPoint(n, m, maze, startX, startY) {
		return 0, 0, nil, Point{}, Point{}, fmt.Errorf("начальная точка (%d, %d) некорректна", startX, startY)
	}
	if !isValidPoint(n, m, maze, endX, endY) {
		return 0, 0, nil, Point{}, Point{}, fmt.Errorf("конечная точка (%d, %d) некорректна", endX, endY)
	}

	return n, m, maze, Point{startX, startY}, Point{endX, endY}, nil
}

func isValidPoint(n, m int, maze [][]int, x, y int) bool {
	return x >= 0 && x < n && y >= 0 && y < m && maze[x][y] != 0
}

func bfs(n, m int, maze [][]int, start, end Point) ([]Point, error) {
	directions := []Point{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
	prev := make([][]Point, n)
	for i := range prev {
		prev[i] = make([]Point, m)
	}

	visited := make([][]bool, n)
	for i := range visited {
		visited[i] = make([]bool, m)
	}

	queue := []Point{start}
	visited[start.x][start.y] = true

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if current == end {
			return reconstructPath(prev, start, end), nil
		}

		for _, dir := range directions {
			nx, ny := current.x+dir.x, current.y+dir.y
			if isValidPoint(n, m, maze, nx, ny) && !visited[nx][ny] {
				visited[nx][ny] = true
				prev[nx][ny] = current
				queue = append(queue, Point{nx, ny})
			}
		}
	}

	return nil, fmt.Errorf("путь не найден")
}

func reconstructPath(prev [][]Point, start, end Point) []Point {
	path := []Point{}
	for current := end; current != start; current = prev[current.x][current.y] {
		path = append([]Point{current}, path...)
	}
	return append([]Point{start}, path...)
}

func main() {
	n, m, maze, start, end, err := input()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Ошибка:", err)
		os.Exit(1)
	}

	path, err := bfs(n, m, maze, start, end)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Ошибка:", err)
		os.Exit(1)
	}

	for _, p := range path {
		fmt.Printf("%d %d\n", p.x, p.y)
	}
	fmt.Println(".")
}
