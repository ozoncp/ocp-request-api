package search

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"github.com/DATA-DOG/go-sqlmock"
	sq "github.com/Masterminds/squirrel"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/ozoncp/ocp-request-api/internal/models"
)

var _ = Describe("Search", func() {

	var (
		search   Searcher
		dbMock   sqlmock.Sqlmock
		mockCtrl *gomock.Controller
		ctx      context.Context
		db       *sql.DB
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		ctx = context.Background()
	})

	AfterEach(func() {
		mockCtrl.Finish()
		defer db.Close()
		if err := dbMock.ExpectationsWereMet(); err != nil {
			Expect(err).ToNot(HaveOccurred())
		}

	})

	Context("Test search", func() {
		JustBeforeEach(func() {
			var err error
			db, dbMock, err = sqlmock.New()
			Expect(err).ToNot(HaveOccurred())
			stmtCache := sq.NewStmtCache(db)

			search = &searcher{
				stmBuilder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar).RunWith(stmtCache),
			}

		})

		It("Simple full text search", func() {
			dbRows := [][]driver.Value{
				{uint64(1), uint64(10), uint64(100), "one"},
				{uint64(2), uint64(20), uint64(200), "two"},
				{uint64(3), uint64(30), uint64(300), "three"},
			}
			expectedRequests := make([]models.Request, 0, len(dbRows))
			returnRows := sqlmock.NewRows([]string{"id", "user_id", "type", "text"})

			for _, row := range dbRows {
				expectedRequests = append(expectedRequests, models.Request{
					Id:     row[0].(uint64),
					UserId: row[1].(uint64),
					Type:   row[2].(uint64),
					Text:   row[3].(string),
				})
				returnRows.AddRow(row...)
			}

			offset, limit := uint64(100), uint64(1000)

			dbMock.ExpectPrepare(
				"SELECT id, user_id, type, text " +
					"FROM requests " +
					"WHERE to_tsvector\\(\\'russian\\', text\\) @@ to_tsquery\\(\\$1\\) " +
					"ORDER BY ts_rank\\(to_tsvector\\(\\'russian\\', text\\), to_tsquery\\(\\$2\\)\\) desc " +
					"LIMIT 1000 " +
					"OFFSET 100",
			).
				ExpectQuery().
				WillReturnRows(returnRows)
			actualRequests, err := search.Search(ctx, "hey", limit, offset)
			Expect(err).ToNot(HaveOccurred())

			Expect(actualRequests).To(Equal(expectedRequests))
		})

	})

})
