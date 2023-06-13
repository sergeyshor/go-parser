package app

import (
	"context"
	"database/sql"
	"go-parser/config"
	"go-parser/pkg/pghelper"
	"log"
	"strings"
	"sync"

	_ "github.com/lib/pq"
	"golang.org/x/net/html"

	"go-parser/internal/dto"
	"go-parser/internal/infrastructure/repository/postgres"
	"go-parser/internal/infrastructure/utils"
)

func Run(cfg *config.Config) {
	// Repository
	connURL := pghelper.GetConnURL(cfg.PG)

	conn, err := sql.Open("postgres", connURL)
	if err != nil {
		log.Println(err)
	}
	defer conn.Close()

	repo := postgres.NewBookRepo(conn)

	// Parsing
	parseUrl := cfg.URL

	pages := utils.GetPagesLinks(parseUrl)

	queue := make(chan []*dto.Book, 1)

	var wg sync.WaitGroup
	wg.Add(len(pages))
	for _, pageUrl := range pages {
		go func(pageUrl string) {
			page, err := utils.GetHtmlPage(pageUrl)
			if err != nil {
				log.Println(err)
			}

			doc, err := html.Parse(strings.NewReader(string(page)))
			if err != nil {
				log.Println(err)
			}

			productDiv := utils.FindNode(doc, "tg-productgrid")

			products := utils.GetNodesByClass(productDiv, "tg-postbook")

			books := utils.ParseBookCards(products)

			queue <- books
		}(pageUrl)
	}

	go func() {
		for books := range queue {
			err = repo.InsertBooks(context.Background(), books)
			if err != nil {
				log.Println(err)
			}
			wg.Done()
		}
	}()

	wg.Wait()
}
