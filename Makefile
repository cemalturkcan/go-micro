PRODUCT_CATALOG=product-catalog
PRODUCT_CATEGORY=product-category

PRODUCT_CATALOG_PATH=./services/product-catalog
PRODUCT_CATEGORY_PATH=./services/product-category

AIR_CMD=air

run-product-catalog:
	@echo "Running Service 1"
	cd $(PRODUCT_CATALOG_PATH) && $(AIR_CMD)

run-product-category:
	@echo "Running Service 2"
	cd $(PRODUCT_CATEGORY_PATH) && $(AIR_CMD)


watch:
	@echo "Running services..."
	(cd $(PRODUCT_CATALOG_PATH) && $(AIR_CMD)) & \
	(cd $(PRODUCT_CATEGORY_PATH) && $(AIR_CMD)) &
	wait
	@echo "Services are running..."
