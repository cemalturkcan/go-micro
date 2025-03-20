PRODUCT_CATALOG=product-catalog
PRODUCT_CATEGORY=product-category

PRODUCT_CATALOG_PATH=./services/product-catalog
PRODUCT_CATEGORY_PATH=./services/product-category

AIR_CMD=air
TIDY_CMD=go mod tidy

run-product-catalog:
	@echo "Running Service 1"
	cd $(PRODUCT_CATALOG_PATH) && $(AIR_CMD)

run-product-category:
	@echo "Running Service 2"
	cd $(PRODUCT_CATEGORY_PATH) && $(AIR_CMD)

tidy:
	@echo "Running go mod tidy for all services..."
	cd $(PRODUCT_CATALOG_PATH) && $(TIDY_CMD)
	cd $(PRODUCT_CATEGORY_PATH) && $(TIDY_CMD)
	@echo "go mod tidy completed for all services."


watch:
	@echo "Running services..."
	(cd $(PRODUCT_CATALOG_PATH) && $(AIR_CMD)) & \
	(cd $(PRODUCT_CATEGORY_PATH) && $(AIR_CMD)) &
	wait
	@echo "Services are running..."
