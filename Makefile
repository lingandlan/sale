# =============================================================================
# 太积堂 - 根目录 Makefile
# =============================================================================

.PHONY: help sync-worktrees

help: ## Display help
	@echo "Targets: sync-worktrees"

# 同步核心配置文件到所有 worktree
SYNC_FILES = .claude/CLAUDE.md shop-pc/vite.config.ts
WORKTREES = .claude/worktrees/alpha .claude/worktrees/beta

sync-worktrees: ## Sync core config files to all worktrees
	@for wt in $(WORKTREES); do \
		if [ -d "$$wt" ]; then \
			echo "=== Syncing to $$wt ==="; \
			for f in $(SYNC_FILES); do \
				if [ -f "$$wt/$$f" ]; then \
					cp "$$f" "$$wt/$$f"; \
					echo "  synced $$f"; \
				fi; \
			done; \
		fi; \
	done
	@echo "Done."
