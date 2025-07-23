/**
 * Godin Hot Reload Client
 * Handles hot reload and hot refresh functionality via WebSocket
 */

class GodinHotReload {
    constructor() {
        this.websocket = null;
        this.wsUrl = this.getWebSocketURL();
        this.reconnectAttempts = 0;
        this.maxReconnectAttempts = 10;
        this.reconnectDelay = 1000;
        this.isEnabled = true;

        this.init();
    }

    init() {
        // Only initialize in development mode
        if (!this.isDevelopmentMode()) {
            return;
        }

        console.log('🔥 Godin Hot Reload initialized');
        this.connectWebSocket();
        this.setupVisibilityHandler();
    }

    isDevelopmentMode() {
        // Check if we're in development mode
        return window.location.hostname === 'localhost' ||
               window.location.hostname === '127.0.0.1' ||
               window.location.hostname.startsWith('192.168.') ||
               window.location.hostname.startsWith('10.') ||
               window.location.hostname.includes('dev');
    }

    getWebSocketURL() {
        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        const host = window.location.host;
        return `${protocol}//${host}/ws`;
    }

    connectWebSocket() {
        if (this.websocket && this.websocket.readyState === WebSocket.OPEN) {
            return;
        }

        console.log('🔌 Connecting to hot reload WebSocket:', this.wsUrl);

        try {
            this.websocket = new WebSocket(this.wsUrl);

            this.websocket.onopen = (event) => {
                console.log('✅ Hot reload WebSocket connected');
                this.reconnectAttempts = 0;
                this.onWebSocketOpen(event);
            };

            this.websocket.onmessage = (event) => {
                this.onWebSocketMessage(event);
            };

            this.websocket.onclose = (event) => {
                console.log('❌ Hot reload WebSocket disconnected');
                this.onWebSocketClose(event);
            };

            this.websocket.onerror = (event) => {
                console.error('🚨 Hot reload WebSocket error:', event);
            };
        } catch (error) {
            console.error('❌ Failed to create hot reload WebSocket:', error);
            this.scheduleReconnect();
        }
    }

    onWebSocketOpen(event) {
        // Subscribe to hot reload channels
        this.subscribe('hot-reload');
        this.subscribe('hot-refresh');

        // Show hot reload status
        this.showStatus('🔥 Hot reload active', 'success');
    }

    onWebSocketMessage(event) {
        try {
            const message = JSON.parse(event.data);
            this.handleHotReloadMessage(message);
        } catch (error) {
            console.error('❌ Failed to parse hot reload message:', error);
        }
    }

    onWebSocketClose(event) {
        this.showStatus('🔌 Hot reload disconnected', 'warning');
        this.scheduleReconnect();
    }

    scheduleReconnect() {
        if (this.reconnectAttempts < this.maxReconnectAttempts) {
            const delay = this.reconnectDelay * Math.pow(1.5, this.reconnectAttempts);
            setTimeout(() => {
                this.reconnectAttempts++;
                console.log(`🔄 Attempting to reconnect (${this.reconnectAttempts}/${this.maxReconnectAttempts})`);
                this.connectWebSocket();
            }, delay);
        } else {
            console.error('❌ Max reconnection attempts reached');
            this.showStatus('❌ Hot reload failed', 'error');
        }
    }

    subscribe(channel) {
        if (this.websocket && this.websocket.readyState === WebSocket.OPEN) {
            const message = {
                type: 'subscribe',
                channel: channel
            };
            this.websocket.send(JSON.stringify(message));
        }
    }

    handleHotReloadMessage(message) {
        console.log('📨 Hot reload message:', message);

        switch (message.type) {
            case 'hot-reload':
                this.handleHotReload(message);
                break;
            case 'hot-refresh':
                this.handleHotRefresh(message);
                break;
            default:
                console.log('🤷 Unknown hot reload message type:', message.type);
        }
    }

    handleHotReload(message) {
        console.log('🔥 Hot reload triggered - reloading page...');
        this.showStatus('🔥 Reloading...', 'info');

        // Save current state before reload
        this.saveCurrentState();

        // Add a small delay to ensure the server is ready
        setTimeout(() => {
            window.location.reload();
        }, 1000);
    }

    handleHotRefresh(message) {
        console.log('🔄 Hot refresh triggered - refreshing content...');
        this.showStatus('🔄 Refreshing...', 'info');

        // Refresh CSS files
        this.refreshCSS();

        // Refresh JavaScript files (if needed)
        this.refreshJS();

        // Trigger HTMX refresh for dynamic content
        this.refreshHTMX();

        // Refresh any state-dependent elements
        this.refreshStateElements();

        setTimeout(() => {
            this.showStatus('✅ Refreshed', 'success');
        }, 1000);
    }

    refreshCSS() {
        const links = document.querySelectorAll('link[rel="stylesheet"]');
        links.forEach(link => {
            const href = link.href;
            const url = new URL(href);
            url.searchParams.set('_t', Date.now());
            link.href = url.toString();
        });
    }

    refreshJS() {
        // For now, we'll just log this - full JS refresh is complex
        console.log('🔄 JavaScript refresh (CSS and HTMX refreshed)');
    }

