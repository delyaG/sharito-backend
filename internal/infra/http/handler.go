package http

import (
	"backend/internal/domain"
	"backend/internal/infra/http/viewmodels"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"net/http"
	"strconv"
	"time"
)

const productCountOnPage int = 10

func (a *adapter) wrap(handler func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := handler(w, r); err != nil {
			a.logger.WithFields(generateFields(r)).WithError(err).Error("Error handling request")
		}
	}
}

func (a *adapter) sayHello(w http.ResponseWriter, r *http.Request) error {
	for k, values := range r.Header {
		fmt.Print(k, ": ")
		for _, v := range values {
			fmt.Print(v, ", ")
		}
	}

	if _, err := w.Write([]byte("Hello!")); err != nil {
		return jError(w, err)
	}

	w.WriteHeader(http.StatusOK)
	return nil
}

func (a *adapter) register(w http.ResponseWriter, r *http.Request) error {
	var req viewmodels.User
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		a.logger.WithError(err).Error("Error while decoding request body!")
		return jError(w, domain.ErrInvalidInputData)
	}

	// check password existence
	if req.Password == nil {
		a.logger.WithError(domain.ErrInvalidInputData).Error("There is no password!")
		return jError(w, domain.ErrInvalidInputData)
	}

	user := req.Domain()

	token, err := a.service.Register(user)
	if err != nil {
		return jError(w, err)
	}

	w.Header().Set("X-Auth", token)
	return j(w, http.StatusOK, struct {
		Token string `json:"token"`
	}{Token: token})
}

func (a *adapter) login(w http.ResponseWriter, r *http.Request) error {
	var req viewmodels.User
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		a.logger.WithError(err).Error("Error while decoding request body!")
		return jError(w, domain.ErrInvalidInputData)
	}

	// check password existence
	if req.Password == nil {
		a.logger.WithError(domain.ErrInvalidInputData).Error("There is no password!")
		return jError(w, domain.ErrInvalidInputData)
	}

	user := req.Domain()

	token, err := a.service.Login(user)
	if err != nil {
		return jError(w, err)
	}

	w.Header().Set("X-Auth", token)
	return j(w, http.StatusOK, struct {
		Token string `json:"token"`
	}{Token: token})
}

func (a *adapter) getUser(w http.ResponseWriter, r *http.Request) error {
	user, err := a.service.GetUser(r.Context())
	if err != nil {
		return jError(w, err)
	}

	var res viewmodels.User
	res.ViewModel(user)
	return j(w, http.StatusOK, res)
}

func (a *adapter) addProduct(w http.ResponseWriter, r *http.Request) error {
	var req viewmodels.Product
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		a.logger.WithError(err).Error("Error while decoding request body!")
		return jError(w, domain.ErrInvalidInputData)
	}

	productID, err := a.service.AddProduct(r.Context(), req.Domain())
	if err != nil {
		return jError(w, err)
	}

	return j(w, http.StatusOK, struct {
		ProductID int `json:"product_id"`
	}{ProductID: productID})
}

func (a *adapter) getProduct(w http.ResponseWriter, r *http.Request) error {
	productIDStr := chi.URLParam(r, "product_id")
	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		a.logger.WithError(err).Error("product_id is not int")
		return jError(w, domain.ErrInvalidInputData)
	}

	product, user, err := a.service.GetProductAndOwnerUserByProductID(productID)
	if err != nil {
		return jError(w, err)
	}

	var res viewmodels.ProductWithUser
	res.ViewModel(product, user)
	return j(w, http.StatusOK, res)
}

func (a *adapter) getProducts(w http.ResponseWriter, r *http.Request) error {
	var page int
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err != nil {
			a.logger.WithError(domain.ErrInvalidInputData).Error("cannot parse 'page' query param")
			return jError(w, domain.ErrInvalidInputData)
		} else {
			page = p - 1
		}
	}

	if page < 0 {
		return jError(w, domain.ErrInvalidInputData)
	}

	search := r.URL.Query().Get("search")

	products, count, err := a.service.GetProductsWithPagination(page, productCountOnPage, search)
	if err != nil {
		return jError(w, err)
	}

	var res viewmodels.ProductsWithCount
	res.ViewModel(products, count)
	return j(w, http.StatusOK, res)
}

func (a *adapter) rentProduct(w http.ResponseWriter, r *http.Request) error {
	productIDStr := chi.URLParam(r, "product_id")
	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		a.logger.WithError(err).Error("product_id is not int")
		return jError(w, domain.ErrInvalidInputData)
	}

	var from, to time.Time
	if dateStr := r.URL.Query().Get("from"); dateStr != "" {
		if date, err := time.Parse("2006-01-02 15:04", dateStr); err != nil {
			a.logger.WithError(domain.ErrInvalidInputData).Error("cannot parse 'from' query param")
			return jError(w, domain.ErrInvalidInputData)
		} else {
			from = date
		}
	}

	if dateStr := r.URL.Query().Get("to"); dateStr != "" {
		if date, err := time.Parse("2006-01-02 15:04", dateStr); err != nil {
			a.logger.WithError(domain.ErrInvalidInputData).Error("cannot parse 'to' query param")
			return jError(w, domain.ErrInvalidInputData)
		} else {
			to = date
		}
	}

	if err := a.service.RentProduct(r.Context(), productID, from, to); err != nil {
		return jError(w, err)
	}

	w.WriteHeader(http.StatusOK)
	return nil
}

func (a *adapter) getOrders(w http.ResponseWriter, r *http.Request) error {
	var isMine bool
	if isMeStr := r.URL.Query().Get("mine"); isMeStr != "" {
		if mm, err := strconv.ParseBool(isMeStr); err != nil {
			a.logger.WithError(domain.ErrInvalidInputData).Error("cannot parse 'mine' query param")
			return jError(w, domain.ErrInvalidInputData)
		} else {
			isMine = mm
		}
	}

	orders, err := a.service.GetOrders(r.Context(), isMine)
	if err != nil {
		return jError(w, err)
	}

	res := viewmodels.Orders{}
	res.ViewModel(orders)
	return j(w, http.StatusOK, res)
}
