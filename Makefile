regenerate:
	@bash ./scripts/fetch-latest-registry.sh && \
	echo 'Regenerating...\n' && \
	go run ./cmd/registry-gen || \
	echo "Skipped code generation (no changes)."
