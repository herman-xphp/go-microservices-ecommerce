.PHONY: up down restart logs ps \
	git-commit git-push git-feature git-hotfix git-release \
	git-merge-develop git-finish-release git-quick

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

# ==================== Git Commands (GitFlow) ====================
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

# Create feature branch FROM DEVELOP: make git-feature name="auth-service"
git-feature:
	@if [ -z "$(name)" ]; then \
		echo "❌ Error: Please provide a branch name with name=\"branch-name\""; \
		exit 1; \
	fi
	git checkout develop
	git pull origin develop
	git checkout -b feature/$(name)
	@echo "✅ Created feature/$(name) from develop"

# Create hotfix branch FROM MAIN: make git-hotfix name="fix-login"
git-hotfix:
	@if [ -z "$(name)" ]; then \
		echo "❌ Error: Please provide a branch name with name=\"branch-name\""; \
		exit 1; \
	fi
	git checkout main
	git pull origin main
	git checkout -b hotfix/$(name)
	@echo "✅ Created hotfix/$(name) from main"

# Create release branch FROM DEVELOP: make git-release version="1.0.0"
git-release:
	@if [ -z "$(version)" ]; then \
		echo "❌ Error: Please provide a version with version=\"x.x.x\""; \
		exit 1; \
	fi
	git checkout develop
	git pull origin develop
	git checkout -b release/v$(version)
	@echo "✅ Created release/v$(version) from develop"

# Merge current feature branch to DEVELOP
git-merge-develop:
	@CURRENT_BRANCH=$$(git branch --show-current); \
	if [ "$$CURRENT_BRANCH" = "develop" ]; then \
		echo "❌ Error: Already on develop branch"; \
		exit 1; \
	fi; \
	git checkout develop && \
	git pull origin develop && \
	git merge $$CURRENT_BRANCH && \
	echo "✅ Merged $$CURRENT_BRANCH into develop"

# Finish release: merge release branch to main AND develop, then tag
git-finish-release:
	@CURRENT_BRANCH=$$(git branch --show-current); \
	if echo "$$CURRENT_BRANCH" | grep -q "^release/"; then \
		VERSION=$$(echo $$CURRENT_BRANCH | sed 's/release\///'); \
		git checkout main && \
		git pull origin main && \
		git merge $$CURRENT_BRANCH && \
		git tag -a $$VERSION -m "Release $$VERSION" && \
		git checkout develop && \
		git pull origin develop && \
		git merge $$CURRENT_BRANCH && \
		echo "✅ Finished release $$VERSION (merged to main & develop, tagged)"; \
	else \
		echo "❌ Error: Not on a release branch (current: $$CURRENT_BRANCH)"; \
		exit 1; \
	fi

# Quick commit and push: make git-quick m="your message"
git-quick: git-commit git-push
