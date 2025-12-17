.PHONY: up down restart logs ps \
	git-commit git-push git-feature git-hotfix git-release git-merge-main

# ==================== Docker Commands ====================
up:
	docker-compose up -d

down:
	docker-compose down

restart: down up

logs:
	docker-compose logs -f

ps:
	docker-compose ps

# ==================== Git Commands ====================
# Auto commit with message: make git-commit m="your message"
git-commit:
	@if [ -z "$(m)" ]; then \
		echo "❌ Error: Please provide a commit message with m=\"your message\""; \
		exit 1; \
	fi
	git add .
	git commit -m "$(m)"
	@echo "✅ Committed: $(m)"

# Push current branch
git-push:
	git push -u origin $$(git branch --show-current)
	@echo "✅ Pushed to origin/$$(git branch --show-current)"

# Create feature branch: make git-feature name="auth-service"
git-feature:
	@if [ -z "$(name)" ]; then \
		echo "❌ Error: Please provide a branch name with name=\"branch-name\""; \
		exit 1; \
	fi
	git checkout -b feature/$(name)
	@echo "✅ Created and switched to feature/$(name)"

# Create hotfix branch: make git-hotfix name="fix-login"
git-hotfix:
	@if [ -z "$(name)" ]; then \
		echo "❌ Error: Please provide a branch name with name=\"branch-name\""; \
		exit 1; \
	fi
	git checkout main
	git pull origin main
	git checkout -b hotfix/$(name)
	@echo "✅ Created and switched to hotfix/$(name)"

# Create release branch: make git-release version="1.0.0"
git-release:
	@if [ -z "$(version)" ]; then \
		echo "❌ Error: Please provide a version with version=\"x.x.x\""; \
		exit 1; \
	fi
	git checkout main
	git pull origin main
	git checkout -b release/v$(version)
	@echo "✅ Created and switched to release/v$(version)"

# Merge current branch to main
git-merge-main:
	@CURRENT_BRANCH=$$(git branch --show-current); \
	if [ "$$CURRENT_BRANCH" = "main" ]; then \
		echo "❌ Error: Already on main branch"; \
		exit 1; \
	fi; \
	git checkout main && \
	git pull origin main && \
	git merge $$CURRENT_BRANCH && \
	echo "✅ Merged $$CURRENT_BRANCH into main"

# Quick commit and push: make git-quick m="your message"
git-quick: git-commit git-push
