/**
 * Godin State Manager
 * Handles client-side state management for ValueListener widgets
 * Provides WebSocket integration with fallback polling for real-time state updates
 */

class GodinStateManager {
    constructor(options = {}) {
        this.websocket = null;
        this.wsUrl = this.getWebSocketURL();
        this.subscriptions = new Map(); // notifierId -> [callbacks]
        this.valueListeners = new Map(); // notifierId -> current value
        this.domElements = new Map(); // notifierId -> [DOM elements]
        
        // Configuration
        this.config = {
            reconnectAttempts: 0,
            maxReconnectAttempts: options.maxReconnectAttempts || 5,
            reconnectDelay: options.reconnectDelay || 1000,
            pollingInterval: options.pollingInterval || 5000,
            enableFallbackPolling: options.enableFallbackPolling !== false,
            debug: options.debug || false,
            ...options
        };
        
        // Polling fallback
        this.pollingTimer = null;
        this.pollingEndpoints = new Map(); // notifierId -> polling endpoint
        
        this.init();
    }
    
    init() {
        if (document.readyState === 'loading') {
            document.addEventListener('DOMContentLoaded', () => this.onDOMReady());
        } else {
            this.onDOMReady();
        }
    }
    
    onDOMReady() {
        this.log('GodinStateManager initialized');
        
        // Initialize WebSocket connection
        this.connectWebSocket();
        
        // Scan for existing ValueListener elements
        this.scanForValueListeners();
        
        // Setup mutation observer for dynamic content
        this.setupMutationObserver();
    }
    
    // WebSocket Management
    getWebSocketURL() {
        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        const host = window.location.host;
        return `${protocol}//${host}/ws/state`;
    }
    
    connectWebSocket() {
        if (this.websocket && this.websocket.readyState === WebSocket.OPEN) {
            return;
        }

        this.log('Attempting to connect to WebSocket:', this.wsUrl);

        try {
            this.websocket = new WebSocket(this.wsUrl);

            this.websocket.onopen = (event) => {
                this.log('WebSocket connected successfully');
                this.config.reconnectAttempts = 0;
                this.onWebSocketOpen(event);
            };
            
            this.websocket.onmessage = (event) => {
                this.onWebSocketMessage(event);
            };
            
            this.websocket.onclose = (event) => {
                this.log('WebSocket disconnected');
                this.onWebSocketClose(event);
            };
            
            this.websocket.onerror = (event) => {
                this.log('WebSocket error:', event);
                this.onWebSocketError(event);
            };
        } catch (error) {
            this.log('Failed to create WebSocket connection:', error);
            this.startFallbackPolling();
        }
    }
    
    onWebSocketOpen(event) {
        // Stop fallback polling if WebSocket is connected
        this.stopFallbackPolling();
        
        // Re-subscribe to all active subscriptions
        this.subscriptions.forEach((callbacks, notifierId) => {
            this.sendWebSocketMessage({
                type: 'subscribe',
                notifier_id: notifierId
            });
        });
    }
    
    onWebSocketMessage(event) {
        try {
            const message = JSON.parse(event.data);
            this.handleWebSocketMessage(message);
        } catch (error) {
            this.log('Failed to parse WebSocket message:', error);
        }
    }
    
    onWebSocketClose(event) {
        if (this.config.reconnectAttempts < this.config.maxReconnectAttempts) {
            const delay = this.config.reconnectDelay * Math.pow(2, this.config.reconnectAttempts);
            setTimeout(() => {
                this.config.reconnectAttempts++;
                this.log(`Attempting to reconnect (${this.config.reconnectAttempts}/${this.config.maxReconnectAttempts})`);
                this.connectWebSocket();
            }, delay);
        } else {
            this.log('Max reconnection attempts reached, starting fallback polling');
            this.startFallbackPolling();
        }
    }
    
    onWebSocketError(event) {
        this.log('WebSocket error occurred');
        // Start fallback polling on error
        this.startFallbackPolling();
    }
    
    handleWebSocketMessage(message) {
        switch (message.type) {
            case 'value_change':
                this.handleValueChange(message);
                break;
            case 'subscription_confirmed':
                this.log(`Subscription confirmed for notifier: ${message.notifier_id}`);
                break;
            case 'error':
                this.log('WebSocket error message:', message.error);
                break;
            default:
                this.log('Unknown WebSocket message type:', message.type);
        }
    }
    
    handleValueChange(message) {
        const { notifier_id, value, old_value, timestamp } = message;
        
        this.log(`Value change received for ${notifier_id}:`, { old_value, value });
        
        // Update stored value
        this.valueListeners.set(notifier_id, value);
        
        // Notify all callbacks
        const callbacks = this.subscriptions.get(notifier_id) || [];
        callbacks.forEach(callback => {
            try {
                callback(value, old_value);
            } catch (error) {
                this.log('Error in value change callback:', error);
            }
        });
        
        // Update DOM elements
        this.updateDOMElements(notifier_id, value);
        
        // Dispatch custom event
        this.dispatchStateChangeEvent(notifier_id, value, old_value);
    }
    
