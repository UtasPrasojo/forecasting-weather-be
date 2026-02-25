run: build
	@echo "Running..."
	@./$(BINARY)

# Live reload using Air
watch:
	@if command -v air > /dev/null; then \
		air; \
	else \
		read -p "Air not installed. Install now? [Y/n] " choice; \
		if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
			go install github.com/cosmtrek/air@latest; \
			air; \
		else \
			echo "Air not installed. Exiting..."; \
			exit 1; \
		fi; \
	fi
