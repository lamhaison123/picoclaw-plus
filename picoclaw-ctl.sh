#!/bin/bash
# PicoClaw Control Script - Helper for managing PicoClaw service

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Functions
print_header() {
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${BLUE}  PicoClaw Control Panel${NC}"
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
}

print_success() {
    echo -e "${GREEN}✓${NC} $1"
}

print_error() {
    echo -e "${RED}✗${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}⚠${NC} $1"
}

print_info() {
    echo -e "${BLUE}ℹ${NC} $1"
}

check_root() {
    if [ "$EUID" -ne 0 ]; then 
        print_error "Please run as root (use sudo)"
        exit 1
    fi
}

check_service() {
    if ! systemctl list-unit-files | grep -q "picoclaw.service"; then
        print_error "PicoClaw service not installed"
        echo "Run: sudo bash install-systemd.sh"
        exit 1
    fi
}

show_status() {
    print_header
    echo ""
    
    # Service status
    if systemctl is-active --quiet picoclaw; then
        print_success "Service is running"
    else
        print_error "Service is stopped"
    fi
    
    # Enabled status
    if systemctl is-enabled --quiet picoclaw; then
        print_success "Auto-start enabled"
    else
        print_warning "Auto-start disabled"
    fi
    
    echo ""
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    systemctl status picoclaw --no-pager -l
}

start_service() {
    check_root
    check_service
    
    print_info "Starting PicoClaw service..."
    systemctl start picoclaw
    sleep 2
    
    if systemctl is-active --quiet picoclaw; then
        print_success "Service started successfully"
        echo ""
        print_info "View logs: sudo journalctl -u picoclaw -f"
    else
        print_error "Failed to start service"
        echo ""
        print_info "Check logs: sudo journalctl -u picoclaw -n 50"
        exit 1
    fi
}

stop_service() {
    check_root
    check_service
    
    print_info "Stopping PicoClaw service..."
    systemctl stop picoclaw
    sleep 1
    
    if ! systemctl is-active --quiet picoclaw; then
        print_success "Service stopped successfully"
    else
        print_error "Failed to stop service"
        exit 1
    fi
}

restart_service() {
    check_root
    check_service
    
    print_info "Restarting PicoClaw service..."
    systemctl restart picoclaw
    sleep 2
    
    if systemctl is-active --quiet picoclaw; then
        print_success "Service restarted successfully"
        echo ""
        print_info "View logs: sudo journalctl -u picoclaw -f"
    else
        print_error "Failed to restart service"
        echo ""
        print_info "Check logs: sudo journalctl -u picoclaw -n 50"
        exit 1
    fi
}

enable_service() {
    check_root
    check_service
    
    print_info "Enabling auto-start..."
    systemctl enable picoclaw
    print_success "Auto-start enabled"
}

disable_service() {
    check_root
    check_service
    
    print_info "Disabling auto-start..."
    systemctl disable picoclaw
    print_success "Auto-start disabled"
}

show_logs() {
    check_service
    
    local lines="${1:-50}"
    print_info "Showing last $lines log lines..."
    echo ""
    journalctl -u picoclaw -n "$lines" --no-pager
}

follow_logs() {
    check_service
    
    print_info "Following logs (Ctrl+C to stop)..."
    echo ""
    journalctl -u picoclaw -f
}

show_help() {
    print_header
    echo ""
    echo "Usage: $0 <command> [options]"
    echo ""
    echo "Commands:"
    echo "  status              Show service status"
    echo "  start               Start the service"
    echo "  stop                Stop the service"
    echo "  restart             Restart the service"
    echo "  enable              Enable auto-start on boot"
    echo "  disable             Disable auto-start"
    echo "  logs [lines]        Show last N log lines (default: 50)"
    echo "  follow              Follow logs in real-time"
    echo "  help                Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 status           # Check service status"
    echo "  $0 start            # Start PicoClaw"
    echo "  $0 logs 100         # Show last 100 log lines"
    echo "  $0 follow           # Follow logs"
    echo ""
}

# Main
case "${1:-}" in
    status)
        show_status
        ;;
    start)
        start_service
        ;;
    stop)
        stop_service
        ;;
    restart)
        restart_service
        ;;
    enable)
        enable_service
        ;;
    disable)
        disable_service
        ;;
    logs)
        show_logs "${2:-50}"
        ;;
    follow)
        follow_logs
        ;;
    help|--help|-h)
        show_help
        ;;
    *)
        print_error "Unknown command: ${1:-}"
        echo ""
        show_help
        exit 1
        ;;
esac
