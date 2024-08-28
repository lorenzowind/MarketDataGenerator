package main

import (
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

type Ordem struct {
	ID         string
	Preco      float64
	Quantidade int
	Tipo       string
	Horario    time.Time
}

type LivroDeOfertas struct {
	Ordens []Ordem
}

func gerarOrdem(precoBase, spread float64, tipo string, horario time.Time) Ordem {
	preco := precoBase
	if tipo == "compra" {
		preco -= rand.Float64() * spread
	} else {
		preco += rand.Float64() * spread
	}
	quantidade := rand.Intn(100) + 1
	return Ordem{
		ID:         fmt.Sprintf("%d", rand.Intn(100000)),
		Preco:      preco,
		Quantidade: quantidade,
		Tipo:       tipo,
		Horario:    horario,
	}
}

func simularPregao(precoBase, spread float64, intervalos int) LivroDeOfertas {
	inicioPregao := time.Date(2024, 8, 22, 9, 0, 0, 0, time.Local)
	fimPregao := time.Date(2024, 8, 22, 17, 0, 0, 0, time.Local)
	duracaoPregao := fimPregao.Sub(inicioPregao)

	intervaloTempo := duracaoPregao / time.Duration(intervalos)
	livro := LivroDeOfertas{}

	for i := 0; i < intervalos; i++ {
		horario := inicioPregao.Add(time.Duration(i) * intervaloTempo)
		livro.Ordens = append(livro.Ordens, gerarOrdem(precoBase, spread, "compra", horario))
		livro.Ordens = append(livro.Ordens, gerarOrdem(precoBase, spread, "venda", horario))
	}

	return livro
}

func salvarCSV(livro LivroDeOfertas, nomeArquivo string) error {
	file, err := os.Create(nomeArquivo)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"ID", "Tipo", "Preco", "Quantidade", "Horario"})

	for _, ordem := range livro.Ordens {
		writer.Write([]string{
			ordem.ID,
			ordem.Tipo,
			strconv.FormatFloat(ordem.Preco, 'f', 2, 64),
			strconv.Itoa(ordem.Quantidade),
			ordem.Horario.Format("2006-01-02 15:04:05"),
		})
	}

	return nil
}

func main() {
	precoBase := 100.0
	spread := 1.0
	intervalos := 480

	livro := simularPregao(precoBase, spread, intervalos)
	if err := salvarCSV(livro, "livro_de_ofertas.csv"); err != nil {
		fmt.Println("Erro ao salvar CSV:", err)
	}
}