    sendWebSocketMessage(message) {
        if (this.websocket && this.websocket.readyState === WebSocket.OPEN) {
            this.websocket.send(JSON.stringify(message));
            return true;
        }
        return false;
    }
    
    // Subscription Management
    subscribe(notifierId, callback, pollingEndpoint = null) {
        if (!this.subscriptions.has(notifierId)) {
            this.subscriptions.set(notifierId, []);
        }
        
        this.subscriptions.get(notifierId).push(callback);
        
        // Store polling endpoint for fallback
        if (pollingEndpoint) {
            this.pollingEndpoints.set(notifierId, pollingEndpoint);
        }
        
        // Subscribe via WebSocket if connected
        if (!this.sendWebSocketMessage({
            type: 'subscribe',
            notifier_id: notifierId
        })) {
            // WebSocket not available, ensure polling is active
            this.startFallbackPolling();
        }
        
        this.log(`Subscribed to notifier: ${notifierId}`);
        
        return () => this.unsubscribe(notifierId, callback);
    }
    
    unsubscribe(notifierId, callback) {
        const callbacks = this.subscriptions.get(notifierId);
        if (callbacks) {
            const index = callbacks.indexOf(callback);
            if (index > -1) {
                callbacks.splice(index, 1);
                
                // If no more callbacks, unsubscribe completely
                if (callbacks.length === 0) {
                    this.subscriptions.delete(notifierId);
                    this.pollingEndpoints.delete(notifierId);
                    
                    this.sendWebSocketMessage({
                        type: 'unsubscribe',
                        notifier_id: notifierId
                    });
                }
            }
        }
        
        this.log(`Unsubscribed from notifier: ${notifierId}`);
    }
    
    // DOM Management
    updateDOMElements(notifierId, value) {
        const elements = this.domElements.get(notifierId) || [];
        
        elements.forEach(element => {
            try {
                this.updateSingleElement(element, value);
            } catch (error) {
                this.log('Error updating DOM element:', error);
            }
        });
    }
    
    updateSingleElement(element, value) {
        const updateType = element.getAttribute('data-value-listener-update') || 'content';
        
        switch (updateType) {
            case 'content':
                element.textContent = this.formatValue(value);
                break;
            case 'html':
                element.innerHTML = this.formatValue(value);
                break;
            case 'attribute':
                const attrName = element.getAttribute('data-value-listener-attr') || 'value';
                element.setAttribute(attrName, value);
                break;
            case 'class':
                const className = element.getAttribute('data-value-listener-class');
                if (className) {
                    element.classList.toggle(className, !!value);
                }
                break;
            case 'style':
                const styleProp = element.getAttribute('data-value-listener-style');
                if (styleProp) {
                    element.style[styleProp] = value;
                }
                break;
            case 'custom':
                // Trigger custom update event
                const customEvent = new CustomEvent('godin:valueUpdate', {
                    detail: { element, value }
                });
                element.dispatchEvent(customEvent);
                break;
            default:
                element.textContent = this.formatValue(value);
        }
        
        // Add updated class for CSS animations
        element.classList.add('godin-value-updated');
        setTimeout(() => {
            element.classList.remove('godin-value-updated');
        }, 300);
    }
    
    formatValue(value) {
        if (value === null || value === undefined) {
            return '';
        }
        
        if (typeof value === 'object') {
            return JSON.stringify(value);
        }
        
        return String(value);
    }
    
    registerDOMElement(notifierId, element) {
        if (!this.domElements.has(notifierId)) {
            this.domElements.set(notifierId, []);
        }
        
        const elements = this.domElements.get(notifierId);
        if (!elements.includes(element)) {
            elements.push(element);
        }
    }
    
    unregisterDOMElement(notifierId, element) {
        const elements = this.domElements.get(notifierId);
        if (elements) {
            const index = elements.indexOf(element);
            if (index > -1) {
                elements.splice(index, 1);
                
                if (elements.length === 0) {
                    this.domElements.delete(notifierId);
                }
            }
        }
    }
    
    // Fallback Polling
    startFallbackPolling() {
        if (!this.config.enableFallbackPolling || this.pollingTimer) {
            return;
        }
        
        this.log('Starting fallback polling');
        
        this.pollingTimer = setInterval(() => {
            this.pollForUpdates();
        }, this.config.pollingInterval);
    }
    
    stopFallbackPolling() {
        if (this.pollingTimer) {
            this.log('Stopping fallback polling');
            clearInterval(this.pollingTimer);
            this.pollingTimer = null;
        }
    }
    
    async pollForUpdates() {
        const promises = [];
        
        this.pollingEndpoints.forEach((endpoint, notifierId) => {
            promises.push(this.pollSingleNotifier(notifierId, endpoint));
        });
        
        try {
            await Promise.allSettled(promises);
        } catch (error) {
            this.log('Error during polling:', error);
        }
    }
    
