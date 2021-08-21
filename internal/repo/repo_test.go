package repo

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	sq "github.com/Masterminds/squirrel"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/ozoncp/ocp-request-api/internal/models"
)

var _ = Describe("Repo", func() {

	var (
		rep      Repo
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

	Context("Adding items with no errors. Will not return any remains.", func() {
		JustBeforeEach(func() {
			var err error
			db, dbMock, err = sqlmock.New()
			Expect(err).ToNot(HaveOccurred())
			stmtCache := sq.NewStmtCache(db)

			rep = &repo{
				stmBuilder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar).RunWith(stmtCache),
			}

		})

		It("Add single request. Expect new ID generated.", func() {
			newReq := models.Request{
				Id:     0,
				UserId: 1,
				Type:   2,
				Text:   "one",
			}
			expectedNewId := uint64(1)
			returnRows := sqlmock.NewRows([]string{"id"}).AddRow(expectedNewId)
			dbMock.ExpectPrepare(
				"INSERT INTO requests \\(user_id,type,text\\) VALUES \\(\\$1,\\$2,\\$3\\) RETURNING id",
			).
				ExpectQuery().
				WithArgs(newReq.UserId, newReq.Type, newReq.Text).
				WillReturnRows(returnRows)

			newId, err := rep.Add(ctx, newReq)
			Expect(err).ToNot(HaveOccurred())

			Expect(newId).To(Equal(expectedNewId))
		})

		It("Add many requests into repository", func() {
			requests := []models.Request{
				{
					UserId: 10,
					Type:   100,
					Text:   "one",
				},
				{
					UserId: 20,
					Type:   200,
					Text:   "two",
				},
				{
					UserId: 30,
					Type:   300,
					Text:   "three",
				},
			}
			expectedQueryArgs := make([]driver.Value, 0, len(requests)*3)

			for _, req := range requests {
				expectedQueryArgs = append(expectedQueryArgs, req.UserId, req.Type, req.Text)
			}
			expctedNewIds := []uint64{1, 2, 3}

			dbMock.ExpectPrepare(
				"INSERT INTO requests \\(user_id,type,text\\) VALUES \\(\\$1,\\$2,\\$3\\),\\(\\$4,\\$5,\\$6\\),\\(\\$7,\\$8,\\$9\\) RETURNING id",
			).
				ExpectQuery().
				WithArgs(expectedQueryArgs...).
				WillReturnRows(sqlmock.NewRows([]string{"id"}).
					AddRow(1).
					AddRow(2).
					AddRow(3))

			newIds, err := rep.AddMany(ctx, requests)
			Expect(err).ToNot(HaveOccurred())
			Expect(newIds).To(Equal(expctedNewIds))

		})

		It("Fetch requests from database", func() {
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
				"SELECT id, user_id, type, text FROM requests LIMIT 1000 OFFSET 100",
			).
				ExpectQuery().
				WillReturnRows(returnRows)
			actualRequests, err := rep.List(ctx, limit, offset)
			Expect(err).ToNot(HaveOccurred())

			Expect(actualRequests).To(Equal(expectedRequests))
		})

		It("Remove request that is exists", func() {
			reqId := uint64(100)
			res := sqlmock.NewResult(0, 1)

			dbMock.ExpectPrepare(
				"DELETE FROM requests WHERE id = \\$1",
			).
				ExpectExec().
				WithArgs(reqId).
				WillReturnResult(res)

			err := rep.Remove(ctx, reqId)
			Expect(err).ToNot(HaveOccurred())
		})

		It("Remove request that is not exists", func() {
			reqId := uint64(100)
			res := sqlmock.NewResult(0, 0)

			dbMock.ExpectPrepare(
				"DELETE FROM requests WHERE id = \\$1",
			).
				ExpectExec().
				WithArgs(reqId).
				WillReturnResult(res)

			err := rep.Remove(ctx, reqId)
			Expect(err).To(Equal(NotFound))
		})

		It("Return single request that is exists", func() {
			reqId := uint64(1)
			expectedReq := models.NewRequest(reqId, 10, 100, "one")

			returnRows := sqlmock.
				NewRows([]string{"id", "user_id", "type", "text"}).
				AddRow(expectedReq.Id, expectedReq.UserId, expectedReq.Type, expectedReq.Text)

			dbMock.ExpectPrepare(
				"SELECT id, user_id, type, text FROM requests WHERE id = \\$1",
			).
				ExpectQuery().
				WithArgs(reqId).
				WillReturnRows(returnRows)

			actualReq, err := rep.Describe(ctx, reqId)
			Expect(err).ToNot(HaveOccurred())
			Expect(actualReq).To(Equal(&expectedReq))
		})

		It("Return single request that is not exists", func() {
			reqId := uint64(1)

			returnRows := sqlmock.
				NewRows([]string{"id", "user_id", "type", "text"})

			dbMock.ExpectPrepare(
				"SELECT id, user_id, type, text FROM requests WHERE id = \\$1",
			).
				ExpectQuery().
				WithArgs(reqId).
				WillReturnRows(returnRows)

			actualReq, err := rep.Describe(ctx, reqId)
			Expect(err).To(Equal(NotFound))
			Expect(actualReq).To(BeNil())
		})

		It("Pops up error on List()", func() {

			offset, limit := uint64(100), uint64(1000)
			expectedError := errors.New("test")
			dbMock.ExpectPrepare(
				"SELECT id, user_id, type, text FROM requests LIMIT 1000 OFFSET 100",
			).
				ExpectQuery().
				WillReturnError(expectedError)
			_, err := rep.List(ctx, limit, offset)
			Expect(err).To(Equal(expectedError))
		})

		It("Pops up error on Add()", func() {

			newReq := models.Request{
				Id:     0,
				UserId: 1,
				Type:   2,
				Text:   "one",
			}
			expectedError := errors.New("test")
			dbMock.ExpectPrepare(
				"INSERT INTO requests \\(user_id,type,text\\) VALUES \\(\\$1,\\$2,\\$3\\) RETURNING id",
			).
				ExpectQuery().
				WithArgs(newReq.UserId, newReq.Type, newReq.Text).
				WillReturnError(expectedError)

			newId, err := rep.Add(ctx, newReq)
			Expect(err).To(Equal(expectedError))
			Expect(newId).To(Equal(uint64(0)))

		})

		It("Pops up error on AddMany()", func() {

			newReq := models.Request{
				Id:     0,
				UserId: 1,
				Type:   2,
				Text:   "one",
			}
			expectedError := errors.New("test")
			dbMock.ExpectPrepare(
				"INSERT INTO requests \\(user_id,type,text\\) VALUES \\(\\$1,\\$2,\\$3\\) RETURNING id",
			).
				ExpectQuery().
				WithArgs(newReq.UserId, newReq.Type, newReq.Text).
				WillReturnError(expectedError)

			_, err := rep.AddMany(ctx, []models.Request{newReq})
			Expect(err).To(Equal(expectedError))

		})

		It("Pops up error on Remove()", func() {

			reqId := uint64(100)
			expectedError := errors.New("test")
			dbMock.ExpectPrepare(
				"DELETE FROM requests WHERE id = \\$1",
			).
				ExpectExec().
				WithArgs(reqId).
				WillReturnError(expectedError)

			err := rep.Remove(ctx, reqId)
			Expect(err).To(Equal(expectedError))
		})

		It("Update request that is exists", func() {
			req := models.NewRequest(1, 10, 100, "one")
			res := sqlmock.NewResult(0, 1)

			dbMock.ExpectPrepare(
				"UPDATE requests SET user_id = \\$1, type = \\$2, text = \\$3 WHERE id = \\$4",
			).
				ExpectExec().
				WithArgs(req.UserId, req.Type, req.Text, req.Id).
				WillReturnResult(res)

			err := rep.Update(ctx, req)
			Expect(err).ToNot(Equal(NotFound))
		})

		It("Update request that is not exists", func() {
			req := models.NewRequest(1, 10, 100, "one")
			res := sqlmock.NewResult(0, 0)

			dbMock.ExpectPrepare(
				"UPDATE requests SET user_id = \\$1, type = \\$2, text = \\$3 WHERE id = \\$4",
			).
				ExpectExec().
				WithArgs(req.UserId, req.Type, req.Text, req.Id).
				WillReturnResult(res)

			err := rep.Update(ctx, req)
			Expect(err).To(Equal(NotFound))
		})

	})

})
