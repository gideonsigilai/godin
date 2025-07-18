/**
 * Godin Framework JavaScript Runtime
 * Provides client-side functionality for HTMX integration and WebSocket support
 */

class GodinFramework {
    constructor() {
        this.websocket = null;
        this.wsUrl = this.getWebSocketURL();
        this.reconnectAttempts = 0;
        this.maxReconnectAttempts = 5;
        this.reconnectDelay = 1000;
        this.subscriptions = new Map();
        
        this.init();
    }
    
    init() {
        // Initialize when DOM is ready
        if (document.readyState === 'loading') {
            document.addEventListener('DOMContentLoaded', () => this.onDOMReady());
        } else {
            this.onDOMReady();
        }
    }
    
    onDOMReady() {
        console.log('Godin Framework initialized');
        
        // Initialize WebSocket connection
        this.connectWebSocket();
        
        // Setup HTMX event listeners
        this.setupHTMXListeners();
        
        // Setup UI event listeners
        this.setupUIListeners();
    }
    
    // WebSocket Management
    getWebSocketURL() {
        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        const host = window.location.host;
        return `${protocol}//${host}/ws`;
    }
    
    connectWebSocket() {
        if (this.websocket && this.websocket.readyState === WebSocket.OPEN) {
            return;
        }
        
        try {
            this.websocket = new WebSocket(this.wsUrl);
            
            this.websocket.onopen = (event) => {
                console.log('WebSocket connected');
                this.reconnectAttempts = 0;
                this.onWebSocketOpen(event);
            };
            
            this.websocket.onmessage = (event) => {
                this.onWebSocketMessage(event);
            };
            
            this.websocket.onclose = (event) => {
                console.log('WebSocket disconnected');
                this.onWebSocketClose(event);
            };
            
            this.websocket.onerror = (event) => {
                console.error('WebSocket error:', event);
                this.onWebSocketError(event);
            };
        } catch (error) {
            console.error('Failed to create WebSocket connection:', error);
        }
    }
    
    onWebSocketOpen(event) {
        // Send any queued subscriptions
        this.subscriptions.forEach((callback, channel) => {
            this.subscribe(channel, callback);
        });
    }
    
    onWebSocketMessage(event) {
        try {
            const message = JSON.parse(event.data);
            this.handleWebSocketMessage(message);
        } catch (error) {
            console.error('Failed to parse WebSocket message:', error);
        }
    }
    
    onWebSocketClose(event) {
        if (this.reconnectAttempts < this.maxReconnectAttempts) {
            setTimeout(() => {
                this.reconnectAttempts++;
                console.log(`Attempting to reconnect (${this.reconnectAttempts}/${this.maxReconnectAttempts})`);
                this.connectWebSocket();
            }, this.reconnectDelay * this.reconnectAttempts);
        }
    }
    
    onWebSocketError(event) {
        // Handle WebSocket errors
    }
    
    handleWebSocketMessage(message) {
        switch (message.type) {
            case 'broadcast':
                this.handleBroadcast(message);
                break;
            case 'pong':
                // Handle ping/pong
                break;
            default:
                console.log('Unknown WebSocket message type:', message.type);
        }
    }
    
    handleBroadcast(message) {
        const callback = this.subscriptions.get(message.channel);
        if (callback) {
            callback(message.data);
        }

        // Handle state changes for automatic UI updates
        if (message.channel.startsWith('state:')) {
            this.handleStateChange(message.channel, message.data);
        }

        // Trigger custom event
        const event = new CustomEvent('godin:broadcast', {
            detail: {
                channel: message.channel,
                data: message.data
            }
        });
        document.dispatchEvent(event);
    }

    handleStateChange(channel, data) {
        const stateKey = channel.replace('state:', '');

        // Find all elements that depend on this state key
        const stateElements = document.querySelectorAll(`[data-state-key="${stateKey}"]`);

        stateElements.forEach(element => {
            // Trigger HTMX refresh for state-dependent elements
            if (element.hasAttribute('hx-get')) {
                htmx.trigger(element, 'refresh');
            } else if (element.hasAttribute('data-state-endpoint')) {
                // Fetch updated content for this state-dependent element
                const endpoint = element.getAttribute('data-state-endpoint');
                fetch(endpoint)
                    .then(response => response.text())
                    .then(html => {
                        element.innerHTML = html;
                        this.initializeComponents(element);
                    })
                    .catch(error => console.error('Error updating state element:', error));
            }
        });

        // Trigger custom state change event
        const stateEvent = new CustomEvent('godin:stateChange', {
            detail: {
                key: stateKey,
                value: data.value
            }
        });
        document.dispatchEvent(stateEvent);
    }
    
    subscribe(channel, callback) {
        this.subscriptions.set(channel, callback);
        
        if (this.websocket && this.websocket.readyState === WebSocket.OPEN) {
            this.websocket.send(JSON.stringify({
                type: 'subscribe',
                channel: channel
            }));
        }
    }
    