    async pollSingleNotifier(notifierId, endpoint) {
        try {
            const response = await fetch(endpoint, {
                method: 'GET',
                headers: {
                    'Accept': 'application/json',
                    'Cache-Control': 'no-cache'
                }
            });
            
            if (!response.ok) {
                throw new Error(`HTTP ${response.status}`);
            }
            
            const data = await response.json();
            const currentValue = this.valueListeners.get(notifierId);
            
            // Check if value has changed
            if (JSON.stringify(data.value) !== JSON.stringify(currentValue)) {
                this.handleValueChange({
                    notifier_id: notifierId,
                    value: data.value,
                    old_value: currentValue,
                    timestamp: Date.now()
                });
            }
        } catch (error) {
            this.log(`Polling error for ${notifierId}:`, error);
        }
    }
    
    // DOM Scanning
    scanForValueListeners() {
        const elements = document.querySelectorAll('[data-value-listener-id]');
        
        elements.forEach(element => {
            const notifierId = element.getAttribute('data-value-listener-id');
            const pollingEndpoint = element.getAttribute('data-value-listener-endpoint');
            
            if (notifierId) {
                this.registerDOMElement(notifierId, element);
                
                // Auto-subscribe if not already subscribed
                if (!this.subscriptions.has(notifierId)) {
                    this.subscribe(notifierId, () => {
                        // Empty callback - DOM updates are handled separately
                    }, pollingEndpoint);
                }
            }
        });
    }
    
    setupMutationObserver() {
        const observer = new MutationObserver((mutations) => {
            mutations.forEach((mutation) => {
                mutation.addedNodes.forEach((node) => {
                    if (node.nodeType === Node.ELEMENT_NODE) {
                        // Check if the added node itself is a value listener
                        if (node.hasAttribute && node.hasAttribute('data-value-listener-id')) {
                            const notifierId = node.getAttribute('data-value-listener-id');
                            const pollingEndpoint = node.getAttribute('data-value-listener-endpoint');
                            this.registerDOMElement(notifierId, node);
                            
                            if (!this.subscriptions.has(notifierId)) {
                                this.subscribe(notifierId, () => {}, pollingEndpoint);
                            }
                        }
                        
                        // Check for value listeners in child nodes
                        const childListeners = node.querySelectorAll('[data-value-listener-id]');
                        childListeners.forEach(child => {
                            const notifierId = child.getAttribute('data-value-listener-id');
                            const pollingEndpoint = child.getAttribute('data-value-listener-endpoint');
                            this.registerDOMElement(notifierId, child);
                            
                            if (!this.subscriptions.has(notifierId)) {
                                this.subscribe(notifierId, () => {}, pollingEndpoint);
                            }
                        });
                    }
                });
                
                mutation.removedNodes.forEach((node) => {
                    if (node.nodeType === Node.ELEMENT_NODE) {
                        // Clean up removed value listeners
                        if (node.hasAttribute && node.hasAttribute('data-value-listener-id')) {
                            const notifierId = node.getAttribute('data-value-listener-id');
                            this.unregisterDOMElement(notifierId, node);
                        }
                        
                        // Clean up child value listeners
                        const childListeners = node.querySelectorAll('[data-value-listener-id]');
                        childListeners.forEach(child => {
                            const notifierId = child.getAttribute('data-value-listener-id');
                            this.unregisterDOMElement(notifierId, child);
                        });
                    }
                });
            });
        });
        
        observer.observe(document.body, {
            childList: true,
            subtree: true
        });
    }
    
    // Event Management
    dispatchStateChangeEvent(notifierId, newValue, oldValue) {
        const event = new CustomEvent('godin:stateChange', {
            detail: {
                notifierId,
                newValue,
                oldValue,
                timestamp: Date.now()
            }
        });
        
        document.dispatchEvent(event);
    }
    
    // Public API Methods
    getValue(notifierId) {
        return this.valueListeners.get(notifierId);
    }
    
    setValue(notifierId, value) {
        // Send value change to server
        this.sendWebSocketMessage({
            type: 'set_value',
            notifier_id: notifierId,
            value: value
        });
    }
    
    // Utility Methods
    log(...args) {
        if (this.config.debug) {
            console.log('[GodinStateManager]', ...args);
        }
    }
    
    // Cleanup
    destroy() {
        this.stopFallbackPolling();
        
        if (this.websocket) {
            this.websocket.close();
        }
        
        this.subscriptions.clear();
        this.valueListeners.clear();
        this.domElements.clear();
        this.pollingEndpoints.clear();
    }
}

// Initialize global instance
window.GodinStateManager = window.GodinStateManager || new GodinStateManager({
    debug: window.location.hostname === 'localhost' || window.location.hostname === '127.0.0.1'
});

// Export for module systems
if (typeof module !== 'undefined' && module.exports) {
    module.exports = GodinStateManager;
}