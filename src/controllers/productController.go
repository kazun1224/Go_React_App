package controllers

import (
	"admin/src/database"
	"admin/src/models"
	"context"
	"encoding/json"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

var (
	products_frontend = "products_frontend"
	products_backend  = "products_backend"
)

func Products(ctx *fiber.Ctx) error {
	var products []models.Product

	// 全てのプロダクトを取得
	database.DB.Find(&products)

	return ctx.JSON(products)
}

func CreateProducts(ctx *fiber.Ctx) error {
	var product models.Product

	// リクエストデータをパースする
	if err := ctx.BodyParser(&product); err != nil {
		return err
	}

	// プロダクト取得
	database.DB.Create(&product)

	// キャッシュ削除(goroutineで呼び出す)
	go database.ClearCache(products_frontend, products_backend)

	return ctx.JSON(product)
}

func GetProduct(ctx *fiber.Ctx) error {
	// リクエストからIDを取得
	id, _ := strconv.Atoi(ctx.Params("id"))

	var product models.Product
	product.ID = uint(id)

	// プロダクト検索
	database.DB.Find(&product)

	return ctx.JSON(product)
}

func UpdateProduct(ctx *fiber.Ctx) error {
	// リクエストからIDを取得
	id, _ := strconv.Atoi(ctx.Params("id"))

	product := models.Product{}
	product.ID = uint(id)

	if err := ctx.BodyParser(&product); err != nil {
		return err
	}

	// プロダクト更新
	database.DB.Model(&product).Updates(&product)

	// キャッシュ削除(goroutineで呼び出す)
	go database.ClearCache(products_frontend, products_backend)

	return ctx.JSON(product)
}

func DeleteProduct(ctx *fiber.Ctx) error {
	// リクエストからIDを取得
	id, _ := strconv.Atoi(ctx.Params("id"))

	product := models.Product{}
	product.ID = uint(id)

	// プロダクト削除
	database.DB.Delete(&product)

	// キャッシュ削除(goroutineで呼び出す)
	go database.ClearCache(products_frontend, products_backend)

	return nil
}

func ProductFrontend(ctx *fiber.Ctx) error {
	var products []models.Product
	var c = context.Background()
	expiredTime := 30 * time.Minute

	// products_frontend keyでRedisからデータを取得
	result, err := database.Cache.Get(c, products_frontend).Result()
	if err != nil {
		database.DB.Find(&products)

		// Redisにデータを格納するため、エンコードする
		productBytes, err := json.Marshal(&products)
		if err != nil {
			panic(err)
		}

		// products_fronend keyでRedisにデータを格納
		err = database.Cache.Set(c, products_frontend, productBytes, expiredTime).Err()
		if err != nil {
			panic(err)
		}
	} else {
		// デコードする
		json.Unmarshal([]byte(result), &products)
	}
	return ctx.JSON(products)
}

func ProductBackend(ctx *fiber.Ctx) error {
	var products []models.Product
	var c = context.Background()
	expiredTime := 30 * time.Minute

	// キャッシュ操作
	result, err := database.Cache.Get(c, products_backend).Result()
	if err != nil {
		database.DB.Find(&products)

		productBytes, err := json.Marshal(&products)
		if err != nil {
			panic(err)
		}

		database.Cache.Set(c, products_backend, productBytes, expiredTime).Err()

	} else {
		json.Unmarshal([]byte(result), &products)
	}

	var searchProducts []models.Product

	// 検索
	// urlの?q=XXXXから文字列を取得
	if q := ctx.Query("q"); q != "" {
		// 大文字小文字の区別をなくすため、全て小文字扱いにする
		lower := strings.ToLower(q)
		for _, product := range products {
			// 検索条件1: Title
			if strings.Contains(strings.ToLower(product.Title), lower) ||
				// 検索条件2: Description
				strings.Contains(strings.ToLower(product.Description), lower) {
				searchProducts = append(searchProducts, product)
			}
		}
	} else {
		// 検索しない場合は、全てのデータを返却
		searchProducts = products
	}

	// ソート
	if sortParam := ctx.Query("sort"); sortParam != "" {
		sortLower := strings.ToLower(sortParam)
		if sortLower == "asc" {
			sort.Slice(searchProducts, func(i, j int) bool {
				return searchProducts[i].Price < searchProducts[j].Price
			})
		} else if sortLower == "desc" {
			sort.Slice(searchProducts, func(i, j int) bool {
				return searchProducts[i].Price > searchProducts[j].Price
			})
		}
	}

	// ページネーション
	var total = len(searchProducts)
	// デフォルトは"1"ページ
	page, _ := strconv.Atoi(ctx.Query("page", "1"))
	// 1ページ最大9個の商品
	perPage := 9

	var data []models.Product

	if total <= page*perPage && total >= (page-1)*perPage {
		data = searchProducts[(page-1)*perPage : total]
	} else if total >= page*perPage {
		data = searchProducts[(page-1)*perPage : page*perPage]
	} else {
		data = []models.Product{}
	}

	// 1ページ目 -> 0 ~ 8
	// 2パージ目 -> 9 ~ 17
	return ctx.JSON(fiber.Map{
		"data":      data,
		"total":     total,
		"page":      page,
		"last_page": total/perPage + 1,
	})
}