    unsubscribe(channel) {
        this.subscriptions.delete(channel);
        
        if (this.websocket && this.websocket.readyState === WebSocket.OPEN) {
            this.websocket.send(JSON.stringify({
                type: 'unsubscribe',
                channel: channel
            }));
        }
    }
    
    // HTMX Integration
    setupHTMXListeners() {
        // Listen for HTMX events
        document.addEventListener('htmx:beforeRequest', (event) => {
            this.onHTMXBeforeRequest(event);
        });
        
        document.addEventListener('htmx:afterRequest', (event) => {
            this.onHTMXAfterRequest(event);
        });
        
        document.addEventListener('htmx:afterSwap', (event) => {
            this.onHTMXAfterSwap(event);
        });
    }
    
    onHTMXBeforeRequest(event) {
        // Add loading indicators
        const target = event.target;
        target.classList.add('godin-loading');
    }
    
    onHTMXAfterRequest(event) {
        // Remove loading indicators
        const target = event.target;
        target.classList.remove('godin-loading');
    }
    
    onHTMXAfterSwap(event) {
        // Re-initialize any new components
        this.initializeComponents(event.target);
    }
    
    // UI Event Listeners
    setupUIListeners() {
        // Handle drawer toggles
        document.addEventListener('click', (event) => {
            if (event.target.matches('[data-godin-drawer-toggle]')) {
                const drawerId = event.target.getAttribute('data-godin-drawer-toggle');
                this.toggleDrawer(drawerId);
            }
        });
        
        // Handle dialog close
        document.addEventListener('click', (event) => {
            if (event.target.matches('.godin-dialog-overlay')) {
                this.closeDialog(event.target.nextElementSibling);
            }
        });
        
        // Handle tab switching
        document.addEventListener('click', (event) => {
            if (event.target.matches('.godin-tab')) {
                this.switchTab(event.target);
            }
        });
    }
    
    // UI Component Methods
    toggleDrawer(drawerId) {
        const drawer = document.getElementById(drawerId);
        if (drawer) {
            drawer.classList.toggle('open');
        }
    }
    
    openDialog(dialogId) {
        const dialog = document.getElementById(dialogId);
        if (dialog) {
            dialog.style.display = 'block';
        }
    }
    
    closeDialog(dialog) {
        if (dialog) {
            dialog.style.display = 'none';
        }
    }
    
    switchTab(tabElement) {
        const tabBar = tabElement.closest('.godin-tab-bar');
        const tabs = tabBar.querySelectorAll('.godin-tab');
        
        // Remove active class from all tabs
        tabs.forEach(tab => tab.classList.remove('active'));
        
        // Add active class to clicked tab
        tabElement.classList.add('active');
        
        // Trigger custom event
        const event = new CustomEvent('godin:tabSwitch', {
            detail: {
                tab: tabElement,
                index: Array.from(tabs).indexOf(tabElement)
            }
        });
        document.dispatchEvent(event);
    }
    
    showSnackbar(message, type = 'info', duration = 3000) {
        const snackbar = document.createElement('div');
        snackbar.className = `godin-snackbar godin-snackbar-${type}`;
        snackbar.textContent = message;
        
        document.body.appendChild(snackbar);
        
        // Auto-remove after duration
        setTimeout(() => {
            if (snackbar.parentNode) {
                snackbar.parentNode.removeChild(snackbar);
            }
        }, duration);
    }
    
    initializeComponents(container = document) {
        // Initialize tooltips
        const tooltips = container.querySelectorAll('.godin-tooltip');
        tooltips.forEach(tooltip => this.initializeTooltip(tooltip));
        
        // Initialize progress indicators
        const progressBars = container.querySelectorAll('.godin-progress-linear');
        progressBars.forEach(bar => this.initializeProgressBar(bar));
    }
    
    initializeTooltip(tooltip) {
        // Tooltip initialization logic
    }
    
    initializeProgressBar(progressBar) {
        // Progress bar initialization logic
    }
    
    // Utility Methods
    debounce(func, wait) {
        let timeout;
        return function executedFunction(...args) {
            const later = () => {
                clearTimeout(timeout);
                func(...args);
            };
            clearTimeout(timeout);
            timeout = setTimeout(later, wait);
        };
    }
    
    throttle(func, limit) {
        let inThrottle;
        return function() {
            const args = arguments;
            const context = this;
            if (!inThrottle) {
                func.apply(context, args);
                inThrottle = true;
                setTimeout(() => inThrottle = false, limit);
            }
        };
    }
}

// Initialize Godin Framework
window.Godin = new GodinFramework();

// Export for module systems
if (typeof module !== 'undefined' && module.exports) {
    module.exports = GodinFramework;
}
