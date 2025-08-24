#!/bin/bash
set -euo pipefail

# Monotonix PRテスト用スクリプト
# Usage:
#   ./scripts/test-monotonix-pr.sh <pr-number>    # 指定したPRのブランチでテスト
#   ./scripts/test-monotonix-pr.sh --cleanup      # テスト後のクリーンアップ

WORKFLOWS_DIR=".github/workflows"
SWITCH_SCRIPT="./scripts/switch-monotonix-ref.sh"

# 色付き出力
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1" >&2
}

step() {
    echo -e "${BLUE}==>${NC} $1"
}

# PRの情報を取得
get_pr_info() {
    local pr_number="$1"
    
    step "Fetching PR #$pr_number information from yuya-takeyama/monotonix"
    
    # gh CLIでPR情報を取得
    local pr_info=$(gh pr view "$pr_number" --repo yuya-takeyama/monotonix --json headRefName,state,mergeable 2>/dev/null)
    
    if [[ -z "$pr_info" ]]; then
        error "Failed to fetch PR #$pr_number"
        exit 1
    fi
    
    local branch=$(echo "$pr_info" | jq -r .headRefName)
    local state=$(echo "$pr_info" | jq -r .state)
    
    info "PR #$pr_number: $branch (state: $state)"
    echo "$branch"
}

# PRのブランチでテスト
test_pr() {
    local pr_number="$1"
    
    # PR情報を取得
    local branch=$(get_pr_info "$pr_number")
    
    if [[ -z "$branch" ]]; then
        error "Failed to get branch name for PR #$pr_number"
        exit 1
    fi
    
    step "Switching monotonix actions ref to branch: $branch"
    "$SWITCH_SCRIPT" "$branch"
    
    step "Creating test commit"
    git add "$WORKFLOWS_DIR"/*.yaml "$WORKFLOWS_DIR"/*.yml 2>/dev/null || true
    
    if git diff --cached --quiet; then
        warn "No changes to commit"
    else
        git commit -m "test: monotonix PR #$pr_number ($branch)

Testing https://github.com/yuya-takeyama/monotonix/pull/$pr_number

This is a temporary commit for testing purposes.
Will be reverted after testing is complete.

Co-Authored-By: Claude <noreply@anthropic.com>"
        
        info "✅ Test commit created"
        echo ""
        info "Next steps:"
        echo "  1. Push to trigger CI: git push"
        echo "  2. Monitor CI results"
        echo "  3. After testing, run: $0 --cleanup"
    fi
}

# クリーンアップ
cleanup() {
    step "Cleaning up test changes"
    
    # デフォルトに戻す
    "$SWITCH_SCRIPT" --reset
    
    # 変更をコミット
    git add "$WORKFLOWS_DIR"/*.yaml "$WORKFLOWS_DIR"/*.yml 2>/dev/null || true
    
    if git diff --cached --quiet; then
        warn "No changes to commit"
    else
        git commit -m "revert: restore default monotonix actions ref

Reverting test changes and restoring default configuration.

Co-Authored-By: Claude <noreply@anthropic.com>"
        
        info "✅ Cleanup commit created"
        echo ""
        info "Next step: Push to apply cleanup: git push"
    fi
}

# メイン処理
main() {
    if [[ $# -eq 0 ]]; then
        error "Usage: $0 <pr-number> | --cleanup"
        echo ""
        echo "Options:"
        echo "  <pr-number>  Test with specified PR from yuya-takeyama/monotonix"
        echo "  --cleanup    Revert test changes and restore defaults"
        echo ""
        echo "Examples:"
        echo "  $0 151        # Test PR #151"
        echo "  $0 --cleanup  # Cleanup after testing"
        exit 1
    fi
    
    case "$1" in
        --cleanup)
            cleanup
            ;;
        -*)
            error "Unknown option: $1"
            exit 1
            ;;
        *)
            if [[ "$1" =~ ^[0-9]+$ ]]; then
                test_pr "$1"
            else
                error "Invalid PR number: $1"
                exit 1
            fi
            ;;
    esac
}

main "$@"