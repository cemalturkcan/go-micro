SERVICE_1=product-catalog

SERVICE_1_PATH=./services/product-catalog

AIR_CMD=air

run-service-1:
	@echo "Running Service 1"
	cd $(SERVICE_1_PATH) && $(AIR_CMD)


watch:
	@echo "Starting both services concurrently..."
	(cd $(SERVICE_1_PATH) && $(AIR_CMD)) &
	wait
	@echo "All services are running"