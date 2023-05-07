package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/zaytal/eshop_go/app/models"

	"github.com/unrolled/render"
)

func (server *Server) Products(w http.ResponseWriter, r *http.Request) {
	renderer := render.New(render.Options{
		Layout: "layout",
	})

	q := r.URL.Query()
	page, _ := strconv.Atoi(q.Get("page"))
	if page <= 0 {
		page = 1
	}

	perPage := 9

	productModel := models.Product{}
	products, totalRows, err := productModel.GetProducts(server.DB, perPage, page)
	if err != nil {
		return
	}

	pagination, _ := GetPaginationLinks(server.AppConfig, PaginationParams{
		Path:        "products",
		TotalRows:   int32(totalRows),
		PerPage:     int32(perPage),
		CurrentPage: int32(page),
	})

	_ = renderer.HTML(w, http.StatusOK, "products", map[string]interface{}{
		"products":   products,
		"pagination": pagination,
	})
}

func (server *Server) GetProductBySlug(w http.ResponseWriter, r *http.Request) {
	renderer := render.New(render.Options{
		Layout: "layout",
	})

	vars := mux.Vars(r)

	if vars["slug"] == "" {
		return
	}

	productModel := models.Product{}
	product, err := productModel.FindBySlug(server.DB, vars["slug"])
	if err != nil {
		return
	}

	_ = renderer.HTML(w, http.StatusOK, "product", map[string]interface{}{
		"product": product,
	})
}

func (server *Server) ApiGetProducts(w http.ResponseWriter, r *http.Request) {
	var productsResponse struct {
		Products   *[]models.Product
		Pagination PaginationLinks
	}

	q := r.URL.Query()
	page, _ := strconv.Atoi(q.Get("page"))
	if page <= 0 {
		page = 1
	}

	perPage := 9

	productModel := models.Product{}
	products, totalRows, err := productModel.GetProducts(server.DB, perPage, page)
	if err != nil {
		return
	}

	productsResponse.Pagination, err = GetPaginationLinks(server.AppConfig, PaginationParams{
		Path:        "products",
		TotalRows:   int32(totalRows),
		PerPage:     int32(perPage),
		CurrentPage: int32(page),
	})
	if err != nil {
		//TODO log
		return
	}
	productsResponse.Products = products

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(productsResponse)
	if err != nil {
		//TODO log
		return
	}
}
