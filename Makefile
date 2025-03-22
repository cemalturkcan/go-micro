PRODUCT_CATALOG_PATH=./services/product-catalog
PRODUCT_CATEGORY_PATH=./services/product-category

AUTHENTICATION_PATH=./services/authentication
API_GATEWAY_PATH=./services/api-gateway


AIR_CMD=air
TIDY_CMD=go mod tidy

run-product-catalog:
	@echo "Running Service 1"
	cd $(PRODUCT_CATALOG_PATH) && $(AIR_CMD)

run-product-category:
	@echo "Running Service 2"
	cd $(PRODUCT_CATEGORY_PATH) && $(AIR_CMD)


run-authentication:
	@echo "Running Service 3"
	cd $(AUTHENTICATION_PATH) && $(AIR_CMD)

run-api-gateway:
	@echo "Running Service 4"
	cd $(API_GATEWAY_PATH) && $(AIR_CMD)


tidy:
	@echo "Running go mod tidy for all services..."
	cd $(PRODUCT_CATALOG_PATH) && $(TIDY_CMD)
	cd $(PRODUCT_CATEGORY_PATH) && $(TIDY_CMD)
	cd $(AUTHENTICATION_PATH) && $(TIDY_CMD)
	cd $(API_GATEWAY_PATH) && $(TIDY_CMD)
	@echo "go mod tidy completed for all services."


watch:
	@echo "Running services..."
	(cd $(PRODUCT_CATALOG_PATH) && $(AIR_CMD)) & \
	(cd $(PRODUCT_CATEGORY_PATH) && $(AIR_CMD)) &
	(cd $(AUTHENTICATION_PATH) && $(AIR_CMD)) & \
	(cd $(API_GATEWAY_PATH) && $(AIR_CMD)) &
	wait
	@echo "Services are running..."


#inside common directory run make compile-proto
compile-proto:
	cd common/grpc && make compile-protos