    refreshHTMX() {
        // Trigger HTMX refresh for elements with hx-get
        if (typeof htmx !== 'undefined') {
            const elements = document.querySelectorAll('[hx-get]');
            elements.forEach(element => {
                htmx.trigger(element, 'refresh');
            });
        }
    }

    showStatus(message, type = 'info') {
        // Create or update status indicator
        let indicator = document.getElementById('godin-hot-reload-status');
        if (!indicator) {
            indicator = document.createElement('div');
            indicator.id = 'godin-hot-reload-status';
            indicator.style.cssText = `
                position: fixed;
                top: 10px;
                right: 10px;
                padding: 8px 12px;
                border-radius: 4px;
                font-family: monospace;
                font-size: 12px;
                z-index: 10000;
                transition: all 0.3s ease;
                pointer-events: none;
            `;
            document.body.appendChild(indicator);
        }

        // Set colors based on type
        const colors = {
            success: { bg: '#10b981', text: '#ffffff' },
            warning: { bg: '#f59e0b', text: '#ffffff' },
            error: { bg: '#ef4444', text: '#ffffff' },
            info: { bg: '#3b82f6', text: '#ffffff' }
        };

        const color = colors[type] || colors.info;
        indicator.style.backgroundColor = color.bg;
        indicator.style.color = color.text;
        indicator.textContent = message;

        // Auto-hide success messages
        if (type === 'success') {
            setTimeout(() => {
                if (indicator.textContent === message) {
                    indicator.style.opacity = '0';
                    setTimeout(() => {
                        if (indicator.style.opacity === '0') {
                            indicator.remove();
                        }
                    }, 300);
                }
            }, 2000);
        }
    }

    saveCurrentState() {
        try {
            // Save form data
            const forms = document.querySelectorAll('form');
            const formData = {};
            forms.forEach((form, index) => {
                const data = new FormData(form);
                formData[`form_${index}`] = Object.fromEntries(data);
            });

            // Save input values
            const inputs = document.querySelectorAll('input, textarea, select');
            const inputData = {};
            inputs.forEach((input, index) => {
                if (input.id || input.name) {
                    const key = input.id || input.name || `input_${index}`;
                    inputData[key] = input.value;
                }
            });

            // Save scroll position
            const scrollData = {
                x: window.scrollX,
                y: window.scrollY
            };

            // Store in sessionStorage
            sessionStorage.setItem('godin_hot_reload_state', JSON.stringify({
                forms: formData,
                inputs: inputData,
                scroll: scrollData,
                timestamp: Date.now()
            }));

            console.log('💾 Current state saved for hot reload');
        } catch (error) {
            console.warn('⚠️ Failed to save state:', error);
        }
    }

    restoreCurrentState() {
        try {
            const savedState = sessionStorage.getItem('godin_hot_reload_state');
            if (!savedState) return;

            const state = JSON.parse(savedState);

            // Only restore if saved recently (within 10 seconds)
            if (Date.now() - state.timestamp > 10000) {
                sessionStorage.removeItem('godin_hot_reload_state');
                return;
            }

            // Restore input values
            if (state.inputs) {
                Object.entries(state.inputs).forEach(([key, value]) => {
                    const input = document.getElementById(key) || document.querySelector(`[name="${key}"]`);
                    if (input && input.value !== value) {
                        input.value = value;
                        // Trigger change event
                        input.dispatchEvent(new Event('change', { bubbles: true }));
                    }
                });
            }

            // Restore scroll position
            if (state.scroll) {
                window.scrollTo(state.scroll.x, state.scroll.y);
            }

            // Clean up
            sessionStorage.removeItem('godin_hot_reload_state');
            console.log('🔄 State restored after hot reload');
        } catch (error) {
            console.warn('⚠️ Failed to restore state:', error);
        }
    }

    refreshStateElements() {
        // Refresh elements that depend on server state
        const stateElements = document.querySelectorAll('[data-state-key], [hx-get]');
        stateElements.forEach(element => {
            if (element.hasAttribute('hx-get')) {
                // Trigger HTMX refresh
                if (typeof htmx !== 'undefined') {
                    htmx.trigger(element, 'refresh');
                }
            }
        });
    }

    setupVisibilityHandler() {
        // Reconnect when page becomes visible (handles browser sleep/wake)
        document.addEventListener('visibilitychange', () => {
            if (!document.hidden && this.websocket && this.websocket.readyState !== WebSocket.OPEN) {
                console.log('🔄 Page visible, reconnecting hot reload...');
                this.connectWebSocket();
            }
        });

        // Restore state on page load (after hot reload)
        if (document.readyState === 'complete') {
            this.restoreCurrentState();
        } else {
            window.addEventListener('load', () => {
                this.restoreCurrentState();
            });
        }
    }
}

// Initialize hot reload when DOM is ready
if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', () => {
        window.godinHotReload = new GodinHotReload();
    });
} else {
    window.godinHotReload = new GodinHotReload();
}

// Export for manual control
window.GodinHotReload = GodinHotReload;