#!/bin/bash
set -euo pipefail

# Monotonix Actions ref切り替えスクリプト
# Usage:
#   ./scripts/switch-monotonix-ref.sh <ref>          # 指定したrefに切り替え
#   ./scripts/switch-monotonix-ref.sh --reset        # デフォルト(最新リリース)に戻す
#   ./scripts/switch-monotonix-ref.sh --current      # 現在のrefを表示

WORKFLOWS_DIR=".github/workflows"
MONOTONIX_REPO="yuya-takeyama/monotonix"

# 色付き出力
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
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

# クリーンアップ処理
cleanup_bak_files() {
    find "$WORKFLOWS_DIR" -name "*.bak" -delete 2>/dev/null || true
}

# 異常終了時のクリーンアップ
trap cleanup_bak_files EXIT

# 最新リリースの情報を取得
get_latest_release() {
    local release_info=$(gh release view --repo "$MONOTONIX_REPO" --json tagName,targetCommitish 2>/dev/null)
    
    if [[ -z "$release_info" ]]; then
        error "Failed to fetch latest release information"
        exit 1
    fi
    
    local tag=$(echo "$release_info" | jq -r .tagName)
    local commit=$(echo "$release_info" | jq -r .targetCommitish)
    
    # commitがSHAでない場合（ブランチ名の場合）、実際のSHAを取得
    if [[ ! "$commit" =~ ^[a-f0-9]{40}$ ]]; then
        commit=$(gh api "repos/$MONOTONIX_REPO/commits/$tag" --jq .sha 2>/dev/null | cut -c1-40)
    fi
    
    echo "$commit $tag"
}

# 現在のrefを取得
get_current_ref() {
    local workflow_file="$WORKFLOWS_DIR/build-docker.yaml"
    if [[ ! -f "$workflow_file" ]]; then
        error "Workflow file not found: $workflow_file"
        return 1
    fi
    
    # ブランチ名もコミットSHAもマッチするように修正
    grep -oE 'yuya-takeyama/monotonix/actions/[^@]+@[^[:space:]#]+' "$workflow_file" | head -1 | cut -d'@' -f2
}

# 現在のrefを表示
show_current() {
    local current_ref=$(get_current_ref)
    if [[ -z "$current_ref" ]]; then
        error "Failed to get current ref"
        exit 1
    fi
    
    # 最新リリースと比較
    local latest_info=$(get_latest_release)
    local latest_ref=$(echo "$latest_info" | cut -d' ' -f1)
    local latest_tag=$(echo "$latest_info" | cut -d' ' -f2)
    
    if [[ "$current_ref" == "$latest_ref" ]]; then
        info "Current ref: $current_ref (latest release: $latest_tag)"
    else
        info "Current ref: $current_ref"
        info "Latest release: $latest_ref ($latest_tag)"
    fi
}

# refを切り替え
switch_ref() {
    local new_ref="$1"
    local version_comment="${2:-}"  # オプション: バージョンコメント
    local current_ref=$(get_current_ref)
    
    if [[ -z "$current_ref" ]]; then
        error "Failed to get current ref"
        exit 1
    fi
    
    if [[ "$current_ref" == "$new_ref" ]]; then
        warn "Already using ref: $new_ref"
        return 0
    fi
    
    info "Switching from $current_ref to $new_ref"
    
    # すべてのworkflowファイルを更新
    for workflow_file in "$WORKFLOWS_DIR"/*.yaml "$WORKFLOWS_DIR"/*.yml; do
        if [[ -f "$workflow_file" ]]; then
            local filename=$(basename "$workflow_file")
            
            # monotonix actionsが含まれているか確認
            if grep -q "yuya-takeyama/monotonix/actions" "$workflow_file"; then
                info "Updating $filename"
                
                # refとコメント部分を置き換え
                # macOSとLinuxの両方で動作するように調整
                if [[ -n "$version_comment" ]]; then
                    # バージョンコメント付きで置き換え（コメントあり・なし両方を1つの正規表現で処理）
                    sed -i.bak -E "s|(yuya-takeyama/monotonix/actions/[^@]+)@[^[:space:]#]+([[:space:]]*#[^$]*)?|\\1@$new_ref # $version_comment|g" "$workflow_file"
                else
                    # コメントなしで置き換え（既存のコメントも削除）
                    sed -i.bak -E "s|(yuya-takeyama/monotonix/actions/[^@]+)@[^[:space:]#]+([[:space:]]*#[^$]*)?|\\1@$new_ref|g" "$workflow_file"
                fi
            fi
        fi
    done
    
    info "✅ Successfully switched to ref: $new_ref"
}

# デフォルト（最新リリース）に戻す
reset_to_default() {
    local latest_info=$(get_latest_release)
    local latest_ref=$(echo "$latest_info" | cut -d' ' -f1)
    local latest_tag=$(echo "$latest_info" | cut -d' ' -f2)
    
    info "Resetting to latest release: $latest_ref ($latest_tag)"
    switch_ref "$latest_ref" "$latest_tag"
}

# メイン処理
main() {
    if [[ $# -eq 0 ]]; then
        error "Usage: $0 <ref> | --reset | --current"
        echo ""
        echo "Options:"
        echo "  <ref>      Switch to specified ref (commit SHA or branch/tag)"
        echo "  --reset    Reset to latest release"
        echo "  --current  Show current ref"
        echo ""
        echo "Examples:"
        echo "  $0 feat/metadata-support"
        echo "  $0 abc123def456"
        echo "  $0 --reset"
        exit 1
    fi
    
    case "$1" in
        --reset)
            reset_to_default
            ;;
        --current)
            show_current
            ;;
        -*)
            error "Unknown option: $1"
            exit 1
            ;;
        *)
            switch_ref "$1"
            ;;
    esac
}

main "$@"
