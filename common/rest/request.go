package rest

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"net/url"
	"strconv"
	"strings"
)

const MaxLimit = 100

type IdRequest struct {
	Id int64 `json:"id" validate:"required"`
}

func IdRequestListToIdList(list []IdRequest) []int64 {
	var idList []int64
	for _, id := range list {
		idList = append(idList, id.Id)
	}
	return idList
}

type DeleteRequest struct {
	IdList []int64 `json:"idList" validate:"required"`
}

func GetPageParams(c *fiber.Ctx) (int, int) {
	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil {
		limit = 10
	}

	if limit > MaxLimit {
		limit = MaxLimit
	}
	offset, err := strconv.Atoi(c.Query("offset", "0"))
	if err != nil {
		offset = 0
	}
	return limit, offset
}

func GetPageParamsWithSort(c *fiber.Ctx, validColumns *map[string]bool) (int, int, string) {
	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil {
		limit = 10
	}
	if limit > MaxLimit {
		limit = MaxLimit
	}
	offset, err := strconv.Atoi(c.Query("offset", "0"))
	if err != nil {
		offset = 0
	}
	order := checkAndGetOrder(c.Query("sort", "id,desc"), validColumns)
	return limit, offset, order
}

func checkAndGetOrder(sort string, validColumns *map[string]bool) string {
	value, err := url.QueryUnescape(sort)
	if err != nil {
		return ""
	}
	orders := strings.Split(value, ":")
	var validOrders []string
	for _, o := range orders {
		column := strings.Split(o, ",")[0]
		if _, ok := (*validColumns)[column]; ok {
			o = strings.Replace(o, ",", " ", 1)
			validOrders = append(validOrders, o)
		}
	}
	return strings.Join(validOrders, ",")
}

var validate = validator.New(validator.WithRequiredStructEnabled())

func validateRequest(req interface{}) error {
	if err := validate.Struct(req); err != nil {
		return err
	}
	return nil
}

func SetBodyAndValidate(c *fiber.Ctx, body interface{}) error {
	err := c.BodyParser(body)
	if err != nil {
		return err
	}
	err = validateRequest(body)
	if err != nil {
		return err
	}
	return nil
}
