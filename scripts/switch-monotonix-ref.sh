#!/bin/bash
set -euo pipefail

# Monotonix Actions ref切り替えスクリプト
# Usage:
#   ./scripts/switch-monotonix-ref.sh <ref>          # 指定したrefに切り替え
#   ./scripts/switch-monotonix-ref.sh --reset        # デフォルト(v0.0.4)に戻す
#   ./scripts/switch-monotonix-ref.sh --current      # 現在のrefを表示

DEFAULT_REF="1ee58090547501c1b691407d21db1dee9374de7e" # v0.0.4
DEFAULT_VERSION="v0.0.4"
WORKFLOWS_DIR=".github/workflows"

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

# 現在のrefを取得
get_current_ref() {
    local workflow_file="$WORKFLOWS_DIR/build-docker.yaml"
    if [[ ! -f "$workflow_file" ]]; then
        error "Workflow file not found: $workflow_file"
        return 1
    fi
    
    grep -oE 'yuya-takeyama/monotonix/actions/[^@]+@[a-f0-9]+' "$workflow_file" | head -1 | cut -d'@' -f2
}

# 現在のrefを表示
show_current() {
    local current_ref=$(get_current_ref)
    if [[ -z "$current_ref" ]]; then
        error "Failed to get current ref"
        exit 1
    fi
    
    if [[ "$current_ref" == "$DEFAULT_REF" ]]; then
        info "Current ref: $current_ref (default: $DEFAULT_VERSION)"
    else
        info "Current ref: $current_ref"
    fi
}

# refを切り替え
switch_ref() {
    local new_ref="$1"
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
                sed -i.bak "s|yuya-takeyama/monotonix/actions/\\([^@]*\\)@[a-f0-9]*|yuya-takeyama/monotonix/actions/\\1@$new_ref|g" "$workflow_file"
                rm "${workflow_file}.bak"
            fi
        fi
    done
    
    info "✅ Successfully switched to ref: $new_ref"
}

# デフォルトに戻す
reset_to_default() {
    info "Resetting to default ref: $DEFAULT_REF ($DEFAULT_VERSION)"
    switch_ref "$DEFAULT_REF"
}

# メイン処理
main() {
    if [[ $# -eq 0 ]]; then
        error "Usage: $0 <ref> | --reset | --current"
        echo ""
        echo "Options:"
        echo "  <ref>      Switch to specified ref (commit SHA or branch/tag)"
        echo "  --reset    Reset to default ref ($DEFAULT_VERSION)"
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