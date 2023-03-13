package pets

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"

	"kata-peya/internal/http/handler"
	"kata-peya/internal/http/response"
	petRepository "kata-peya/internal/pet/repository"
	petUseCase "kata-peya/internal/pet/usecase"
	"kata-peya/internal/test/integration"

	sq "github.com/Masterminds/squirrel"
	"github.com/cucumber/godog"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"github.com/testcontainers/testcontainers-go"
)

type feature struct {
	container testcontainers.Container
	db        *sql.DB
	server    *echo.Echo
	resp      *httptest.ResponseRecorder
	handler   *handler.Pets
}

func (f *feature) StartSuite() {
	fmt.Println("INITIALIZE PETS SUITE")
	f.db, f.container = integration.New()
}

func (f *feature) Reset() {
	fmt.Println("INITIALIZE PETS SCENARIO")

	if _, err := f.db.Exec("delete from pets;"); err != nil {
		log.Fatal(err)
		return
	}

	f.resp = httptest.NewRecorder()

	f.server = echo.New()

	repo := petRepository.NewPetMysqlRepository(f.db)
	uc := petUseCase.NewUseCase(repo)
	f.handler = handler.NewPetsHandler(uc)

}

func (f *feature) iSendRequestTo(method, endpoint string) error {

	req, err := http.NewRequest(method, endpoint, nil)
	if err != nil {
		return err
	}

	c := f.server.NewContext(req, f.resp)
	if err := f.handler.GetAll(c); err != nil {
		return err
	}

	return nil
}

func (f *feature) theAreRegisteredPets(t *godog.Table) error {

	for i := 1; i < len(t.Rows); i++ {
		var vaccines *string
		row := t.Rows[i]
		id, _ := strconv.Atoi(row.Cells[0].Value)
		age, _ := strconv.Atoi(row.Cells[3].Value)

		if row.Cells[2].Value != "" {
			vaccines = &row.Cells[2].Value
		}

		_, err := sq.StatementBuilder.
			RunWith(f.db).
			PlaceholderFormat(sq.Question).
			Insert("pets").
			Columns("id", "name", "vaccines", "age_months").
			Values(id, row.Cells[1].Value, vaccines, age).
			ExecContext(context.Background())
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Info("pets registered in database")
	return nil
}

func (f *feature) theResponseShouldMatchJson(jsonData *godog.DocString) error {

	var (
		expected []response.Pet
		actual   []response.Pet
	)

	if err := json.Unmarshal([]byte(jsonData.Content), &expected); err != nil {
		return err
	}

	actualJson := f.resp.Body.Bytes()
	if err := json.Unmarshal(actualJson, &actual); err != nil {
		return err
	}

	if !reflect.DeepEqual(expected, actual) {
		return fmt.Errorf("the response json not match, expected %s and got %s", jsonData.Content, string(actualJson))
	}

	return nil
}

func (f *feature) theResponseStatusCodeShouldBe(code int) error {
	if f.resp.Code != code {
		return fmt.Errorf("the response status code expected is %d and get %d with response %s", code, f.resp.Code, f.resp.Body.String())
	}
	return nil
}

func (f *feature) Teardown() {
	if err := f.container.Terminate(context.Background()); err != nil {
		log.Error(err)
	}

	fmt.Println("END PETS SUITE")
}
